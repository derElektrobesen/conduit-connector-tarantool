package tarantool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/derElektroBesen/conduit-connector-tarantool/internal/tntlogger"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
	vshardrouter "github.com/tarantool/go-vshard-router/v2"
)

var (
	errUnsupportedOperation  = fmt.Errorf("operation isn't supported")
	errUnsupportedCollection = fmt.Errorf("unsupported collection")
	errBadKey                = fmt.Errorf("unexpected key format")
)

// some value => bucket id.
type shardFunction func(string) uint64

type Destination struct {
	sdk.UnimplementedDestination

	config   DestinationConfig
	router   *vshardrouter.Router
	bucketID shardFunction
}

func NewDestination() sdk.Destination {
	// Create Destination and wrap it in the default middleware.
	return sdk.DestinationWithMiddleware(&Destination{})
}

func (d *Destination) Config() sdk.DestinationConfig {
	return &d.config
}

func (d *Destination) Open(ctx context.Context) error {
	// Open is called after Configure to signal the plugin it can prepare to
	// start writing records. If needed, the plugin should open connections in
	// this function.

	topology, err := d.config.makeTopology()
	if err != nil {
		return fmt.Errorf("unable to make tarantool cluster topology: %w", err)
	}

	router, err := vshardrouter.NewRouter(ctx, vshardrouter.Config{
		Loggerf: tntlogger.NewTntLogger(sdk.Logger(ctx)),

		DiscoveryTimeout: time.Minute,
		DiscoveryMode:    vshardrouter.DiscoveryModeOn,

		TotalBucketCount: d.config.TotalBuckets,

		User:             d.config.User,
		Password:         d.config.Password,
		TopologyProvider: topology,
	})
	if err != nil {
		return fmt.Errorf("unable to create vshard router: %w", err)
	}

	d.router = router

	d.bucketID = func(v string) uint64 {
		// default shard function
		return vshardrouter.BucketIDStrCRC32(v, router.BucketCount())
	}

	return nil
}

func (d *Destination) Write(ctx context.Context, recs []opencdc.Record) (int, error) {
	// Write writes len(r) records from r to the destination right away without
	// caching. It should return the number of records written from r
	// (0 <= n <= len(r)) and any error encountered that caused the write to
	// stop early. Write must return a non-nil error if it returns n < len(r).

	for i, rec := range recs {
		err := d.unmarshalRecord(ctx, &rec)
		if err != nil {
			return i, fmt.Errorf("can't unmarshal record: %w", err)
		}

		replicaset, err := d.replicaset(ctx, rec)
		if err != nil {
			return i, fmt.Errorf("can't get replicaset: %w", err)
		}

		data, err := d.makeTarantoolTuple(rec)
		if err != nil {
			return i, fmt.Errorf("can't make tarantool tuple: %w", err)
		}

		collection, err := rec.Metadata.GetCollection()
		if err != nil {
			return i, fmt.Errorf("can't understand collection name: %w", err)
		}

		switch rec.Operation {
		case opencdc.OperationSnapshot, opencdc.OperationCreate:
			if err := d.insertRecord(ctx, replicaset, collection, data); err != nil {
				return i, err
			}
		case opencdc.OperationUpdate:
			if err := d.updateRecord(ctx, rec); err != nil {
				return i, err
			}
		case opencdc.OperationDelete:
			if err := d.deleteRecord(ctx, rec); err != nil {
				return i, err
			}
		default:
			return i, errUnsupportedOperation
		}
	}

	return len(recs), nil
}

func (d *Destination) Teardown(ctx context.Context) error {
	replicasets := d.router.RouteAll()

	var err error

	lg := sdk.Logger(ctx)
	for name, r := range replicasets {
		lg.Debug().Str("replicaset", name).Msg("closing connections")
		p := r.Pooler()

		if errs := p.Close(); len(errs) > 0 {
			lg.Error().
				Errs("errors", errs).
				Str("replicaset", name).
				Msg("unable to close connection")

			err = errs[0]
		}
	}

	return err
}

func (d *Destination) getKey(data opencdc.Data) (string, error) {
	key, err := d.structuredDataFormatter(data)
	if err != nil {
		return "", fmt.Errorf("unable to parse data: %w", err)
	}

	if len(key) > 1 {
		return "", fmt.Errorf("key %+v: %w", key, errBadKey)
	}

	for _, v := range key {
		return fmt.Sprintf("%v", v), nil
	}

	return "", fmt.Errorf("no key found: %w", errBadKey)
}

func (d *Destination) getBucketID(ctx context.Context, rec opencdc.Record) (uint64, error) {
	collection, err := rec.Metadata.GetCollection()
	if err != nil {
		return 0, fmt.Errorf("unable to get record collection: %w", err)
	}

	key, err := d.getKey(rec.Key)
	if err != nil && !errors.Is(err, errBadKey) {
		return 0, fmt.Errorf("unable to get key: %w", err)
	}

	col := d.config.Collections[collection]
	switch {
	case key == "" && col.GetShardingKey == nil:
		return 0, fmt.Errorf("unable to get sharding key: %w", err)

	case col.GetShardingKey != nil:
		key, err = col.GetShardingKey(rec)
		if err != nil {
			return 0, fmt.Errorf("bad sharding key: %w", err)
		}
	}

	if key == "" {
		return 0, fmt.Errorf("unable to get sharding key: %w", errBadKey)
	}

	return d.bucketID(key), nil
}

func (d *Destination) replicaset(ctx context.Context, rec opencdc.Record) (pool.Pooler, error) {
	bucket, err := d.getBucketID(ctx, rec)
	if err != nil {
		return nil, fmt.Errorf("unable to get bucket id: %w", err)
	}

	replicaset, err := d.router.Route(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("unable to find route for bucket %d: %w", bucket, err)
	}

	return replicaset.Pooler(), nil
}

func (d *Destination) unmarshalRecord(ctx context.Context, rec *opencdc.Record) error {
	var err error
	rec.Key, err = d.structuredDataFormatter(rec.Key)
	if err != nil {
		return fmt.Errorf("unable to unmarshal key: %w", err)
	}

	lg := sdk.Logger(ctx)
	e := lg.Trace().Any("key", rec.Key)

	defer e.Msg("got record")

	rec.Payload.After, err = d.structuredDataFormatter(rec.Payload.After)
	if err != nil {
		return fmt.Errorf("unable to unmarshal payload: %w", err)
	}

	e = e.Any("payload", rec.Payload.After)

	return nil
}

func (d *Destination) structuredDataFormatter(data opencdc.Data) (opencdc.StructuredData, error) {
	if data == nil {
		return opencdc.StructuredData{}, nil
	}
	if sdata, ok := data.(opencdc.StructuredData); ok {
		return sdata, nil
	}
	raw := data.Bytes()
	if len(raw) == 0 {
		return opencdc.StructuredData{}, nil
	}

	m := make(map[string]interface{})
	err := json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *Destination) makeTarantoolTuple(rec opencdc.Record) ([]any, error) {
	data, err := d.structuredDataFormatter(rec.Payload.After)
	return []any{data}, err
}

func (d *Destination) insertRecord(
	ctx context.Context,
	replicaset pool.Pooler,
	collection string,
	data []any,
) error {
	req := tarantool.NewInsertRequest(collection).
		Context(ctx).
		Tuple(data)

	_, err := replicaset.Do(req, pool.RW).Get()
	if err != nil {
		return fmt.Errorf("unable to insert tuple: %w", err)
	}

	return nil
}

func (d *Destination) updateRecord(ctx context.Context, rec opencdc.Record) error {
	return nil
}

func (d *Destination) deleteRecord(ctx context.Context, rec opencdc.Record) error {
	return nil
}

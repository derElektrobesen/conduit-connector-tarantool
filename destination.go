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

	var g = map[opencdc.Operation]func(opencdc.Record) (string, []any, error){
		opencdc.OperationSnapshot: makePutRequest,
		opencdc.OperationCreate:   makePutRequest,
		opencdc.OperationUpdate:   makePutRequest,
		opencdc.OperationDelete:   makeDeleteRequest,
	}

	for i, rec := range recs {
		generator := g[rec.Operation]
		if generator == nil {
			return i, errUnsupportedOperation
		}

		if err := d.modify(ctx, rec, generator); err != nil {
			return i, fmt.Errorf("unable to submit modification: %w", err)
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

func getKey(data opencdc.Data) (string, error) {
	key, err := formatData(data)
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

	key, err := getKey(rec.Key)
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

func formatData(data opencdc.Data) (opencdc.StructuredData, error) {
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

func makePutRequest(rec opencdc.Record) (string, []any, error) {
	// Following function should exists and
	// tarantool user shoud have enough privileges to call it.
	// That's required because there is no simple method to
	// create a tuple with specified format.
	//
	// function put_tuple(space, tuple_map)
	//     return box.space[space]:put(box.space[space]:frommap(tuple_map))
	// end

	tup, err := formatData(rec.Payload.After)
	return "put_tuple", []any{tup}, err
}

func makeDeleteRequest(rec opencdc.Record) (string, []any, error) {
	// Following function should exists and
	// tarantool user shoud have enough privileges to call it.
	// This function is just a sugar: it simplifies go-code and
	// doesn't cause any problems because of put_tuple function is required.
	//
	// function delete_tuple(space, key)
	//     return box.space[space]:delete(key)
	// end

	key, err := getKey(rec.Key)
	return "delete_tuple", []any{key}, err
}

func (d *Destination) modify(
	ctx context.Context,
	rec opencdc.Record,
	generate func(opencdc.Record) (string, []any, error),
) error {
	collection, err := rec.Metadata.GetCollection()
	if err != nil {
		return fmt.Errorf("can't understand collection name: %w", err)
	}

	tup := []any{collection}

	method, data, err := generate(rec)
	if err != nil {
		return fmt.Errorf("unable to generate tuple: %w", err)
	}

	bucketID, err := d.getBucketID(ctx, rec)
	if err != nil {
		return fmt.Errorf("unable to get bucket id: %w", err)
	}

	resp, err := d.router.Call(
		ctx,
		bucketID,
		vshardrouter.CallModeRW,
		method,
		append(tup, data...),
		vshardrouter.CallOpts{Timeout: time.Second},
	)

	if err != nil {
		return fmt.Errorf("unable to call router: %w", err)
	}

	res, err := resp.Get()
	if err != nil {
		return fmt.Errorf("unable to decode resp: %w", err)
	}

	sdk.Logger(ctx).Trace().Any("resp", res).Str("method", method).Msg("got response")

	return nil
}

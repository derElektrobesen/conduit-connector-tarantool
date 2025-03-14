package tarantool

import (
	"context"
	"fmt"
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/derElektroBesen/conduit-connector-tarantool/internal/tntlogger"
	vshardrouter "github.com/tarantool/go-vshard-router"
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
		return vshardrouter.BucketIDStrCRC32(v, router.RouterBucketCount())
	}

	return nil
}

func (d *Destination) Write(_ context.Context, _ []opencdc.Record) (int, error) {
	// Write writes len(r) records from r to the destination right away without
	// caching. It should return the number of records written from r
	// (0 <= n <= len(r)) and any error encountered that caused the write to
	// stop early. Write must return a non-nil error if it returns n < len(r).
	return 0, nil
}

func (d *Destination) Teardown(_ context.Context) error {
	// Teardown signals to the plugin that all records were written and there
	// will be no more calls to any other function. After Teardown returns, the
	// plugin should be ready for a graceful shutdown.
	return nil
}

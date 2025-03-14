package tarantool

import (
	"context"
	"fmt"
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/derElektroBesen/conduit-connector-tarantool/internal/tntlogger"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	vshardrouter "github.com/tarantool/go-vshard-router"
	"github.com/tarantool/go-vshard-router/providers/static"
)

// some value => bucket id
type shardFunction func(string) uint64

type Destination struct {
	sdk.UnimplementedDestination

	config   DestinationConfig
	router   *vshardrouter.Router
	bucketID shardFunction
}

type ReplicaConfig struct {
	Addr string `yaml:"addr"`
	Name string `yaml:"name"`
	UUID string `yaml:"uuid"`
}

type ReplicasetConfig struct {
	Name     string          `yaml:"name"`
	UUID     string          `yaml:"uuid"`
	Replicas []ReplicaConfig `yaml:"replicas"`
}

// XXX: extra spaces are required to correctly generate readme
type DestinationConfig struct {
	sdk.DefaultDestinationMiddleware

	// Config includes parameters that are the same in the source and destination.
	Config

	// replicasets specifies tarantool cluster configuration
	//
	// Example configuration:
	//
	//   replicasets: |-
	//
	//     - name: "replicaset_1"                             # optional
	//
	//       uuid: "ff69c808-039f-4478-9f28-27a487b3d1d3"     # optional
	//
	//       replicas:                                        # *required*
	//
	//         - addr: "127.0.0.1:1001"                       # *required*
	//
	//           name: "1_1"                                  # optional
	//
	//           uuid: "15299fc8-fb53-44dc-9a84-2332fad9687c" # optional
	//
	//         - addr: "127.0.0.1:1002"
	//
	//           name: "1_2"
	//
	//           uuid: "15299fc8-fb53-44dc-9a84-2332fad9687c"
	//
	//     - name: "replicaset_2"
	//
	//       uuid: "29bbdcc4-2aa8-475f-b660-4fdc20dd5052"
	//
	//       replicas:
	//
	//         - addr: "127.0.0.1:1003"
	//
	//           name: "2_1"
	//
	//         - addr: "127.0.0.1:1004"
	//
	//           name: "2_2"
	ReplicasetsYaml string             `json:"replicasets" validate:"required"`
	Replicasets     []ReplicasetConfig `json:"-"`

	// total_buckets specifies a number of buckets in tarantool cluster
	TotalBuckets uint64 `json:"total_buckets" validate:"required,greater-than=0"`

	// shard_function specifies a shard function to be used during sharding.
	// At now only default sharding function is supported.
	ShardFunction string `json:"shard_function" validate:"inclusion=default" default:"default"`

	// user specifies username to access tarantool
	User string `json:"user" validate:"required"`

	// password specifies a password to access tarantool
	Password string `json:"password" validate:"required"`
}

func (s *DestinationConfig) Validate(context.Context) error {
	err := yaml.Unmarshal([]byte(s.ReplicasetsYaml), &s.Replicasets)
	if err != nil {
		return fmt.Errorf("unable to parse replicasets configuration: %w", err)
	}

	if len(s.Replicasets) == 0 {
		return fmt.Errorf("no replicasets found")
	}

	for i, c := range s.Replicasets {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("unable to validate replicaset %d: %w", i, err)
		}
	}

	return nil
}

func (c ReplicasetConfig) Validate() error {
	if len(c.Replicas) == 0 {
		return fmt.Errorf("no replicas found")
	}

	for i, r := range c.Replicas {
		if r.Addr == "" {
			return fmt.Errorf("no address found for replica %d", i)
		}
	}

	return nil
}

func NewDestination() sdk.Destination {
	// Create Destination and wrap it in the default middleware.
	return sdk.DestinationWithMiddleware(&Destination{})
}

func (d *Destination) Config() sdk.DestinationConfig {
	return &d.config
}

func (d *Destination) makeTopology() (vshardrouter.TopologyProvider, error) {
	topology := make(map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo)

	for i, r := range d.config.Replicasets {
		id, err := uuid.Parse(r.UUID)
		if err != nil {
			return nil, fmt.Errorf("bad UUID in replicaset %d: %q. %w",
				i, r.UUID, err)
		}

		replicasetInfo := vshardrouter.ReplicasetInfo{
			Name: r.Name,
			UUID: id,
		}

		replicas := make([]vshardrouter.InstanceInfo, 0, len(r.Replicas))
		for j, instance := range r.Replicas {
			id, err := uuid.Parse(instance.UUID)
			if err != nil {
				return nil, fmt.Errorf("bad UUID in instance %d in replicaset %d: %q. %w",
					j, i, instance.UUID, err)
			}

			replicas = append(replicas, vshardrouter.InstanceInfo{
				Name: instance.Name,
				Addr: instance.Addr,
				UUID: id,
			})
		}

		topology[replicasetInfo] = replicas
	}

	return static.NewProvider(topology), nil
}

func (d *Destination) Open(ctx context.Context) error {
	// Open is called after Configure to signal the plugin it can prepare to
	// start writing records. If needed, the plugin should open connections in
	// this function.

	topology, err := d.makeTopology()
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

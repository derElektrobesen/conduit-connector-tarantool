package tarantool

import (
	"context"
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	vshardrouter "github.com/tarantool/go-vshard-router"
	"github.com/tarantool/go-vshard-router/providers/static"
	"gopkg.in/yaml.v2"
)

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

// XXX: extra spaces are required to correctly generate readme.
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

var (
	errNoReplicas    = fmt.Errorf("no replicas")
	errNoReplicasets = fmt.Errorf("no replicasets found")
	errNoReplicaAddr = fmt.Errorf("replica addr not set")
)

func (c *DestinationConfig) Validate(context.Context) error {
	err := yaml.Unmarshal([]byte(c.ReplicasetsYaml), &c.Replicasets)
	if err != nil {
		return fmt.Errorf("unable to parse replicasets configuration: %w", err)
	}

	if len(c.Replicasets) == 0 {
		return errNoReplicasets
	}

	for i, c := range c.Replicasets {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("unable to validate replicaset %d: %w", i, err)
		}
	}

	return nil
}

func (c ReplicasetConfig) Validate() error {
	if len(c.Replicas) == 0 {
		return errNoReplicas
	}

	for i, r := range c.Replicas {
		if r.Addr == "" {
			return fmt.Errorf("%w for replica %d", errNoReplicaAddr, i)
		}
	}

	return nil
}

func (c *DestinationConfig) makeTopology() (vshardrouter.TopologyProvider, error) {
	topology := make(map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo)

	for i, r := range c.Replicasets {
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

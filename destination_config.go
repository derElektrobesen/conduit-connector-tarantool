package tarantool

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/conduitio/conduit-commons/opencdc"
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

type ShardKeyFn func(opencdc.Record) (string, error)

type CollectionConfig struct {
	// sharding_key specifies a field used to calculate bucket ID.
	// Key could be a text template or a static value.
	// Static value could be used if there is only one tarantool shard
	// present.
	//
	// See [Referencing record fields](https://conduit.io/docs/using/processors/referencing-fields)
	// for more info.
	ShardingKey    string     `json:"sharding_key" validate:"required" default="{{ .Key }}"`
	GetShardingKey ShardKeyFn `json:"-"`
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
	//     - name: "replicaset_1"                             # required
	//
	//       uuid: "ff69c808-039f-4478-9f28-27a487b3d1d3"     # required
	//
	//       replicas:                                        # *required*
	//
	//         - addr: "127.0.0.1:1001"                       # *required*
	//
	//           name: "1_1"                                  # optional
	//
	//           uuid: "15299fc8-fb53-44dc-9a84-2332fad9687c" # required
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

	// total_buckets specifies a number of buckets in tarantool cluster.
	TotalBuckets uint64 `json:"total_buckets" validate:"required,greater-than=0"`

	// shard_function specifies a shard function to be used during sharding.
	// At now only default sharding function is supported.
	ShardFunction string `json:"shard_function" validate:"inclusion=default" default:"default"`

	// user specifies username to access tarantool.
	User string `json:"user" validate:"required"`

	// password specifies a password to access tarantool.
	Password string `json:"password" validate:"required"`

	// collection specifies per-collection configuration.
	Collections map[string]CollectionConfig `json:"collection" validate:"required"`
}

var (
	errNoReplicas         = fmt.Errorf("no replicas")
	errNoReplicasets      = fmt.Errorf("no replicasets found")
	errNoReplicaAddr      = fmt.Errorf("replica addr not set")
	errReplicasetNotNamed = fmt.Errorf("replicaset not named")
	errUUIDRequired       = fmt.Errorf("uuid is required")
)

func (c *DestinationConfig) Validate(ctx context.Context) error {
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

	for k, v := range c.Collections {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("unable to validate collection %q: %w", k, err)
		}
	}

	return c.DefaultDestinationMiddleware.Validate(ctx)
}

func (c *CollectionConfig) Validate() error {
	var err error
	c.GetShardingKey, err = c.shardKeyFunction()

	if err != nil {
		return fmt.Errorf("bad sharding key %q: %w", c.ShardingKey, err)
	}

	return nil
}

func (c ReplicasetConfig) Validate() error {
	if len(c.Replicas) == 0 {
		return errNoReplicas
	}

	if c.UUID == "" {
		return errUUIDRequired
	}

	if c.Name == "" {
		return errReplicasetNotNamed
	}

	for i, r := range c.Replicas {
		if r.Addr == "" {
			return fmt.Errorf("%w for replica %d", errNoReplicaAddr, i)
		}

		if r.UUID == "" {
			return fmt.Errorf("%w for replica %d", errUUIDRequired, i)
		}
	}

	return nil
}

// shardKeyFunction returns a function that determines the sharding key for each record individually.
func (c *CollectionConfig) shardKeyFunction() (f ShardKeyFn, err error) {
	// Not a template, i.e. it's a static table name
	if !strings.HasPrefix(c.ShardingKey, "{{") && !strings.HasSuffix(c.ShardingKey, "}}") {
		return func(_ opencdc.Record) (string, error) {
			return c.ShardingKey, nil
		}, nil
	}

	// Try to parse the sharding key
	t, err := template.New("sharding_key").Funcs(sprig.FuncMap()).Parse(c.ShardingKey)
	if err != nil {
		return nil, fmt.Errorf("sharding key is neither a valid static value nor a valid Go template: %w", err)
	}

	var buf bytes.Buffer
	return func(r opencdc.Record) (string, error) {
		buf.Reset()
		if err := t.Execute(&buf, r); err != nil {
			return "", fmt.Errorf("failed to execute sharding key template: %w", err)
		}
		return buf.String(), nil
	}, nil
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

	p := static.NewProvider(topology)
	return p, p.Validate()
}

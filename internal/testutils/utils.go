package testutils

import (
	"context"
	"testing"

	mysql "github.com/conduitio-labs/conduit-connector-mysql"
	"github.com/conduitio-labs/conduit-connector-mysql/common"
	sdk "github.com/conduitio/conduit-connector-sdk"
	tarantool "github.com/derElektroBesen/conduit-connector-tarantool"
	"github.com/matryer/is"
	"github.com/rs/zerolog"
)

var replicasetsConfig = `
- name: storage_1
  uuid: cbf06940-0790-498b-948d-042b62cf3d29
  replicas:
    - addr: 127.0.0.1:3301
      name: storage_1_a
      uuid: 8a274925-a26d-47fc-9e1b-af88ce939412
    - addr: 127.0.0.1:3302
      name: storage_1_b
      uuid: 3de2e3e1-9ebe-4d0d-abb1-26d301b84633
    - addr: 127.0.0.1:3303
      name: storage_1_c
      uuid: 3de2e3e1-9ebe-4d0d-abb1-26d301b84635
- name: storage_2
  uuid: ac522f65-aa94-4134-9f64-51ee384f1a54
  replicas:
    - addr: 127.0.0.1:3311
      name: storage_2_a
      uuid: 1e02ae8a-afc0-4e91-ba34-843a356b8ed7
    - addr: 127.0.0.1:3312
      name: storage_2_b
      uuid: 001688c3-66f8-4a31-8e19-036c17d489c2
    - addr: 127.0.0.1:3313
      name: storage_2_c
      uuid: c23516d5-22de-4ef4-8918-73d52e7661e2
`

var mysqlDSN = "root:pass@tcp(127.0.0.1:3326)/mPOP?charset=utf8mb4&parseTime=True&loc=Local"

func NewDestination(ctx context.Context, is *is.I) (sdk.Destination, func()) {
	return NewDestinationWithConfig(ctx, is, NewDestinationConfig())
}

func NewDestinationWithConfig(
	ctx context.Context,
	is *is.I,
	config tarantool.DestinationConfig,
) (sdk.Destination, func()) {
	is.Helper()

	dest := tarantool.Destination{}
	cfg := dest.Config().(*tarantool.DestinationConfig)

	*cfg = config
	is.NoErr(cfg.Validate(ctx))
	is.NoErr(dest.Open(ctx))

	return &dest, func() {
		is.Helper()
		is.NoErr(dest.Teardown(ctx))
	}
}

func NewDestinationConfig() tarantool.DestinationConfig {
	return tarantool.DestinationConfig{
		ReplicasetsYaml: replicasetsConfig,
		TotalBuckets:    10000,
		ShardFunction:   "default",
		User:            "user",
		Password:        "pass",
	}
}

func NewSource(ctx context.Context, is *is.I, tables ...string) (sdk.Source, func()) {
	is.Helper()

	src := mysql.Source{}
	cfg := src.Config().(*common.SourceConfig)

	cfg.DSN = mysqlDSN
	cfg.Tables = tables

	is.NoErr(cfg.Validate(ctx))
	is.NoErr(src.Open(ctx, nil))

	return &src, func() {
		is.Helper()
		is.NoErr(src.Teardown(ctx))
	}
}

func NewContext(t *testing.T) context.Context {
	lg := zerolog.New(
		zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t)),
	).Level(zerolog.TraceLevel)
	return lg.WithContext(context.Background())
}

package tarantool_test

import (
	"context"
	"testing"
	"time"

	mysql "github.com/conduitio-labs/conduit-connector-mysql"
	"github.com/conduitio-labs/conduit-connector-mysql/common"
	"github.com/conduitio/conduit-commons/opencdc"
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

func newContext(t *testing.T) context.Context {
	lg := zerolog.New(
		zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t)),
	).Level(zerolog.TraceLevel)
	return lg.WithContext(context.Background())
}

func newDestination(ctx context.Context, is *is.I) (sdk.Destination, func()) {
	is.Helper()

	dest := tarantool.Destination{}
	cfg := dest.Config().(*tarantool.DestinationConfig)

	*cfg = tarantool.DestinationConfig{
		ReplicasetsYaml: replicasetsConfig,
		TotalBuckets:    10000,
		ShardFunction:   "default",
		User:            "user",
		Password:        "pass",
	}

	is.NoErr(cfg.Validate(ctx))
	is.NoErr(dest.Open(ctx))

	return &dest, func() {
		is.Helper()
		is.NoErr(dest.Teardown(ctx))
	}
}

func newSource(ctx context.Context, is *is.I, tables ...string) (sdk.Source, func()) {
	is.Helper()

	src := mysql.Source{}
	cfg := src.Config().(*common.SourceConfig)

	cfg.DSN = "root:pass@tcp(127.0.0.1:3326)/mPOP"
	cfg.Tables = tables

	is.NoErr(cfg.Validate(ctx))
	is.NoErr(src.Open(ctx, nil))

	return &src, func() {
		is.Helper()
		is.NoErr(src.Teardown(ctx))
	}
}

func TestWrite(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithTimeout(newContext(t), 5*time.Second)
	defer cancel()

	src, srcTeardown := newSource(ctx, is, "user")
	defer srcTeardown()

	dst, dstTeardown := newDestination(ctx, is)
	defer dstTeardown()

	rec, err := src.Read(ctx)
	is.NoErr(src.Ack(ctx, nil))
	is.NoErr(err)

	n, err := dst.Write(ctx, []opencdc.Record{rec})
	is.NoErr(err)
	is.Equal(1, n)
}

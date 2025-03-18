package tarantool_test

import (
	"os"
	"testing"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/derElektroBesen/conduit-connector-tarantool/internal/testutils"
	"github.com/matryer/is"
	"gorm.io/gorm"
)

var g *gorm.DB

func TestMain(t *testing.M) {
	var teardown func()
	g, teardown = testutils.Gorm()
	defer teardown()

	os.Exit(t.Run())
}

func TestWrite(t *testing.T) {
	is := is.New(t)
	ctx := testutils.NewContext(t)

	src, srcTeardown := testutils.NewSource(ctx, is, "user")
	defer srcTeardown()

	dst, dstTeardown := testutils.NewDestination(ctx, is)
	defer dstTeardown()

	u := testutils.NewUsers(10)
	g.Create(u)

	rec, err := src.Read(ctx)
	is.NoErr(src.Ack(ctx, nil))
	is.NoErr(err)

	n, err := dst.Write(ctx, []opencdc.Record{rec})
	is.NoErr(err)
	is.Equal(1, n)
}

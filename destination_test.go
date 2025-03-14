package tarantool_test

import (
	"context"
	"testing"

	tarantool "github.com/derElektroBesen/conduit-connector-tarantool"
	"github.com/matryer/is"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := tarantool.NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}

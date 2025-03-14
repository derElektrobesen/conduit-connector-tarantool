package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
	tarantool "github.com/derElektroBesen/conduit-connector-tarantool"
)

func main() {
	sdk.Serve(tarantool.Connector)
}

package main

import (
	tarantool "github.com/derElektroBesen/conduit-connector-tarantool"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func main() {
	sdk.Serve(tarantool.Connector)
}

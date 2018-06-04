package main

import (
	"fmt"
	"testing"

	"github.com/ToQoz/gopwt/assert"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

func TestServerStart(t *testing.T) {
	go launch(":8888")

	client := jsonrpc2.NewHTTPClient(fmt.Sprintf("http://%s/rpc", ":8888"))
	defer client.Close()

	var reply NameRes
	err := client.Call("LSPHandler.FullName", NameArg{"Takatoshi", "Matsumoto"}, &reply)
	assert.Require(t, err == nil)
	assert.OK(t, reply == "Takatoshi Matsumoto")
}

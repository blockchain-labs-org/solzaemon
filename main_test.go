package main

import (
	"context"
	"flag"
	"net"
	"os"
	"testing"

	"github.com/ToQoz/gopwt"
	"github.com/blockchain-labs-org/solzaemon/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

var (
	handler *lsp.Handler
	addr    string
)

func TestMain(m *testing.M) {
	flag.Parse()
	gopwt.Empower()

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	addr = l.Addr().String()
	l.Close()

	handler = lsp.NewHandler()
	go func() {
		err := launch(handler, addr)
		if err != nil {
			panic(err)
		}
	}()

	os.Exit(m.Run())
}

func dialServer(addr string) (*jsonrpc2.Conn, error) {
	conn, err := (&net.Dialer{}).Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), jsonrpc2.HandlerWithError(func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request) (interface{}, error) {
		// no-op
		return nil, nil
	})), nil
}

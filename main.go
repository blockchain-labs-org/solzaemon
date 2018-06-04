package main

import (
	"context"
	"log"
	"net"

	"github.com/blockchain-labs-org/solzaemon/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

var connOpt = []jsonrpc2.ConnOpt{}

func main() {
	log.Printf("listen :8080")
	err := launch(lsp.NewHandler(), ":8080")
	if err != nil {
		panic(err)
	}
}

func launch(handler *lsp.Handler, addr string) error {
	h := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
		return handler.Handle(ctx, conn, req)
	})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), h, connOpt...)
	}
}

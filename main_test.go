package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ToQoz/gopwt"
	"github.com/blockchain-labs-org/solzaemon/langserver"
	"github.com/sourcegraph/jsonrpc2"
)

var (
	handler *langserver.Handler
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

	handler = langserver.NewHandler()
	go func() {
		err := launch(handler, addr)
		if err != nil {
			panic(err)
		}
	}()
	if err := waitServer(addr); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func waitServer(addr string) error {
	wait := 0
	for {
		_, err := (&net.Dialer{}).Dial("tcp", addr)
		if err == nil {
			return nil
		}
		wait++
		if wait > 10 {
			return err
		}
		fmt.Println("waiting server...", err)
		time.Sleep(100 * time.Duration(wait) * time.Millisecond)
	}
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

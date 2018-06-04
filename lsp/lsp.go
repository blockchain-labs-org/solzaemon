package main

import (
	"log"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

type LSPHandler struct{}

type NameArg struct{ Fname, Lname string }
type NameRes string

// FullName concats first name and last name
func (*LSPHandler) FullName(t NameArg, res *NameRes) error {
	*res = NameRes(t.Fname + " " + t.Lname)
	return nil
}

func init() {
	rpc.Register(&LSPHandler{})
}

func main() {
	log.Printf("listen :8080")
	launch(":8080")
}

func launch(addr string) {
	http.Handle("/rpc", jsonrpc2.HTTPHandler(nil))
	http.ListenAndServe(addr, nil)
}

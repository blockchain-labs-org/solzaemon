package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	protocol "github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

type Handler struct {
	Mu   sync.Mutex
	Docs map[protocol.DocumentURI][]byte
}

func NewHandler() *Handler {
	return &Handler{
		Mu:   sync.Mutex{},
		Docs: map[protocol.DocumentURI][]byte{},
	}
}

func (h *Handler) Handle(ctx context.Context, conn jsonrpc2.JSONRPC2, req *jsonrpc2.Request) (result interface{}, err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("%v", perr)
		}
	}()
	switch req.Method {
	case "textDocument/didOpen":
		var params protocol.DidOpenTextDocumentParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		changed, err := do(params.TextDocument.URI, func() error {
			h.Mu.Lock()
			h.Docs[params.TextDocument.URI] = []byte(params.TextDocument.Text)
			h.Mu.Unlock()
			return nil
		})
		if changed {
			// clear cache
		}
		return params.TextDocument.URI, err
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	}
}

func do(uri protocol.DocumentURI, op func() error) (bool, error) {
	err := op()
	if err != nil {
		return true, err
	}
	return false, nil
}

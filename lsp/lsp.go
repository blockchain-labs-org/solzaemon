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
	case "initialize":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "initialized":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "shutdown":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "exit":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "$/cancelRequest":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/hover":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/definition":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/typeDefinition":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/xdefinition":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/references":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/implementation":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/documentSymbol":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/signatureHelp":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/formatting":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "workspace/symbol":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "workspace/xreferences":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/didOpen":
		var params protocol.DidOpenTextDocumentParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		changed, err := do(params.TextDocument.URI, func() error {
			h.setDocString(params.TextDocument.URI, params.TextDocument.Text)
			return nil
		})
		if changed {
			// clear cache
		}
		return params.TextDocument.URI, err
	case "textDocument/didChange":
		var params protocol.DidChangeTextDocumentParams
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
		return h.handleTextDocumentDidChange(params)
	case "textDocument/didClose":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	case "textDocument/didSave":
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
	}
}

func (h *Handler) setDocString(uri protocol.DocumentURI, doc string) {
	h.setDoc(uri, []byte(doc))
}

func (h *Handler) setDoc(uri protocol.DocumentURI, doc []byte) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	h.Docs[uri] = doc
}

func (h *Handler) getDoc(uri protocol.DocumentURI) ([]byte, bool) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	doc, found := h.Docs[uri]
	return doc, found
}

func do(uri protocol.DocumentURI, op func() error) (bool, error) {
	err := op()
	if err != nil {
		return true, err
	}
	return false, nil
}

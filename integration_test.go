package main

import (
	"context"
	"testing"

	"github.com/ToQoz/gopwt/assert"
	protocol "github.com/sourcegraph/go-langserver/pkg/lsp"
)

func TestTextDocument_didOpen(t *testing.T) {
	client, err := dialServer(addr)
	assert.Require(t, err == nil)

	defer client.Close()

	var reply protocol.DocumentURI
	params := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:  "uri",
			Text: "code",
		},
	}
	err = client.Call(context.Background(), "textDocument/didOpen", params, &reply)
	assert.Require(t, err == nil)
	assert.OK(t, reply == "uri")
	assert.OK(t, string(handler.Docs["uri"]) == "code")
}

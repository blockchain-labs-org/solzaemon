package lsp

import (
	"testing"

	"github.com/ToQoz/gopwt/assert"
	protocol "github.com/sourcegraph/go-langserver/pkg/lsp"
)

func TestHandleTextDocumentDidOpen(t *testing.T) {
	handler := NewHandler()
	params := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:  "code",
			Text: "func() {}",
		},
	}
	_, err := handler.handleTextDocumentDidOpen(params)
	assert.Require(t, err == nil)
	assert.OK(t, string(handler.Docs["code"]) == "func() {}")
}

func TestHandleTextDocumentDidChange(t *testing.T) {
	handler := NewHandler()
	handler.Docs["code"] = []byte(`func A() {
	fmt.Println("X")
	os.Exit(0)
}`)
	params := protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: "code",
			},
			Version: 0,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			protocol.TextDocumentContentChangeEvent{
				Range: &protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 0},
				},
				Text: "X",
			},
			protocol.TextDocumentContentChangeEvent{
				Range: &protocol.Range{
					Start: protocol.Position{Line: 0, Character: 1},
					End:   protocol.Position{Line: 0, Character: 1},
				},
				Text: "Y",
			},
			protocol.TextDocumentContentChangeEvent{
				Range: &protocol.Range{
					Start: protocol.Position{Line: 0, Character: 2},
					End:   protocol.Position{Line: 0, Character: 3},
				},
				Text: "L",
			},
			protocol.TextDocumentContentChangeEvent{
				Range: &protocol.Range{
					Start: protocol.Position{Line: 1, Character: 13},
					End:   protocol.Position{Line: 1, Character: 16},
				},
				Text: "^Y^",
			},
			protocol.TextDocumentContentChangeEvent{
				Range: &protocol.Range{
					Start: protocol.Position{Line: 2, Character: 1},
				},
				RangeLength: 2,
				Text:        "myos",
			},
		},
	}
	_, err := handler.handleTextDocumentDidChange(params)
	assert.Require(t, err == nil)

	assert.OK(t, string(handler.Docs["code"]) == `XYLunc A() {
	fmt.Println(^Y^)
	myos.Exit(0)
}`)
}

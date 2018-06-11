package langserver

import (
	"bytes"
	"fmt"

	"github.com/blockchain-labs-org/solzaemon/parser"
	"github.com/blockchain-labs-org/solzaemon/token"
	protocol "github.com/sourcegraph/go-langserver/pkg/lsp"
)

func (h *Handler) handleTextDocumentDefinition(params protocol.TextDocumentPositionParams) ([]protocol.Location, error) {
	contents, found := h.getDoc(params.TextDocument.URI)
	if !found {
		return nil, fmt.Errorf("received textDocument/definition for unknown file %q", params.TextDocument.URI)
	}

	f := token.NewFile()
	p, err := parser.Parse(f, []rune(string(contents)))
	if err != nil {
		panic(err) // FIXME
	}

	d, err := definition(p, token.Pos(f.Offset(params.Position.Line, params.Position.Character)))
	if err != nil {
		panic(err)
	}
	loc := protocol.Location{
		URI: "code",
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      f.Line(int(d)),
				Character: f.Character(int(d)),
			},
		},
	}
	locs := []protocol.Location{loc}
	return locs, nil
}

func (h *Handler) handleTextDocumentDidOpen(params protocol.DidOpenTextDocumentParams) (protocol.DocumentURI, error) {
	h.setDocString(params.TextDocument.URI, params.TextDocument.Text)
	return params.TextDocument.URI, nil
}

func (h *Handler) handleTextDocumentDidChange(params protocol.DidChangeTextDocumentParams) (protocol.DocumentURI, error) {
	contents, found := h.getDoc(params.TextDocument.URI)
	if !found {
		return params.TextDocument.URI, fmt.Errorf("received textDocument/didChange for unknown file %q", params.TextDocument.URI)
	}

	// The content changes descibe single state changes to the document.
	// So if there are two content changes c1 and c2 for a document in state S10 then c1 move the document to S11 and c2 to S12.
	// https://github.com/Microsoft/language-server-protocol/commit/fcb32f98317a1d37c798ca7309bb42ad8749d81d
	for _, change := range params.ContentChanges {
		sp := change.Range.Start
		ep := change.Range.End

		line := 0
		col := 0
		soffset := 0
		eoffset := 0
		startFound := false
		rangelen := int(change.RangeLength)
		for _, b := range contents {
			if line == sp.Line && col == sp.Character {
				startFound = true
			}
			if !startFound {
				if (line == sp.Line && col > sp.Character) || line > sp.Line {
					return params.TextDocument.URI, fmt.Errorf("received textDocument/didChange for invalid start position %#v on %s", sp, params.TextDocument.URI)
				}
				soffset++
			} else {
				if rangelen != 0 || (line == ep.Line && col == ep.Character) {
					eoffset += soffset + rangelen
					if soffset < 0 || eoffset >= len(contents) || eoffset < soffset {
						return params.TextDocument.URI, fmt.Errorf("received textDocument/didChange for out of range %#v on %s", change.Range, params.TextDocument.URI)
					}

					b := &bytes.Buffer{}
					b.Grow(soffset + len(change.Text) + len(contents) - eoffset)
					b.Write(contents[:soffset])
					b.WriteString(change.Text)
					b.Write(contents[eoffset:])
					contents = b.Bytes()
					break
				}
				if (line == ep.Line && col > ep.Character) || line > ep.Line {
					return params.TextDocument.URI, fmt.Errorf("received textDocument/didChange for invalid end position %#v on %s", ep, params.TextDocument.URI)
				}

				eoffset++
			}

			if b == '\n' {
				line++
				col = 0
			} else {
				col++
			}
		}
	}

	h.setDoc(params.TextDocument.URI, contents)
	return params.TextDocument.URI, nil
}

package lsp

import (
	"reflect"
	"testing"

	"github.com/ToQoz/gopwt/assert"
	protocol "github.com/sourcegraph/go-langserver/pkg/lsp"
)

func TestHandleTextDocumentDefinition(t *testing.T) {
	handler := NewHandler()
	handler.Docs["code"] = []byte(`pragma solidity ^0.4.23;
import "../token/ERC20/StandardToken.sol";

contract SimpleToken is StandardToken {
	string public constant name = "SimpleToken";
	string public constant symbol = "SIM";
	uint8 public constant decimals = 18;

	uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(decimals));

	constructor() public {
		totalSupply_ = INITIAL_SUPPLY;
		balances[msg.sender] = INITIAL_SUPPLY;
		emit Transfer(0x0, msg.sender, INITIAL_SUPPLY);
	}
}`)
	params := protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "code",
		},
		Position: protocol.Position{
			Line:      11,
			Character: 17,
		},
	}
	locs, err := handler.handleTextDocumentDefinition(params)
	assert.Require(t, err == nil)
	assert.Require(t, len(locs) == 1)
	assert.OK(t, locs[0].URI == "code")
	assert.OK(t, reflect.DeepEqual(locs[0].Range.Start, protocol.Position{Line: 8, Character: 25}))
	assert.OK(t, reflect.DeepEqual(locs[0].Range.End, protocol.Position{Line: 8, Character: 25}))
}

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

package scanner

import (
	"testing"

	"github.com/ToQoz/gopwt/assert"
	"github.com/blockchain-labs-org/solzaemon/token"
)

func TestScan(t *testing.T) {
	{
		// pragma solidity ^0.4.23
		s := NewScanner([]rune(`pragma solidity ^0.4.23;
import "../token/ERC20/StandardToken.sol";

contract SimpleToken is StandardToken {
	string public constant name = "SimpleToken";
	string public constant symbol = "SIM";
	uint8 public constant decimals = 18;

	uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(decimals));

	constructor() public {
		totalSupply_ = INITIAL_SUPPLY;
		balances[msg.sender] = INITIAL_SUPPLY;
	}
}`))
		expected := []struct {
			tok token.Token
			lit string
		}{
			{token.IDENT, `pragma`}, {token.IDENT, `solidity`}, {token.XOR, `^`}, {token.INT, `0`}, {token.PERIOD, `.`}, {token.INT, `4`}, {token.PERIOD, `.`}, {token.INT, `23`}, {token.SEMICOLON, `;`},
			{token.IDENT, `import`}, {token.STRING, `"../token/ERC20/StandardToken.sol"`}, {token.SEMICOLON, `;`},
			{token.IDENT, `contract`}, {token.IDENT, `SimpleToken`}, {token.IDENT, `is`}, {token.IDENT, `StandardToken`}, {token.LBRACE, `{`},
			{token.IDENT, `string`}, {token.IDENT, `public`}, {token.IDENT, `constant`}, {token.IDENT, `name`}, {token.ASSIGN, `=`}, {token.STRING, `"SimpleToken"`}, {token.SEMICOLON, `;`},
			{token.IDENT, `string`}, {token.IDENT, `public`}, {token.IDENT, `constant`}, {token.IDENT, `symbol`}, {token.ASSIGN, `=`}, {token.STRING, `"SIM"`}, {token.SEMICOLON, `;`},
			{token.IDENT, `uint8`}, {token.IDENT, `public`}, {token.IDENT, `constant`}, {token.IDENT, `decimals`}, {token.ASSIGN, `=`}, {token.INT, `18`}, {token.SEMICOLON, `;`},
			{token.IDENT, `uint256`}, {token.IDENT, `public`}, {token.IDENT, `constant`}, {token.IDENT, `INITIAL_SUPPLY`}, {token.ASSIGN, `=`},
			/* - */ {token.INT, `10000`}, {token.MUL, "*"}, {token.LPAREN, "("}, {token.INT, "10"}, {token.POW, "**"},
			/* - */ {token.IDENT, "uint256"}, {token.LPAREN, "("}, {token.IDENT, "decimals"}, {token.RPAREN, ")"}, {token.RPAREN, ")"}, {token.SEMICOLON, `;`},
			{token.IDENT, `constructor`}, {token.LPAREN, "("}, {token.RPAREN, ")"}, {token.IDENT, `public`}, {token.LBRACE, `{`},
			{token.IDENT, `totalSupply_`}, {token.ASSIGN, "="}, {token.IDENT, "INITIAL_SUPPLY"}, {token.SEMICOLON, `;`},
			{token.IDENT, `balances`}, {token.LBRACK, "["}, {token.IDENT, "msg"}, {token.PERIOD, `.`}, {token.IDENT, "sender"}, {token.RBRACK, "]"}, {token.ASSIGN, "="}, {token.IDENT, "INITIAL_SUPPLY"}, {token.SEMICOLON, `;`},
			{token.RBRACE, `}`},
		}
		for _, e := range expected {
			tok, lit := s.Scan()
			assert.Require(t, lit == e.lit)
			assert.Require(t, tok == e.tok)
		}
	}
}

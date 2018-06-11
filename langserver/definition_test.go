package langserver

import (
	"testing"

	"github.com/ToQoz/gopwt/assert"
	"github.com/blockchain-labs-org/solzaemon/parser"
	"github.com/blockchain-labs-org/solzaemon/token"
)

func TestDefinition_StateVarToStateVar(t *testing.T) {
	f := token.NewFile()
	got, err := parser.Parse(f, []rune(`pragma solidity ^0.4.23;
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

	def, err := definition(got, token.Pos(f.Offset(8, len(`	uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(d`))))
	assert.Require(t, err == nil)
	assert.OK(t, f.Line(int(def)) == 7)
	assert.OK(t, f.Character(int(def)) == len(`	uint8 public constant d`))
}

func TestDefinitin_FuncBodyToStateVar(t *testing.T) {
	f := token.NewFile()
	got, err := parser.Parse(f, []rune(`pragma solidity ^0.4.23;
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

	def, err := definition(got, token.Pos(f.Offset(11, len(`		totalSupply_ = I`))))
	assert.Require(t, err == nil)
	assert.OK(t, f.Line(int(def)) == 8)
	assert.OK(t, f.Character(int(def)) == len(`	uint256 public constant I`))
}

func TestDefinition_FuncBodyToFuncDef(t *testing.T) {
	f := token.NewFile()
	got, err := parser.Parse(f, []rune(`pragma solidity ^0.4.23;
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

	function a() public {
		b();
	}

	function b() public {
	}
}`))

	def, err := definition(got, token.Pos(f.Offset(16, len(`		b`))))
	assert.Require(t, err == nil)
	assert.OK(t, f.Line(int(def)) == 19)
	assert.OK(t, f.Character(int(def)) == len(`	function b`))
}

func TestDefinition_FuncLocalVar(t *testing.T) {
	f := token.NewFile()
	got, err := parser.Parse(f, []rune(`pragma solidity ^0.4.23;
import "../token/ERC20/StandardToken.sol";

contract SimpleToken is StandardToken {
	string public constant name = "SimpleToken";
	string public constant symbol = "SIM";
	uint8 public constant decimals = 18;
	uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(decimals));

	constructor() public {
		totalSupply_ = INITIAL_SUPPLY;
		totalSupply2_ = totalSupply_ * 2;
	}
}`))

	def, err := definition(got, token.Pos(f.Offset(12, len(`		totalSupply2_ = t`))))
	assert.Require(t, err == nil)
	assert.OK(t, f.Line(int(def)) == 11)
	assert.OK(t, f.Character(int(def)) == len(`		t`))
}

func TestDefinition_Contract(t *testing.T) {
	f := token.NewFile()
	got, err := parser.Parse(f, []rune(`pragma solidity ^0.4.23;
import "../token/ERC20/StandardToken.sol";

contract B is A {
}

contract A is StandardToken {
}`))

	def, err := definition(got, token.Pos(f.Offset(4, len(`contract B is A`))))
	assert.Require(t, err == nil)
	assert.OK(t, f.Line(int(def)) == 7)
	assert.OK(t, f.Character(int(def)) == len(`contract A`))
}

func TestDefinition_UndefinedVar(t *testing.T) {
	f := token.NewFile()
	got, err := parser.Parse(f, []rune(`pragma solidity ^0.4.23;
import "../token/ERC20/StandardToken.sol";

contract SimpleToken is StandardToken {
	string public constant name = "SimpleToken";
	string public constant symbol = "SIM";
	uint8 public constant decimals = 18;
	uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(decimals));

	constructor() public {
		totalSupply_ = INITIAL_SUPPLY;
		totalSupply2_ = totalSupplyUndefined_ * 2;
	}
}`))
	assert.Require(t, err == nil)

	{
		_, err := definition(got, token.Pos(f.Offset(12, len(`		totalSupply2_ = t`))))
		assert.Require(t, err.Error() == `definition of totalSupplyUndefined_ is not found in scope`)
	}
}

func TestDefinition_UnknownPosition(t *testing.T) {
	f := token.NewFile()
	got, err := parser.Parse(f, []rune(`pragma solidity ^0.4.23;
import "../token/ERC20/StandardToken.sol";

contract SimpleToken is StandardToken {
	string public constant name = "SimpleToken";
	string public constant symbol = "SIM";
	uint8 public constant decimals = 18;
	uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(decimals));

	constructor() public {
		totalSupply_ = INITIAL_SUPPLY;
		totalSupply2_ = totalSupplyUndefined_ * 2;
	}
}`))
	assert.Require(t, err == nil)

	{
		_, err := definition(got, token.Pos(f.Offset(50, 5)))
		assert.Require(t, err == unknownPosition)
	}
}

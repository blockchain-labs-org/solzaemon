package parser

import (
	"testing"

	"github.com/ToQoz/gopwt/assert"
	"github.com/blockchain-labs-org/solzaemon/ast"
	"github.com/blockchain-labs-org/solzaemon/token"
)

func TestParseERC20SimpleToken(t *testing.T) {
	got, err := Parse(token.NewFile(), []rune(`pragma solidity ^0.4.23;
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
	assert.Require(t, err == nil)
	// pragma
	assert.OK(t, got.PragmaDirective.Name.Name == "solidity")
	assert.OK(t, got.PragmaDirective.Value == "^0.4.23")
	// imports
	assert.Require(t, len(got.ImportDirectives) == 1)
	assert.OK(t, got.ImportDirectives[0].Path == `"../token/ERC20/StandardToken.sol"`)
	// contract
	assert.Require(t, len(got.ContractDefinition) == 1)
	// contract/state-vars
	// contract/state-vars/1
	assert.Require(t, len(got.ContractDefinition[0].StateVariableDeclarations) == 4)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[0].Typ.Name == "string")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[0].Visibility == "public")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[0].IsConstant == true)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[0].Name.Name == "name")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[0].Rhs.(*ast.BasicLit).Value == `"SimpleToken"`)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[0].Rhs.(*ast.BasicLit).Kind == token.STRING)
	// contract/state-vars/2
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[1].Typ.Name == "string")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[1].Visibility == "public")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[1].IsConstant == true)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[1].Name.Name == "symbol")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[1].Rhs.(*ast.BasicLit).Value == `"SIM"`)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[1].Rhs.(*ast.BasicLit).Kind == token.STRING)
	// contract/state-vars/3
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[2].Typ.Name == "uint8")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[2].Visibility == "public")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[2].IsConstant == true)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[2].Name.Name == "decimals")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[2].Rhs.(*ast.BasicLit).Kind == token.INT)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[2].Rhs.(*ast.BasicLit).Value == `18`)
	// contract/state-vars/4
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Typ.Name == "uint256")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Visibility == "public")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].IsConstant == true)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Name.Name == "INITIAL_SUPPLY")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Rhs.(*ast.BinaryExpr).X.(*ast.BasicLit).Kind == token.INT)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Rhs.(*ast.BinaryExpr).X.(*ast.BasicLit).Value == "10000")
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Rhs.(*ast.BinaryExpr).Op == token.MUL)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Rhs.(*ast.BinaryExpr).Y.(*ast.ParenExpr).X.(*ast.BinaryExpr).X.(*ast.BasicLit).Kind == token.INT)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Rhs.(*ast.BinaryExpr).Y.(*ast.ParenExpr).X.(*ast.BinaryExpr).X.(*ast.BasicLit).Value == "10")
	assert.Require(t, len(got.ContractDefinition[0].StateVariableDeclarations[3].Rhs.(*ast.BinaryExpr).Y.(*ast.ParenExpr).X.(*ast.BinaryExpr).Y.(*ast.CallExpr).Args) == 1)
	assert.OK(t, got.ContractDefinition[0].StateVariableDeclarations[3].Rhs.(*ast.BinaryExpr).Y.(*ast.ParenExpr).X.(*ast.BinaryExpr).Y.(*ast.CallExpr).Args[0].(*ast.Ident).Name == "decimals")
	// contract/function-defs
	assert.Require(t, len(got.ContractDefinition[0].FunctionDefinitions) == 1)
}

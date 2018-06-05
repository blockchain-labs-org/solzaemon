package ast

import (
	"github.com/blockchain-labs-org/solzaemon/token"
	protocol "github.com/sourcegraph/go-langserver/pkg/lsp"
)

type Node interface{}
type Stmt interface{}
type Expr interface{}

type Program struct {
	PragmaDirective    *PragmaDirective
	ImportDirectives   []*ImportDirective
	ContractDefinition []*ContractPart
}

type PragmaDirective struct {
	Name  *Ident
	Value string
}

type ImportDirective struct {
	Path protocol.DocumentURI
}

type ContractPart struct {
	StateVariableDeclarations []*StateVariableDeclaration
	FunctionDefinitions       []*FunctionDefinition
}

type StateVariableDeclaration struct {
	Name       *Ident
	Typ        *Ident
	Rhs        Expr
	IsConstant bool
	Visibility string
}

type FunctionDefinition struct {
	Name       *Ident
	Args       []Expr
	Visibility string
	Modifiers  []*Modifier
	Returns    Returns
	Block      []Stmt
}

type Modifier struct{}
type Returns struct {
	Typs []*Ident
}

type Ident struct {
	Name string
}

type BinaryExpr struct {
	X  Expr
	Op token.Token
	Y  Expr
}

type IndexExpr struct {
	X     Expr
	Index Expr
}

type SelectorExpr struct {
	X   Expr
	Sel Expr
}

type ParenExpr struct {
	Lparen token.Token
	X      Expr
	Rparen token.Token
}

type BasicLit struct {
	Kind  token.Token
	Value string
}

type AssignStmt struct {
	Lhs []Expr
	Rhs []Expr
}

type CallExpr struct {
	Fun    Expr
	Lparen token.Token
	Args   []Expr
	Rparen token.Token
}

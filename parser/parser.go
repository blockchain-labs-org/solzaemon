package parser

import (
	"fmt"

	"github.com/blockchain-labs-org/solzaemon/ast"
	"github.com/blockchain-labs-org/solzaemon/scanner"
	"github.com/blockchain-labs-org/solzaemon/token"
	protocol "github.com/sourcegraph/go-langserver/pkg/lsp"
)

func Parse(f *token.File, src []rune) (*ast.Program, error) {
	p := &Parser{}
	p.scanner = scanner.NewScanner(f, src)
	p.node = &ast.Program{}
	return p.parse()
}

type Parser struct {
	// DefPos to ReferencePos list
	scanner *scanner.Scanner
	node    ast.Node
	offset  token.Pos
	tok     token.Token
	lit     string
}

func (p *Parser) parse() (*ast.Program, error) {
	program := &ast.Program{}
	for {
		p.next()
		if p.tok == 0 {
			break
		}
		if p.tok == token.IDENT {
			switch p.lit {
			case "pragma":
				directive, err := p.parsePragma()
				if err != nil {
					return nil, err
				}
				program.PragmaDirective = directive
			case "import":
				imp, err := p.parseImport()
				if err != nil {
					return nil, err
				}
				program.ImportDirectives = append(program.ImportDirectives, imp)
			case "contract":
				contract, err := p.parseContract()
				if err != nil {
					return nil, err
				}
				program.ContractDefinition = append(program.ContractDefinition, contract)
			}
		}
	}
	return program, nil
}

func (p *Parser) parsePragma() (*ast.PragmaDirective, error) {
	p.next()
	name := p.lit
	p.next()
	val := ""
	for {
		if p.tok == token.SEMICOLON {
			break
		}
		val += p.lit
		p.next()
	}
	return &ast.PragmaDirective{
		Name:  &ast.Ident{Name: name, NamePos: p.offset},
		Value: val,
	}, nil
}

func (p *Parser) parseImport() (*ast.ImportDirective, error) {
	p.next()
	return &ast.ImportDirective{Path: protocol.DocumentURI(p.lit)}, nil
}

func (p *Parser) parseContract() (*ast.ContractPart, error) {
	part := &ast.ContractPart{}
	p.next()
	p.expect(token.IDENT)
	part.Name = p.parseIdent()

	if p.tok == token.IDENT && p.lit == "is" {
		p.next()
		for {
			part.Inherits = append(part.Inherits, &ast.Ident{Name: p.lit, NamePos: p.offset})
			p.next()
			if p.tok != token.COMMA {
				break
			}
		}
	}
	p.expect(token.LBRACE)

	for {
		p.next()
		if p.tok == 0 {
			return part, nil
		}

		switch p.tok {
		case token.RBRACE:
			return part, nil
		case token.IDENT:
			switch p.lit {
			case "constructor", "function":
				fn, err := p.parseFunction()
				if err != nil {
					return nil, err
				}
				part.FunctionDefinitions = append(part.FunctionDefinitions, fn)
			default:
				stateVar, err := p.parseStateVariable()
				if err != nil {
					return nil, err
				}
				part.StateVariableDeclarations = append(part.StateVariableDeclarations, stateVar)
			}
		}
	}
	return part, nil
}

func (p *Parser) parseStateVariable() (*ast.StateVariableDeclaration, error) {
	stateVar := &ast.StateVariableDeclaration{}
	stateVar.Typ = &ast.Ident{Name: p.lit, NamePos: p.offset}

done:
	for {
		p.next()
		switch p.tok {
		case token.IDENT:
			switch p.lit {
			case "constant":
				stateVar.IsConstant = true
			case "public", "internal", "private":
				stateVar.Visibility = p.lit
			default:
				stateVar.Name = &ast.Ident{Name: p.lit, NamePos: p.offset}
				break done
			}
		}
	}

	p.expectNext(token.ASSIGN)
	p.next()

	rhs := p.parseBinaryExpr()
	stateVar.Rhs = rhs

	return stateVar, nil
}

func (p *Parser) next() {
	p.offset, p.tok, p.lit = p.scanner.Scan()
}

func (p *Parser) parseFunction() (*ast.FunctionDefinition, error) {
	functionDef := &ast.FunctionDefinition{}
	if p.lit != "constructor" {
		p.next()
	}
	functionDef.Name = &ast.Ident{Name: p.lit, NamePos: p.offset}
	p.expectNext(token.LPAREN)
	// NOTE: args is not supported yet
	p.expectNext(token.RPAREN)

	p.next()
	if p.tok == token.IDENT {
		if p.lit == "public" || p.lit == "internal" || p.lit == "private" {
			functionDef.Visibility = p.lit
		}

		p.expectNext(token.LBRACE)
	} else {
		p.expect(token.LBRACE)
		functionDef.Visibility = "public"
	}

	// stmts
	for {
		p.next()
		expr := p.parseBinaryExpr()
		if expr == nil {
			break
		}
		functionDef.Block = append(functionDef.Block, expr)
	}

	return functionDef, nil
}

func (p *Parser) parseBinaryExpr() ast.Expr {
	if p.tok == '0' {
		return nil
	}
	if p.tok == token.SEMICOLON {
		p.next()
		return nil
	}

	x := p.parseUnaryExpr()
	if x == nil {
		return nil
	}
	for {
		if p.tok == '0' {
			return x
		}
		if p.tok == token.SEMICOLON {
			return x
		}
		switch p.tok {
		case token.ASSIGN, token.ADD, token.MUL, token.POW:
			// save tok/pos for op
			op := p.tok
			opPos := p.offset
			p.next()
			y := p.parseBinaryExpr()
			x = &ast.BinaryExpr{
				X:     x,
				Op:    op,
				OpPos: opPos,
				Y:     y,
			}
		default:
			return x
		}
	}
}

func (p *Parser) parseUnaryExpr() ast.Expr {
	// NOTE: unary expr is not supported yet
	return p.parsePrimaryExpr()
}

func (p *Parser) parsePrimaryExpr() ast.Expr {
	// NOTE: primary expr is not supported yet
	x := p.parseOperand()
	switch p.tok {
	case token.PERIOD:
		p.next()
		switch p.tok {
		case token.IDENT:
			x = p.parseSelector(x)
		}
	case token.LPAREN:
		x = p.parseCallOrConversion(x)
	case token.LBRACK:
		x = p.parseIndexExpr(x)
	}
	return x
}

func (p *Parser) parseCallOrConversion(x ast.Expr) ast.Expr {
	var args []ast.Expr
	p.expect(token.LPAREN)
	lparen := p.offset
	p.next()
	for p.tok != token.RPAREN && p.tok != 0 {
		args = append(args, p.parseUnaryExpr())
		p.next()
	}
	p.expect(token.RPAREN)
	rparen := p.offset
	return &ast.CallExpr{Fun: x, Lparen: lparen, Args: args, Rparen: rparen}
}

func (p *Parser) parseSelector(x ast.Expr) ast.Expr {
	sel := p.parseIdent()
	return &ast.SelectorExpr{X: x, Sel: sel}
}

func (p *Parser) parseIndexExpr(x ast.Expr) ast.Expr {
	idxExpr := &ast.IndexExpr{}
	idxExpr.X = x

	p.expect(token.LBRACK)
	idxExpr.Rbrack = p.offset
	p.next()
	idxExpr.Index = p.parseBinaryExpr()

	p.expect(token.RBRACK)
	idxExpr.Rbrack = p.offset
	p.next()

	return idxExpr
}

func (p *Parser) parseOperand() ast.Expr {
	switch p.tok {
	case token.IDENT:
		return p.parseIdent()
	case token.INT, token.STRING:
		return p.parseBasicLit()
	case token.LPAREN:
		p.next()
		x := p.parseBinaryExpr()
		p.expect(token.RPAREN)
		rparen := p.offset
		p.next()
		return &ast.ParenExpr{Lparen: p.offset, X: x, Rparen: rparen}
	case token.RBRACE: // hotfix
		return nil
	}
	panic("parseOperand: BadExpr: " + string(p.tok) + p.lit)
}

func (p *Parser) parseIdent() *ast.Ident {
	name := &ast.Ident{Name: p.lit, NamePos: p.offset}
	p.next()
	return name
}

func (p *Parser) parseBasicLit() ast.Expr {
	lit := &ast.BasicLit{Kind: p.tok, Value: p.lit, ValuePos: p.offset}
	p.next()
	return lit
}

func (p *Parser) expect(expect token.Token) (token.Token, error) {
	if p.tok != expect {
		return 0, fmt.Errorf("expect token %v but got %v", expect, p.tok)
	}
	return p.tok, nil
}

func (p *Parser) expectNext(expect token.Token) (token.Token, error) {
	p.next()
	return p.expect(expect)
}

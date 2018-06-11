package langserver

import (
	"errors"
	"fmt"

	"github.com/blockchain-labs-org/solzaemon/ast"
	"github.com/blockchain-labs-org/solzaemon/token"
)

var unknownPosition = errors.New("unknown position")

func definition(p *ast.Program, pos token.Pos) (token.Pos, error) {
	ret, found, err := newDefinitionFinder(p).lookup(pos)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, unknownPosition
	}

	return ret, nil
}

type scope struct {
	outer   *scope
	objects map[string]token.Pos
}

func (s *scope) lookup(name string) (token.Pos, bool) {
	ss := s
	for {
		if pos, ok := ss.objects[name]; ok {
			return pos, true
		}
		if ss.outer == nil {
			break
		}
		ss = ss.outer
	}
	return 0, false
}

func newScope(outer *scope) *scope {
	return &scope{outer: outer, objects: map[string]token.Pos{}}
}

type definitionFinder struct {
	scope *scope
	node  ast.Expr
}

func newDefinitionFinder(src ast.Expr) *definitionFinder {
	return &definitionFinder{
		scope: newScope(nil),
		node:  src,
	}
}

func (f *definitionFinder) lookup(pos token.Pos) (token.Pos, bool, error) {
	switch f.node.(type) {
	case *ast.Program:
		n := f.node.(*ast.Program)
		for _, def := range n.ContractDefinition {
			f.scope.objects[def.Name.Name] = def.Name.NamePos
		}

		for _, def := range n.ContractDefinition {
			for _, inherit := range def.Inherits {
				f.node = inherit
				ret, found, err := f.lookup(pos)
				if err != nil {
					return 0, false, err
				}
				if found {
					return ret, true, nil
				}
			}
			f.node = def

			f.scope = newScope(f.scope)
			ret, found, err := f.lookup(pos)
			if err != nil {
				return 0, false, err
			}
			if found {
				return ret, true, nil
			}
			f.scope = f.scope.outer
		}

		return 0, false, nil
	case *ast.ContractPart:
		n := f.node.(*ast.ContractPart)

		for _, def := range n.StateVariableDeclarations {
			f.scope.objects[def.Name.Name] = def.Name.NamePos
		}
		for _, def := range n.FunctionDefinitions {
			f.scope.objects[def.Name.Name] = def.Name.NamePos
		}

		for _, def := range n.StateVariableDeclarations {
			f.node = def
			ret, found, err := f.lookup(pos)
			if err != nil {
				return 0, false, err
			}
			if found {
				return ret, true, nil
			}
		}
		for _, def := range n.FunctionDefinitions {
			f.node = def
			f.scope = newScope(f.scope)
			ret, found, err := f.lookup(pos)
			if err != nil {
				return 0, false, err
			}
			if found {
				return ret, true, nil
			}
			f.scope = f.scope.outer
		}

		return 0, false, nil
	case *ast.StateVariableDeclaration:
		n := f.node.(*ast.StateVariableDeclaration)
		f.node = n.Name
		ret, found, err := f.lookup(pos)
		if err != nil {
			return 0, false, err
		}
		if found {
			return ret, true, nil
		}
		f.node = n.Rhs
		return f.lookup(pos)
	case *ast.FunctionDefinition:
		n := f.node.(*ast.FunctionDefinition)
		for _, stmt := range n.Block {
			f.node = stmt
			str, found, err := f.lookup(pos)
			if err != nil {
				return 0, false, err
			}
			if found {
				return str, found, nil
			}
		}

		return 0, false, nil
	case *ast.CallExpr:
		n := f.node.(*ast.CallExpr)
		f.node = n.Fun
		str, found, err := f.lookup(pos)
		if err != nil {
			return 0, false, err
		}
		if found {
			return str, true, nil
		}
		for _, n := range n.Args {
			f.node = n
			str, found, err := f.lookup(pos)
			if err != nil {
				return 0, false, err
			}
			if found {
				return str, true, nil
			}
		}
		return 0, false, nil
	case *ast.BinaryExpr:
		n := f.node.(*ast.BinaryExpr)
		if n.Op == token.ASSIGN {
			switch n.X.(type) {
			case *ast.Ident:
				name := n.X.(*ast.Ident)
				f.scope.objects[name.Name] = name.NamePos
			}
		}

		f.node = n.X
		str, found, err := f.lookup(pos)
		if err != nil {
			return 0, false, err
		}
		if found {
			return str, true, nil
		}
		f.node = n.Y
		return f.lookup(pos)
	case *ast.ParenExpr:
		n := f.node.(*ast.ParenExpr)
		f.node = n.X
		return f.lookup(pos)
	case *ast.IndexExpr:
		n := f.node.(*ast.IndexExpr)
		f.node = n.X
		str, found, err := f.lookup(pos)
		if err != nil {
			return 0, false, err
		}
		if found {
			return str, true, nil
		}
		f.node = n.Index
		return f.lookup(pos)
	case *ast.SelectorExpr:
		n := f.node.(*ast.SelectorExpr)
		f.node = n.X
		str, found, err := f.lookup(pos)
		if err != nil {
			return 0, false, err
		}
		if found {
			return str, true, nil
		}
		f.node = n.Sel
		return f.lookup(pos)
	case *ast.Ident:
		n := f.node.(*ast.Ident)
		if n.NamePos == pos {
			if ret, ok := f.scope.lookup(n.Name); ok {
				return ret, true, nil
			}
			return 0, false, fmt.Errorf("definition of %s is not found in scope", n.Name)
		}

		return 0, false, nil
	case *ast.BasicLit:
		// ignore unnamed
		return 0, false, nil
	default:
		fmt.Printf("%#v\n", f.node)
		panic("unexpected node")
	}
}

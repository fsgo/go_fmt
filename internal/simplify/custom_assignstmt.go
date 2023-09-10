// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/10

package simplify

import (
	"go/ast"
	"go/token"
)

type cAssignStmt struct {
	customBase
	Node *ast.AssignStmt
}

func (c *cAssignStmt) doFix() {
	c.numIncDec()
	c.chanReceive()
	c.mapRead()
}

// id+=1 -> id++
// id-=1 -> id--
func (c *cAssignStmt) numIncDec() {
	node := c.Node
	if len(node.Rhs) != 1 {
		return
	}
	var newTok token.Token

	switch node.Tok {
	case token.ADD_ASSIGN:
		newTok = token.INC
	case token.SUB_ASSIGN:
		newTok = token.DEC
	default:
		return
	}

	rh, ok := node.Rhs[0].(*ast.BasicLit)
	if !ok {
		return
	}
	if rh.Value != "1" {
		return
	}
	node.TokPos--
	node.Rhs = nil
	node.Tok = newTok
}

// _ = <-ch  ->  <-ch
func (c *cAssignStmt) chanReceive() {
	node := c.Node
	if node.Tok != token.ASSIGN || (len(node.Rhs) != 1 || len(node.Lhs) != 1) || len(node.Lhs) != 1 {
		return
	}
	x, ok1 := node.Lhs[0].(*ast.Ident)
	if !ok1 {
		return
	}
	if x.Name != "_" {
		return
	}
	y, ok2 := node.Rhs[0].(*ast.UnaryExpr)
	if !ok2 {
		return
	}
	if y.Op != token.ARROW {
		return
	}
	c1 := &ast.ExprStmt{
		X: y,
	}
	c.Cursor.Replace(c1)
}

// x, _ := someMap["key"] -> x:=someMap["key"]
func (c *cAssignStmt) mapRead() {
	node := c.Node
	if node.Tok != token.DEFINE || (len(node.Lhs) != 2 || len(node.Rhs) != 1) || len(node.Rhs) != 1 {
		return
	}
	x1, ok1 := node.Lhs[1].(*ast.Ident)
	if !ok1 || x1.Name != "_" {
		return
	}
	_, ok2 := node.Rhs[0].(*ast.IndexExpr)
	if !ok2 {
		return
	}
	node.Lhs = node.Lhs[0:1]
}

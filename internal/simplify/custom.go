// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// customSimplify 自定义的简化规则
func customSimplify(f *ast.File) {
	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		// log.Println("c.Name()", c.Name())
		switch vt := c.Node().(type) {
		case *ast.AssignStmt:
			newCustomApply(c).fixAssignStmt(vt)
		case *ast.BinaryExpr:
			newCustomApply(c).fixBinaryExpr(c, vt)
		}
		return true
	})
}

func newCustomApply(c *astutil.Cursor) *customApply {
	return &customApply{
		Cursor: c,
	}
}

type customApply struct {
	Cursor *astutil.Cursor
}

// id+=1 -> id++
// id-=1 -> id--
func (c *customApply) fixAssignStmt(vt *ast.AssignStmt) {
	if len(vt.Rhs) != 1 {
		return
	}
	var newTok token.Token

	switch vt.Tok {
	case token.ADD_ASSIGN:
		newTok = token.INC
	case token.SUB_ASSIGN:
		newTok = token.DEC
	default:
		return
	}

	rh, ok := vt.Rhs[0].(*ast.BasicLit)
	if !ok {
		return
	}
	if rh.Value != "1" {
		return
	}
	vt.TokPos--
	vt.Rhs = nil
	vt.Tok = newTok
}

func (c *customApply) fixBinaryExpr(cu *astutil.Cursor, cond *ast.BinaryExpr) {
	y, ok2 := cond.Y.(*ast.Ident)
	if !ok2 {
		return
	}

	var isTrue bool
	switch y.Name {
	case "true":
		isTrue = true
	case "false":
		isTrue = false
	default:
		return
	}

	setCond := func(c ast.Expr) {
		cu.Replace(c)
	}

	if cond.Op == token.EQL {
		if isTrue { // if b==true  -> if b
			setCond(cond.X)
		} else { // if b==false -> if !b
			c1 := &ast.UnaryExpr{
				Op: token.NOT,
				X:  cond.X,
			}
			setCond(c1)
		}
		return
	}
	if cond.Op == token.NEQ {
		if isTrue { // if b!=true -> if !b
			c1 := &ast.UnaryExpr{
				Op: token.NOT,
				X:  cond.X,
			}
			setCond(c1)
		} else { // if b!=false -> if b
			setCond(cond.X)
		}
	}
}

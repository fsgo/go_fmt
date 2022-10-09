// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"go/ast"
	"go/token"
)

var ins = &custom{}

// customSimplify 自定义的简化规则
func customSimplify(f *ast.File) {
	ast.Inspect(f, func(node ast.Node) bool {
		switch vt := node.(type) {
		case *ast.AssignStmt:
			ins.fixAssignStmt(vt)
		case *ast.IfStmt:
			ins.fixIfStmt(vt)
		}
		return true
	})
}

type custom struct {
}

// id+=1 -> id++
// id-=1 -> id--
func (c custom) fixAssignStmt(vt *ast.AssignStmt) {
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

func (c custom) fixIfStmt(vt *ast.IfStmt) {
	cond, ok := vt.Cond.(*ast.BinaryExpr)
	if !ok {
		return
	}
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

	if cond.Op == token.EQL {
		if isTrue { // if b==true  -> if b
			vt.Cond = cond.X
		} else { // if b==false -> if !b
			vt.Cond = &ast.UnaryExpr{
				Op: token.NOT,
				X:  cond.X,
			}
		}
		return
	}
	if cond.Op == token.NEQ {
		if isTrue { // if b!=true -> if !b
			vt.Cond = &ast.UnaryExpr{
				Op: token.NOT,
				X:  cond.X,
			}
		} else { // if b!=false -> if b
			vt.Cond = cond.X
		}
	}
}

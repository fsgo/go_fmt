// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/29

package simplify

import (
	"go/ast"
)

// isSimpleCondExpr 是否简单的 bool 表达式条件
func isSimpleCondExpr(e ast.Expr) bool {
	// if a==b
	if _, ok1 := e.(*ast.Ident); ok1 {
		return true
	}

	// if a 、if !a
	if ub, ok2 := e.(*ast.UnaryExpr); ok2 {
		switch ub.X.(type) {
		case *ast.Ident, // !a
			*ast.SelectorExpr: // !user.a
			return true
		default:
			return false
		}
	}

	// if user.ok
	if _, ok3 := e.(*ast.SelectorExpr); ok3 {
		return true
	}

	c, ok := e.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return (isExpVarBasic(c.X) || isSimpleCallExpr(c.X)) && (isExpVarBasic(c.Y) || isSimpleCallExpr(c.Y))
}

func isSimpleCallExpr(e ast.Expr) bool {
	// u.id > num 的 u.id
	// u.detail.id > num 的 u.detail.id
	// (*u).num >num 的 (*u).num
	if _, ok1 := e.(*ast.SelectorExpr); ok1 {
		return true
	}
	ce, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	fn, ok4 := ce.Fun.(*ast.Ident)
	if !ok4 || (fn.Name != "len" && fn.Name != "cap") || len(ce.Args) != 1 {
		return false
	}
	return isExpVarBasic(ce.Args[0])
}

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/26

package simplify

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func isIntExprValue(v ast.Expr) bool {
	if isBasicLitKind(v, token.INT) {
		return true
	}
	_, ok := v.(*ast.BasicLit) // fmt.Println(8)
	if ok {
		return true
	}
	cv, ok1 := v.(*ast.CallExpr) // fmt.Println(int16(8))
	if !ok1 {
		return false
	}
	fn, ok2 := cv.Fun.(*ast.Ident)
	if ok2 {
		switch fn.Name {
		case "int8",
			"int16",
			"int32",
			"int64",
			"uint8",
			"uint16",
			"uint32",
			"uint64":
			return true
		default:
			return false
		}
	}
	return false
}

func toBoolIdentExpr(v ast.Expr) (*ast.Ident, bool) {
	vi, ok := v.(*ast.Ident)
	if !ok {
		return nil, false
	}
	switch vi.Name {
	case "true", "false":
		return vi, true
	default:
		return nil, false
	}
}

func isBasicLitValue(n ast.Expr, kind token.Token, val string) bool {
	nv, ok := n.(*ast.BasicLit)
	if !ok {
		return false
	}
	return nv.Value == val && nv.Kind == kind
}

func isBasicLitKind(n ast.Expr, kind token.Token) bool {
	nv, ok := n.(*ast.BasicLit)
	if !ok {
		return false
	}
	return nv.Kind == kind
}

// isFunAny 判断是否任意的方法调用
// fnNames: regexp.Compile
func isFunAny(fn ast.Expr, fnNames ...string) bool {
	for i := 0; i < len(fnNames); i++ {
		info := strings.Split(fnNames[i], ".")
		if len(info) != 2 {
			panic(fmt.Sprintf("invalid fuc name %q, expect like x.y", fnNames[i]))
		}
		if isFun(fn, info[0], info[1]) {
			return true
		}
	}
	return false
}

func isFun(fn ast.Expr, pkg string, name string) bool {
	fun, ok2 := fn.(*ast.SelectorExpr)
	if !ok2 {
		return false
	}
	if fun.Sel.Name != name {
		return false
	}
	x, ok3 := fun.X.(*ast.Ident)
	return !(!ok3 || x.Name != pkg)
}

func isExpVarBasic(e ast.Expr) bool {
	_, ok1 := e.(*ast.BasicLit)
	if ok1 {
		return true
	}
	_, ok2 := e.(*ast.Ident)
	return ok2
}

// 将条件取非值
// 如  a > b  ->  a <= b
func conditionToNot(ne ast.Expr) ast.Expr {
	uc, ok1 := ne.(*ast.UnaryExpr)
	if ok1 {
		switch uc.Op {
		case token.NOT:
			return uc.X
		}
	}
	bc, ok2 := ne.(*ast.BinaryExpr)
	if ok2 {
		switch bc.Op {
		case token.EQL: // a == b => a != b
			bc.Op = token.NEQ
		case token.NEQ: // a != b => a == b
			bc.Op = token.EQL

		case token.GEQ: // a >= b => a < b
			bc.Op = token.LSS
		case token.GTR: // a > b => a <= b
			bc.Op = token.LEQ

		case token.LEQ: // a <= b => a > b
			bc.Op = token.GTR
		case token.LSS: // a < b => a  >= b
			bc.Op = token.GEQ
		default:
			goto end
		}
		return ne
	}

end:
	return &ast.UnaryExpr{
		Op: token.NOT,
		X:  ne,
	}
}

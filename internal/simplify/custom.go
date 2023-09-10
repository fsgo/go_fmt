// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/fsgo/go_fmt/internal/common"
)

// customSimplify 自定义的简化规则
func customSimplify(req *common.Request) {
	pre := func(c *astutil.Cursor) bool {
		base := customBase{
			req:    req,
			Cursor: c,
		}
		switch vt := c.Node().(type) {
		case *ast.BinaryExpr:
			(&cBinaryExpr{
				customBase: base,
				Node:       vt,
			}).doFix()
		case *ast.IfStmt:
			(&cIfStmt{
				customBase: base,
				Node:       vt,
			}).doFix()
		case *ast.ForStmt:
			(&cForStmt{
				customBase: base,
				Node:       vt,
			}).doFix()
		}
		return true
	}
	post := func(c *astutil.Cursor) bool {
		base := customBase{
			req:    req,
			Cursor: c,
		}
		// log.Println("c.Name()", c.Name())
		switch vt := c.Node().(type) {
		case *ast.AssignStmt:
			(&cAssignStmt{
				customBase: base,
				Node:       vt,
			}).doFix()
		case *ast.BinaryExpr:
			(&cBinaryExpr{
				customBase: base,
				Node:       vt,
			}).doFix()
		case *ast.CallExpr:
			(&cCallExpr{
				customBase: base,
				Node:       vt,
			}).doFix()
		case *ast.FuncDecl:
			(&cFuncDecl{
				customBase: base,
			}).fixFuncDecl(vt)
		case *ast.FuncLit:
			(&cFuncDecl{
				customBase: base,
			}).fixFuncLit(vt)
		}
		return true
	}
	astutil.Apply(req.AstFile, pre, post)
}

type customBase struct {
	req    *common.Request
	Cursor *astutil.Cursor
}

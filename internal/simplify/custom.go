// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xpasser"
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
		case *ast.BlockStmt:
			(&cBlockStmt{
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

func (cb customBase) isBasicKind(node ast.Expr, kind types.BasicKind) bool {
	vtp, err := xpasser.TypeOf(cb.req, node)
	if err != nil {
		return false
	}
	vb, ok := vtp.(*types.Basic)
	if !ok {
		return false
	}
	return vb.Kind() == kind
}

// isMinMaxNameAvailable 判断 min/max 在 node 所在作用域中是否可用。
//
// 当名称未被声明，或者名称解析到 Go 内置作用域中的同名标识符时，
// 返回 true；当名称被局部变量、常量、函数、类型等其他声明遮蔽时，
// 返回 false。
//
// 例如：
//   - min 未被声明：true（使用 Go 内置 min）
//   - min 是 Go 内置函数：true
//   - min 是局部变量：false
//   - min 是局部常量：false
//   - min 是自定义函数：false
func (cb customBase) isMinMaxNameAvailable(node ast.Node, name string) bool {
	if node == nil || name == "" {
		return false
	}
	switch name {
	case "min", "max":
	default:
		return false
	}

	info, err := xpasser.TypesInfo(cb.req, node)
	if err != nil || info == nil {
		return false
	}

	// 找到包含 node 的最内层作用域。
	var scope *types.Scope
	var minSpan token.Pos

	for n, s := range info.Scopes {
		if n.Pos() > node.Pos() || n.End() < node.End() {
			continue
		}
		span := n.End() - n.Pos()

		if scope == nil || span < minSpan {
			scope = s
			minSpan = span
		}
	}
	if scope == nil {
		return false
	}

	// 查找当前作用域及其外层作用域。
	_, obj := scope.LookupParent(name, node.Pos())
	return obj == nil || obj == types.Universe.Lookup(name)
}

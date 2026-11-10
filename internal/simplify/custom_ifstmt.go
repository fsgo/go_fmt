// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/10

package simplify

import (
	"go/ast"
	"go/token"
)

type cIfStmt struct {
	customBase
	Node *ast.IfStmt
}

func (c *cIfStmt) doFix() {
	if c.ifReturnMinMax1(true, false) {
		return
	}
	if c.ifReturnMinMax1(false, false) {
		return
	}
	if c.ifReturnMinMax1(true, true) {
		return
	}
	if c.ifReturnMinMax1(false, true) {
		return
	}
	c.ifReturnNoElse()
}

// 简化1：
//
//	if a<b {
//		return a
//	}else{
//		return b
//	}
//
// 为：
// return min(a,b)
//
// 简化2：
//
//	if a>b {
//		return a
//	}else{
//		return b
//	}
//
// 为：
// return max(a,b)
//
// 简化3:
//
//	if b<a {
//		return a
//	}else{
//		return b
//	}
//
// 为：
// return max(a,b)
//
// 简化4：
//
//	if b>a {
//		return a
//	}else{
//		return b
//	}
//
// 为：
// return min(a,b)
//
// 没有 else 的 return 的实现在 custom_blockstmt.go
//
// 参数：
// isMin: 是 min 还是 max
// reverse: 反向检查，即将 min 变为 max，max->min
func (c *cIfStmt) ifReturnMinMax1(isMin bool, reverse bool) bool {
	var opStr string
	var op token.Token
	if reverse {
		if isMin {
			opStr = "max"
			op = token.LSS
		} else {
			opStr = "min"
			op = token.GTR
		}
	} else {
		// 正常情况
		if isMin {
			opStr = "min"
			op = token.LSS
		} else {
			opStr = "max"
			op = token.GTR
		}
	}

	if !c.isMinMaxNameAvailable(c.Node, opStr) {
		return false
	}
	node := c.Node
	if node.Else == nil || len(node.Body.List) != 1 {
		return false
	}
	be, ok1 := node.Cond.(*ast.BinaryExpr)
	if !ok1 || be.Op != op {
		return false
	}
	var x, y ast.Expr
	if reverse {
		x = be.Y
		y = be.X
	} else {
		x = be.X
		y = be.Y
	}

	body, ok2 := node.Body.List[0].(*ast.ReturnStmt)
	if !ok2 || len(body.Results) != 1 || !isSameIdent(body.Results[0], x) {
		return false
	}

	el, ok3 := node.Else.(*ast.BlockStmt)
	if !ok3 || len(el.List) != 1 {
		return false
	}
	elr, ok4 := el.List[0].(*ast.ReturnStmt)
	if !ok4 || len(elr.Results) != 1 || !isSameIdent(elr.Results[0], y) {
		return false
	}

	nn := &ast.ReturnStmt{
		Return: elr.Return,
		Results: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.Ident{
					Name: opStr,
				},
				Args: []ast.Expr{
					be.X,
					be.Y,
				},
			},
		},
	}
	c.Cursor.Replace(nn)
	return true
}

// 简化 if else:
//
//	if cond{
//		  	// do something
//		  return // 必须存在
//	}else{
//			// do something
//	}
//
// 去掉 else 部分
func (c *cIfStmt) ifReturnNoElse() {
	node := c.Node

	if node.Else == nil || node.Body == nil || len(node.Body.List) == 0 {
		return
	}

	if node.Init != nil {
		//  if _,ok:=hello();ok {
		return
	}

	if !isBlockStmtReturn(node.Body) {
		return
	}

	if c.Cursor.Name() != "List" {
		return
	}

	curNode, ok3 := node.Else.(*ast.IfStmt)
	for ok3 && curNode.Cond != nil {
		if !isBlockStmtReturn(curNode.Body) {
			return
		}
		curNode, ok3 = curNode.Else.(*ast.IfStmt)
	}

	stElse, ok2 := node.Else.(*ast.BlockStmt)
	if !ok2 {
		return
	}

	for i := len(stElse.List) - 1; i >= 0; i-- {
		c.Cursor.InsertAfter(stElse.List[i])
	}
	node.Else = nil
}

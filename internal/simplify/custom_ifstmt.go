// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/10

package simplify

import "go/ast"

type cIfStmt struct {
	customBase
	Node *ast.IfStmt
}

func (c *cIfStmt) doFix() {
	c.ifReturnNoElse()
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

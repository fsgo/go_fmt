// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/10

package simplify

import (
	"go/ast"
	"go/token"
)

type cForStmt struct {
	customBase
	Node *ast.ForStmt
}

func (c *cForStmt) doFix() {
	c.loopBreak()
}

func (c *cForStmt) loopBreak() {
	node := c.Node

	// for 循环已经有条件
	// 如：for ok
	if node.Cond != nil {
		return
	}
	if node.Body == nil || len(node.Body.List) == 0 {
		return
	}

	first, ok1 := node.Body.List[0].(*ast.IfStmt)
	if !ok1 || first.Init != nil || first.Body == nil || len(first.Body.List) != 1 {
		return
	}

	ifb, ok2 := first.Body.List[0].(*ast.BranchStmt)
	if !ok2 || ifb.Tok != token.BREAK {
		return
	}

	node.Cond = conditionToNot(first.Cond)
	node.Body.List = node.Body.List[1:]
}

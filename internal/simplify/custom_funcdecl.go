// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/10

package simplify

import "go/ast"

type cFuncDecl struct {
	customBase
}

func (c *cFuncDecl) fixFuncDecl(node *ast.FuncDecl) {
	// c.shortReturnBool(node.Type, node.Body)
}

func (c *cFuncDecl) fixFuncLit(node *ast.FuncLit) {
	// c.shortReturnBool(node.Type, node.Body)
}

// 对应case: custom13.go.input
//
//nolint:gocyclo
//lint:ignore U1000 临时忽略
func (c *cFuncDecl) shortReturnBool(ft *ast.FuncType, funBody *ast.BlockStmt) {
	if funBody == nil || len(funBody.List) < 2 {
		// 函数已经很简单
		return
	}

	ret := ft.Results
	if ret == nil || len(ret.List) != 1 {
		return
	}
	// 函数的返回值只有一个并且是 bool 类型，如：
	// func ok() bool
	r1Type, ok := ret.List[0].Type.(*ast.Ident)
	if !ok || r1Type.Name != "bool" {
		return
	}
	st2If, ok2 := funBody.List[len(funBody.List)-2].(*ast.IfStmt)
	if !ok2 || st2If.Body == nil || len(st2If.Body.List) != 1 {
		return
	}

	if st2If.Init != nil {
		return
	}

	if _, ok21 := toBoolIdentExpr(st2If.Cond); ok21 {
		// 忽略：这种可能是临时调试代码
		// if true {}
		// if false {}
		return
	}

	st2rt, ok21 := st2If.Body.List[0].(*ast.ReturnStmt)
	if !ok21 || len(st2rt.Results) != 1 {
		return
	}

	st2Bi, ok22 := toBoolIdentExpr(st2rt.Results[0])
	if !ok22 {
		return
	}

	st1, ok1 := funBody.List[len(funBody.List)-1].(*ast.ReturnStmt)
	if !ok1 || len(st1.Results) != 1 {
		// 语法错误，忽略
		return
	}
	st1Bi, ok3 := toBoolIdentExpr(st1.Results[0])
	if !ok3 {
		return
	}

	if st2Bi.Name == st1Bi.Name {
		// 可能存在 bug (返回一样的值):
		//   if ok1(){
		//      return true
		//   }
		//   return true
		return
	}

	var cond ast.Expr

	if st2Bi.Name == "true" {
		// if cond(){
		// 	return true
		// }
		// return false
		cond = st2If.Cond
	} else {
		// if cond(){
		// 	return false
		// }
		// return true
		cond = conditionToNot(st2If.Cond)
	}
	stNew := &ast.ReturnStmt{
		Return: st2If.Pos(),
		Results: []ast.Expr{
			cond,
		},
	}
	funBody.List[len(funBody.List)-2] = stNew
	funBody.List = funBody.List[0 : len(funBody.List)-1]
}

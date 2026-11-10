package simplify

import (
	"go/ast"
	"go/token"
)

type cBlockStmt struct {
	customBase
	Node *ast.BlockStmt
}

func (c *cBlockStmt) doFix() {
	if c.ifReturnMinMax2(true, false) {
		return
	}
	if c.ifReturnMinMax2(false, false) {
		return
	}
	if c.ifReturnMinMax2(false, true) {
		return
	}
	if c.ifReturnMinMax2(true, true) {
		return
	}
}

// 简化1：
//
//	if a < b {
//		return a
//	}
//	return b
//
// 为：
// return min(a,b)
//
// 简化2：
//
//	if a > b {
//		return a
//	}
//	return b
//
// 为：
// return max(a,b)
func (c *cBlockStmt) ifReturnMinMax2(isMin bool, reverse bool) bool {
	node := c.Node
	if len(node.List) < 2 {
		return false
	}

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

	node1 := node.List[len(node.List)-1]
	node2 := node.List[len(node.List)-2]

	if !c.isMinMaxNameAvailable(node2, opStr) {
		return false
	}

	// 从最后一个节点开始
	tailReturn, ok1 := node1.(*ast.ReturnStmt)
	if !ok1 || len(tailReturn.Results) != 1 {
		return false
	}

	// 判断  if 语句:
	// if a<b{ return a}
	tail2, ok2 := node2.(*ast.IfStmt)
	if !ok2 || len(tail2.Body.List) != 1 || tail2.Else != nil || tail2.Init != nil {
		return false
	}

	be, ok3 := tail2.Cond.(*ast.BinaryExpr)
	if !ok3 || be.Op != op {
		return false
	}
	ifReturn, ok4 := tail2.Body.List[0].(*ast.ReturnStmt)
	if !ok4 || len(ifReturn.Results) != 1 {
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

	if !isSameIdent(ifReturn.Results[0], x) || !isSameIdent(tailReturn.Results[0], y) {
		return false
	}
	nn := node
	nn.List = nn.List[:len(nn.List)-2]

	nr := &ast.ReturnStmt{
		Return: tailReturn.Return,
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

	nn.List = append(nn.List, nr)
	c.Cursor.Replace(nn)
	return true
}

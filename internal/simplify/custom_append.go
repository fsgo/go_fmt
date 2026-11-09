package simplify

import (
	"go/ast"
)

// cCallAppend 修复 append 相关的
// 由  cCallExpr 调用
// 测试文件为 custom21.append.input
type cCallAppend struct {
	customBase
	Node *ast.CallExpr
}

func (c *cCallAppend) start() {
	c.nilSliceClone()
	c.emptySliceClone()
}

// append([]int(nil), a...) --> slice.Clone(a)
func (c *cCallAppend) nilSliceClone() {
	if len(c.Node.Args) != 2 {
		return
	}
	arg1, ok1 := c.Node.Args[0].(*ast.CallExpr)
	if !ok1 || len(arg1.Args) != 1 || !isIdentName(arg1.Args[0], "nil") {
		return
	}

	_, ok2 := arg1.Fun.(*ast.ArrayType)
	if !ok2 {
		return
	}
	nn := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "slices",
			},
			Sel: &ast.Ident{
				Name: "Clone",
			},
		},
		Args: []ast.Expr{
			c.Node.Args[1],
		},
	}
	c.Cursor.Replace(nn)

	pkgReplace(c.req.FSet, c.req.AstFile, "", "slices")
}

// append([]int{}, a...) --> slice.Clone(a)
func (c *cCallAppend) emptySliceClone() {
	if len(c.Node.Args) != 2 {
		return
	}
	arg1, ok1 := c.Node.Args[0].(*ast.CompositeLit)
	if !ok1 || len(arg1.Elts) > 0 {
		return
	}
	_, ok2 := arg1.Type.(*ast.ArrayType)
	if !ok2 {
		return
	}
	nn := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "slices",
			},
			Sel: &ast.Ident{
				Name: "Clone",
			},
		},
		Args: []ast.Expr{
			c.Node.Args[1],
		},
	}
	c.Cursor.Replace(nn)

	pkgReplace(c.req.FSet, c.req.AstFile, "", "slices")
}

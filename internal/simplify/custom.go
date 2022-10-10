// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// customSimplify 自定义的简化规则
func customSimplify(f *ast.File) {
	astutil.Apply(f, nil, func(c *astutil.Cursor) bool {
		// log.Println("c.Name()", c.Name())
		switch vt := c.Node().(type) {
		case *ast.AssignStmt:
			newCustomApply(f, c).fixAssignStmt(vt)
		case *ast.BinaryExpr:
			newCustomApply(f, c).fixBinaryExpr(vt)
		case *ast.CallExpr:
			newCustomApply(f, c).fixCallExpr(vt)
		}
		return true
	})
}

func newCustomApply(f *ast.File, c *astutil.Cursor) *customApply {
	return &customApply{
		Cursor: c,
		f:      f,
	}
}

type customApply struct {
	f      *ast.File
	Cursor *astutil.Cursor
}

// id+=1 -> id++
// id-=1 -> id--
func (c *customApply) fixAssignStmt(vt *ast.AssignStmt) {
	if len(vt.Rhs) != 1 {
		return
	}
	var newTok token.Token

	switch vt.Tok {
	case token.ADD_ASSIGN:
		newTok = token.INC
	case token.SUB_ASSIGN:
		newTok = token.DEC
	default:
		return
	}

	rh, ok := vt.Rhs[0].(*ast.BasicLit)
	if !ok {
		return
	}
	if rh.Value != "1" {
		return
	}
	vt.TokPos--
	vt.Rhs = nil
	vt.Tok = newTok
}

func (c *customApply) fixBinaryExpr(cond *ast.BinaryExpr) {
	c.trueFalse(cond)

	c.stringsCount0(cond)
	c.bytesCount0(cond)

	c.stringsIndex1(cond)
	c.bytesIndex1(cond)
}

func (c *customApply) trueFalse(cond *ast.BinaryExpr) {
	y, ok2 := cond.Y.(*ast.Ident)
	if !ok2 {
		return
	}

	var isTrue bool
	switch y.Name {
	case "true":
		isTrue = true
	case "false":
		isTrue = false
	default:
		return
	}

	setCond := func(e ast.Expr) {
		c.Cursor.Replace(e)
	}

	if cond.Op == token.EQL {
		if isTrue { // if b==true  -> if b
			setCond(cond.X)
		} else { // if b==false -> if !b
			c1 := &ast.UnaryExpr{
				Op: token.NOT,
				X:  cond.X,
			}
			setCond(c1)
		}
		return
	}
	if cond.Op == token.NEQ {
		if isTrue { // if b!=true -> if !b
			c1 := &ast.UnaryExpr{
				Op: token.NOT,
				X:  cond.X,
			}
			setCond(c1)
		} else { // if b!=false -> if b
			setCond(cond.X)
		}
	}
}

// strings.Count(s,"a") == 0   -> !strings.Contains(s,"a")
// strings.Count(s,"a") <= 0   -> !strings.Contains(s,"a")
// strings.Count(s,"a") < 1    -> !strings.Contains(s,"a")
// strings.Count(s,"a") > 0    -> strings.Contains(s,"a")
// strings.Count(s,"a") != 0   -> strings.Contains(s,"a")
func (c *customApply) stringsCount0(cond *ast.BinaryExpr) {
	c.stringsBytesCount0(cond, "strings")
}

func (c *customApply) bytesCount0(cond *ast.BinaryExpr) {
	c.stringsBytesCount0(cond, "bytes")
}

func (c *customApply) stringsBytesCount0(cond *ast.BinaryExpr, pkg string) {
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, pkg, "Count") {
		return
	}
	isVal0 := isBasicLit(cond.Y, token.INT, "0")
	isVal1 := !isVal0 && isBasicLit(cond.Y, token.INT, "1")
	if !isVal0 && !isVal1 {
		return
	}
	if !astutil.UsesImport(c.f, pkg) {
		return
	}

	fun := x.Fun.(*ast.SelectorExpr)

	// strings.Contains(s,"a")
	if (isVal0 && cond.Op == token.GTR) || // // strings.Count(s,"a") > 0
		(isVal0 && cond.Op == token.NEQ) { // strings.Count(s,"a") != 0
		fun.Sel.Name = "Contains"
		c.Cursor.Replace(cond.X)
		return
	}

	// !strings.Contains(s,"a")
	if (isVal0 && cond.Op == token.EQL) || // strings.Count(s,"a") == 0
		(isVal0 && cond.Op == token.LEQ) || // strings.Count(s,"a") <= 0
		(isVal1 && cond.Op == token.LSS) { //  strings.Count(s,"a") < 1
		fun.Sel.Name = "Contains"
		c1 := &ast.UnaryExpr{
			Op: token.NOT,
			X:  x,
		}
		c.Cursor.Replace(c1)
	}
}

// strings.Index(s,"a") == -1   ->  !strings.Contains(s,"a")
// strings.Index(s,"a") <= -1   ->  !strings.Contains(s,"a")
// strings.Index(s,"a") != -1   ->  strings.Contains(s,"a")
// strings.Index(s,"a") >  -1   ->  strings.Contains(s,"a")
// strings.Index(s,"a") >=  0   ->  strings.Contains(s,"a")
// strings.Index(s,"a") <   0   ->  !strings.Contains(s,"a")
func (c *customApply) stringsIndex1(cond *ast.BinaryExpr) {
	c.stringsBytesIndex1(cond, "strings")
}

func (c *customApply) bytesIndex1(cond *ast.BinaryExpr) {
	c.stringsBytesIndex1(cond, "bytes")
}

func (c *customApply) stringsBytesIndex1(cond *ast.BinaryExpr, pkg string) {
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, pkg, "Index") {
		return
	}

	//  判断的值是 0
	isVal0 := isBasicLit(cond.Y, token.INT, "0")

	//  判断的值是 -1
	var isValSub1 bool

	if !isVal0 {
		y, ok2 := cond.Y.(*ast.UnaryExpr)
		if !ok2 {
			return
		}
		//  判断的值是 -1
		isValSub1 = y.Op == token.SUB && isBasicLit(y.X, token.INT, "1")
	}

	if !isVal0 && !isValSub1 {
		return
	}

	if !astutil.UsesImport(c.f, pkg) {
		return
	}

	fun := x.Fun.(*ast.SelectorExpr)
	if (isValSub1 && cond.Op == token.NEQ) || // strings.Index(s,"a") != -1
		(isValSub1 && cond.Op == token.GTR) || // strings.Index(s,"a") >  -1
		(isVal0 && cond.Op == token.GEQ) { // strings.Index(s,"a") >=  0
		fun.Sel.Name = "Contains"
		c.Cursor.Replace(cond.X)
		return
	}

	if (isValSub1 && cond.Op == token.EQL) || //  strings.Index(s,"a") == -1
		(isValSub1 && cond.Op == token.LEQ) || // strings.Index(s,"a") <= -1
		(isVal0 && cond.Op == token.LSS) { // strings.Index(s,"a") <   0
		fun.Sel.Name = "Contains"
		c1 := &ast.UnaryExpr{
			Op: token.NOT,
			X:  x,
		}
		c.Cursor.Replace(c1)
		return
	}
}

func (c *customApply) fixCallExpr(node *ast.CallExpr) {
	c.stringsReplace(node)
}

// strings.Replace(s,"a","b",-1) -> strings.ReplaceAll(s,"a","b")
func (c *customApply) stringsReplace(node *ast.CallExpr) {
	if len(node.Args) != 4 {
		return
	}
	arg3, ok0 := node.Args[3].(*ast.UnaryExpr)
	if !ok0 || arg3.Op != token.SUB {
		return
	}
	if !isBasicLit(arg3.X, token.INT, "1") {
		return
	}
	if !astutil.UsesImport(c.f, "strings") {
		return
	}
	if !isFun(node.Fun, "strings", "Replace") {
		return
	}
	fun := node.Fun.(*ast.SelectorExpr)
	node.Args = node.Args[:3]
	fun.Sel.Name = "ReplaceAll"
}

func isBasicLit(n ast.Expr, kind token.Token, val string) bool {
	nv, ok := n.(*ast.BasicLit)
	if !ok {
		return false
	}
	return nv.Value == val && nv.Kind == kind
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
	if !ok3 || x.Name != pkg {
		return false
	}
	return true
}

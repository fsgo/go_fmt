// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/fsgo/go_fmt/internal/common"
)

// customSimplify 自定义的简化规则
func customSimplify(req *common.Request) {
	astutil.Apply(req.AstFile, nil, func(c *astutil.Cursor) bool {
		// log.Println("c.Name()", c.Name())
		switch vt := c.Node().(type) {
		case *ast.AssignStmt:
			newCustomApply(req, c).fixAssignStmt(vt)
		case *ast.BinaryExpr:
			newCustomApply(req, c).fixBinaryExpr(vt)
		case *ast.CallExpr:
			newCustomApply(req, c).fixCallExpr(vt)
		}
		return true
	})
}

func newCustomApply(req *common.Request, c *astutil.Cursor) *customApply {
	return &customApply{
		Cursor: c,
		req:    req,
	}
}

type customApply struct {
	req    *common.Request
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

	c.stringsCompare0(cond)
	c.bytesCompare0(cond)
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
	if !astutil.UsesImport(c.req.AstFile, pkg) {
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

	if !astutil.UsesImport(c.req.AstFile, pkg) {
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

// bytes.Compare(s,a) == 0 -> bytes.Equal(s,a)
// bytes.Compare(s,a) != 0 -> !bytes.Equal(s,a)
func (c *customApply) bytesCompare0(cond *ast.BinaryExpr) {
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, "bytes", "Compare") {
		return
	}
	if !isBasicLit(cond.Y, token.INT, "0") {
		return
	}

	if !astutil.UsesImport(c.req.AstFile, "bytes") {
		return
	}

	fun := x.Fun.(*ast.SelectorExpr)

	switch cond.Op {
	case token.EQL: // bytes.Compare(s,a) == 0
		fun.Sel.Name = "Equal"
		c.Cursor.Replace(cond.X)
	case token.NEQ: // bytes.Compare(s,a) != 0
		fun.Sel.Name = "Equal"
		c1 := &ast.UnaryExpr{
			Op: token.NOT,
			X:  x,
		}
		c.Cursor.Replace(c1)
	}
}

// strings.Compare(s,a) == 0 -> s==a
// strings.Compare(s,a) != 0 -> s!=a
func (c *customApply) stringsCompare0(cond *ast.BinaryExpr) {
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, "strings", "Compare") {
		return
	}
	if !isBasicLit(cond.Y, token.INT, "0") {
		return
	}

	switch cond.Op {
	case token.EQL,
		token.NEQ:
		c1 := &ast.BinaryExpr{
			X:  x.Args[0],
			Op: cond.Op,
			Y:  x.Args[1],
		}
		c.Cursor.Replace(c1)
		if !astutil.UsesImport(c.req.AstFile, "strings") {
			astutil.DeleteImport(c.req.FSet, c.req.AstFile, "strings")
		}
	default:
		return
	}
}

func (c *customApply) fixCallExpr(node *ast.CallExpr) {
	c.stringsReplace(node)
	c.timeNowSub(node)
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
	if !astutil.UsesImport(c.req.AstFile, "strings") {
		return
	}
	if !isFun(node.Fun, "strings", "Replace") {
		return
	}
	fun := node.Fun.(*ast.SelectorExpr)
	node.Args = node.Args[:3]
	fun.Sel.Name = "ReplaceAll"
}

// time.Now().Sub(n) -> time.Since(n)
func (c *customApply)timeNowSub(node *ast.CallExpr){
	if len(node.Args)!=1{
		return
	}
	x1, ok1 := node.Fun.(*ast.SelectorExpr)
	if !ok1 {
		return
	}
	if x1.Sel.Name!="Sub"{
		return
	}
	x2,ok2:=x1.X.(*ast.CallExpr)
	if !ok2{
		return
	}
	if !isFun(x2.Fun, "time", "Now") {
		return
	}
	if !astutil.UsesImport(c.req.AstFile, "time") {
		return
	}
	
	c1:=&ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "time",
			},
			Sel: &ast.Ident{
				Name: "Since",
			},
		},
		Args: node.Args,
	}
	c.Cursor.Replace(c1)
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

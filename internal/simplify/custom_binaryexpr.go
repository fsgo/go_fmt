// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/10

package simplify

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/fsgo/go_fmt/internal/xpasser"
)

type cBinaryExpr struct {
	customBase
	Node *ast.BinaryExpr
}

func (c *cBinaryExpr) doFix() {
	c.yadoCond0()
	c.yadoCond1()

	c.trueFalse()

	c.stringsCount0()
	c.bytesCount0()

	c.stringsIndex1()
	c.bytesIndex1()

	c.stringsCompare0()
	c.bytesCompare0()

	c.checkSliceNilLen()

	c.sortXY()
}

// "a" == val -> val == "a"
func (c *cBinaryExpr) yadoCond0() {
	cond := c.Node
	_, ok := cond.X.(*ast.BasicLit)
	if !ok {
		return
	}
	if _, ok2 := cond.Y.(*ast.BasicLit); ok2 {
		return
	}
	c.switchBinaryExprXY()
}

func (c *cBinaryExpr) switchBinaryExprXY() {
	cond := c.Node
	switch cond.Op {
	case token.EQL: // "a" == val
	// do nothing
	case token.NEQ: // "a" != val
	// do nothing
	case token.GEQ: // 1 >= val  -> val <= 1
		cond.Op = token.LEQ
	case token.LEQ: // 1 <= val  -> val >= 1
		cond.Op = token.GEQ
	case token.GTR: // 1 > val   -> val < 1
		cond.Op = token.LSS
	case token.LSS: // 1 < val  -> val > 1
		cond.Op = token.GTR
	default:
		return
	}

	x := cond.X

	cond.X = cond.Y
	cond.Y = x
}

// true == val -> val == true
func (c *cBinaryExpr) yadoCond1() {
	cond := c.Node
	x, ok := cond.X.(*ast.Ident)
	if !ok || (x.Name != "true" && x.Name != "false") {
		return
	}

	if _, ok2 := cond.Y.(*ast.BasicLit); ok2 {
		return
	}
	c.switchBinaryExprXY()
}

// if val==true   --> if val
// if val==false  --> if !val
// if val!=true   --> if !val
// if val!=false  --> if val
//
// testcase: custom1.go.input
func (c *cBinaryExpr) trueFalse() {
	cond := c.Node
	y, ok2 := cond.Y.(*ast.Ident)
	if !ok2 {
		return
	}

	// 判断 val 是否是 bool 类型的，如不是则不应该处理
	if vtp, err := xpasser.TypeOf(c.req, cond.X); err == nil {
		vb, ok3 := vtp.(*types.Basic)
		if !ok3 || vb.Info() != types.IsBoolean {
			return
		}
	} else {
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

// https://staticcheck.io/docs/checks#S1009
// if x != nil && len(x) != 0 {}
// =>
// if len(x) != 0 {}
//
//nolint:gocyclo
func (c *cBinaryExpr) checkSliceNilLen() {
	cond := c.Node
	if cond.Op != token.LAND {
		return
	}
	x, ok1 := cond.X.(*ast.BinaryExpr)
	// 检查是否  x != nil
	if !ok1 || x.Op != token.NEQ || !isIdentName(x.Y, "nil") {
		return
	}

	y, ok2 := cond.Y.(*ast.BinaryExpr)
	if !ok2 {
		return
	}
	yx, ok3 := y.X.(*ast.CallExpr)
	// 检查是否 len(x)
	if !ok3 || !isIdentName(yx.Fun, "len") || len(yx.Args) != 1 {
		return
	}
	// 检查是否同样的变量
	if !nameExprEq(x.X, yx.Args[0]) {
		return
	}

	// 检查其他的情况：
	// len(x) >= z、len(x) > z
	//  z >=0 即可
	checkCond := func() bool {
		if y.Op == token.GEQ || y.Op == token.GTR {
			v1, ok4 := y.Y.(*ast.BasicLit)
			if !ok4 || v1.Kind != token.INT {
				return false
			}
			vi, _ := strconv.Atoi(v1.Value)
			return vi >= 0
		}
		return false
	}

	if condYIs(y, token.NEQ, "0") || // !=0
		condYIs(y, token.GEQ, "0") || // >= 0
		condYIs(y, token.GTR, "0") ||
		checkCond() { // > 0
		c.Cursor.Replace(cond.Y)
	}
}

// strings.Count(s,"a") == 0   -> !strings.Contains(s,"a")
// strings.Count(s,"a") <= 0   -> !strings.Contains(s,"a")
// strings.Count(s,"a") < 1    -> !strings.Contains(s,"a")
// strings.Count(s,"a") > 0    -> strings.Contains(s,"a")
// strings.Count(s,"a") != 0   -> strings.Contains(s,"a")
func (c *cBinaryExpr) stringsCount0() {
	c.stringsBytesCount0("strings")
}

func (c *cBinaryExpr) bytesCount0() {
	c.stringsBytesCount0("bytes")
}

//nolint:gocyclo
func (c *cBinaryExpr) stringsBytesCount0(pkg string) {
	cond := c.Node
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, pkg, "Count") {
		return
	}
	isVal0 := isBasicLitValue(cond.Y, token.INT, "0")
	isVal1 := !isVal0 && isBasicLitValue(cond.Y, token.INT, "1")
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
func (c *cBinaryExpr) stringsIndex1() {
	c.stringsBytesIndex1("strings")
}

func (c *cBinaryExpr) bytesIndex1() {
	c.stringsBytesIndex1("bytes")
}

//nolint:gocyclo
func (c *cBinaryExpr) stringsBytesIndex1(pkg string) {
	cond := c.Node
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, pkg, "Index") {
		return
	}

	//  判断的值是 0
	isVal0 := isBasicLitValue(cond.Y, token.INT, "0")

	//  判断的值是 -1
	var isValSub1 bool

	if !isVal0 {
		y, ok2 := cond.Y.(*ast.UnaryExpr)
		if !ok2 {
			return
		}
		//  判断的值是 -1
		isValSub1 = y.Op == token.SUB && isBasicLitValue(y.X, token.INT, "1")
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
func (c *cBinaryExpr) bytesCompare0() {
	cond := c.Node
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, "bytes", "Compare") {
		return
	}
	if !isBasicLitValue(cond.Y, token.INT, "0") {
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

// 策略还存在瑕疵
const enableSortXY = false

// sortXY
//
//	ok1() && "a"=="b"       ->  "a"=="b" && ok1()
//	ok1() && len("a") > 0   ->  len("a") > 0 && ok1()
func (c *cBinaryExpr) sortXY() {
	cond := c.Node
	if !enableSortXY {
		return
	}
	if cond.Op != token.LAND && cond.Op != token.LOR {
		return
	}
	if isSimpleCondExpr(cond.X) { // 左边已经是简单条件，则跳过
		return
	}
	if !isSimpleCondExpr(cond.Y) { // 右边是复杂条件，则跳过
		return
	}
	ny := cond.X
	nx := cond.Y
	// todo 调整树结构
	if xv, ok := cond.X.(*ast.BinaryExpr); ok && xv.Op == cond.Op {
		nx = &ast.BinaryExpr{
			Op: cond.Op,
			X:  cond.Y,
			Y:  xv.X,
		}
		ny = xv.Y
	}
	cond.X = nx
	cond.Y = ny
}

// strings.Compare(s,a) == 0 -> s==a
// strings.Compare(s,a) != 0 -> s!=a
func (c *cBinaryExpr) stringsCompare0() {
	cond := c.Node
	x, ok1 := cond.X.(*ast.CallExpr)
	if !ok1 {
		return
	}
	if !isFun(x.Fun, "strings", "Compare") {
		return
	}
	if !isBasicLitValue(cond.Y, token.INT, "0") {
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

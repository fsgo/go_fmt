// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/fsgo/go_fmt/internal/common"
)

// customSimplify 自定义的简化规则
func customSimplify(req *common.Request) {
	pre := func(c *astutil.Cursor) bool {
		switch vt := c.Node().(type) {
		case *ast.BinaryExpr:
			newCustomApply(req, c).fixBinaryExpr(vt)
		case *ast.IfStmt:
			newCustomApply(req, c).fixIfStmt(vt)
		}
		return true
	}
	post := func(c *astutil.Cursor) bool {
		// log.Println("c.Name()", c.Name())
		switch vt := c.Node().(type) {
		case *ast.AssignStmt:
			newCustomApply(req, c).fixAssignStmt(vt)
		case *ast.BinaryExpr:
			newCustomApply(req, c).fixBinaryExpr(vt)
		case *ast.CallExpr:
			newCustomApply(req, c).fixCallExpr(vt)
		case *ast.FuncDecl:
			newCustomApply(req, c).fixFuncDecl(vt)
		case *ast.FuncLit:
			newCustomApply(req, c).fixFuncLit(vt)
		}
		return true
	}
	astutil.Apply(req.AstFile, pre, post)
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

func (c *customApply) fixAssignStmt(vt *ast.AssignStmt) {
	c.numIncDec(vt)
	c.chanReceive(vt)
	c.mapRead(vt)
}

// id+=1 -> id++
// id-=1 -> id--
func (c *customApply) numIncDec(vt *ast.AssignStmt) {
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

// _ = <-ch  ->  <-ch
func (c *customApply) chanReceive(node *ast.AssignStmt) {
	if node.Tok != token.ASSIGN || (len(node.Rhs) != 1 || len(node.Lhs) != 1) || len(node.Lhs) != 1 {
		return
	}
	x, ok1 := node.Lhs[0].(*ast.Ident)
	if !ok1 {
		return
	}
	if x.Name != "_" {
		return
	}
	y, ok2 := node.Rhs[0].(*ast.UnaryExpr)
	if !ok2 {
		return
	}
	if y.Op != token.ARROW {
		return
	}
	c1 := &ast.ExprStmt{
		X: y,
	}
	c.Cursor.Replace(c1)
}

// x, _ := someMap["key"] -> x:=someMap["key"]
func (c *customApply) mapRead(node *ast.AssignStmt) {
	if node.Tok != token.DEFINE || (len(node.Lhs) != 2 || len(node.Rhs) != 1) || len(node.Rhs) != 1 {
		return
	}
	x1, ok1 := node.Lhs[1].(*ast.Ident)
	if !ok1 || x1.Name != "_" {
		return
	}
	_, ok2 := node.Rhs[0].(*ast.IndexExpr)
	if !ok2 {
		return
	}
	node.Lhs = node.Lhs[0:1]
}

func (c *customApply) fixBinaryExpr(cond *ast.BinaryExpr) {
	c.yadoCond0(cond)
	c.yadoCond1(cond)

	c.trueFalse(cond)

	c.stringsCount0(cond)
	c.bytesCount0(cond)

	c.stringsIndex1(cond)
	c.bytesIndex1(cond)

	c.stringsCompare0(cond)
	c.bytesCompare0(cond)

	c.sortXY(cond)
}

// "a" == val -> val == "a"
func (c *customApply) yadoCond0(cond *ast.BinaryExpr) {
	_, ok := cond.X.(*ast.BasicLit)
	if !ok {
		return
	}
	if _, ok2 := cond.Y.(*ast.BasicLit); ok2 {
		return
	}
	c.switchBinaryExprXY(cond)
}

func (c *customApply) switchBinaryExprXY(cond *ast.BinaryExpr) {
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
func (c *customApply) yadoCond1(cond *ast.BinaryExpr) {
	x, ok := cond.X.(*ast.Ident)
	if !ok || (x.Name != "true" && x.Name != "false") {
		return
	}

	if _, ok2 := cond.Y.(*ast.BasicLit); ok2 {
		return
	}
	c.switchBinaryExprXY(cond)
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

//nolint:gocyclo
func (c *customApply) stringsBytesCount0(cond *ast.BinaryExpr, pkg string) {
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
func (c *customApply) stringsIndex1(cond *ast.BinaryExpr) {
	c.stringsBytesIndex1(cond, "strings")
}

func (c *customApply) bytesIndex1(cond *ast.BinaryExpr) {
	c.stringsBytesIndex1(cond, "bytes")
}

//nolint:gocyclo
func (c *customApply) stringsBytesIndex1(cond *ast.BinaryExpr, pkg string) {
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
func (c *customApply) bytesCompare0(cond *ast.BinaryExpr) {
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
func (c *customApply) sortXY(cond *ast.BinaryExpr) {
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
func (c *customApply) stringsCompare0(cond *ast.BinaryExpr) {
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

func (c *customApply) fixCallExpr(node *ast.CallExpr) {
	c.regexpRawString(node)

	c.stringsReplace(node)
	c.timeNowSub(node)
	c.timeSubNow(node)
	c.fmtErrorf(node)
	c.xPrintf(node)
}

// 提高正则的可读性
// see https://staticcheck.io/docs/checks#S1007
// regexp.Compile("\\A(\\w+) profile: total \\d+\\n\\z")
// ->
// regexp.Compile(`\A(\w+) profile: total \d+\n\z`)
func (c *customApply) regexpRawString(node *ast.CallExpr) {
	if !astutil.UsesImport(c.req.AstFile, "regexp") {
		return
	}

	if isFunAny(node.Fun, "regexp.Compile", "regexp.MustCompile") && len(node.Args) == 1 {
		c.nodeRawString(node.Args[0])
		return
	}
}

func (c *customApply) nodeRawString(node ast.Node) {
	ab, ok := node.(*ast.BasicLit)
	if !ok || ab.Kind != token.STRING {
		return
	}
	if !strings.HasSuffix(ab.Value, `"`) {
		// already a raw string
		return
	}
	if !strings.Contains(ab.Value, `\\`) {
		return
	}
	if strings.Contains(ab.Value, "`") {
		return
	}
	raw, err := strconv.Unquote(ab.Value)
	if err != nil {
		return
	}
	if raw == ab.Value {
		return
	}
	ab.Value = "`" + raw + "`"
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
	if !isBasicLitValue(arg3.X, token.INT, "1") {
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
func (c *customApply) timeNowSub(node *ast.CallExpr) {
	if len(node.Args) != 1 {
		return
	}
	x1, ok1 := node.Fun.(*ast.SelectorExpr)
	if !ok1 {
		return
	}
	if x1.Sel.Name != "Sub" {
		return
	}
	x2, ok2 := x1.X.(*ast.CallExpr)
	if !ok2 {
		return
	}
	if !isFun(x2.Fun, "time", "Now") {
		return
	}
	if !astutil.UsesImport(c.req.AstFile, "time") {
		return
	}

	c1 := &ast.CallExpr{
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

// t.Sub(time.Now()) -> time.Until(t)
func (c *customApply) timeSubNow(node *ast.CallExpr) {
	if len(node.Args) != 1 {
		return
	}
	arg, ok := node.Args[0].(*ast.CallExpr)
	if !ok {
		return
	}
	if !isFun(arg.Fun, "time", "Now") {
		return
	}
	fn, ok2 := node.Fun.(*ast.SelectorExpr)
	if !ok2 {
		return
	}
	if fn.Sel.Name != "Sub" {
		return
	}
	if !astutil.UsesImport(c.req.AstFile, "time") {
		return
	}
	c1 := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "time",
			},
			Sel: &ast.Ident{
				Name: "Until",
			},
		},
		Args: []ast.Expr{
			fn.X,
		},
	}
	c.Cursor.Replace(c1)
}

// fmt.Errorf("abc") -> errors.New("def")
func (c *customApply) fmtErrorf(node *ast.CallExpr) {
	if !isFun(node.Fun, "fmt", "Errorf") {
		return
	}
	if len(node.Args) != 1 {
		return
	}

	// 只处理 fmt.Errorf("abc")
	// 而不处理
	// var msg="abc"
	// fmt.Errorf(msg)
	arg, ok := node.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}
	if arg.Kind == token.STRING && strings.Contains(arg.Value, "%") {
		return
	}

	if !astutil.UsesImport(c.req.AstFile, "fmt") {
		return
	}

	c1 := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "errors",
			},
			Sel: &ast.Ident{
				Name: "New",
			},
		},
		Args: []ast.Expr{
			arg,
		},
	}
	c.Cursor.Replace(c1)
	pkgReplace(c.req.FSet, c.req.AstFile, "fmt", "errors")
}

func (c *customApply) xPrintf(node *ast.CallExpr) {
	c.xPrintfByPkg(node, "fmt", "Printf", "Print")
	c.xPrintfByPkg(node, "log", "Printf", "Print")

	c.xPrintfByPkg(node, "log", "Fatalf", "Fatal")
	c.xPrintfByPkg(node, "log", "Panicf", "Panic")

	c.xFprintfByPkg(node, "fmt", "Fprintf", "Fprint")
}

// fmt.Printf("abc") -> fmt.Print("abc")
// fmt.Printf("%s","abc") -> fmt.Print("abc")
// log.Printf("abc") -> log.Print("abc")
// log.Fatalf("abc") -> log.Fatal("abc")
// log.Panicf("abc") -> log.Panic("abc")
func (c *customApply) xPrintfByPkg(node *ast.CallExpr, pkg string, fnOld string, fnNew string) {
	if !isFun(node.Fun, pkg, fnOld) {
		return
	}

	if len(node.Args) == 1 {
		arg, ok := node.Args[0].(*ast.BasicLit)
		if !ok {
			return
		}
		if arg.Kind == token.STRING && strings.Contains(arg.Value, "%") {
			return
		}
		if !astutil.UsesImport(c.req.AstFile, pkg) {
			return
		}
		fn := node.Fun.(*ast.SelectorExpr)
		fn.Sel.Name = fnNew
		return
	}

	if len(node.Args) == 2 {
		arg0, ok := node.Args[0].(*ast.BasicLit)
		if !ok {
			return
		}
		if arg0.Kind != token.STRING {
			return
		}
		arg1 := node.Args[1]
		if (arg0.Value == `"%s"` && isBasicLitKind(arg1, token.STRING)) ||
			(arg0.Value == `"%d"` && isIntExprValue(arg1)) {
			fn := node.Fun.(*ast.SelectorExpr)
			fn.Sel.Name = fnNew
			node.Args = node.Args[1:]
			return
		}
		return
	}
}

// fmt.Fprintf(os.Stderr,"abc") -> fmt.Fprint(os.Stderr,"abc")
func (c *customApply) xFprintfByPkg(node *ast.CallExpr, pkg string, fnOld string, fnNew string) {
	if !isFun(node.Fun, pkg, fnOld) {
		return
	}
	if len(node.Args) != 2 {
		return
	}
	arg1, ok := node.Args[1].(*ast.BasicLit)
	if !ok {
		return
	}
	if arg1.Kind == token.STRING && strings.Contains(arg1.Value, "%") {
		return
	}
	if !astutil.UsesImport(c.req.AstFile, pkg) {
		return
	}
	fn := node.Fun.(*ast.SelectorExpr)
	fn.Sel.Name = fnNew
}

func (c *customApply) fixFuncDecl(node *ast.FuncDecl) {
	c.shortReturnBool(node.Type, node.Body)
}

func (c *customApply) fixFuncLit(node *ast.FuncLit) {
	c.shortReturnBool(node.Type, node.Body)
}

// 对应case: custom13.go.input
//
//nolint:gocyclo
func (c *customApply) shortReturnBool(ft *ast.FuncType, funBody *ast.BlockStmt) {
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
	if funBody == nil || len(funBody.List) < 2 {
		// 函数已经很简单
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

func (c *customApply) fixIfStmt(node *ast.IfStmt) {
	c.ifReturnNoElse(node)
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
func (c *customApply) ifReturnNoElse(node *ast.IfStmt) {
	if node.Else == nil {
		return
	}

	if node.Init != nil {
		//  if _,ok:=hello();ok {
		return
	}

	if c.Cursor.Name() != "List" {
		return
	}

	if node.Body == nil || len(node.Body.List) == 0 {
		return
	}

	_, ok1 := node.Body.List[len(node.Body.List)-1].(*ast.ReturnStmt)
	if !ok1 {
		// 必须是  if cond{  return}
		return
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

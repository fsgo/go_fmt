// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/10

package simplify

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/fsgo/go_fmt/internal/xpasser"
)

type cCallExpr struct {
	customBase
	Node *ast.CallExpr
}

func (c *cCallExpr) doFix() {
	c.regexpRawString()

	c.stringsReplace()
	c.timeNowSub()
	c.timeSubNow()

	c.errorsNewFmt()
	c.fmtErrorf()

	c.xPrintf()

	c.sortSlice()

	c.writeFmtSprintf()

	c.fmtSprintfInt()

	c.fmtSprintfStrings()
}

// 提高正则的可读性
// see https://staticcheck.io/docs/checks#S1007
// regexp.Compile("\\A(\\w+) profile: total \\d+\\n\\z")
// ->
// regexp.Compile(`\A(\w+) profile: total \d+\n\z`)
func (c *cCallExpr) regexpRawString() {
	if !astutil.UsesImport(c.req.AstFile, "regexp") {
		return
	}
	node := c.Node
	if isFunAny(node.Fun, "regexp.Compile", "regexp.MustCompile") && len(node.Args) == 1 {
		c.nodeRawString(node.Args[0])
		return
	}
}

func (c *cCallExpr) nodeRawString(node ast.Node) {
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
func (c *cCallExpr) stringsReplace() {
	node := c.Node
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
func (c *cCallExpr) timeNowSub() {
	node := c.Node
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
func (c *cCallExpr) timeSubNow() {
	node := c.Node
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
func (c *cCallExpr) fmtErrorf() {
	node := c.Node
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

// errors.New(fmt.Sprintf(...)) => fmt.Errorf(...)
// 测试用例：
// custom14.go.input
// custom15.go.input
func (c *cCallExpr) errorsNewFmt() {
	node := c.Node
	if !isFun(node.Fun, "errors", "New") {
		return
	}
	if len(node.Args) != 1 {
		return
	}

	a1, ok1 := node.Args[0].(*ast.CallExpr)
	if !ok1 {
		return
	}
	if len(a1.Args) == 0 {
		// 语法错误
		return
	}

	if !isFun(a1.Fun, "fmt", "Sprintf") {
		return
	}

	if len(a1.Args) == 1 {
		arg, ok := a1.Args[0].(*ast.BasicLit)
		if ok && arg.Kind == token.STRING && !strings.Contains(arg.Value, "%") {
			node.Args = a1.Args
			pkgReplace(c.req.FSet, c.req.AstFile, "fmt", "errors")
			return
		}
	}

	nn := &ast.CallExpr{
		Lparen: node.Lparen,
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "fmt",
			},
			Sel: &ast.Ident{
				Name: "Errorf",
			},
		},
		Args: a1.Args,
	}
	c.Cursor.Replace(nn)
	pkgReplace(c.req.FSet, c.req.AstFile, "errors", "fmt")
}

func (c *cCallExpr) xPrintf() {
	c.xPrintfByPkg("fmt", "Printf", "Print")
	c.xPrintfByPkg("log", "Printf", "Print")

	c.xPrintfByPkg("log", "Fatalf", "Fatal")
	c.xPrintfByPkg("log", "Panicf", "Panic")

	c.xFprintfByPkg("fmt", "Fprintf", "Fprint")
}

// fmt.Printf("abc") -> fmt.Print("abc")
// fmt.Printf("%s","abc") -> fmt.Print("abc")
// log.Printf("abc") -> log.Print("abc")
// log.Fatalf("abc") -> log.Fatal("abc")
// log.Panicf("abc") -> log.Panic("abc")
func (c *cCallExpr) xPrintfByPkg(pkg string, fnOld string, fnNew string) {
	node := c.Node
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
func (c *cCallExpr) xFprintfByPkg(pkg string, fnOld string, fnNew string) {
	node := c.Node
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

// sort.Sort(sort.StringSlice(x)) => sort.Strings(x)
func (c *cCallExpr) sortSlice() {
	node := c.Node
	if !isFun(node.Fun, "sort", "Sort") {
		return
	}
	if len(node.Args) != 1 {
		return
	}

	xfun, ok0 := node.Fun.(*ast.SelectorExpr)
	if !ok0 {
		return
	}

	arg, ok1 := node.Args[0].(*ast.CallExpr)
	if !ok1 {
		return
	}

	if isFun(arg.Fun, "sort", "StringSlice") {
		// sort.Sort(sort.StringSlice(x)) => sort.Strings(x)
		xfun.Sel.Name = "Strings"
		node.Args = arg.Args
		return
	}

	if isFun(arg.Fun, "sort", "Float64Slice") {
		// sort.Sort(sort.Float64Slice(x)) => sort.Float64s(x)
		xfun.Sel.Name = "Float64s"
		node.Args = arg.Args
		return
	}

	if isFun(arg.Fun, "sort", "IntSlice") {
		// sort.Sort(sort.IntSlice(x)) => sort.Ints(x)
		xfun.Sel.Name = "Ints"
		node.Args = arg.Args
		return
	}
}

// bf.Write([]byte(fmt.Sprintf("hello %d",1)))
// =>
// fmt.Fprintf(bf,"hello %d",1)
// 性能：1 => 1.37
// ut: custom4.go.input
func (c *cCallExpr) writeFmtSprintf() {
	if !astutil.UsesImport(c.req.AstFile, "fmt") {
		return
	}
	if len(c.Node.Args) != 1 {
		return
	}
	fun, ok1 := c.Node.Fun.(*ast.SelectorExpr)
	if !ok1 || fun.Sel == nil || fun.Sel.Name != "Write" {
		return
	}

	arg, ok2 := c.Node.Args[0].(*ast.CallExpr)
	if !ok2 || len(arg.Args) != 1 {
		return
	}

	af, ok3 := arg.Fun.(*ast.ArrayType)
	if !ok3 || !isIdentName(af.Elt, "byte") {
		return
	}

	aa, ok4 := arg.Args[0].(*ast.CallExpr)
	if !ok4 {
		return
	}
	if !isFun(aa.Fun, "fmt", "Sprintf") {
		return
	}

	var na []ast.Expr
	na = append(na, fun.X)
	na = append(na, aa.Args...)

	n := &ast.CallExpr{
		Lparen: c.Node.Lparen,
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "fmt",
			},
			Sel: &ast.Ident{
				Name: "Fprintf",
			},
		},
		Args: na,
	}
	c.Cursor.Replace(n)
}

// fmt.Sprintf("%d",123) -> strconv.Atoi(123)
//
// b:=int64(789)
// strconv.FormatInt(b, 10)
//
// c:=uint32(456)
// strconv.FormatUint(uint64(c), 10)
//
// fmt.Sprintf("%d",int8(1))  --> 	_ = strconv.FormatInt(int64(int8(1)), 10)
// int64(int8(1) 是符合预期的
//
// custom19_fmtSprintf.go.input
func (c *cCallExpr) fmtSprintfInt() {
	node := c.Node
	if !isFun(node.Fun, "fmt", "Sprintf") {
		return
	}

	arg1, ok := node.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}

	switch arg1.Value {
	default:
		return
	case `"%d"`, `"%v"`:
	}

	vtp, err := xpasser.TypeOf(c.req, node.Args[1])
	if err != nil {
		return
	}
	vb, ok3 := vtp.(*types.Basic)
	if !ok3 {
		return
	}

	doReplace := func(node ast.Node) {
		c.Cursor.Replace(node)
		pkgReplace(c.req.FSet, c.req.AstFile, "fmt", "strconv")
	}

	if vb.Kind() == types.Int {
		n := &ast.CallExpr{
			Lparen: c.Node.Lparen,
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "strconv",
				},
				Sel: &ast.Ident{
					Name: "Itoa",
				},
			},
			Args: []ast.Expr{node.Args[1]},
		}
		doReplace(n)
		return
	}

	formatInt := &ast.SelectorExpr{
		X: &ast.Ident{
			Name: "strconv",
		},
		Sel: &ast.Ident{
			Name: "FormatInt",
		},
	}
	blit10 := &ast.BasicLit{
		Kind:  token.INT,
		Value: "10",
	}

	formatUint := &ast.SelectorExpr{
		X: &ast.Ident{
			Name: "strconv",
		},
		Sel: &ast.Ident{
			Name: "FormatUint",
		},
	}

	switch vb.Kind() {
	case types.Int8, types.Int16, types.Int32:
		n := &ast.CallExpr{
			Lparen: c.Node.Lparen,
			Fun:    formatInt,
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: "int64",
					},
					Args: []ast.Expr{
						node.Args[1],
					},
				},
				blit10,
			},
		}
		doReplace(n)
		return
	case types.Int64:
		n := &ast.CallExpr{
			Lparen: c.Node.Lparen,
			Fun:    formatInt,
			Args: []ast.Expr{
				node.Args[1],
				blit10,
			},
		}
		doReplace(n)
		return
	case types.Uint8, types.Uint16, types.Uint32, types.Uint:
		n := &ast.CallExpr{
			Lparen: c.Node.Lparen,
			Fun:    formatUint,
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.Ident{
						Name: "uint64",
					},
					Args: []ast.Expr{
						node.Args[1],
					},
				},
				blit10,
			},
		}
		doReplace(n)
		return
	case types.Uint64:
		n := &ast.CallExpr{
			Lparen: c.Node.Lparen,
			Fun:    formatUint,
			Args: []ast.Expr{
				node.Args[1],
				blit10,
			},
		}
		doReplace(n)
		return
	}
}

// fmt.Sprintf("a %s b","hello") --> "a " + "hello" + " b"
func (c *cCallExpr) fmtSprintfStrings() {
	const enable = false
	if !enable {
		return
	}
	// todo
	node := c.Node

	if len(node.Args) < 2 || len(node.Args) > 5 {
		return
	}

	if !isFun(node.Fun, "fmt", "Sprintf") {
		return
	}

	fmtLayout, ok1 := node.Args[0].(*ast.BasicLit)
	if !ok1 || !strings.Contains(fmtLayout.Value, "%s") || len(fmtLayout.Value) <= 4 {
		return
	}

	list := strings.Split(fmtLayout.Value[1:len(fmtLayout.Value)-1], "%s")
	for i := 0; i < len(list); i++ {
		if strings.Contains(list[i], "%") {
			return
		}
	}

	for i := 1; i < len(node.Args); i++ {
		if !c.isBasicKind(node.Args[i], types.String) {
			return
		}
	}

	qt := string(fmtLayout.Value[0])

	items := make([]ast.Expr, 0, len(list)*2)
	for i := 0; i < len(list); i++ {
		if list[i] != "" {
			val := &ast.BasicLit{
				Kind:  token.STRING,
				Value: qt + list[i] + qt,
			}
			items = append(items, val)
		}
		if i == len(list)-1 {
			break
		}
		items = append(items, node.Args[i+1])
	}
	nn := stringExprJoin(items)

	c.Cursor.Replace(nn)
	pkgReplace(c.req.FSet, c.req.AstFile, "fmt", "fmt")
}

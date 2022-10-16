// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/18

package common

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"sort"
)

// NewTestRequest 给测试场景使用的，创建一个新的 request 对象
func NewTestRequest(fileName string) *Request {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	opt := NewDefaultOptions()
	return &Request{
		FileName: fileName,
		Opt:      *opt,
		FSet:     fset,
		AstFile:  f,
	}
}

// Request 一次格式化的请求
type Request struct {
	FSet      *token.FileSet
	AstFile   *ast.File
	tokenLine *TokenLine
	FileName  string
	Opt       Options
}

// TokenLine 获取 TokenLine
func (req *Request) TokenLine() *TokenLine {
	if req.tokenLine == nil {
		req.tokenLine = &TokenLine{
			file: req.FSet.File(req.AstFile.Pos()),
		}
	}
	return req.tokenLine
}

// FormatFile 将 AstFile 格式化、得到源码
func (req *Request) FormatFile() ([]byte, error) {
	return req.Opt.Source(req.FSet, req.AstFile)
}

// Save 保存文件
func (req *Request) Save(name string) error {
	code, err := req.FormatFile()
	if err != nil {
		return err
	}
	return os.WriteFile(name, code, 0644)
}

// ReParse 重新解析
func (req *Request) ReParse() error {
	fs, f, err := req.reParse()
	if err == nil {
		req.FSet = fs
		req.AstFile = f
	}
	return err
}

// MustReParse 重新解析，若失败会 panic
func (req *Request) MustReParse() {
	if err := req.ReParse(); err != nil {
		panic(fmt.Errorf("reParse %q failed: %w", req.FileName, err))
	}
}

func (req *Request) reParse() (*token.FileSet, *ast.File, error) {
	req.TokenLine().Execute()
	code, err := req.FormatFile()
	if err != nil {
		return nil, nil, err
	}
	req.tokenLine = nil
	return ParseOneFile(req.FileName, code)
}

// Clone reParser it and return a new Request
func (req *Request) Clone() *Request {
	c := &Request{
		FileName: req.FileName,
		Opt:      req.Opt,
	}
	fs, f, err := req.reParse()
	if err != nil {
		panic(fmt.Sprintf("reParser %q failed: %v", req.FileName, err))
	}
	c.FSet = fs
	c.AstFile = f
	return c
}

// TokenLine 记录对文件的换行的处理
type TokenLine struct {
	file       *token.File
	addPos     []int
	deleteLine []int
}

// AddLine 在指定位置添加新行
func (tf *TokenLine) AddLine(depth int, at token.Pos) {
	if Debug {
		DebugPrintln(depth+1, "AddLine", tf.file.Position(at), "atPos=", at)
	}
	tf.addPos = append(tf.addPos, tf.file.Offset(at))
}

// DeleteLine 删除指定位置的新行
func (tf *TokenLine) DeleteLine(depth int, line int) {
	if line < 1 {
		panic(fmt.Sprintf("invalid line number %d (should be >= 1)", line))
	}
	max := tf.file.LineCount()
	if line > max {
		panic(fmt.Sprintf("invalid line number %d (should be < %d)", line, max))
	}
	if Debug {
		DebugPrintln(depth+1, "DeleteLine=", line, "lineStart=", tf.file.LineStart(line))
	}
	tf.deleteLine = append(tf.deleteLine, line)
}

// Execute 将 Add、 Delete 的结果生效
func (tf *TokenLine) Execute() {
	if len(tf.addPos) == 0 && len(tf.deleteLine) == 0 {
		return
	}
	defer func() {
		tf.addPos = nil
		tf.deleteLine = nil
	}()
	lines := tokenFileLines(tf.file)
	lines = intSliceDelete(lines, tf.deleteLine...)

	lm := make(map[int]struct{}, len(lines))
	for _, v := range lines {
		lm[v] = struct{}{}
	}
	for _, offset := range tf.addPos {
		lm[offset] = struct{}{}
	}

	result := make([]int, 0, len(lm))
	for k := range lm {
		if k >= tf.file.Size() {
			continue
		}
		result = append(result, k)
	}
	sort.Ints(result)
	if !tf.file.SetLines(result) {
		panic(fmt.Sprintf("SetLines failed, size=%d, lines=%v", tf.file.Size(), result))
	}
}

func tokenFileLines(f *token.File) []int {
	field := reflect.ValueOf(f).Elem().FieldByName("lines")
	total := field.Len()
	lines := make([]int, 0, total)
	for i := 0; i < total; i++ {
		cur := int(field.Index(i).Int())
		lines = append(lines, cur)
	}
	return lines
}

func intSliceDelete(lines []int, delete ...int) []int {
	dm := make(map[int]bool, len(delete))
	for _, v := range delete {
		dm[v] = true
	}
	result := make([]int, 0, len(lines))
	for line, v := range lines {
		if !dm[line] {
			result = append(result, v)
		}
	}
	return result
}

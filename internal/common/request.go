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
	"sync"

	"golang.org/x/mod/semver"
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
	FSet    *token.FileSet
	AstFile *ast.File

	// tokenLine 用于处理 AstFile 的对应的 TokenLine
	tokenLine *TokenLine

	// FileName 文件名
	FileName string

	// goVersion 所属 go module(go.mod) 里的 定义的 go 版本
	goVersion string

	// Opt 处理的参数
	Opt Options

	goVersionOnce sync.Once

	directives    directives
	directiveOnce sync.Once
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

// GoVersionGEQ 判断模块 Go 的版本是否 >= 指定版本
// version: 版本号，如 1.19
func (req *Request) GoVersionGEQ(version string) bool {
	req.goVersionOnce.Do(func() {
		req.goVersion = goVersionByFile(req.FileName)
	})
	return semver.Compare("v"+req.goVersion, "v"+version) >= 0
}

// NoFormat 判断一个节点是否不需要执行格式化
func (req *Request) NoFormat(node ast.Node) bool {
	return req.HasDirective(node, "no_fmt")
}

// HasDirective 判断一个节点是否有指定的指令
func (req *Request) HasDirective(node ast.Node, name string) bool {
	ds := req.getDirectives()
	items0 := ds.ByNode(req.AstFile)
	if items0.Has(name) {
		return true
	}
	items := ds.ByNode(node)
	return items.Has(name)
}

func (req *Request) getDirectives() directives {
	req.directiveOnce.Do(func() {
		req.directives = parserDirectives(req.AstFile, req.FSet)
	})
	return req.directives
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
	maxLine := tf.file.LineCount()
	if line > maxLine {
		panic(fmt.Sprintf("invalid line number %d (should be < %d)", line, maxLine))
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

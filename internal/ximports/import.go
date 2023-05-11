// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package ximports

import (
	/*
	* 多行注释1
	 */ /*
	* 多行注释2
	 */ /* 多行注释3 */ "bytes"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"path/filepath"
	"regexp" // after regexp
	// on strconv
	"strconv" // strconv 后面
	"strings"
	// internal common 上面
	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/fsgo/fstypes"
)

// FormatImports 格式化 import 部分
// 多段import 会合并成一段
// 默认按照3段：系统库、第三方库、当前项目库
// 单独的注释行会保留
// 不会自动去除没使用的 import
//
//nolint:funlen,gocyclo
func FormatImports(req *common.Request) (out []byte, err error) {
	_opts := req.Opt.Clone()
	opts := *_opts
	err = req.ReParse()
	if err != nil {
		return nil, err
	}
	src, err := req.FormatFile()
	if err != nil {
		return nil, err
	}
	fileName := req.FileName
	file := req.AstFile

	// 将分组的 import 合并为一组，方便后续处理
	// 已知若 import 区域出现单独行的注释将不正确
	if len(file.Imports) < 2 {
		return src, nil
	}

	if err = findSubModules(fileName, &opts); err != nil {
		return nil, err
	}

	if opts.Trace {
		log.Println("ThirdParty Modules:", opts.ThirdModules)
	}

	var importDecls []*ast.GenDecl

	// 原始的 import 部分，未排序的
	// 保存的是原始的分组情况
	originImports := make(map[int][]*importDecl)

	var nextID int
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)

		// import "C" 是 Cgo 的，不可以处理
		if !ok || gen.Tok != token.IMPORT || declImports(gen, "C") {
			continue
		}
		importDecls = append(importDecls, gen)

		// 这个是整段的 import 语句
		importSrc := src[decl.Pos()-1 : decl.End()]
		lines, err := parserImportSrc(importSrc)
		if err != nil {
			panic(err.Error())
		}

		if _, has := originImports[nextID]; !has {
			originImports[nextID] = []*importDecl{}
		}

		var myImportDecls []*importDecl

		myImportDecls = originImports[nextID]

		myImportDecls = append(myImportDecls, lines...)

		originImports[nextID] = myImportDecls

		// 若将多段 import merge 到一个里,nextID 不变化
		// 下一段 import 将和之前的 merge 到一起
		if !opts.MergeImports {
			nextID++
		}
	}

	if opts.Trace {
		log.Println("originImports:", originImports)
	}

	var buf bytes.Buffer
	var start int
	var end int
	for i := 0; i < len(importDecls); i++ {
		decl := importDecls[i]

		buf.Write(src[start : decl.Pos()-1])

		importNew := formatImportDecls(originImports[i], opts)
		if len(importNew) > 0 {
			buf.Write(importNew)
		}
		start = int(decl.End())
		end = int(decl.End())
		if decl.Rparen.IsValid() {
			end = int(decl.Rparen)
		}
	}

	buf.Write(src[end:])
	code := cleanSpecCode(buf.Bytes())
	code, err = opts.Format(code)
	if err != nil {
		return nil, fmt.Errorf("format.Source failed: %w, code=\n%s", err, buf.String())
	}
	return code, err
}

// parserImportSrc 直接对 import 部分的源代码进行解析
func parserImportSrc(src []byte) (lines []*importDecl, err error) {
	body := bytes.TrimSpace(src[len("import"):])
	body = bytes.TrimLeft(body, "(")
	body = bytes.TrimRight(body, ")")

	// fmt.Println("body:", string(body))

	var decl *importDecl

	var checkDecl = func() {
		if decl == nil {
			decl = new(importDecl)
		}
	}

	var inCmt bool
	for _, lineBf := range bytes.Split(body, []byte("\n")) {
		bf := bytes.TrimSpace(lineBf)

		if len(bf) == 0 {
			continue
		}

		/*
		 * nihao
		 */

		if !inCmt {
			if isImportPathLine(bf) {
				line := string(bf)
				checkDecl()

				decl.Path = line
				lines = append(lines, decl)

				decl = nil
				continue
			}
		}

		if bytes.HasPrefix(bf, []byte("//")) {
			checkDecl()
			decl.AddComment(bf)
			// 注释掉的 import path 的情况
			if decl.CommentHasImportPath() {
				lines = append(lines, decl)
				// fmt.Println("CommentHasImportPath",decl)
				decl = nil
			}

			continue
		}

		if bytes.HasPrefix(bf, []byte("/*")) {
			inCmt = true
		}

		checkDecl()
		decl.AddComment(bf)

		if bytes.HasSuffix(bf, []byte("*/")) {
			inCmt = false
		}
	}

	if decl != nil {
		lines = append(lines, decl)
	}

	return lines, nil
}

// declImports reports whether gen contains an import of path.
// Taken from golang.org/x/tools/ast/astutil.
func declImports(gen *ast.GenDecl, path string) bool {
	if gen.Tok != token.IMPORT {
		return false
	}
	for _, spec := range gen.Specs {
		impspec := spec.(*ast.ImportSpec)
		if importPath(impspec) == path {
			return true
		}
	}
	return false
}

func importPath(s ast.Spec) string {
	t, err := strconv.Unquote(s.(*ast.ImportSpec).Path.Value)
	if err == nil {
		return t
	}
	return ""
}

// import 的 Name 部分
// 如 import(
//
//	redis "github.com/xxx"
//	_ "github.com/xxx"
//	. "github.com/xxx"
//
// )
// 上面的 "redis"、"_"、和 "." 都是
var importNameReg = regexp.MustCompile(`^([a-zA-Z_]?[a-zA-Z0-9_]*)|(\.)$`)

var importPathReg = regexp.MustCompile(`^[a-zA-Z_0-9\-\/\.]+$`)

func isImportPathLine(bf []byte) bool {
	line := bytes.TrimSpace(bf)
	if len(line) == 0 {
		return false
	}

	if bytes.ContainsAny(line, "\n\r") {
		return false
	}

	if !isImportPathHeader(line[0]) {
		return false
	}

	// 暂时未支持"`"
	if bytes.Count(line, []byte(`"`)) < 2 {
		return false
	}

	arr := bytes.SplitN(line, []byte(`"`), 3)
	name := bytes.TrimSpace(arr[0])

	if len(name) > 0 && !importNameReg.Match(name) {
		return false
	}

	// 不包括双引号
	importPath := arr[1]

	if !importPathReg.Match(importPath) {
		return false
	}

	cmt := bytes.TrimSpace(arr[2])
	if len(cmt) > 0 {
		// 若不是注释
		if !bytes.HasPrefix(cmt, []byte("/")) {
			return false
		}
	}

	return true
}

func isImportPathHeader(first byte) bool {
	return first == '"' ||
		first == '.' ||
		first == '_' ||
		(first >= 'A' && first <= 'Z') ||
		(first >= 'a' && first <= 'z')
}

func findSubModules(fileName string, opts *common.Options) error {
	// 这种情况是由于当前代码并没有 go.mod，所以需要忽略掉
	if strings.HasPrefix(opts.LocalModule, "_") {
		return nil
	}
	mp, err := common.FindGoModPath(fileName)
	if err != nil {
		return err
	}
	ms, err := common.ListModules(filepath.Dir(mp))
	if err != nil {
		return err
	}
	for i := 0; i < len(ms); i++ {
		module := ms[i]
		if module == opts.LocalModule {
			continue
		}
		if !fstypes.SliceHas(opts.ThirdModules, module) {
			opts.ThirdModules = append(opts.ThirdModules, module)
		}
	}
	return nil
}

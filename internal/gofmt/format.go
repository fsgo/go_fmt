/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package gofmt

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

	"golang.org/x/tools/imports"
)

// FormatFile 输出格式化的go代码
func Format(fileName string, src []byte, localPrefix string) ([]byte, error) {
	fileSet := token.NewFileSet()
	parserMode := parser.Mode(0) | parser.ParseComments
	file, err := parser.ParseFile(fileSet, fileName, src, parserMode)
	if err != nil {
		return nil, err
	}

	resetImportDecls(fileSet, file)

	// ast.Print(fileSet, file)

	fix(fileSet, file)

	var buf bytes.Buffer
	printerMode := printer.UseSpaces
	printerMode |= printer.TabIndent
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: 8}

	if err := printConfig.Fprint(&buf, fileSet, file); err != nil {
		return nil, err
	}

	imports.LocalPrefix = localPrefix

	opt := &imports.Options{
		Fragment:  true,
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}
	return imports.Process(fileName, buf.Bytes(), opt)
}

func resetImportDecls(fileSet *token.FileSet, f *ast.File) {
	// 将分组的import 合并为一组，方便后续处理
	// 已知若import 区域出现单独行的注释将不正确
	var firstLine int
	for _, ip := range f.Imports {
		p := ip.Pos()
		if firstLine == 0 {
			firstLine = fileSet.Position(p).Line + 1
		}
		fileSet.File(p).MergeLine(firstLine)
	}
}

func fix(fileSet *token.FileSet, f *ast.File) {
	if len(f.Decls) <= 1 {
		return
	}
	for _, cms := range f.Comments {
		for _, cm := range cms.List {
			if strings.HasPrefix(cm.Text, "//") && !strings.HasPrefix(cm.Text, "// ") {
				cm.Text = "// " + strings.TrimSpace(cm.Text[2:])
			}
		}
	}
}

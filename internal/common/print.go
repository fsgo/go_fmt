// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package common

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
)

// PrintCode 输出代码
func PrintCode(fileSet *token.FileSet, file *ast.File) (out []byte, err error) {
	var buf bytes.Buffer
	printerMode := printer.UseSpaces
	printerMode |= printer.TabIndent
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: 8}

	if err := printConfig.Fprint(&buf, fileSet, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FormatSource 格式化源代码
func FormatSource(src []byte) ([]byte, error) {
	fset, f, err := ParseOneFile("code.go", src)
	if err != nil {
		return nil, err
	}
	return PrintCode(fset, f)
}

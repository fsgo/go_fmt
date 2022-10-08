// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package common

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ParseOneFile 解析为 astFile
func ParseOneFile(fileName string, src []byte) (*token.FileSet, *ast.File, error) {
	fileSet := token.NewFileSet()
	parserMode := parser.Mode(0) | parser.ParseComments
	file, err := parser.ParseFile(fileSet, fileName, src, parserMode)
	return fileSet, file, err
}

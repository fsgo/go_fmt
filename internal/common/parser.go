// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package common

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

// ParseOneFile 解析为 astFile
func ParseOneFile(fileName string, src []byte) (*token.FileSet, *ast.File, error) {
	fileSet := token.NewFileSet()
	parserMode := parser.Mode(0) | parser.ParseComments
	file, err := parser.ParseFile(fileSet, fileName, src, parserMode)
	return fileSet, file, err
}

// ParseFile 解析文件，优先尝试使用 ParseDir
func ParseFile(fileName string, src []byte) (*token.FileSet, *ast.File, error) {
	fileSet := token.NewFileSet()
	var file *ast.File
	parserMode := parser.Mode(0) | parser.ParseComments
	pkgs, err := parser.ParseDir(fileSet, filepath.Dir(fileName), nil, parserMode)
	if err == nil {
		for _, pkg := range pkgs {
			for name, f := range pkg.Files {
				if name == fileName {
					file = f
				}
			}
		}
	}
	if file == nil {
		fileSet = token.NewFileSet()
		file, err = parser.ParseFile(fileSet, fileName, src, parserMode)
	}
	return fileSet, file, err
}

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/18

package common

import (
	"go/ast"
	"go/token"
)

// Request 一次格式化的请求
type Request struct {
	FileName string
	FSet     *token.FileSet
	AstFile  *ast.File
	Opt      Options
}

// FormatFile 将 AstFile 格式化、得到源码
func (req *Request) FormatFile() ([]byte, error) {
	return req.Opt.Source(req.FSet, req.AstFile)
}

// ReParse 重新解析
func (req *Request) ReParse() (*token.FileSet, *ast.File, error) {
	code, err := req.FormatFile()
	if err != nil {
		return nil, nil, err
	}
	return ParseOneFile(req.FileName, code)
}

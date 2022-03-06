// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/16

package gofmt

import (
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/localmodule"
	"github.com/fsgo/go_fmt/internal/simplify"
	// a "github.com/fsgo/go_fmt/internal/ximports"
	"github.com/fsgo/go_fmt/internal/ximports"
)

// Options 别名
type Options = common.Options

// Format 输出格式化的go代码
func Format(fileName string, src []byte, opts *Options) ([]byte, error) {
	options := opts.Clone()

	if src == nil {
		var err error
		src, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
	}

	if common.DoNotEdit(src) {
		return src, nil
	}

	module, err := localmodule.Get(options, fileName)
	if err != nil {
		return nil, err
	}

	if options.Trace {
		fmt.Println("[go.module]", fileName, "--->", module)
	}

	options.LocalModule = module

	outImports, errImports := ximports.FormatImports(fileName, src, options)
	if errImports != nil {
		return nil, errImports
	}

	src = outImports

	fileSet, file, err := common.ParseFile(fileName, src)
	if err != nil {
		return nil, err
	}

	// ast.Print(fileSet, file)

	fix(fileSet, file, src, options)

	return common.PrintCode(fileSet, file)
}

func fix(fileSet *token.FileSet, file *ast.File, src []byte, options *Options) {
	if options.Simplify {
		simplify.Format(file)
	}
	FormatComments(fileSet, file, options)
}

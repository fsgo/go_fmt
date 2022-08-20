// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/16

package gofmt

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/localmodule"
	"github.com/fsgo/go_fmt/internal/simplify"
	"github.com/fsgo/go_fmt/internal/ximports"
)

// Options 别名
type Options = common.Options

// Format 输出格式化的go代码
func Format(fileName string, src []byte, opts *Options) ([]byte, error) {
	options := opts.Clone()

	if src == nil {
		var err error
		src, err = os.ReadFile(fileName)
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

	file, err = fix(fileSet, file, options)
	if err != nil {
		return nil, err
	}
	return common.PrintCode(fileSet, file)
}

func fix(fileSet *token.FileSet, file *ast.File, opt *Options) (*ast.File, error) {
	if opt.Simplify {
		simplify.Format(file)
	}
	FormatComments(fileSet, file, opt)
	if len(opt.RewriteRules) > 0 {
		var err error
		file, err = simplify.Rewrites(file, opt.RewriteRules)
		if err != nil {
			return nil, err
		}
	}
	if opt.RewriteWithBuildIn {
		return simplify.Rewrites(file, simplify.BuildInRewriteRules())
	}
	return file, nil
}

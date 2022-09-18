// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/16

package gofmt

import (
	"fmt"
	"go/ast"
	"log"
	"os"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/localmodule"
	"github.com/fsgo/go_fmt/internal/simplify"
	"github.com/fsgo/go_fmt/internal/ximports"
	"github.com/fsgo/go_fmt/internal/xpasser"
)

// Options 别名
type Options = common.Options

// Format 输出格式化的go代码
func Format(fileName string, src []byte, opts *Options) (code []byte, formatted bool, err error) {
	options := opts.Clone()
	if src == nil {
		src, err = os.ReadFile(fileName)
		if err != nil {
			return nil, false, err
		}
	}

	if common.DoNotEdit(fileName, src) {
		return src, false, nil
	}

	module, err := localmodule.Get(options, fileName)
	if err != nil {
		return nil, false, err
	}

	if options.Trace {
		log.Println("[go.module]", fileName, "--->", module)
	}

	options.LocalModule = module

	file, err := xpasser.ParserFile(fileName, src)

	if err != nil {
		return nil, false, err
	}
	// ast.Print(fileSet, file)

	code, err = fix(fileName, file, options)
	if err != nil {
		return nil, false, err
	}
	return code, true, err
}

func fix(fileName string, file *ast.File, opt *Options) ([]byte, error) {
	// ast.Print(fileSet,file)
	fset := xpasser.Default.FSet
	if opt.Simplify {
		simplify.Format(fset, file)
	}

	FormatComments(fset, file, opt)

	if len(opt.RewriteRules) > 0 {
		var err error
		file, err = simplify.Rewrites(file, opt.RewriteRules)
		if err != nil {
			return nil, fmt.Errorf("rewrite failed: %w", err)
		}
	}

	if opt.RewriteWithBuildIn {
		var err error
		file, err = simplify.Rewrites(file, simplify.BuildInRewriteRules())
		if err != nil {
			return nil, fmt.Errorf("rewrite with build in rules failed: %w", err)
		}
	}

	// if opt.FieldAlignment == 1 {
	// 	fieldalignment.Run(fset, file, true)
	// }

	code, err := ximports.FormatImports(fileName, file, opt)
	if err != nil {
		return nil, fmt.Errorf("format import failed: %w", err)
	}
	code, err = common.FormatSource(code)
	if err != nil {
		return nil, fmt.Errorf("reformat failed: %w, code=\n%s", err, code)
	}
	return code, nil
}

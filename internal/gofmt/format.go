// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/16

package gofmt

import (
	"fmt"
	"log"
	"os"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/localmodule"
	"github.com/fsgo/go_fmt/internal/simplify"
	"github.com/fsgo/go_fmt/internal/xanalysis"
	"github.com/fsgo/go_fmt/internal/ximports"
	"github.com/fsgo/go_fmt/internal/xpasser"
)

// Options 别名
type Options = common.Options

// Format 输出格式化的 go 代码
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

	module, err := localmodule.Get(*options, fileName)
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

	req := &common.Request{
		FileName: fileName,
		FSet:     xpasser.Default.FSet,
		AstFile:  file,
		Opt:      *options,
	}
	// ast.Print(fileSet, file)

	code, err = fix(req)
	if err != nil {
		return nil, false, err
	}
	return code, true, err
}

func fix(req *common.Request) ([]byte, error) {
	// ast.Print(fileSet,file)
	if req.Opt.Simplify {
		simplify.Format(req)
	}
	FormatComments(req)

	if len(req.Opt.RewriteRules) > 0 {
		var err error
		file, err := simplify.Rewrites(req, req.Opt.RewriteRules)
		if err != nil {
			return nil, fmt.Errorf("rewrite failed: %w", err)
		}
		req.AstFile = file
	}

	if req.Opt.RewriteWithBuildIn {
		var err error
		file, err := simplify.Rewrites(req, simplify.BuildInRewriteRules())
		if err != nil {
			return nil, fmt.Errorf("rewrite with build in rules failed: %w", err)
		}
		req.AstFile = file
	}

	// if opt.FieldAlignment == 1 {
	// 	fieldalignment.Run(fset, file, true)
	// }

	xanalysis.Format(req)

	code, err := ximports.FormatImports(req)
	if err != nil {
		return nil, fmt.Errorf("format import failed: %w", err)
	}
	code, err = common.FormatSource(code)
	if err != nil {
		return nil, fmt.Errorf("reformat failed: %w, code=\n%s", err, code)
	}
	return code, nil
}

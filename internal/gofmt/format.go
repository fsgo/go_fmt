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
	"github.com/fsgo/go_fmt/internal/xast"
	"github.com/fsgo/go_fmt/internal/ximports"
	"github.com/fsgo/go_fmt/internal/xpasser"
)

// Options 别名
type Options = common.Options

// Format 输出格式化的 go 代码
func Format(fileName string, src []byte, opts *Options) (code []byte, formatted bool, err error) {
	defer func() {
		if re := recover(); re != nil {
			panic(fmt.Sprintf("when Format %q", fileName))
		}
	}()
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
	xanalysis.Format(req)

	// ast.Print(fileSet,file)
	if req.Opt.Simplify {
		simplify.Format(req)
	}
	FormatComments(req)

	if len(req.Opt.RewriteRules) > 0 {
		err := simplify.Rewrites(req, req.Opt.RewriteRules)
		if err != nil {
			return nil, fmt.Errorf("rewrite failed: %w", err)
		}
	}

	if req.Opt.RewriteWithBuildIn {
		err := simplify.Rewrites(req, simplify.BuildInRewriteRules())
		if err != nil {
			return nil, fmt.Errorf("rewrite with build in rules failed: %w", err)
		}
	}

	// if opt.FieldAlignment == 1 {
	// 	fieldalignment.Run(fset, file, true)
	// }

	if req.Opt.Extra {
		xast.Format(req)
	}

	code, err := ximports.FormatImports(req)
	if err != nil {
		return nil, fmt.Errorf("format import failed: %w", err)
	}
	code, err = req.Opt.Format(code)
	if err != nil {
		return nil, fmt.Errorf("reformat failed: %w, code=\n%s", err, code)
	}
	return code, nil
}

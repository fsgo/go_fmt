/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package gofmt

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/localmodule"
	"github.com/fsgo/go_fmt/internal/ximports"
)

type Options = common.Options

// Format 输出格式化的go代码
func Format(fileName string, src []byte, options *Options) ([]byte, error) {
	localPrefix, err := localmodule.Get(options.LocalPrefix, fileName)
	if err != nil {
		return nil, err
	}

	if options.Trace {
		fmt.Println("fileName--->", fileName)
	}

	options.LocalPrefix = localPrefix

	outImports, errImports := ximports.FormatImports(fileName, src, options, nil)
	if errImports != nil {
		return nil, errImports
	}

	src = outImports

	fileSet, file, err := common.ParseFile(fileName, src)
	if err != nil {
		return nil, err
	}

	// ast.Print(fileSet, file)

	fix(fileSet, file, src)

	return common.PrintCode(fileSet, file)

	// opt := &imports.Options{
	// Fragment:  true,
	// Comments:  true,
	// TabIndent: options.TabIndent,
	// TabWidth:  options.TabWidth,
	// }
	//
	// return imports.Process(fileName, buf.Bytes(), opt)
}

func fix(fileSet *token.FileSet, file *ast.File, src []byte) {
	FormatComments(fileSet, file)
}

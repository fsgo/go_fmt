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
	"io/ioutil"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/localmodule"
	// a "github.com/fsgo/go_fmt/internal/ximports"
	"github.com/fsgo/go_fmt/internal/ximports"
)

// Options 别名
type Options = common.Options

// Format 输出格式化的go代码
func Format(fileName string, src []byte, options *Options) ([]byte, error) {
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

	localPrefix, err := localmodule.Get(options.LocalPrefix, fileName)
	if err != nil {
		return nil, err
	}

	if options.Trace {
		fmt.Println("fileName--->", fileName)
	}

	options.LocalPrefix = localPrefix

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

	fix(fileSet, file, src)

	return common.PrintCode(fileSet, file)
}

func fix(fileSet *token.FileSet, file *ast.File, src []byte) {
	FormatComments(fileSet, file)
}

/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/1
 */

package ximports

import (
	"fmt"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/fsgo/go_fmt/internal/common"
)

// Clean 未使用的import
func Clean(fileName string, src []byte) ([]byte, error) {
	fileSet, file, err := common.ParseFile(fileName, src)
	if err != nil {
		return nil, err
	}
	if len(file.Imports) < 1 {
		return src, nil
	}

	for _, imp := range file.Imports {
		if !astutil.UsesImport(file, imp.Path.Value) {
			suc := astutil.DeleteImport(fileSet, file, imp.Path.Value)
			fmt.Println("clean:", imp.Path.Value, suc)
		}
	}
	return common.PrintCode(fileSet, file)
}

// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/16

package main

import (
	"fmt"
	"os"

	"github.com/fsgo/go_fmt/gofmtapi"
	"github.com/fsgo/go_fmt/internal/common"
)

func main() {
	gf := gofmtapi.NewFormatter()
	opt := gofmtapi.NewOptions()
	opt.BindFlags()

	err := gf.Execute(opt)

	if err != nil {
		fmt.Fprint(os.Stderr, common.ConsoleRed(err.Error())+"\n")
		os.Exit(2)
	}
}

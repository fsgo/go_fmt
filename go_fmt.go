// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/16

package main

import (
	"fmt"
	"os"

	"github.com/fsgo/go_fmt/gofmtapi"
)

func main() {
	gf := gofmtapi.NewFormatter()
	opt := gofmtapi.NewOptions()
	opt.BindFlags()

	err := gf.Execute(opt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}
}

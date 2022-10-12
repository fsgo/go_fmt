// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/12

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

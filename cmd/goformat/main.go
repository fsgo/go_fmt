// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/29

package main

import (
	"fmt"
	"os"

	"github.com/fsgo/go_fmt/gofmtapi"
	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/gofmt"
)

func main() {
	gf := gofmtapi.NewFormatter()
	opt := gofmtapi.NewOptions()
	opt.Write = false
	gofmt.OnExecute = func(opt *gofmt.Options) error {
		if len(opt.Files) == 1 && opt.Files[0] == common.NameGitChange {
			opt.Files[0] = common.NameSTDIN
			opt.Write = false
		}
		return nil
	}
	opt.BindFlags()

	err := gf.Execute(opt)
	if err != nil {
		fmt.Fprint(os.Stderr, common.ConsoleRed(err.Error())+"\n")
		os.Exit(2)
	}
}

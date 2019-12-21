/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package main

import (
	"github.com/fsgo/go_fmt/internal/gofmt"
)

func main() {
	gf := gofmt.NewGoFmt()
	gf.BindFlags()
	gf.Execute()
}

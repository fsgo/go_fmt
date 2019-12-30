/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package gofmt

// Options 选项
type Options struct {
	Trace bool

	TabIndent bool

	TabWidth int

	LocalPrefix string

	Write bool

	Files []string
}

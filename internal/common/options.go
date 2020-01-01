/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/1
 */

package common

// Options 选项
type Options struct {
	Trace bool

	TabIndent bool

	TabWidth int

	LocalPrefix string

	Write bool

	Files []string
}

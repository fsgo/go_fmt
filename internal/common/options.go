/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/1
 */

package common

// ImportGroupFn import 排序逻辑
type ImportGroupFn func(importPath string, opt *Options) int

// Options 选项
type Options struct {
	Trace bool

	TabIndent bool

	TabWidth int

	LocalPrefix string

	Write bool

	// 待处理的文件列表
	Files []string

	ImportGroupFn ImportGroupFn

	// 是否将多段import 合并为一个
	MergeImports bool
}

/*
 * Copyright(C) 2021 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2021/1/15
 */

package gofmtapi

import (
	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/gofmt"
)

// Formatter 格式化工具类
type Formatter = gofmt.Formatter

// Options 参数
type Options = common.Options

// ImportGroupFunc import 排序逻辑
type ImportGroupFunc = common.ImportGroupFunc

// NewOptions 生成默认的 option
func NewOptions() *Options {
	return common.NewDefaultOptions()
}

// NewFormatter 生成新 Formatter 对象
func NewFormatter() *Formatter {
	return gofmt.NewFormatter()
}

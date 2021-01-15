/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/1
 */

package common

import (
	"flag"
	"fmt"
	"os"
)

// ImportGroupFunc import 排序逻辑
type ImportGroupFunc func(importPath string, opt *Options) int

// Options 选项
type Options struct {
	Trace bool

	TabIndent bool

	TabWidth int

	LocalPrefix string

	Write bool

	// 待处理的文件列表
	Files []string

	ImportGroupFn ImportGroupFunc

	// 是否将多段 import 合并为一个
	MergeImports bool
}

// NewDefaultOptions 生成默认的 options
func NewDefaultOptions() *Options {
	return &Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalPrefix:  "auto",
		Write:        true,
		MergeImports: true,
	}
}

// AllGoFiles 获取所有的待格式化的 .go 文件
func (opt *Options) AllGoFiles() ([]string, error) {
	if len(opt.Files) == 0 {
		return nil, fmt.Errorf("opt.Files cannot empty")
	}

	var files []string
	var err error

	for _, name := range opt.Files {
		if name == "" {
			continue
		}
		var tmpList []string

		switch name {
		case "./...":
			tmpList, err = allGoFiles("./")
		case "git_change":
			tmpList, err = filesGitDirChange()
		default:
			info, errStat := os.Stat(name)
			if errStat != nil {
				return nil, err
			}
			if info.IsDir() {
				tmpList, err = allGoFiles(name)
			} else {
				// 若属实传入 文件名 可以不用检查是否是.go文件
				// 在一些特殊场景可能会有用
				if len(opt.Files) == 1 || isGoFileName(name) {
					tmpList = []string{name}
				} else {
					err = fmt.Errorf("%q is not .go file", name)
				}
			}
		}

		if err != nil {
			return nil, err
		}

		if len(tmpList) > 0 {
			files = append(files, tmpList...)
		}
	}
	return files, nil
}

// BindFlags 绑定参数信息
func (opt *Options) BindFlags() {
	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.BoolVar(&opt.Write, "w", true, "write result to (source) file instead of stdout")
	commandLine.StringVar(&opt.LocalPrefix, "local", "auto", "put imports beginning with this string as 3rd-party packages")
	commandLine.BoolVar(&opt.Trace, "trace", false, "show trace infos")
	commandLine.BoolVar(&opt.MergeImports, "mi", false, "merge imports into one")

	commandLine.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [path ...]\n", cmd)
		commandLine.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nsite :    github.com/fsgo/go_fmt\n")
		fmt.Fprintf(os.Stderr, "version:  %s\n", Version)
		os.Exit(2)
	}

	if err := commandLine.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "parser commandLine with error: %s\n", err.Error())
		os.Exit(2)
	}

	opt.Files = commandLine.Args()

	if len(opt.Files) == 0 {
		opt.Files = []string{"git_change"}
	}
}

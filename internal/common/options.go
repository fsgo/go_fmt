// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package common

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// ImportGroupFunc import 排序逻辑
type ImportGroupFunc func(importPath string, opt *Options) int

// Options 选项
type Options struct {
	Trace bool

	TabIndent bool

	TabWidth int

	// LocalModule 当前代码所在的 module
	// 对应其 go.mod 文件中的 module 的值
	LocalModule string

	// ThirdModules 可选，第三方模块列表
	//
	// 是为解决这种情况：
	// LocalModule = github.com/test
	// 但是其子目录有其他的子模块，如：
	// github.com/test/hello/say
	// github.com/test/world
	// 这个时候，在 github.com/test 里的代码，应该将 github.com/test/hello/say 归为第三方模块代码的分组
	ThirdModules Modules

	// Write 是否直接将格式化后的内容写入文件
	Write bool

	// Simplify  是否简化代码
	Simplify bool

	// DisplayDiff  是否只检查是否已格式化，
	// 当值为 true 时，会强制设置 Write=false
	DisplayDiff bool

	// DisplayFormat 输出 DisplayDiff 的格式，默认为 text，还可以是 json
	DisplayFormat string

	// 待处理的文件列表
	Files []string

	ImportGroupFn ImportGroupFunc

	// 是否将多段 import 合并为一个
	MergeImports bool

	// SingleLineCopyright 是否将 copyright 的多行注释格式化为单行注释
	SingleLineCopyright bool

	// import 分组的排序规则,可选
	// 总共 可分为 3 组，分别是 标准库(简称 s)，第三方库(简称 t)，模块自身(简称 c)
	// stc: 默认的排序规则
	// sct: Go 源码中的排序规则
	ImportGroupRule string
}

// ImportGroupType import 分组类型
type ImportGroupType byte

const (
	// ImportGroupGoStandard 标准库(简称 s)
	ImportGroupGoStandard ImportGroupType = 's'

	// ImportGroupThirdParty 第三方库(简称 t)
	ImportGroupThirdParty ImportGroupType = 't'

	// ImportGroupCurrentModule 模块自身(简称 c)
	ImportGroupCurrentModule ImportGroupType = 'c'
)

var defaultImportGroupRule = "stc"

// GetImportGroup 读取 import 分组的排序
func (opt *Options) GetImportGroup(t ImportGroupType) int {
	if len(opt.ImportGroupRule) == 0 {
		return strings.Index(defaultImportGroupRule, string(t))
	}
	return strings.Index(opt.ImportGroupRule, string(t))
}

// Check 简称 option 是否正确
func (opt *Options) Check() error {
	if opt.DisplayDiff {
		opt.Write = false
	}
	if len(opt.ImportGroupRule) > 0 {
		if len(opt.ImportGroupRule) != 3 {
			return fmt.Errorf("invalid ig %q", opt.ImportGroupRule)
		}
		for i := 0; i < len(opt.ImportGroupRule); i++ {
			switch opt.ImportGroupRule[i] {
			case 's', 't', 'c':
			default:
				return fmt.Errorf("invalid ig %q", opt.ImportGroupRule)
			}
		}
	}
	return nil
}

// NewDefaultOptions 生成默认的 options
func NewDefaultOptions() *Options {
	return &Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalModule:  "auto",
		Write:        true,
		MergeImports: true,
		Simplify:     true,
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
	commandLine.BoolVar(&opt.DisplayDiff, "d", false, "display diffs instead of rewriting files")
	commandLine.StringVar(&opt.DisplayFormat, "df", "text", "display diffs format, support: text, json")
	commandLine.BoolVar(&opt.Simplify, "s", true, "simplify code")
	commandLine.StringVar(&opt.LocalModule, "local", "auto", "put imports beginning with this string as 3rd-party packages")
	commandLine.BoolVar(&opt.Trace, "trace", false, "show trace infos")
	commandLine.BoolVar(&opt.MergeImports, "mi", false, "merge imports into one")
	commandLine.BoolVar(&opt.SingleLineCopyright, "slcr", false, "multiline copyright to single-line")
	commandLine.StringVar(&opt.ImportGroupRule, "ig", defaultImportGroupRule, `import group sort rule,
stc: Go Standard pkgs, Third Party pkgs, Current ModuleByFile pkg
sct: Go Standard pkgs, Current ModuleByFile pkg, Third Party pkgs
`)

	commandLine.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [path ...]\n", cmd)
		commandLine.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")

		titleFormat := "%15s : %s\n"
		fmt.Fprintf(os.Stderr, titleFormat, "build with", runtime.Version())
		fmt.Fprintf(os.Stderr, titleFormat, "site", "https://github.com/fsgo/go_fmt")
		fmt.Fprintf(os.Stderr, titleFormat, "check update", "go install github.com/fsgo/go_fmt@master")
		fmt.Fprintf(os.Stderr, titleFormat, "version", Version)
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

// Clone 当执行 format 的时候，每个文件都 clone 一份
func (opt *Options) Clone() *Options {
	o1 := &Options{}
	*o1 = *opt
	return o1
}

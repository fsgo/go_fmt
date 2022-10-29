// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package common

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"runtime"
	"strings"
)

// ImportGroupFunc import 排序逻辑
type ImportGroupFunc func(importPath string, opt Options) int

// Options 选项
type Options struct {
	// ImportGroupFn 排序函数，可选，若不为空，将不会使用默认内置的规则
	ImportGroupFn ImportGroupFunc

	// DisplayFormat 输出 DisplayDiff 的格式，默认为 text，还可以是 json
	DisplayFormat string

	// import 分组的排序规则,可选
	// 总共 可分为 3 组，分别是 标准库(简称 s)，第三方库(简称 t)，模块自身(简称 c)
	// stc: 默认的排序规则
	// sct: Go 源码中的排序规则
	ImportGroupRule string

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

	// 待处理的文件列表
	Files []string

	// 重写、简化代码的规则，可选
	RewriteRules []string

	TabWidth int

	// FieldAlignment struct 字段对齐优化，可选，默认 0
	// 可选值:
	// 1-对发现的进行修正，同时打印日志
	// 2-只打印出需优化的日志信息
	// FieldAlignment int

	// Write 是否直接将格式化后的内容写入文件
	Write bool

	// Simplify  是否简化代码
	Simplify bool

	// DisplayDiff  是否只检查是否已格式化，
	// 当值为 true 时，会强制设置 Write=false
	DisplayDiff bool

	Trace bool

	// 是否将多段 import 合并为一个
	MergeImports bool

	// SingleLineCopyright 是否将 copyright 的多行注释格式化为单行注释
	SingleLineCopyright bool

	TabIndent bool

	// 是否使用内置的 rewrite 规则简化代码，可选，默认 false
	RewriteWithBuildIn bool

	// Extra 更多额外的、高级的格式化规则
	Extra bool
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
		TabIndent:          true,
		TabWidth:           8,
		LocalModule:        "auto",
		Write:              true,
		MergeImports:       true,
		Simplify:           true,
		Extra:              true,
		RewriteWithBuildIn: true,
	}
}

const (
	// NameSTDIN 特殊的文件名，用于标志从 stdin 读取代码
	NameSTDIN = "stdin"

	// NameGitChange 特殊的文件名，表示查找所有当前有修改的文件
	NameGitChange = "git_change"
)

// AllGoFiles 获取所有的待格式化的 .go 文件
func (opt *Options) AllGoFiles() ([]string, error) {
	if len(opt.Files) == 0 {
		return nil, errors.New("opt.Files cannot empty")
	}

	var files []string
	var err error

	for _, name := range opt.Files {
		if len(name) == 0 {
			continue
		}
		var tmpList []string

		switch name {
		case "./...":
			tmpList, err = allGoFiles("./")
		case NameGitChange:
			tmpList, err = filesGitDirChange()
			if err != nil {
				if opt.Trace {
					log.Println("git_change:", err.Error())
				}
				// 若获取 git 变化的文件失败，则获取当前目录下所有文件
				tmpList, err = allGoFiles("./")
			}
		case NameSTDIN:
			files = append(files, name)
			opt.Write = false
			goto end
		default:
			info, errStat := os.Stat(name)
			if errStat != nil {
				return nil, err
			}
			if info.IsDir() {
				tmpList, err = allGoFiles(name)
			} else { // 若属实传入 文件名 可以不用检查是否是.go文件
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
end:
	return files, nil
}

// BindFlags 绑定参数信息
func (opt *Options) BindFlags() {
	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.BoolVar(&opt.Write, "w", opt.Write, "write result to (source) file instead of stdout")
	commandLine.BoolVar(&opt.DisplayDiff, "d", false, "display diffs instead of rewriting files")
	commandLine.StringVar(&opt.DisplayFormat, "df", "text", "display diffs format, support: text, json")
	commandLine.BoolVar(&opt.Simplify, "s", opt.Simplify, "simplify code")
	commandLine.StringVar(&opt.LocalModule, "local", "auto",
		`current package path, will put imports beginning with this string as 3rd-party packages.
by default, it will got from 'go.mod' file.
`,
	)
	commandLine.BoolVar(&opt.Trace, "trace", false, "show trace messages")
	commandLine.BoolVar(&opt.Extra, "e", opt.Extra, "enable extra rules")
	commandLine.BoolVar(&opt.MergeImports, "mi", flagEnvBool("GORGEOUS_MI", opt.MergeImports),
		`merge imports into one section.
with env 'GORGEOUS_MI=false' to set default value as false`)
	commandLine.BoolVar(&opt.SingleLineCopyright, "slcr", flagEnvBool("GORGEOUS_SLCR", opt.SingleLineCopyright),
		`multiline copyright to single-line
with env 'GORGEOUS_SLCR=true' to set default value as true`)
	commandLine.StringVar(&opt.ImportGroupRule, "ig", defaultImportGroupRule, `import groups sort rule,
stc: Go Standard package, Third Party package, Current package
sct: Go Standard package, Current package, Third Party package
`)
	var rewriteRule stringSlice
	commandLine.Var(&rewriteRule, "r", `rewrite rule (e.g., 'a[b:len(a)] -> a[b:]')
or a file path for rewrite rules (like -rr)
`)

	commandLine.BoolVar(&opt.RewriteWithBuildIn, "rr", flagEnvBool("GORGEOUS_RR", opt.RewriteWithBuildIn),
		`rewrite with build in rules:
`+strings.Join(BuildInRewriteRules(), "\n")+`

with env 'GORGEOUS_RR=false' to set default value as false
`)

	commandLine.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [path ...]\n", cmd)
		commandLine.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")

		titleFormat := "%15s : %s\n"
		fmt.Fprintf(os.Stderr, titleFormat, "build with", runtime.Version())
		fmt.Fprintf(os.Stderr, titleFormat, "site", "https://github.com/fsgo/go_fmt")
		fmt.Fprintf(os.Stderr, titleFormat, "check update", "go install github.com/fsgo/go_fmt/cmd/gorgeous@master")
		fmt.Fprintf(os.Stderr, titleFormat, "version", Version)
		os.Exit(2)
	}

	if err := commandLine.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "parser commandLine with error: %s\n", err.Error())
		os.Exit(2)
	}
	opt.RewriteRules = rewriteRule
	opt.Files = commandLine.Args()

	if len(opt.Files) == 0 {
		opt.Files = []string{NameGitChange}
	}
}

// Clone 当执行 format 的时候，每个文件都 clone 一份
func (opt *Options) Clone() *Options {
	o1 := &Options{}
	*o1 = *opt
	return o1
}

func (opt *Options) getTabWidth() int {
	if opt.TabWidth > 0 {
		return opt.TabWidth
	}
	return 8
}

const printerNormalizeNumbers = 1 << 30

// Source 格式化文件
func (opt *Options) Source(fileSet *token.FileSet, file *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	printerMode := printer.Mode(0) | printer.UseSpaces
	if opt.TabIndent {
		printerMode |= printer.TabIndent
	}
	printerMode |= printerNormalizeNumbers

	printConfig := &printer.Config{
		Mode:     printerMode,
		Tabwidth: opt.getTabWidth(),
	}

	if err := printConfig.Fprint(&buf, fileSet, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Format 重新格式化代码
func (opt *Options) Format(src []byte) ([]byte, error) {
	fset, f, err := ParseOneFile("tmp.go", src)
	if err != nil {
		return nil, err
	}
	return opt.Source(fset, f)
}

func flagEnvBool(key string, def bool) bool {
	switch os.Getenv(key) {
	case "true":
		return true
	case "false":
		return false
	default:
		return def
	}
}

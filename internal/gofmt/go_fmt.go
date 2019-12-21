/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/21
 */

package gofmt

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// NewGoFmt 创建一个新的带默认格式化规则的格式化实例
func NewGoFmt() *GoFmt {
	return &GoFmt{
		Options: &Options{
			TabIndent:   true,
			TabWidth:    8,
			LocalPrefix: "auto",
			Write:       true,
		},
	}
}

// GoFmt 代码格式化实例
type GoFmt struct {
	Options *Options
}

// BindFlags 绑定参数信息
func (gf *GoFmt) BindFlags() {
	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.BoolVar(&gf.Options.Write, "w", true, "write result to (source) file instead of stdout")
	commandLine.StringVar(&gf.Options.LocalPrefix, "local", "auto", "put imports beginning with this string after 3rd-party packages; comma-separated list")

	commandLine.Usage = func() {
		cmd := os.Args[0]
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [path ...]\n", cmd)
		commandLine.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nsite :    github.com/fsgo/go_fmt\n")
		fmt.Fprintf(os.Stderr, "version:  %s\n", Version)
		os.Exit(2)
	}

	commandLine.Parse(os.Args[1:])

	gf.Options.Files = commandLine.Args()

	if len(gf.Options.Files) == 0 {
		gf.Options.Files = []string{"git_change"}
	}
}

// Execute 执行代码格式化
func (gf *GoFmt) Execute() {
	err := gf.execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}
}

func (gf *GoFmt) execute() error {
	files, err := gf.ParserOptionsFiles()
	if err != nil {
		return err
	}

	var errTotal int

	for _, fileName := range files {

		var out []byte
		var change bool
		var errFmt error

		if gf.Options.Write {
			change, errFmt = gf.FormatAndWriteFile(fileName)
		} else {
			out, change, errFmt = gf.Format(fileName, nil)
		}
		if errFmt != nil {
			errTotal++
		}

		gf.printFmtResult(fileName, change, errFmt)

		if len(out) > 0 {
			fmt.Println(string(out))
		}
	}
	return nil
}

func (gf *GoFmt) printFmtResult(fileName string, change bool, err error) {
	var consoleColorTag = 0x1B
	if change {
		fmt.Fprintf(os.Stdout, "%c[31m rewrited : %s%c[0m\n", consoleColorTag, fileName, consoleColorTag)
	} else {
		fmt.Fprintf(os.Stderr, "%c[32m unchange : %s%c[0m\n", consoleColorTag, fileName, consoleColorTag)
	}
}

// ParserOptionsFiles 通过解析Options.Files值获取要进行格式化的文件列表
func (gf *GoFmt) ParserOptionsFiles() ([]string, error) {
	if len(gf.Options.Files) == 0 {
		return nil, fmt.Errorf("Options.Files cannot empty")
	}

	var files []string
	var err error

	for _, name := range gf.Options.Files {
		if name == "" {
			continue
		}
		var tmpList []string
		switch name {
		case "./...":
			tmpList, err = filesCurrentDirAll()
		case "git_change":
			tmpList, err = filesGitDirChange()
		default:
			// 若属实传入 文件名 可以不用检查是否是.go文件
			// 在一些特殊场景可能会有用
			if len(gf.Options.Files) == 1 || isGoFileName(name) {
				tmpList = []string{name}
			} else {
				err = fmt.Errorf("%q is not .go file", name)
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

// Format 格式化文件，获取格式化后的内容
func (gf *GoFmt) Format(fileName string, src []byte) (out []byte, change bool, err error) {
	if src == nil {
		src, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, false, err
		}
	}

	out, err = Format(fileName, src, gf.Options)

	if err != nil {
		return nil, false, err
	}
	change = !bytes.Equal(src, out)
	return out, change, nil
}

// FormatAndWriteFile 格式化并写入文件
func (gf *GoFmt) FormatAndWriteFile(fileName string) (bool, error) {
	out, change, err := gf.Format(fileName, nil)

	if err != nil {
		return false, err
	}

	if !change {
		return change, nil
	}

	if gf.Options.Write {
		err = ioutil.WriteFile(fileName, out, 0)
	}
	return change, err
}

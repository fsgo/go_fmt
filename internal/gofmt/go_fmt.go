// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/21

package gofmt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

// NewFormatter 创建一个新的带默认格式化规则的格式化实例
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Formatter 代码格式化实例
type Formatter struct {
	// PrintResult 用于格式化过程中，打印结果
	PrintResult func(fileName string, change bool, err error)
}

// Execute 执行代码格式化
func (gf *Formatter) Execute(opt *Options) error {
	return gf.execute(opt)
}

func (gf *Formatter) execute(opt *Options) error {
	files, err := opt.AllGoFiles()
	if err != nil {
		return err
	}
	var errTotal int

	for _, fileName := range files {

		var out []byte
		var change bool
		var errFmt error

		if opt.Write {
			change, errFmt = gf.FormatAndWriteFile(fileName, opt)
		} else {
			out, change, errFmt = gf.Format(fileName, nil, opt)
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

func (gf *Formatter) printFmtResult(fileName string, change bool, err error) {

	if gf.PrintResult != nil {
		gf.PrintResult(fileName, change, err)
		return
	}

	var consoleColorTag = 0x1B
	if change {
		fmt.Fprintf(os.Stderr, "%c[31m rewrited : %s%c[0m\n", consoleColorTag, fileName, consoleColorTag)
	} else {
		fmt.Fprintf(os.Stderr, "%c[32m unchange : %s%c[0m\n", consoleColorTag, fileName, consoleColorTag)
	}
}

// Format 格式化文件，获取格式化后的内容
func (gf *Formatter) Format(fileName string, src []byte, opt *Options) (out []byte, change bool, err error) {
	if len(src) == 0 {
		src, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, false, err
		}
	}
	out, err = Format(fileName, src, opt)

	if err != nil {
		return nil, false, err
	}
	change = !bytes.Equal(src, out)
	return out, change, nil
}

// FormatAndWriteFile 格式化并写入文件
func (gf *Formatter) FormatAndWriteFile(fileName string, opt *Options) (bool, error) {
	out, change, err := gf.Format(fileName, nil, opt)

	if err != nil {
		return false, err
	}

	if !change {
		return change, nil
	}

	if opt.Write {
		err = ioutil.WriteFile(fileName, out, 0)
	}
	return change, err
}

// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/21

package gofmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/fsgo/go_fmt/internal/common"
)

// NewFormatter 创建一个新的带默认格式化规则的格式化实例
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Formatter 代码格式化实例
type Formatter struct {
	// PrintResult 用于格式化过程中，打印结果
	PrintResult func(fileName string, change bool, err error)

	diffs diffResults
}

// Execute 执行代码格式化
func (ft *Formatter) Execute(opt *Options) error {
	if e := opt.Check(); e != nil {
		return e
	}
	return ft.execute(opt)
}

func (ft *Formatter) execute(opt *Options) error {
	files, err := opt.AllGoFiles()
	if err != nil {
		return err
	}
	ft.diffs = nil

	ch := make(chan bool, 20) // 控制并发

	var wg sync.WaitGroup
	var failNum int64
	var changeNum int64
	var mu sync.Mutex
	for i := 0; i < len(files); i++ {
		wg.Add(1)
		fileName := files[i]
		ch <- true

		// 并发，同时对多个文件进行格式化
		go func() {
			defer func() {
				<-ch
				wg.Done()
			}()
			change, err2 := ft.doFormat(opt, fileName)

			mu.Lock()
			defer mu.Unlock()

			if err2 != nil {
				failNum++
				err = err2
			}
			if change {
				changeNum++
			}
		}()
	}
	wg.Wait()
	close(ch)

	if len(ft.diffs) > 0 {
		ft.diffs.Output(opt.DisplayFormat)
	}
	mu.Lock()
	defer mu.Unlock()
	if err != nil {
		return err
	}
	if failNum > 0 {
		return fmt.Errorf("%d files format failed", failNum)
	}
	if opt.DisplayDiff && changeNum > 0 {
		return fmt.Errorf("%d files sholud be formated", changeNum)
	}
	return nil
}

func (ft *Formatter) printFmtResult(fileName string, change bool, event string, err error) {
	if ft.PrintResult != nil {
		ft.PrintResult(fileName, change, err)
		return
	}
	txt := fmt.Sprintf(" %8s : %s", event, fileName)
	if change {
		txt = common.ConsoleRed(txt)
	} else {
		txt = common.ConsoleGreen(txt)
	}
	fmt.Fprint(os.Stderr, txt, "\n")
}

func (ft *Formatter) execCallBack(opt *Options, fileName string, originSrc []byte, prettySrc []byte, err error) {
	if bytes.Equal(originSrc, prettySrc) || !opt.DisplayDiff {
		return
	}
	result := common.Diff(string(originSrc), string(prettySrc), opt.Trace)
	if result == nil {
		return
	}

	df := &diffResult{
		File:   fileName,
		Diffs:  result.Detail(),
		result: result,
	}
	if err != nil {
		df.Error = err.Error()
	}
	ft.diffs = append(ft.diffs, df)
}

func (ft *Formatter) doFormat(opt *Options, fileName string) (bool, error) {
	originSrc, prettySrc, err := ft.Format(fileName, nil, opt)
	ft.execCallBack(opt, fileName, originSrc, prettySrc, err)
	changed := !bytes.Equal(originSrc, prettySrc)
	if err != nil {
		ft.printFmtResult(fileName, true, "error", err)
		return changed, err
	}

	if !changed {
		ft.printFmtResult(fileName, false, "pretty", nil)
		return changed, nil
	}

	if opt.Write {
		err = os.WriteFile(fileName, prettySrc, 0)
		ft.printFmtResult(fileName, true, "rewrote", err)
		return changed, err
	} else if opt.DisplayDiff {
		ft.printFmtResult(fileName, true, "ugly", err)
	}
	if !opt.DisplayDiff {
		fmt.Println(string(prettySrc))
	}
	return changed, nil
}

// Format 格式化文件，获取格式化后的内容
func (ft *Formatter) Format(fileName string, src []byte, opt *Options) (originSrc []byte, prettySrc []byte, err error) {
	if len(src) == 0 {
		src, err = os.ReadFile(fileName)
		if err != nil {
			return nil, nil, err
		}
	}
	prettySrc, err = Format(fileName, src, opt)

	if err != nil {
		return nil, nil, err
	}
	return src, prettySrc, nil
}

// FormatAndWriteFile 格式化并写入文件
func (ft *Formatter) FormatAndWriteFile(fileName string, opt *Options) (bool, error) {
	originSrc, prettySrc, err := ft.Format(fileName, nil, opt)

	if err != nil {
		return false, err
	}

	if bytes.Equal(originSrc, prettySrc) {
		return false, nil
	}

	if opt.Write {
		err = os.WriteFile(fileName, prettySrc, 0)
	}
	return true, err
}

type diffResult struct {
	File   string
	Error  string
	Diffs  interface{}
	result common.DiffResult
}

type diffResults []*diffResult

func (drs diffResults) Output(format string) {
	if format == "json" {
		bf, _ := json.MarshalIndent(drs, " ", "    ")
		fmt.Fprint(os.Stdout, string(bf), "\n")
		return
	}

	for i := 0; i < len(drs); i++ {
		item := drs[i]
		title := common.ConsoleRed(item.File)
		msg := item.result.String()
		if len(item.Error) > 0 {
			msg += "Error:" + common.ConsoleRed(item.Error)
		}
		fmt.Fprint(os.Stdout, title, "\n", msg, "\n")
	}
}

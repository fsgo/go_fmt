// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/4

package std

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var curPkgs []string
var curErr error

var doOnce sync.Once

// PKGs 获取当前环境所有的标准库
// 值获取第一级目录
func PKGs() ([]string, error) {
	doOnce.Do(func() {
		curPkgs, curErr = currentPKGs()
	})
	return curPkgs, curErr
}

func currentPKGs() ([]string, error) {
	ctx := build.Default
	if len(ctx.GOROOT) == 0 {
		return nil, errors.New("GOROOT is empty")
	}
	cmd := exec.Command("go", "env", "GOROOT")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	out = bytes.TrimSpace(out)
	if len(out) == 0 {
		return nil, errors.New("not found GOROOT")
	}
	var pkgs []string
	stdSrcDir := filepath.Join(string(out), "src")
	gls, err := filepath.Glob(stdSrcDir + "/*")
	if err != nil {
		return nil, err
	}
	for _, fl := range gls {
		info, errStat := os.Stat(fl)

		if errStat != nil {
			return nil, fmt.Errorf("os.Stat(%q),with error:%w", fl, errStat)
		}

		if info.IsDir() {
			name := info.Name()
			if name == "vendor" {
				continue
			}

			pkgs = append(pkgs, name)
		}
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("cannot find any pkgs from dir %q", stdSrcDir)
	}
	return pkgs, nil
}

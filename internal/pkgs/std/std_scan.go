/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/4
 */

package std

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

// PKGs 获取当前环境所有的标准库
// 值获取第一级目录
func PKGs() ([]string, error) {
	ctx := build.Default
	if ctx.GOROOT == "" {
		return nil, fmt.Errorf("GOROOT is empty")
	}
	var pkgs []string
	stdSrcDir := filepath.Join(ctx.GOROOT, "src")
	gls, err := filepath.Glob(stdSrcDir + "/*")
	if err != nil {
		return nil, err
	}
	for _, fl := range gls {
		info, errStat := os.Stat(fl)

		if errStat != nil {
			return nil, fmt.Errorf("os.Stat(%q),with error:%s", fl, errStat)
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

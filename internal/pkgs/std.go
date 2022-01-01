// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/4

package pkgs

import (
	"log"
	"strings"

	"github.com/fsgo/go_fmt/internal/pkgs/std"
)

//go:generate go run cmd/update_std_static.go -out std_static.go

// IsStd 判断是否系统标准库
func IsStd(path string) bool {
	arr := strings.SplitN(path, "/", 2)
	first := arr[0]

	// 判断默认的标准库
	if inSlice(first, stdPKGs) {
		return true
	}

	// 获取当前系统的标准库
	currentStds, err := std.PKGs()
	if err != nil {
		log.Printf("scan std pkgs failed, err=%s\n", err)
		return false
	}

	if inSlice(first, currentStds) {
		return true
	}

	return false
}

func inSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

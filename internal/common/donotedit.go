// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/13

package common

import (
	"bytes"
	"strings"
)

// DoNotEdit 该代码是否不让修改
func DoNotEdit(name string, src []byte) bool {
	if strings.HasSuffix(name, ".pb.go") {
		return true
	}
	return doNotEditBySrc(src)
}

// doNotEditBySrc 该代码是否不让修改，这类一般是由代码工具自动生成
func doNotEditBySrc(src []byte) bool {
	lines := bytes.Split(src, []byte("\n"))
	if len(lines) == 0 {
		return false
	}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if !bytes.HasPrefix(line, []byte("//")) {
			return false
		}
		if bytes.Contains(line, []byte("DO NOT EDIT")) {
			return true
		}
		line1 := bytes.Replace(line, []byte(" "), []byte(""), -1)
		line1 = bytes.ToUpper(line1)
		if bytes.Contains(line1, []byte("DONOTEDIT")) {
			return true
		}
	}
	return false
}

// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/3/13

package common

import (
	"bytes"
)

// DoNotEdit 该代码是否不让修改，一般是有代码工具自动生成
func DoNotEdit(src []byte) bool {
	lines := bytes.Split(src, []byte("\n"))
	if len(lines) == 0 {
		return false
	}

	first := bytes.TrimSpace(lines[0])

	if !bytes.HasPrefix(first, []byte("//")) {
		return false
	}

	first = bytes.Replace(first, []byte(" "), []byte(""), -1)
	first = bytes.ToUpper(first)
	return bytes.Contains(first, []byte("DONOTEDIT"))
}

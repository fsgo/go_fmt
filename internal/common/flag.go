// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/8/20

package common

import (
	"strings"
)

type stringSlice []string

func (ss *stringSlice) String() string {
	return strings.Join(*ss, ";")
}

func (ss *stringSlice) Set(s2 string) error {
	arr := strings.Split(s2, ";")
	result := make([]string, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		line := strings.TrimSpace(arr[i])
		if len(line) > 0 {
			result = append(result, line)
		}
	}
	if len(result) > 0 {
		*ss = result
	}
	return nil
}

// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/8/12

package common

func SliceHas[S ~[]T, T comparable](arr S, values ...T) bool {
	if len(values) == 0 {
		return false
	}
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(values); j++ {
			if arr[i] == values[j] {
				return true
			}
		}
	}
	return false
}

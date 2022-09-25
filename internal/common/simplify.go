// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/24

package common

var rewriteRules = []string{
	`a[b:len(a)] -> a[b:]`,
	`interface{} -> any`,
	`str == ""   -> len(str) > 0`,
}

// BuildInRewriteRules 获取内置的简化规则
func BuildInRewriteRules() []string {
	return rewriteRules
}

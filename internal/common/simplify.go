// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/24

package common

var rewriteRules = []string{
	`a[b:len(a)] -> a[b:]`,
	`interface{} -> any`,
	`a == ""     -> len(a) == 0`,
	`a != ""     -> len(a) != 0`,

	"io/#ioutil.NopCloser -> io.NopCloser",
	"io/#ioutil.ReadAll   -> io.ReadAll",
	"io/#ioutil.ReadFile  -> os.ReadFile",
	"io/#ioutil.TempFile  -> os.CreateTemp",
	"io/#ioutil.TempDir   -> os.MkdirTemp",
	"io/#ioutil.WriteFile -> os.WriteFile",
	"io/#ioutil.Discard   -> io.Discard",

	// "io/#ioutil.ReadDir -> os.WriteFile", #  这两个方法不兼容，不能直接转换
}

// BuildInRewriteRules 获取内置的简化规则
func BuildInRewriteRules() []string {
	return rewriteRules
}

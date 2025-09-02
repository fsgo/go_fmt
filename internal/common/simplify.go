// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/24

package common

var rewriteRules = []string{
	`a[b:len(a)] -> a[b:]`,

	`interface{} -> any                    // go1.18`,

	"io/#ioutil.NopCloser -> io.NopCloser  // go1.16",
	"io/#ioutil.ReadAll   -> io.ReadAll    // go1.16",
	"io/#ioutil.ReadFile  -> os.ReadFile   // go1.16",
	"io/#ioutil.TempFile  -> os.CreateTemp // go1.16",
	"io/#ioutil.TempDir   -> os.MkdirTemp  // go1.16",
	"io/#ioutil.WriteFile -> os.WriteFile  // go1.16",
	"io/#ioutil.Discard   -> io.Discard    // go1.16",

	// "io/#ioutil.ReadDir -> os.ReadDir", #  这两个方法不兼容，不能直接转换

	// "os.ErrInvalid        -> fs.ErrInvalid        // go1.20",
	// "os.ErrPermission     -> fs.ErrPermission     // go1.20",
	// "os.ErrExist          -> fs.ErrExist          // go1.20",
	// "os.ErrNotExist       -> fs.ErrNotExist       // go1.20",
	// "os.ErrClosed         -> fs.ErrClosed         // go1.20",
}

// BuildInRewriteRules 获取内置的简化规则
func BuildInRewriteRules() []string {
	return rewriteRules
}

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/3/6

package common

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/mod/modfile"
)

var goModPathCache = &sync.Map{}

// FindGoModPath 查找文件对应的 go.mod 文件
func FindGoModPath(fileName string) (string, error) {
	ap, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	pd := filepath.Dir(ap)

	if v, has := goModPathCache.Load(pd); has {
		return v.(string), nil
	}

	// 限定最大往上查找 128 次，避免
	for i := 0; i < 128; i++ {
		modPath := filepath.Join(pd, "go.mod")
		info, err := os.Stat(modPath)
		if err == nil && !info.IsDir() {
			goModPathCache.Store(pd, modPath)
			return modPath, nil
		}
		cpd := filepath.Dir(pd)
		if cpd == pd {
			break
		}
		pd = cpd
	}
	return "", fmt.Errorf("cannot found go.mod")
}

// ModuleByFile 解析 go.mod 文件里的  module 的值
func ModuleByFile(goModPath string) (string, error) {
	goModBuf, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}
	module := modfile.ModulePath(goModBuf)
	if len(module) == 0 {
		return "", fmt.Errorf("parser %s failed", goModPath)
	}
	return module, nil
}

// InModule 判断指定 pkg 的 是否属于 module
func InModule(pkg string, module string) bool {
	pkg1 := pkg + "/"
	module1 := module + "/"
	return strings.HasPrefix(pkg1, module1)
}

// Modules 模块列表
type Modules []string

// PkgIn 判断 pkg 是否属于 模块列表范围
func (ms Modules) PkgIn(pkg string) bool {
	for _, m := range ms {
		if InModule(pkg, m) {
			return true
		}
	}
	return false
}

// key：dir，string
// value: Modules
var modulesCache = &sync.Map{}

// ListModules 找到指定目录下的所有子 module
//
// 可能是这样的：
//
//	a.go
//	go.mod
//	+ world (目录)
//		say.go    // 这个和 下面的 hello 就是两个不同的 module
//	+ hello (目录) // 这是一个独立的 module
//		hello.go
//		go.mod
func ListModules(dir string) (Modules, error) {
	if v, has := modulesCache.Load(dir); has {
		return v.(Modules), nil
	}
	root := filepath.Join(dir, "go.mod")
	var result Modules
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || info.Name() != "go.mod" || root == path {
			return nil
		}
		m, err1 := ModuleByFile(path)
		if err1 != nil {
			return err1
		}
		result = append(result, m)
		return nil
	})

	if err == nil {
		modulesCache.Store(dir, result)
	}
	return result, err
}

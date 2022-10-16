// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/8

package simplify

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// fixImport 修正 import
// 主要是针对 Go 官方标记为 Deprecated 的修复规则的补充
// 如 io/ioutil.WriteFile -> os.WriteFile
// 需要将 import 的 io/ioutil 替换为  os
// 也可以是其他非标准库的pkg
func fixImport(pattern, replace *expr, fset *token.FileSet, f *ast.File) {
	oldPkg := pattern.PkgName()
	newPkg := replace.PkgName()
	if len(oldPkg) == 0 || len(newPkg) == 0 || oldPkg == newPkg {
		return
	}
	pkgReplace(fset, f, oldPkg, newPkg)
}

func pkgReplace(fset *token.FileSet, f *ast.File, oldPkg string, newPkg string) {
	astutil.AddImport(fset, f, newPkg)
	if !astutil.UsesImport(f, newPkg) {
		astutil.DeleteImport(fset, f, newPkg)
	}
	if !astutil.UsesImport(f, oldPkg) {
		astutil.DeleteImport(fset, f, oldPkg)
	}
}

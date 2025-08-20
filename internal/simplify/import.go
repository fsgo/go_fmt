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
	// ----------------------
	// 兼容 Go1.22+ ast.Scope 已标记废弃问题
	// 待后续 astutil 升级兼容后可删除此逻辑
	if f.Scope == nil {
		//lint:ignore SA1019 兼容，以后删除
		f.Scope = &ast.Scope{}
	}
	// ----------------------

	astutil.AddImport(fset, f, newPkg)
	if !astutil.UsesImport(f, newPkg) {
		astutil.DeleteImport(fset, f, newPkg)
	}
	if !astutil.UsesImport(f, oldPkg) {
		astutil.DeleteImport(fset, f, oldPkg)
	}
}

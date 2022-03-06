// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package ximports

import (
	"sort"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/pkgs"
)

func defaultImportGroup(importPath string, opt *common.Options) int {
	// 若是纯注释，则排在最上面
	if importPath == "" {
		return 0
	}

	module := opt.LocalModule

	// 当前的路径是否属于配置的 第三方模块列表里的
	if opt.ThirdModules.PkgIn(importPath) {
		return opt.GetImportGroup(common.ImportGroupThirdParty)
	}

	// 判断是否属于当前的项目库
	// 因为当前项目的 模块名(LocalModule) 可能和 系统标准库 的出现同名，所以优先判断
	if common.InModule(importPath, module) {
		return opt.GetImportGroup(common.ImportGroupCurrentModule)
	}

	// 系统标准库
	if pkgs.IsStd(importPath) {
		return opt.GetImportGroup(common.ImportGroupGoStandard)
	}

	// 默认为第三方库
	return opt.GetImportGroup(common.ImportGroupThirdParty)
}

func sortImportDecls(decls []*importDecl, opts *common.Options) importDeclGroups {
	groupFn := opts.ImportGroupFn
	if groupFn == nil {
		groupFn = defaultImportGroup
	}

	result := make([]*importDeclGroup, 0)
	groups := make(map[int]*importDeclGroup)

	for _, decl := range decls {
		num := groupFn(decl.RealPath(), opts)
		group, has := groups[num]
		if !has {
			group = &importDeclGroup{
				Group: num,
			}
			groups[num] = group

			result = append(result, group)
		}
		group.Decls = append(group.Decls, decl)
		group.sort()
	}

	sort.Slice(result, func(i, j int) bool {
		a := result[i]
		b := result[j]
		return a.Group < b.Group
	})

	return result
}

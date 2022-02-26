// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package ximports

import (
	"sort"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/pkgs"
)

func defaultImportGroup(importPath string, opt *common.Options) int {
	// 若是纯注释，则排在最上面
	if importPath == "" {
		return 0
	}

	LocalPrefix := opt.LocalPrefix
	// 本地项目库
	// 因为当前项目的 模块名(LocalPrefix) 可能和 系统标准库 的出现同名，所以优先判断
	if strings.HasPrefix(importPath, LocalPrefix) || strings.TrimSuffix(LocalPrefix, "/") == importPath {
		return opt.GetImportGroup('c')
	}

	// 系统标准库
	if pkgs.IsStd(importPath) {
		return opt.GetImportGroup('s')
	}

	// 第三方项目库
	// if strings.Contains(importPath, ".") {
	// 	return 1
	// }

	// 默认为第三方库
	return opt.GetImportGroup('t')
}

func sortImportDecls(decls []*importDecl, options *common.Options) importDeclGroups {
	groupFn := options.ImportGroupFn
	if groupFn == nil {
		groupFn = defaultImportGroup
	}

	result := make([]*importDeclGroup, 0)
	groups := make(map[int]*importDeclGroup)

	for _, decl := range decls {
		num := groupFn(decl.RealPath(), options)
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

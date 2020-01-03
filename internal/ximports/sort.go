/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/1
 */

package ximports

import (
	"sort"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

func defaultImportGroup(importPath string, opt *common.Options) int {
	LocalPrefix := opt.LocalPrefix
	// 本地项目库
	if strings.HasPrefix(importPath, LocalPrefix) || strings.TrimSuffix(LocalPrefix, "/") == importPath {
		return 2
	}
	// 第三方项目库
	if strings.Contains(importPath, ".") {
		return 1
	}
	//
	return 0
}

func sortImportDecls(decls []*importDecl, options *common.Options) []*importDeclGroup {
	groupFn := options.ImportGroupFn
	if groupFn == nil {
		groupFn = defaultImportGroup
	}

	result := make([]*importDeclGroup, 0)
	groups := make(map[int]*importDeclGroup, 0)

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
	}

	sort.Slice(result, func(i, j int) bool {
		a := result[i]
		b := result[j]
		return a.Group < b.Group
	})

	return result
}

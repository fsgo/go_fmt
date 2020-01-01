/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/1
 */

package ximports

import (
	"sort"
	"strings"
)

type GroupFn func(importPath string, LocalPrefix string) int

type GroupFns []GroupFn

func (fns GroupFns) Group(importPath string, LocalPrefix string) int {
	importPath = strings.TrimLeft(importPath, `"`)
	for _, fn := range fns {
		if num := fn(importPath, LocalPrefix); num >= 0 {
			return num
		}
	}
	//
	return 0
}

var importToGroup = GroupFns{
	// 本地项目库
	func(importPath string, LocalPrefix string) int {
		if strings.HasPrefix(importPath, LocalPrefix) || strings.TrimSuffix(LocalPrefix, "/") == importPath {
			return 2
		}
		return -1
	},

	// 第三方项目库
	func(importPath string, LocalPrefix string) int {
		if strings.Contains(importPath, ".") {
			return 1
		}
		return -1
	},

	// go标准库
	func(importPath string, LocalPrefix string) int {
		if importPath != "" {
			return 0
		}
		return -1
	},
}

func sortImportDecls(decls []*importDecl, groupFns GroupFns, LocalPrefix string) []*importDeclGroup {
	if groupFns == nil {
		groupFns = importToGroup
	}
	result := make([]*importDeclGroup, 0)
	groups := make(map[int]*importDeclGroup, 0)

	for _, decl := range decls {
		num := groupFns.Group(decl.RealPath(), LocalPrefix)
		group, has := groups[num]
		if !has {
			group = &importDeclGroup{}
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

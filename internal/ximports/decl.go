/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/1
 */

package ximports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

type importDecl struct {
	Comments []string
	Path     string
}

func (decl *importDecl) AddComment(bf []byte) {
	v := bytes.TrimSpace(bf)
	if len(v) == 0 {
		return
	}
	decl.Comments = append(decl.Comments, string(v))
}

// CommentHasImportPath 是否有注释 掉的import path
func (decl *importDecl) CommentHasImportPath() bool {
	if len(decl.Comments) == 0 {
		return false
	}

	if decl.realPathFromCmt() != "" {
		return true
	}

	return false
}
func (decl *importDecl) realPathFromCmt() string {
	for _, cmt := range decl.Comments {
		if !strings.HasPrefix(cmt, "//") || len(cmt) < 3 {
			continue
		}
		cmt = strings.TrimSpace(cmt[2:])
		if cmt == "" {
			continue
		}
		if isImportPathLine([]byte(cmt)) {
			return strings.TrimLeft(cmt, `_ `)
		}
	}
	return ""
}

// RealPath 获取实际的import path
// 比如下列 获取到的都是:a.com/aa"
// a "a.com/aa"
// _"a.com/aa"
// "a.com/aa"
func (decl *importDecl) RealPath() string {
	name := strings.TrimLeft(decl.Path, `/*_ `)
	if name == "" {
		name = decl.realPathFromCmt()
	}

	if !strings.Contains(name, `"`) {
		return name
	}
	arr := strings.SplitN(name, `"`, 3)
	return arr[1]
}

type importDeclGroup struct {
	Group int
	Decls []*importDecl
}

func (group *importDeclGroup) sort() {
	if len(group.Decls) < 2 {
		return
	}
	sort.Slice(group.Decls, func(i, j int) bool {
		a := group.Decls[i]
		b := group.Decls[j]
		return a.RealPath() < b.RealPath()
	})
}

func (group *importDeclGroup) Bytes() []byte {
	group.sort()
	var buf bytes.Buffer
	for _, item := range group.Decls {
		for _, cmt := range item.Comments {
			buf.WriteString(cmt)
			buf.WriteString("\n")
		}
		if item.Path != "" {
			buf.WriteString(item.Path)
			buf.WriteString("\n")
		}
	}
	return buf.Bytes()
}

func formatImportDecls(decls []*importDecl, options *common.Options) []byte {
	var buf, buf2 bytes.Buffer

	buf.WriteString("import (\n")
	groups := sortImportDecls(decls, options)

	if options.Trace {
		a, _ := json.MarshalIndent(groups, " ", " ")
		fmt.Println("formatImportDecls:=", string(a))
	}
	for _, group := range groups {
		if group.Group >= 0 {
			buf.Write(group.Bytes())
			// 每个分组使用空行分割
			buf.WriteString("\n")
		} else {
			buf2.Write(group.Bytes())
			// 每个分组使用空行分割
			buf2.WriteString("\n")
		}
	}

	buf.WriteString("\n)\n")

	if buf2.Len() > 0 {
		buf.WriteString("\n")
		buf.Write(buf2.Bytes())
		buf.WriteString("\n")
	}
	return buf.Bytes()
}

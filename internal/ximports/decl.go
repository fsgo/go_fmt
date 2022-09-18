// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/1

package ximports

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

type importDecl struct {
	Path     string
	Comments []string
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

// Bytes 重新序列化
func (decl *importDecl) Bytes() []byte {
	var buf bytes.Buffer
	for _, cmt := range decl.Comments {

		// 对注释中的多个空格替换为一个空格
		cmt = regexp.MustCompile(`\s+`).ReplaceAllString(cmt, " ")

		buf.WriteString(cmt)
		buf.WriteString("\n")
	}
	if decl.Path != "" {
		buf.WriteString(decl.Path)
		buf.WriteString("\n")
	}
	return buf.Bytes()
}

type importDeclGroups []*importDeclGroup

// String 序列化，调试打印用
func (ig importDeclGroups) String() string {
	var buf bytes.Buffer
	for idx, item := range ig {
		buf.WriteString("idx=")
		buf.WriteString(strconv.Itoa(idx))
		buf.WriteString("\n")
		buf.Write(item.Bytes())
	}
	return buf.String()
}

type importDeclGroup struct {
	Decls []*importDecl
	Group int
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
		buf.Write(item.Bytes())
	}
	return buf.Bytes()
}

// 多个 import group 分组直接的分隔符
// 使用这个是由于 source.Format() 方法，对于一个注释的第三方 path,如
// import(
//
//	  "fmt"
//	    // "github.com/fsgo/a"
//	)
//	// 期望是将 注释行的第三方库放入单独的一组，但是 source.Format 会将其中间的空行给去掉
//
// 为了达到预期，故目前这样处理，Format 后再使用 cleanSpecCode 方法将该分隔符删除掉
const importGroupSpit = "\"github.com/fsgo/gofmtgofmtgofmtgofmt\"\n"

// formatImportDecls 格式化 import 的一个分组
// 会对这个分组重新排序
func formatImportDecls(decls []*importDecl, options *common.Options) []byte {
	if len(decls) == 0 {
		return nil
	}
	var buf bytes.Buffer

	buf.WriteString("import (\n")
	groups := sortImportDecls(decls, options)

	if options.Trace {
		a, _ := json.MarshalIndent(groups, " ", " ")
		log.Println("formatImportDecls:=", string(a))
	}

	for _, group := range groups {
		groupCode := group.Bytes()

		// 每个分组使用特定分割
		if len(groupCode) > 0 {
			buf.WriteString("\n")
			buf.WriteString(importGroupSpit)
			buf.WriteString(string(groupCode))
		}
	}
	if !bytes.HasSuffix(buf.Bytes(), []byte("\n")) {
		buf.WriteString("\n")
	}

	buf.WriteString(")\n")

	return buf.Bytes()
}

func cleanSpecCode(src []byte) []byte {
	return bytes.ReplaceAll(src, []byte(importGroupSpit), []byte(""))
}

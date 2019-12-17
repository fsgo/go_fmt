/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package gofmt

import (
	"go/ast"
	"go/token"
	"strings"
)

// FormatComments 对代码注释进行格式化
// 若是空行注释将删除
func FormatComments(fileSet *token.FileSet, f *ast.File) {
	for _, cms := range f.Comments {
		var cmList []*ast.Comment
		for _, cm := range cms.List {
			if strings.HasPrefix(cm.Text, "//") {
				// 单行注释
				cm.Text = "// " + strings.TrimSpace(cm.Text[2:])
			} else if strings.HasPrefix(cm.Text, "/*") {
				// 多行注释
				// 会删除没有注释内容的多行注释
				if !fixMultilineComment(cm) {
					continue
				}
			}
			cmList = append(cmList, cm)
		}

		cms.List = cmList
	}
}

func fixMultilineComment(cm *ast.Comment) (ok bool) {
	txt := strings.TrimSpace(cm.Text)

	// 使用 /* */的多行注释
	txt = strings.TrimSpace(txt[1 : len(txt)-1]) // 去除两边的 "/"
	lines := strings.Split(txt, "\n")

	/* 下面是处理单行注释的逻辑 */
	if len(lines) == 1 {
		// 去掉两边的 *
		txt = strings.TrimSpace(txt[1 : len(txt)-1])
		if txt == "" {
			cm.Text = "/* */"
			return false
		} else {
			cm.Text = strings.Join([]string{"/* ", txt, " */"}, "")
		}
		return true
	}

	/*
	 * 下面是处理多行注释的逻辑
	 */

	var cmtLines []string

	for idx, line := range lines {
		// 去除开头的 "*"
		line = strings.TrimLeft(strings.TrimSpace(line), "*")
		line = strings.TrimSpace(line)
		// 首行 和末行
		if (idx == 0 || idx == len(lines)-1) && line == "" {
			continue
		}
		// 保留中间的空行
		if len(cmtLines) > 0 || line != "" {
			cmtLines = append(cmtLines, line)
		}
	}

	// 去除末尾的 空行注释
	cmtTxt := strings.TrimSpace(strings.Join(cmtLines, "\n"))
	if len(cmtTxt) == 0 {
		cm.Text = "/* */"
		return false
	}

	cmtLines = strings.Split(cmtTxt, "\n")

	var builder strings.Builder
	builder.WriteString("/*\n")
	for _, cmtLine := range cmtLines {
		builder.WriteString(" * ")
		builder.WriteString(cmtLine)
		builder.WriteString("\n")
	}
	builder.WriteString(" */")
	cm.Text = builder.String()

	return true
}

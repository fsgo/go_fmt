// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/16

package gofmt

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

// FormatComments 对代码注释进行格式化
// 若是空行注释将删除
func FormatComments(req *common.Request) {
	f := req.AstFile
	options := req.Opt
	for _, cms := range f.Comments {
		var cmList []*ast.Comment
		for _, cm := range cms.List {
			if strings.HasPrefix(cm.Text, "//") {
				// 单行注释
				// @see https://github.com/golang/go/blob/master/src/runtime/HACKING.md
				// @see go doc compile
				if strings.HasPrefix(cm.Text, "//go:") ||
					strings.HasPrefix(cm.Text, "//line ") ||
					strings.HasPrefix(cm.Text, "//*line ") ||
					strings.HasPrefix(cm.Text, "//nolint") ||
					strings.HasPrefix(cm.Text, "//lint") ||
					strings.HasPrefix(cm.Text, "//export ") {
					// spec comment,ignore
				} else {
					// 若 // 后没有空格则补充
					if !strings.HasPrefix(cm.Text, "// ") && !strings.HasPrefix(cm.Text, "//	") {
						cm.Text = "// " + cm.Text[2:]
					}
				}
			} else if strings.HasPrefix(cm.Text, "/*") {
				// 多行注释
				if !fixMultilineComment(req.FSet, cm, options) {
					continue
				}
			}
			cmList = append(cmList, cm)
		}

		cms.List = cmList
	}
}

// fixMultilineComment 多行注释处理
func fixMultilineComment(fileSet *token.FileSet, cm *ast.Comment, options Options) (ok bool) {
	if options.SingleLineCopyright {
		fixCopyright(fileSet, cm)
	}
	return true

	// txt := strings.TrimSpace(cm.Text)
	//
	// // cgo
	// if strings.Contains(txt, "#include ") ||
	// 	strings.Contains(txt, "#cgo ") {
	// 	return true
	// }
	//
	// // markdown 的注释代码
	// if strings.Contains(txt, "```") {
	// 	return true
	// }
	//
	// // 使用 /* */的多行注释
	// txt = strings.TrimSpace(txt[1 : len(txt)-1]) // 去除两边的 "/"
	// lines := strings.Split(txt, "\n")
	//
	// /* 下面是处理单行注释的逻辑 */
	// if len(lines) == 1 {
	// 	// 去掉两边的 *
	// 	txt = strings.TrimSpace(txt[1 : len(txt)-1])
	// 	if txt == "" {
	// 		cm.Text = "/* */"
	// 		return false
	// 	}
	// 	cm.Text = strings.Join([]string{"/* ", txt, " */"}, "")
	// 	return true
	// }
	//
	// /*
	//  * 下面是处理多行注释的逻辑
	//  */
	//
	// var cmtLines []string
	//
	// for idx, line := range lines {
	// 	// 去除开头的 "*"
	// 	line = strings.TrimLeft(strings.TrimSpace(line), "*")
	// 	line = strings.TrimSpace(line)
	// 	// 首行 和末行
	// 	if (idx == 0 || idx == len(lines)-1) && line == "" {
	// 		continue
	// 	}
	// 	// 保留中间的空行
	// 	if len(cmtLines) > 0 || line != "" {
	// 		cmtLines = append(cmtLines, line)
	// 	}
	// }
	//
	// // 去除末尾的 空行注释
	// cmtTxt := strings.TrimSpace(strings.Join(cmtLines, "\n"))
	// if len(cmtTxt) == 0 {
	// 	cm.Text = "/* */"
	// 	return false
	// }
	//
	// cmtLines = strings.Split(cmtTxt, "\n")
	//
	// var builder strings.Builder
	// builder.WriteString("/*\n")
	// for _, cmtLine := range cmtLines {
	// 	builder.WriteString(" * ")
	// 	builder.WriteString(cmtLine)
	// 	builder.WriteString("\n")
	// }
	// builder.WriteString(" */")
	// cm.Text = builder.String()
	//
	// return true
}

func fixCopyright(fset *token.FileSet, cm *ast.Comment) {
	p := fset.Position(cm.Pos())
	if p.Line != 1 {
		return
	}
	cm.Text = multilineCommentSingle(cm.Text)
}

func multilineCommentSingle(txt string) string {
	txt1 := strings.TrimSpace(txt[1 : len(txt)-1]) // 去除两边的 "/"
	lines := strings.Split(txt1, "\n")

	cmtLines := make([]string, 0, len(lines))
	for _, line := range lines {
		// 去除开头的 "*"
		line = strings.TrimLeft(strings.TrimSpace(line), "*")
		line = strings.TrimSpace(line)
		cmtLines = append(cmtLines, line)
	}
	newTxt := strings.TrimSpace(strings.Join(cmtLines, "\n"))

	lines = strings.Split(newTxt, "\n")
	cmtLines = cmtLines[:0]
	for _, line := range lines {
		cmtLines = append(cmtLines, "// "+line)
	}
	return strings.Join(cmtLines, "\n")
}

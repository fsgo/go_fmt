// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/10/4

package common

import (
	"go/ast"
	"go/token"
	"log"
	"net/url"
	"strings"
)

// 指令前缀，用于在代码中配置执行的条件，如对特定的代码不执行某项格式化处理
// 格式：// gorgeous:directive1,directive2,directive3(param1=1&param2=2)
const directivePrefix = "// gorgeous:"

type directives []directive

func (ds directives) Has(name string) bool {
	for _, d := range ds {
		if d.Name == name {
			return true
		}
	}
	return false
}

func (ds directives) ByNode(node ast.Node) directives {
	var result []directive
	for _, d := range ds {
		if d.Node == node {
			result = append(result, d)
		}
	}
	return result
}

type directive struct {
	Name    string       // 指令名称，必填
	Params  url.Values   // 指令参数，可选
	Comment *ast.Comment // 所属 comment
	Node    ast.Node     // 所属节点
}

func parserDirectiveComment(node ast.Node, cmt *ast.Comment) directives {
	txt := cmt.Text
	if !strings.HasPrefix(txt, directivePrefix) {
		return nil
	}
	var dirs []directive
	txt = strings.TrimPrefix(txt, directivePrefix)
	ss := strings.Split(txt, ",")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if strings.HasSuffix(s, ")") && strings.Contains(s, "(") {
			cmd, args, _ := strings.Cut(s[:len(s)-1], "(")
			vs, err := url.ParseQuery(args)
			if err != nil {
				// todo 打印代码位置
				log.Printf("invalid Directive：%q\n", s)
				continue
			}
			dir := directive{
				Name:    cmd,
				Params:  vs,
				Comment: cmt,
				Node:    node,
			}
			dirs = append(dirs, dir)
		} else {
			dir := directive{
				Name:    s,
				Params:  nil,
				Comment: cmt,
				Node:    node,
			}
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

func parserDirectives(f *ast.File, fset *token.FileSet) directives {
	var list []directive
	cm := ast.NewCommentMap(fset, f, f.Comments)
	for node, cgs := range cm {
		for _, cg := range cgs {
			for _, c := range cg.List {
				items := parserDirectiveComment(node, c)
				list = append(list, items...)
			}
		}
	}
	return list
}

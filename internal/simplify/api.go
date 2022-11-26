// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu (duv123@baidu.com)
// Date: 2022/3/5

package simplify

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

// Format call simplify
//
// rewrite.go 和 simplify.go 来自于 go1.19
func Format(req *common.Request) {
	simplify(req.AstFile)
	customSimplify(req)

	// rewrite 之后需要重新解析，否则 token.FileSet 可能和 ast.File 里的内容不匹配
	// 导致 panic
	req.MustReParse()
}

var ruleMsg1 = "rewrite rule must be of the form 'pattern -> replacement'"
var ruleMsg2 = ruleMsg1 + " or a valid filepath"

// Rewrite 简化代码
//
//nolint:gocyclo
func Rewrite(req *common.Request, rule string) (f *ast.File, e error) {
	defer func() {
		if re := recover(); re != nil {
			st := debug.Stack()
			e = fmt.Errorf("panic when Rewrite, file=%q rule=%q, detail=%v,stack=%s", req.FileName, rule, re, st)
		}
	}()
	rule, cmt := splitRule(rule)
	if len(rule) == 0 {
		return nil, fmt.Errorf("empty rewrite rule: %q", rule)
	}

	callDoRewrite := func(patternExp string, replaceExp string, what string) error {
		f1, err1 := doRewrite(req.FSet, req.AstFile, patternExp, replaceExp, what)
		if err1 != nil {
			return err1
		}
		req.AstFile = f1
		if err2 := req.ReParse(); err2 != nil {
			return err2
		}
		return nil
	}

	ps := strings.Split(rule, "->")
	if len(ps) == 2 {
		if gv := goVersionFromComment(cmt); len(gv) != 0 && !req.GoVersionGEQ(gv) {
			if req.Opt.Trace {
				log.Println("[Rewrite]", rule, "ignored")
			}
			return req.AstFile, nil
		}
		if err3 := callDoRewrite(ps[0], ps[1], rule); err3 != nil {
			return nil, err3
		}
		return req.AstFile, nil
	}

	if rules, err1 := parserRuleFile(rule); err1 == nil {
		for i := 0; i < len(rules); i++ {
			txt, cmt1 := splitRule(rules[i])
			if len(txt) == 0 {
				continue
			}
			ps1 := strings.Split(txt, "->")
			if len(ps1) != 2 {
				return nil, fmt.Errorf(ruleMsg1+",but at %s:%d, rule is %q", rule, i+1, txt)
			}
			if gv := goVersionFromComment(cmt1); len(gv) != 0 && !req.GoVersionGEQ(gv) {
				if req.Opt.Trace {
					log.Println("[Rewrite]", txt, "ignored")
				}
				continue
			}
			if err3 := callDoRewrite(ps1[0], ps1[1], txt); err3 != nil {
				return nil, err3
			}
		}
		return req.AstFile, nil
	}
	return nil, fmt.Errorf(ruleMsg2+", now got %q", rule)
}

var goVersionReg = regexp.MustCompile(`\sgo(1.\d+)`)

// 从注释里解析出 go 版本
// 如： interface{} -> any // go1.19 ，解析得到 1.19
func goVersionFromComment(comment string) string {
	txt := goVersionReg.FindString(comment)
	if len(txt) > 0 {
		return txt[3:]
	}
	return ""
}

// splitRule  将规则和注释分开
//
//	a == "" -> len(a) == 0 // comment
//	a != "" -> len(a) != 0
func splitRule(rule string) (r string, comment string) {
	rule = strings.TrimSpace(rule)
	arr := strings.SplitN(rule, "//", 2)
	if len(arr) == 2 {
		return arr[0], arr[1]
	}
	return rule, ""
}

func doRewrite(fs *token.FileSet, f *ast.File, patternExp string, replaceExp string, what string) (*ast.File, error) {
	patternExp = strings.TrimSpace(patternExp)
	if len(patternExp) == 0 {
		return nil, fmt.Errorf("pattern is empty, %q", what)
	}
	pattern := customParseExpr(patternExp, "pattern")

	replaceExp = strings.TrimSpace(replaceExp)
	if len(replaceExp) == 0 {
		return nil, fmt.Errorf("replacement is empty, %q", what)
	}
	replace := customParseExpr(replaceExp, "replacement")
	result, changed := rewriteFile(fs, pattern.expr, replace.expr, f)
	if changed {
		fixImport(pattern, replace, fs, result)
	}
	return result, nil
}

// Rewrites rewrite with many rules
func Rewrites(req *common.Request, rules []string) error {
	for i := 0; i < len(rules); i++ {
		if len(rules[i]) > 0 {
			f, err := Rewrite(req, rules[i])
			if err != nil {
				return err
			}
			req.AstFile = f
			if err = req.ReParse(); err != nil {
				return err
			}
		}
	}
	return nil
}

// RewriteWithExpr  使用指定的规则替换
func RewriteWithExpr(req *common.Request, pattern, replace ast.Expr) {
	result, _ := rewriteFile(req.FSet, pattern, replace, req.AstFile)
	req.AstFile = result
	req.MustReParse()
}

type expr struct {
	expr   ast.Expr
	input  string // 输入的原始的表达式,如 io/#ioutil.WriteFile
	what   string
	pkgDir string // 包所在目录，是 input 的最后一个 / 之前的部分（包括 /），如 "io/"
}

func (e *expr) PkgName() string {
	v1, ok1 := e.expr.(*ast.SelectorExpr)
	if !ok1 {
		return ""
	}
	xv, ok2 := v1.X.(*ast.Ident)
	if !ok2 {
		return ""
	}
	return e.pkgDir + xv.Name
}

// customParseExpr 解析表达式
// expStr 的 import path 使用 /# 作为分割，并置于前部
// eg:
//
//	io/#ioutil.WriteFile
//	os.WriteFile
func customParseExpr(expStr string, what string) *expr {
	expStr = strings.TrimSpace(expStr)
	e := &expr{
		input: expStr,
		what:  what,
	}
	idx := strings.LastIndex(expStr, "/#")
	if idx > 0 {
		e.pkgDir = expStr[:idx+1]
		e.expr = parseExpr(expStr[idx+2:], what)
	} else {
		e.expr = parseExpr(expStr, what)
	}
	return e
}

func parserRuleFile(name string) ([]string, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu (duv123@baidu.com)
// Date: 2022/3/5

package simplify

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

// Format call simplify
//
// rewrite.go 和 simplify.go 来自于 go1.19
func Format(req *common.Request) {
	simplify(req.AstFile)
}

// Rewrite 简化代码
func Rewrite(req *common.Request, rule string) (*ast.File, error) {
	if len(rule) == 0 {
		return nil, errors.New("empty rewrite rule")
	}
	ps := strings.Split(rule, "->")
	if len(ps) != 2 {
		return nil, fmt.Errorf("rewrite rule must be of the form 'pattern -> replacement', now got %q", rule)
	}
	pattern := customParseExpr(ps[0], "pattern")
	replace := customParseExpr(ps[1], "replacement")

	result := rewriteFile(req.FSet, pattern.expr, replace.expr, req.AstFile)

	fixImport(pattern, replace, req.FSet, result)

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

type expr struct {
	input  string // 输入的原始的表达式,如 io/#ioutil.WriteFile
	what   string
	pkgDir string // 包所在目录，是 input 的最后一个 / 之前的部分（包括 /），如 "io/"
	expr   ast.Expr
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

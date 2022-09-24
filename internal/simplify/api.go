// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu (duv123@baidu.com)
// Date: 2022/3/5

package simplify

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
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
func Rewrite(f *ast.File, rule string) (*ast.File, error) {
	if len(rule) == 0 {
		return nil, errors.New("empty rewrite rule")
	}
	ps := strings.Split(rule, "->")
	if len(ps) != 2 {
		return nil, fmt.Errorf("rewrite rule must be of the form 'pattern -> replacement', now got %q", rule)
	}
	pattern := parseExpr(ps[0], "pattern")
	replace := parseExpr(ps[1], "replacement")

	fSet := token.NewFileSet()
	return rewriteFile(fSet, pattern, replace, f), nil
}

// Rewrites rewrite with many rules
func Rewrites(req *common.Request, rules []string) (*ast.File, error) {
	f := req.AstFile
	var err error
	for i := 0; i < len(rules); i++ {
		if len(rules[i]) > 0 {
			if f, err = Rewrite(f, rules[i]); err != nil {
				return nil, err
			}
		}
	}
	return f, nil
}

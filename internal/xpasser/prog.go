// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/17

package xpasser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type Program struct {
	FSet *token.FileSet
	pkgs []*packages.Package
}

var Default = &Program{
	FSet: token.NewFileSet(),
}

func Load(patterns []string) error {
	conf := packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: true,
		Fset:  Default.FSet,
		// Logf: func(format string, args ...interface{}) {
		// 	log.Printf(format,args...)
		// },
	}
	pkgs, err := packages.Load(&conf, patterns...)
	if err != nil {
		return err
	}
	Default.pkgs = pkgs
	return nil
}

var errFileNotFound = errors.New("file not found in pkgs")

func (pr *Program) findPkg(filename string) (*packages.Package, *ast.File, error) {
	ap, err := filepath.Abs(filename)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < len(pr.pkgs); i++ {
		p := pr.pkgs[i]
		for j, n := range p.CompiledGoFiles {
			if n == ap {
				return p, p.Syntax[j], nil
			}
		}
	}
	return nil, nil, fmt.Errorf("%w: %s", errFileNotFound, filename)
}

func (pr *Program) ParserFile(filename string, src any) (*ast.File, error) {
	_, f, err := pr.findPkg(filename)
	if len(pr.pkgs) == 0 || (err != nil && errors.Is(err, errFileNotFound)) {
		mod := parser.Mode(0) | parser.ParseComments
		return parser.ParseFile(pr.FSet, filename, src, mod)
	}
	return f, err
}

func ParserFile(filename string, src any) (*ast.File, error) {
	return Default.ParserFile(filename, src)
}

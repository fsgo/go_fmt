// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/17

package xpasser

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/tools/go/packages"

	"github.com/fsgo/go_fmt/internal/common"
)

// Program 一个应用
type Program struct {
	FSet *token.FileSet
	pkgs []*packages.Package
}

// Default 默认的应用
var Default = &Program{
	FSet: token.NewFileSet(),
}

// Reset 重置默认环境，用于测试
func Reset() {
	Default = &Program{
		FSet: token.NewFileSet(),
	}
	Overlay = nil
}

// Overlay see packages.Config.Overlay
var Overlay map[string][]byte

// LoadOverlay 加载文件，用于测试
func LoadOverlay(fileName string) error {
	if Overlay == nil {
		Overlay = map[string][]byte{}
	}
	bf, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	Overlay[fileName] = bf
	return nil
}

// Load 加载解析应用
func Load(opt common.Options, patterns []string) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	conf := packages.Config{
		Context: ctx,
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo,
		Tests:   true,
		Fset:    Default.FSet,
		Overlay: Overlay,
		// Logf: func(format string, args ...interface{}) {
		// 	log.Printf(format,args...)
		// },
	}
	if len(patterns) == 0 {
		patterns = []string{"./..."}
	}

	if opt.Trace {
		log.Println("prog.Load_start")
	}

	pkgs, err := packages.Load(&conf, patterns...)
	if err != nil {
		return err
	}
	Default.pkgs = pkgs
	if opt.Trace {
		log.Println("prog.Load_done", "pkgs", pkgs, "patterns=", patterns, "cost=", time.Since(start).String())
	}
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

// FindPackage  查找所属 package
func (pr *Program) FindPackage(filename string) (*packages.Package, error) {
	p, _, err := pr.findPkg(filename)
	return p, err
}

// ParserFile 解析文件
func (pr *Program) ParserFile(filename string, src any) (*ast.File, error) {
	_, f, err := pr.findPkg(filename)
	if len(pr.pkgs) == 0 || (err != nil && errors.Is(err, errFileNotFound)) {
		mod := parser.Mode(0) | parser.ParseComments
		return parser.ParseFile(pr.FSet, filename, src, mod)
	}
	return f, err
}

// ParserFile 解析文件
func ParserFile(filename string, src any) (*ast.File, error) {
	return Default.ParserFile(filename, src)
}

// FindPackage 查找 pkg 信息
func FindPackage(filename string) (*packages.Package, error) {
	return Default.FindPackage(filename)
}

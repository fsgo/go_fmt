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
	"sync"
	"time"

	"golang.org/x/tools/go/packages"

	"github.com/fsgo/go_fmt/internal/common"
)

// Program 一个应用
type Program struct {
	FSet *token.FileSet
	pkgs []*packages.Package
}

func (p *Program) Packages() []*packages.Package {
	return p.pkgs
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
	overlayMux.Lock()
	overlay = nil
	overlayMux.Unlock()
}

var overlayMux sync.Mutex

// Overlay see packages.Config.Overlay
var overlay map[string][]byte

// LoadOverlay 加载文件，用于测试
func LoadOverlay(fileName string, code []byte) error {
	overlayMux.Lock()
	defer overlayMux.Unlock()

	if overlay == nil {
		overlay = map[string][]byte{}
	}
	ap, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	if len(code) == 0 {
		bf, err := os.ReadFile(ap)
		if err != nil {
			return err
		}
		code = bf
	}
	overlay[ap] = code
	return nil
}

func getOverlay() map[string][]byte {
	overlayMux.Lock()
	defer overlayMux.Unlock()
	return overlay
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
		Overlay: getOverlay(),
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
	goFiles := 0
	for i := 0; i < len(pr.pkgs); i++ {
		p := pr.pkgs[i]
		for j, n := range p.CompiledGoFiles {
			goFiles++
			if n == ap {
				if len(p.Syntax) <= j {
					return nil, nil, fmt.Errorf("%w: %s, invalid Syntax, Errors: %v", errFileNotFound, filename, p.Errors)
				}
				return p, p.Syntax[j], nil
			}
		}
	}
	return nil, nil, fmt.Errorf("%w: %s, len(pkgs)=%d, goFiles=%d", errFileNotFound, filename, len(pr.pkgs), goFiles)
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

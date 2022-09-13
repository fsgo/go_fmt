// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/11

package xpasser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/go/loader"
)

// File parser file with Program
type File struct {
	FileName  string
	AstFile   *ast.File
	FileSet   *token.FileSet
	Pkg       *types.Package // type information about the package
	TypesInfo *types.Info    // type information about the syntax trees
}

// Load parser and load Program
func (f *File) Load(src any) error {
	f.FileSet = token.NewFileSet()
	conf := &loader.Config{
		Fset:       f.FileSet,
		ParserMode: parser.Mode(0) | parser.ParseComments,
	}
	file, err := conf.ParseFile(f.FileName, src)
	if err != nil {
		return err
	}
	f.AstFile = file
	return nil
	// todo 优化性能

	prog, err := f.loadProgram(conf)

	if err != nil {
		return err
	}
	if prog != nil {
		p := prog.Package(file.Name.Name)
		f.TypesInfo = &p.Info
	}
	return nil
}

var caches sync.Map

func (f *File) loadProgram(conf *loader.Config) (*loader.Program, error) {
	ap, err := filepath.Abs(f.FileName)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(ap)
	val, _ := caches.Load(dir)
	if p, ok := val.(*loader.Program); ok {
		return p, nil
	}

	files, err := loadFiles(conf, filepath.Dir(f.FileName), f.FileName)
	if err != nil {
		return nil, err
	}

	files = append(files, f.AstFile)
	conf.CreateFromFiles("", files...)
	prog, err := conf.Load()
	if err == nil {
		caches.Store(dir, prog)
	} else {
		caches.Store(dir, err)
	}
	// ignore error
	return prog, nil
}

func loadFiles(conf *loader.Config, dir string, ignore string) ([]*ast.File, error) {
	ms, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	log.Println("ms", ms)
	files := make([]*ast.File, 0, len(ms))
	var mux sync.Mutex
	var wg errgroup.Group
	for i := 0; i < len(ms); i++ {
		fileName := ms[i]
		if filepath.Base(fileName) == filepath.Base(ignore) {
			continue
		}
		wg.Go(func() error {
			f, err1 := conf.ParseFile(fileName, nil)
			if err1 != nil {
				return fmt.Errorf("parser %s failed: %w", err1)
			}
			mux.Lock()
			files = append(files, f)
			mux.Unlock()
			return nil
		})
	}
	err2 := wg.Wait()
	return files, err2
}

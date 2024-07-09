// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/24

package gofmt_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fsgo/fst"
	"golang.org/x/tools/go/packages"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/gofmt"
	"github.com/fsgo/go_fmt/internal/xpasser"
)

func TestFormatter_Execute(t *testing.T) {
	tests := []struct {
		name    string
		caseDir string
	}{
		{
			name:    "app1",
			caseDir: "testdata/app1",
		},
		{
			name:    "app2",
			caseDir: "testdata/app2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("execute case: %s %s", tt.name, tt.caseDir)
			CheckDir(t, tt.caseDir)
		})
	}
}

// CheckDir 检查一个目录
//
// 测试数据放在 dir 下的 case 子目录
func CheckDir(t *testing.T, dir string) {
	// 避免 case 受到当前项目的 go.work 的影响
	t.Setenv("GOWORK", "off")

	defer xpasser.Reset()
	xpasser.Reset()

	caseDir := filepath.Join(dir, "case")
	t.Logf("case data dir is: %s", caseDir)
	info, err := os.Stat(caseDir)
	fst.NoError(t, err)
	fst.True(t, info.IsDir())

	runDir := filepath.Join(dir, "tmp")
	_ = os.RemoveAll(runDir)
	fst.NoError(t, os.MkdirAll(runDir, 0755))
	wants := map[string]string{}
	err = filepath.Walk(caseDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, err1 := filepath.Rel(caseDir, path)
		fst.NoError(t, err1)

		code, err2 := os.ReadFile(path)
		fst.NoError(t, err2)

		to := filepath.Join(runDir, rel)
		if strings.HasSuffix(path, ".want") {
			to = strings.ReplaceAll(rel, ".want", "")
			wants[to] = string(code)
		} else {
			to = strings.ReplaceAll(to, ".input", "")
			_ = os.MkdirAll(filepath.Dir(to), 0755)
			err3 := os.WriteFile(to, code, 0644)
			fst.NoError(t, err3)
		}
		return nil
	})
	fst.NoError(t, err)
	ft := gofmt.NewFormatter()
	opt := &common.Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalModule:  "auto",
		Write:        true,
		MergeImports: false,
		Trace:        true,
		Extra:        true,
		Files:        []string{"./..."},
	}
	pwd, err := os.Getwd()
	fst.NoError(t, err)
	fst.NotEmpty(t, pwd)

	fst.NoError(t, os.Chdir(runDir))
	defer func() {
		_ = os.Chdir(pwd)
		if !t.Failed() {
			_ = os.RemoveAll(runDir)
		}
	}()
	err = ft.Execute(opt)
	fst.NoError(t, err)

	var n int
	packages.Visit(xpasser.Default.Packages(), nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			t.Logf("pkg Errors:%s %s", pkg.Name, err)
			n++
		}
	})
	fst.Equal(t, 0, n)

	for name, want := range wants {
		t.Run(name, func(t *testing.T) {
			bf, err := os.ReadFile(name)
			fst.NoError(t, err)
			fst.Equal(t, want, string(bf))
		})
	}
}

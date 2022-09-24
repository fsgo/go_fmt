// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/24

package gofmt

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatter_Execute(t *testing.T) {
	fromDir := "testdata/app1/tpls"
	runDir := "testdata/app1/dotest"
	testExecute(t, fromDir, runDir)
}

func testExecute(t *testing.T, caseDir string, runDir string) {
	_ = os.RemoveAll(runDir)
	require.NoError(t, os.MkdirAll(runDir, 0755))
	wants := map[string]string{}
	err := filepath.Walk(caseDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(caseDir, path)
		require.NoError(t, err)

		code, err := os.ReadFile(path)
		require.NoError(t, err)

		to := filepath.Join(runDir, rel)
		if strings.HasSuffix(path, ".want") {
			to = strings.ReplaceAll(rel, ".want", "")
			wants[to] = string(code)
		} else {
			to = strings.ReplaceAll(to, ".input", "")
			err = os.WriteFile(to, code, 0644)
			require.NoError(t, err)
		}
		return nil
	})
	require.NoError(t, err)
	ft := NewFormatter()
	opt := &Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalModule:  "auto",
		Write:        true,
		MergeImports: false,
		Trace:        true,
		Files:        []string{"./..."},
	}
	pwd, err := os.Getwd()
	require.NoError(t, err)
	require.NotEmpty(t, pwd)

	require.NoError(t, os.Chdir(runDir))
	defer func() {
		_ = os.Chdir(pwd)
		if !t.Failed() {
			_ = os.RemoveAll(runDir)
		}
	}()
	err = ft.Execute(opt)
	require.NoError(t, err)

	for name, want := range wants {
		t.Run(name, func(t *testing.T) {
			bf, err := os.ReadFile(name)
			require.NoError(t, err)
			require.Equal(t, want, string(bf))
		})
	}

}

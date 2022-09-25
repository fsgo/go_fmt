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
			testExecute(t, tt.caseDir)
		})
	}
}

func testExecute(t *testing.T, caseDir string) {
	runDir := filepath.Join(caseDir, "tmp")
	caseDir = filepath.Join(caseDir, "tpls")
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
			_ = os.MkdirAll(filepath.Dir(to), 0755)
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

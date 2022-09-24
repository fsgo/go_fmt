// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/3

package gofmt

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fsgo/go_fmt/internal/xpasser"
	"github.com/stretchr/testify/require"
)

func TestFormat_rule1(t *testing.T) {
	opt := &Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalModule:  "auto",
		Write:        false,
		MergeImports: true,
	}
	runTest(t, "rule1", opt)
}

func TestFormat_rule2(t *testing.T) {
	opt := &Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalModule:  "auto",
		Write:        false,
		MergeImports: false,
		Simplify:     true,
	}
	runTest(t, "rule2", opt)
}

func TestFormat_rule3(t *testing.T) {
	opt := &Options{
		TabIndent:           true,
		TabWidth:            8,
		LocalModule:         "auto",
		Write:               false,
		MergeImports:        false,
		SingleLineCopyright: true,
	}
	runTest(t, "rule3", opt)
}

func TestFormat_rule4(t *testing.T) {
	opt := &Options{
		TabIndent:           true,
		TabWidth:            8,
		LocalModule:         "cmd",
		Write:               false,
		MergeImports:        false,
		SingleLineCopyright: true,
		ImportGroupRule:     "sct",
	}
	runTest(t, "rule4", opt)
}

func TestFormat_rule5(t *testing.T) {
	// 校验包含有子模块的情况
	opt := &Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalModule:  "auto",
		Write:        false,
		MergeImports: false,
	}
	runTest(t, "rule5", opt)
}

func runTest(t *testing.T, ruleDirName string, opt *Options) {
	rule1Dir := "./testdata/" + ruleDirName + "/"

	var checkFileTotal int

	err := filepath.Walk(rule1Dir+"/input/", func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil || !strings.HasSuffix(path, ".go.txt") {
			return nil
		}
		checkFileTotal++
		t.Run(path, func(t *testing.T) {
			xpasser.Reset()
			defer xpasser.Reset()

			pts := []string{"file=" + path}
			require.NoError(t, xpasser.Load(*opt, pts))

			got, _, err := Format(path, nil, opt)
			require.NoError(t, err)
			got = bytes.TrimSpace(got)
			want, err := os.ReadFile(rule1Dir + "want/" + filepath.Base(path))
			require.NoError(t, err)

			want = bytes.TrimSpace(want)
			fileGot := rule1Dir + "/tmp/" + filepath.Base(path) + ".got"
			fileWant := rule1Dir + "/tmp/" + filepath.Base(path) + ".want"
			if !bytes.Equal(got, want) {
				err1 := os.WriteFile(fileGot, got, 0644)
				err2 := os.WriteFile(fileWant, want, 0644)
				t.Errorf("not eq, fileget=%q (write=%v), filewant=%q (write=%v)", fileGot, err1, fileWant, err2)
			} else {
				_ = os.Remove(fileGot)
				_ = os.Remove(fileWant)
			}
		})
		return nil
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, checkFileTotal, 1)
}

// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/3

package gofmt

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

	err := filepath.Walk(rule1Dir+"/input/", func(path string, info os.FileInfo, err error) error {
		t.Run(path, func(t *testing.T) {
			if err == nil && strings.HasSuffix(path, ".go.txt") {
				checkFileTotal++

				got, err := Format(path, nil, opt)
				if err != nil {
					t.Errorf("Format returl error: %s", err)
					return
				}
				got = bytes.TrimSpace(got)
				want, err := ioutil.ReadFile(rule1Dir + "want/" + filepath.Base(path))
				if err != nil {
					t.Fatalf("ioutil.ReadFile with error: %s", err)
				}
				want = bytes.TrimSpace(want)
				if !bytes.Equal(got, want) {
					fileGot := rule1Dir + "/tmp/" + filepath.Base(path) + ".got"
					fileWant := rule1Dir + "/tmp/" + filepath.Base(path) + ".want"
					err1 := ioutil.WriteFile(fileGot, got, 0644)
					err2 := ioutil.WriteFile(fileWant, want, 0644)
					t.Errorf("not eq, fileget=%q (write=%v), filewant=%q (write=%v)", fileGot, err1, fileWant, err2)
				}
			}
		})
		return nil
	})

	if err != nil {
		t.Errorf("filepath.Walk with error:%s", err.Error())
	}

	if checkFileTotal < 1 {
		t.Fatalf("checkFileTotal==0")
	}
}

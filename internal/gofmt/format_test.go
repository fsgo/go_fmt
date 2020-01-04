/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/3
 */

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
		LocalPrefix:  "auto",
		Write:        false,
		MergeImports: true,
	}
	runTest(t, "rule1", opt)
}

func TestFormat_rule2(t *testing.T) {
	opt := &Options{
		TabIndent:    true,
		TabWidth:     8,
		LocalPrefix:  "auto",
		Write:        false,
		MergeImports: false,
	}
	runTest(t, "rule2", opt)
}

func runTest(t *testing.T, ruleDirName string, opt *Options) {
	rule1Dir := "./testdata/" + ruleDirName + "/"

	var checkFileTotal int

	filepath.Walk(rule1Dir+"/input/", func(path string, info os.FileInfo, err error) error {
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
					t.Errorf("not eq")
					ioutil.WriteFile(rule1Dir+"/tmp/"+filepath.Base(path)+".got", got, 0644)
					ioutil.WriteFile(rule1Dir+"/tmp/"+filepath.Base(path)+".want", want, 0644)
				}
			}
		})
		return nil
	})

	if checkFileTotal < 1 {
		t.Fatalf("checkFileTotal==0")
	}
}

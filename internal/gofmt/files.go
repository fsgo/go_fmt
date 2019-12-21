/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/21
 */

package gofmt

import (
	"os"
	"path/filepath"
	"strings"
)

// filesCurrentDirAll 获取当前目录所有的.go 文件
func filesCurrentDirAll() ([]string, error) {
	var goFiles []string
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err == nil && isGoFile(info) {
			goFiles = append(goFiles, path)
		}
		return err
	})
	return goFiles, err
}

// filesGitDirChange 获取当前git项目目录所有有修改问.go文件
func filesGitDirChange() ([]string, error) {
	gitFiles, err := GitChangeFiles()
	if err != nil {
		return nil, err
	}
	var goFiles []string
	for _, fileName := range gitFiles {
		if isGoFileName(fileName) {
			goFiles = append(goFiles, fileName)
		}
	}
	return goFiles, nil
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func isGoFileName(fileName string) bool {
	info, err := os.Stat(fileName)
	if err != nil {
		return false
	}

	return isGoFile(info)
}

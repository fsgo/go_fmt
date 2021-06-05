// Copyright(C) 2021 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2021/1/15

package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// allGoFiles 获取当前目录所有的.go 文件
func allGoFiles(dir string) ([]string, error) {
	var goFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && isGoFile(info) && !isIgnorePath(path) {
			goFiles = append(goFiles, path)
		}
		return err
	})
	return goFiles, err
}

var ignoreDirNames = []string{
	"testdata",
	"vendor",
}

// isIgnorePath 判断一个路径是否应该忽略掉
func isIgnorePath(path string) bool {
	for _, name := range ignoreDirNames {
		ignoreDir := fmt.Sprintf("%s%c", name, filepath.Separator)
		if strings.HasPrefix(path, ignoreDir) {
			return true
		}
		ignoreDir = fmt.Sprintf("%c%s", filepath.Separator, ignoreDir)
		if strings.Contains(path, ignoreDir) {
			return true
		}
	}
	return false
}

// filesGitDirChange 获取当前git项目目录所有有修改问.go文件
func filesGitDirChange() ([]string, error) {
	gitFiles, err := GitChangeFiles()
	if err != nil {
		return nil, err
	}
	var goFiles []string
	for _, fileName := range gitFiles {
		if isGoFileName(fileName) && !isIgnorePath(fileName) {
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

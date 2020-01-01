/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package gofmt

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GitChangeFiles 获取当前git项目有变更的文件
// 支持如下git状态:
// M  auth/md5_sign.go
// R  utils/counter_test.go -> component/counter/counter_test.go
// A  unittest/internal/monitor/bvar.apis_monitor.data
// ?? internal/gofmt/files.go
func GitChangeFiles() ([]string, error) {
	data, err := exec.Command("git", "status", "-s").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("exec ( git status -s ) with error: %s", err.Error())
	}
	data = bytes.TrimSpace(data)
	lines := bytes.Split(data, []byte("\n"))
	var files []string
	for _, line := range lines {
		// 非.go文件
		if !bytes.Contains(line, []byte(".go")) {
			continue
		}
		// 删除的文件
		if bytes.HasPrefix(line, []byte("D")) {
			continue
		}

		arr := bytes.Split(line, []byte(" "))

		fileName := string(arr[len(arr)-1])

		// 是上一级目录的文件
		if strings.HasPrefix(fileName, "..") {
			continue
		}

		files = append(files, fileName)
	}
	return files, nil
}

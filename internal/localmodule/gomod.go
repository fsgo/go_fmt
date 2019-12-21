/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/21
 */

package localmodule

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
)

// 通过项目的go.mod 文件来获取项目的module值
func detectByGoMod(_ string) (string, error) {
	data, err := exec.Command("go", "env", "GOMOD").CombinedOutput()
	if err != nil {
		return "", err
	}
	goModPath := string(bytes.TrimSpace(data))
	if goModPath == "" {
		return "", fmt.Errorf("go.mod not found")
	}
	goModBuf, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return "", err
	}
	lines := bytes.Split(goModBuf, []byte("\n"))
	if len(lines) < 1 {
		return "", fmt.Errorf("parser %s failed, no contents", goModPath)
	}
	firstLine := lines[0]
	if !bytes.HasPrefix(firstLine, []byte("module ")) {
		return "", fmt.Errorf("parser %s failed, first line not start with 'module '", goModPath)
	}
	module := string(bytes.TrimSpace(firstLine[len("module "):]))
	if module == "" {
		return "", fmt.Errorf("parser %s failed,module value is empty", goModPath)
	}
	return module, nil
}

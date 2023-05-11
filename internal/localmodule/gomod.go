// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/21

package localmodule

import (
	"log"

	"github.com/fsgo/go_fmt/internal/common"
)

// 通过项目的 go.mod 文件来获取项目的 module 值
func detectByGoMod(opt common.Options, fileName string) (string, error) {
	goModPath, err := common.FindGoModPath(fileName)
	if opt.Trace {
		log.Println("detect go.module, file=", fileName, "go.mod=", goModPath, "err=", err)
	}
	if err != nil {
		return "", err
	}
	m, err := common.ModuleByFile(goModPath)
	if opt.Trace {
		log.Println("parser ", goModPath, "module=", m, ", err=", err)
	}
	return m, err
}

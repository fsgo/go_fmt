// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/21

package localmodule

import (
	"fmt"

	"github.com/fsgo/go_fmt/internal/common"
)

// 通过项目的go.mod 文件来获取项目的module值
func detectByGoMod(opt *common.Options, fileName string) (string, error) {
	goModPath, err := common.FindGoModPath(fileName)
	if opt.Trace {
		fmt.Println("detect go.module, file=", fileName, "go.mod=", goModPath, "err=", err)
	}
	if err != nil {
		return "", err
	}
	return common.ModuleByFile(goModPath)
}

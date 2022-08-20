// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/21

package localmodule

import (
	"errors"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

var detectAutoFuncs = []func(opt *common.Options, fileName string) (string, error){
	detectByGoMod, // 最高优先级，使用 go.mod 自动获取
	detectByGoPath,
	detectFailBack,
}

// 当代码没有在一个其他目录的情况下的时候(没有 go.mod，也不在 GOPATH/src 中)
func detectFailBack(opt *common.Options, fileName string) (string, error) {
	return "_unknown_module", nil
}

// Get 自动推断当前项目地址
func Get(opt *common.Options, fileName string) (string, error) {
	if err := loadConfig(); err != nil {
		return "", err
	}
	if opt.LocalModule == "auto" {
		var builder strings.Builder
		for _, fn := range detectAutoFuncs {
			val, err := fn(opt, fileName)
			if err == nil {
				return val, nil
			}
			builder.WriteString(err.Error())
			builder.WriteString(";\n")
		}
		return "", errors.New(builder.String())
	}
	return opt.LocalModule, nil
}

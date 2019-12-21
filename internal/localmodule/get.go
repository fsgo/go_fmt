/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/21
 */

package localmodule

import (
	"errors"
	"strings"
)

var detectAutoFuncs = []func(fileName string) (string, error){
	detectByGoMod, // 最高优先级，使用go.mod 自动获取
	detectByGoPath,
}

// Get 自动推断当前项目地址
func Get(local string, fileName string) (string, error) {
	if err := loadConfig(); err != nil {
		return "", err
	}
	if local == "auto" {
		var builder strings.Builder
		for _, fn := range detectAutoFuncs {
			val, err := fn(fileName)
			if err == nil {
				return val, nil
			}
			builder.WriteString(err.Error())
			builder.WriteString(";\n")
		}
		return "", errors.New(builder.String())
	}
	return local, nil
}

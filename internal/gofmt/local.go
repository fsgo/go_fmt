/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package gofmt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var localInited bool

// DetectLocal 自动推断当前项目地址
func DetectLocal(local string, fileName string) (string, error) {
	if !localInited {
		if err := userLocalConfigLoad(); err != nil {
			return "", err
		}
	}
	localInited = true

	if local == "auto" {
		// 通过文件绝对路径，自动推断当前项目的地址
		// 如当前文件推断出的local 值为 github.com/fsgo/go_fmt
		srcDirName := fmt.Sprintf("%csrc%c", filepath.Separator, filepath.Separator)

		// 若文件地址不是绝对地址，需要对地址进行补全，后续才可使用
		if !filepath.IsAbs(fileName) {
			wd, err := os.Getwd()
			if err != nil {
				return "", err
			}
			fileName = filepath.Join(wd, fileName)
		}

		infos := strings.SplitN(filepath.Clean(fileName), srcDirName, 2)
		if len(infos) != 2 || len(infos[1]) == 0 {
			return "", nil
		}
		dirs := strings.Split(infos[1], string(filepath.Separator))
		level := domainLevel[dirs[0]]

		var paths []string
		paths = append(paths, dirs[0])

		// 若是域名 并且不在未配置规则，则使用默认值
		if level == 0 && strings.Contains(dirs[0], ".") {
			level = defaultDomainLevel
		}

		if level > 0 {
			for i := 0; i < level; i++ {
				paths = append(paths, dirs[i+1])
			}
		}

		return strings.Join(paths, "/"), nil
	}
	return local, nil
}

// 各个域名
var domainLevel = map[string]int{
	"github.com":      2,
	"icode.baidu.com": 3,
	"golang.org":      2,
}

var defaultDomainLevel int = 2

// 加载用户自定义配置,配置内容为json，
// 比如：
// {"abc.com":3}
func userLocalConfigLoad() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	confPath := filepath.Join(home, ".go_fmt", "auto_local.json")
	_, err = os.Stat(confPath)
	if os.IsNotExist(err) {
		return nil
	}
	confBuf, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	var data map[string]int
	if err := json.Unmarshal(confBuf, &data); err != nil {
		return nil
	}

	for k, v := range data {
		domainLevel[k] = v
	}
	return nil
}

/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/21
 */

package localmodule

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	// 域名对应目录层级
	DomainLevel map[string]int
}

var configLoaded bool

// 加载用户自定义配置,配置内容为json，
// 比如：
// {"abc.com":3}
func loadConfig() error {
	// 只需要加载一次
	if configLoaded {
		return nil
	}
	configLoaded = true

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	confPath := filepath.Join(home, ".go_fmt", "local_module.json")
	_, err = os.Stat(confPath)
	// 若配置文件不存在则 跳过
	if os.IsNotExist(err) {
		return nil
	}

	confBuf, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	var conf *config
	if err := json.Unmarshal(confBuf, &conf); err != nil {
		return nil
	}

	if err := parserConfig4GoPath(conf); err != nil {
		return err
	}

	return nil
}

// Copyright(C) 2019 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2019/12/21

package localmodule

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

// detectByGoPath 通过文件绝对路径来自动推断当前项目的模块名
// 项目的workspace 满足如下 https://golang.org/doc/code.html#Workspaces
// | bin/
// |     go_fmt                               # command executable
// | src/
// |    github.com/fsgo/go_fmt/
// |                           .git/          # Git repository metadata
// |                            go_fmt.go     # command source
// 基本原理为：
// 1.获取文件的完整路径
// 2.查找src目录，其后第一个目录为域名
// 3.默认域名再后面2级目录则为项目的模块名，如github上所有项目。
// 4.若比较特殊，不是2级目录，也可以通过配置文件（~/.go_fmt/local_module.json）来设置
func detectByGoPath(opt *common.Options, fileName string) (string, error) {
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
		return "", errors.New("project not in GOPATH")
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

// 各个域名
var domainLevel = map[string]int{
	"github.com":      2,
	"icode.baidu.com": 3,
	"golang.org":      2,
}

var defaultDomainLevel int = 2

func parserConfig4GoPath(conf *config) error {
	if len(conf.DomainLevel) == 0 {
		return nil
	}

	for domain, level := range conf.DomainLevel {
		domainLevel[domain] = level
	}
	return nil
}

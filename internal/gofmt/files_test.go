/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/3
 */

package gofmt

import (
	"strings"
	"testing"
)

func Test_currentDirAllGoFiles(t *testing.T) {
	files, err := currentDirAllGoFiles()
	if err != nil {
		t.Fatalf("currentDirAllGoFiles with error:%s", err)
	}

	for _, fpath := range files {
		if strings.Contains(fpath, "testdata") {
			t.Errorf("fpath=%s should ignored", fpath)
		}
	}
}

/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/4
 */

package gofmt

import (
	"testing"
)

func TestGitChangeFiles(t *testing.T) {
	_, err := GitChangeFiles()
	if err != nil {
		t.Fatalf("GitChangeFiles() with error:%s", err)
	}
}

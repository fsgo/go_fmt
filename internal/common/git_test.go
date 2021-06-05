// Copyright(C) 2021 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2021/1/15

package common

import (
	"testing"
)

func TestGitChangeFiles(t *testing.T) {
	_, err := GitChangeFiles()
	if err != nil {
		t.Fatalf("GitChangeFiles() with error:%s", err)
	}
}

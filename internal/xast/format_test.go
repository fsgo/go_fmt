// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/2

package xast

import (
	"testing"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func TestFormat(t *testing.T) {
	testTestFormat(t, "testdata/case1.go.input")
	testTestFormat(t, "testdata/case2.go.input")
	testTestFormat(t, "testdata/case3.go.input")
	testTestFormat(t, "testdata/case4.go.input")
}

func testTestFormat(t *testing.T, input string) {
	t.Helper()
	xtest.CheckFileAuto(t, input, func(req *common.Request) {
		Format(req)
	})
}

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
	testTestFormat(t, "testdata/case1.go.input", "testdata/case1.go.want")
	testTestFormat(t, "testdata/case2.go.input", "testdata/case2.go.want")
	testTestFormat(t, "testdata/case3.go.input", "testdata/case3.go.want")
}

func testTestFormat(t *testing.T, input string, want string) {
	t.Helper()
	xtest.Check(t, input, want, func(req *common.Request) {
		Format(req)
	})
}

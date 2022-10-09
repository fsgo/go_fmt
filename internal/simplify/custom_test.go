// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"testing"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func Test_customSimplify(t *testing.T) {
	xtest.Check(t, "testdata/custom1.go.input", "testdata/custom1.go.want", func(req *common.Request) {
		customSimplify(req.AstFile)
	})
}

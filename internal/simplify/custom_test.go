// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"fmt"
	"testing"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func Test_customSimplify(t *testing.T) {
	for i := 1; i <= 8; i++ {
		input := fmt.Sprintf("testdata/custom%d.go.input", i)
		want := fmt.Sprintf("testdata/custom%d.go.want", i)
		xtest.Check(t, input, want, func(req *common.Request) {
			customSimplify(req)
		})
	}
}

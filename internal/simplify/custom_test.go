// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/9

package simplify

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func Test_customSimplify(t *testing.T) {
	fs, err := filepath.Glob("testdata/custom*.go.input")
	require.NoError(t, err)
	for _, input := range fs {
		xtest.CheckAuto(t, input, func(req *common.Request) {
			customSimplify(req)
		})
	}
}

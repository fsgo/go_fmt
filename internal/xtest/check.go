// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/2

package xtest

import (
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fsgo/go_fmt/internal/common"
)

// Check 运行测试 case
func Check(t *testing.T, inputFile, wantFile string, do func(req *common.Request)) {
	t.Run(inputFile, func(t *testing.T) {
		t.Helper()
		t.Logf("Check %q", inputFile)

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, inputFile, nil, parser.ParseComments)
		require.NoError(t, err)
		req := &common.Request{
			FileName: inputFile,
			AstFile:  f,
			FSet:     fset,
			Opt:      *common.NewDefaultOptions(),
		}
		do(req)
		req.MustReParse() // 重新格式化
		code, err := req.FormatFile()
		require.NoError(t, err)

		if len(wantFile) > 0 {
			bw, _ := os.ReadFile(wantFile)
			gf := wantFile + ".got"
			if assert.Equal(t, string(bw), string(code)) {
				_ = os.Remove(gf)
			} else {
				t.Logf("want: %s", gf)
				e2 := os.WriteFile(gf, code, 0644)
				require.NoError(t, e2)
			}
			require.NotEmpty(t, bw)
		} else {
			t.Log("wantFile is empty, skipped check result file")
		}
	})
}

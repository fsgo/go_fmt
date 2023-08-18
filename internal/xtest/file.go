// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/2

package xtest

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xpasser"
)

func TestGoFileName(name string) string {
	return strings.TrimSuffix(name, ".input")
}

// CheckFileAuto inputFile 必须以 .input 结尾；wantFile 则自动依据 inputFile 推断
// 如 inputFile=demo.go.input 则 wantFile=demo.go.want
func CheckFileAuto(t *testing.T, inputFile string, do func(req *common.Request)) {
	suf := ".input"
	require.True(t, strings.HasSuffix(inputFile, suf))
	wantFile := inputFile[:len(inputFile)-len(suf)] + ".want"
	CheckFile(t, inputFile, wantFile, do)
}

// CheckFile 运行测试,检查单个 case 文件
func CheckFile(t *testing.T, inputFile, wantFile string, do func(req *common.Request)) {
	opt := common.NewDefaultOptions()
	t.Run(inputFile, func(t *testing.T) {
		t.Helper()
		defer xpasser.Reset()
		t.Logf("Check %q", inputFile)

		fileContent, err := os.ReadFile(inputFile)
		require.NoError(t, err)

		name1 := TestGoFileName(inputFile)

		require.NoError(t, xpasser.LoadOverlay(name1, fileContent))
		require.NoError(t, xpasser.Load(*opt, []string{"file=" + name1}))
		asfFile, err := xpasser.ParserFile(name1, fileContent)
		require.NoError(t, err)

		req := &common.Request{
			FileName: name1,
			AstFile:  asfFile,
			FSet:     xpasser.Default.FSet,
			Opt:      *opt,
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

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/8

package simplify

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func TestFormat(t *testing.T) {
	xtest.Check(t, "testdata/fmt1.go.input", "testdata/fmt1.go.want", func(req *common.Request) {
		Format(req)
	})
}

func TestRewrite(t *testing.T) {
	fn1 := func(req *common.Request) {
		f, err := Rewrite(req, "io/#ioutil.WriteFile -> os.WriteFile")
		require.NoError(t, err)
		req.AstFile = f
	}
	xtest.Check(t, "testdata/rewrite1.go.input", "testdata/rewrite1.go.want", fn1)
	xtest.Check(t, "testdata/rewrite2.go.input", "testdata/rewrite2.go.want", fn1)
}

func TestRewrites(t *testing.T) {
	rules := []string{
		"io/#ioutil.NopCloser -> io.NopCloser",
		"io/#ioutil.ReadAll -> io.ReadAll",
		"io/#ioutil.ReadFile -> os.ReadFile",
		"io/#ioutil.TempFile -> os.CreateTemp",
		"io/#ioutil.TempDir -> os.MkdirTemp",
		"io/#ioutil.WriteFile -> os.WriteFile",
	}
	fn1 := func(req *common.Request) {
		err := Rewrites(req, rules)
		require.NoError(t, err)
	}
	xtest.Check(t, "testdata/rewrite3.go.input", "testdata/rewrite3.go.want", fn1)

	xtest.Check(t, "testdata/rewrite4.go.input", "testdata/rewrite4.go.want", func(req *common.Request) {
		err := Rewrites(req, BuildInRewriteRules())
		require.NoError(t, err)
	})
}
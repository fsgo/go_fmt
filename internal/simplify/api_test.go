// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/8

package simplify

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func TestFormat(t *testing.T) {
	xtest.CheckAuto(t, "testdata/fmt1.go.input", func(req *common.Request) {
		Format(req)
	})
}

func TestRewrite(t *testing.T) {
	fn1 := func(req *common.Request) {
		f, err := Rewrite(req, "io/#ioutil.WriteFile -> os.WriteFile")
		require.NoError(t, err)
		req.AstFile = f
	}
	xtest.CheckAuto(t, "testdata/rewrite1.go.input", fn1)
	xtest.CheckAuto(t, "testdata/rewrite2.go.input", fn1)

	fn2 := func(req *common.Request) {
		f, err := Rewrite(req, "testdata/rules/rule1.txt")
		require.NoError(t, err)
		req.AstFile = f
	}
	xtest.CheckAuto(t, "testdata/rewrite5.go.input", fn2)

	t.Run("invalid rule", func(t *testing.T) {
		req := common.NewTestRequest("testdata/rewrite5.go.input")
		rules := []string{"invalid", "invalid ->", " -> ", "", " "}
		for i := 0; i < len(rules); i++ {
			f, err := Rewrite(req, rules[i])
			require.Error(t, err)
			require.Nil(t, f)
		}
	})

	t.Run("invalid rule file", func(t *testing.T) {
		req := common.NewTestRequest("testdata/rewrite5.go.input")
		f, err := Rewrite(req, "testdata/rules/rule2.txt")
		require.Error(t, err)
		require.Nil(t, f)
	})
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
	xtest.CheckAuto(t, "testdata/rewrite3.go.input", fn1)

	t.Run("build in", func(t *testing.T) {
		buildInCases,err:=filepath.Glob("testdata/buildin*.go.input")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(buildInCases),3)
		for _,bc:=range buildInCases{
			xtest.CheckAuto(t, bc, func(req *common.Request) {
				err := Rewrites(req, BuildInRewriteRules())
				require.NoError(t, err)
			})
		}
	})
}

func Test_goVersionFromComment(t *testing.T) {
	type args struct {
		comment string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{},
			want: "",
		},
		{
			name: "go1.19",
			args: args{
				comment: "// go1.19",
			},
			want: "1.19",
		},
		{
			name: "go1.19 with other",
			args: args{
				comment: "// go1.19 other",
			},
			want: "1.19",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := goVersionFromComment(tt.args.comment); got != tt.want {
				t.Errorf("goVersionFromComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

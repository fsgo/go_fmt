// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/2

package ximports

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fsgo/go_fmt/internal/common"
)

func Test_sortImportDecls(t *testing.T) {
	type args struct {
		decls   []*importDecl
		options common.Options
	}
	tests := []struct {
		name string
		args args
		want importDeclGroups
	}{
		{
			name: "case 1",
			args: args{
				decls: []*importDecl{
					{
						Docs: nil,
						Path: `"github.com/b"`,
					},
					{
						Docs: nil,
						Path: `"a.com/a"`,
					},
					{
						Docs: nil,
						Path: `"github.com/a"`,
					},
					{
						Docs: nil,
						Path: `"fmt"`,
					},
					{
						Docs: []string{
							"//a",
						},
						Path: ``,
					},
					{
						Docs: nil,
						Path: `_ "net/http"`,
					},
				},
				options: common.Options{
					LocalModule: "github.com/a",
				},
			},
			want: importDeclGroups{
				{
					Group: 0,
					Decls: []*importDecl{
						{
							Docs: nil,
							Path: `"fmt"`,
						},
						{
							Docs: []string{
								"//a",
							},
							Path: ``,
						},
						{
							Docs: nil,
							Path: `_ "net/http"`,
						},
					},
				},
				{
					Group: 1,
					Decls: []*importDecl{
						{
							Docs: nil,
							Path: `"github.com/b"`,
						},
						{
							Docs: nil,
							Path: `"a.com/a"`,
						},
					},
				},
				{
					Group: 2,
					Decls: []*importDecl{
						{
							Docs: nil,
							Path: `"github.com/a"`,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortImportDecls(tt.args.decls, tt.args.options).String()
			want := tt.want.String()
			require.Equal(t, want, got)
		})
	}
}

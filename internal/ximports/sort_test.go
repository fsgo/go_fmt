/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/2
 */

package ximports

import (
	"reflect"
	"testing"

	"github.com/fsgo/go_fmt/internal/common"
)

func Test_sortImportDecls(t *testing.T) {
	type args struct {
		decls   []*importDecl
		options *common.Options
	}
	tests := []struct {
		name string
		args args
		want []*importDeclGroup
	}{
		{
			name: "case 1",
			args: args{
				decls: []*importDecl{
					{
						Comments: nil,
						Path:     `"github.com/b"`,
					},
					{
						Comments: nil,
						Path:     `"a.com/a"`,
					},
					{
						Comments: nil,
						Path:     `"github.com/a"`,
					},
					{
						Comments: nil,
						Path:     `"fmt"`,
					},
					{
						Comments: []string{
							"//a",
						},
						Path: ``,
					},
				},
				options: &common.Options{
					LocalPrefix: "github.com/a",
				},
			},
			want: []*importDeclGroup{
				{
					Group: 0,
					Decls: []*importDecl{
						{
							Comments: nil,
							Path:     `"fmt"`,
						},
						{
							Comments: []string{
								"//a",
							},
							Path: ``,
						},
					},
				},
				{
					Group: 1,
					Decls: []*importDecl{
						{
							Comments: nil,
							Path:     `"github.com/b"`,
						},
						{
							Comments: nil,
							Path:     `"a.com/a"`,
						},
					},
				},
				{
					Group: 2,
					Decls: []*importDecl{
						{
							Comments: nil,
							Path:     `"github.com/a"`,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortImportDecls(tt.args.decls, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortImportDecls() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/2

package ximports

import (
	"testing"
)

func Test_importDecl_RealPath(t *testing.T) {
	type fields struct {
		Comments []string
		Path     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "case 1",
			fields: fields{
				Comments: []string{
					`"github.com/fsgo/cache"`,
				},
				Path: `"github.com/fsgo/go_fmt"`,
			},
			want: `github.com/fsgo/go_fmt`,
		},
		{
			name: "case 2",
			fields: fields{
				Comments: []string{
					`"github.com/fsgo/cache"`,
				},
				Path: `go_fmt "github.com/fsgo/go_fmt"`,
			},
			want: `github.com/fsgo/go_fmt`,
		},
		{
			name: "case 3",
			fields: fields{
				Comments: []string{
					`"github.com/fsgo/cache"`,
				},
				Path: `_ "github.com/fsgo/go_fmt"`,
			},
			want: `github.com/fsgo/go_fmt`,
		},
		{
			name: "case 4",
			fields: fields{
				Comments: []string{
					`"github.com/fsgo/cache"`,
				},
				Path: `_"github.com/fsgo/go_fmt"`,
			},
			want: `github.com/fsgo/go_fmt`,
		},
		{
			name: "case 5",
			fields: fields{
				Comments: []string{
					`"github.com/fsgo/cache"`,
				},
				Path: `gofmt"github.com/fsgo/go_fmt"`,
			},
			want: `github.com/fsgo/go_fmt`,
		},
		{
			name: "case 6",
			fields: fields{
				Comments: []string{
					`// "github.com/fsgo/cache"`,
				},
			},
			want: `github.com/fsgo/cache`,
		},
		{
			name: "case 7",
			fields: fields{
				Comments: []string{
					`// 这个是注释`,
					`// 这个是注释`,
					`//"github.com/fsgo/cache"`,
				},
			},
			want: `github.com/fsgo/cache`,
		},
		{
			name: "case 8",
			fields: fields{
				Comments: []string{
					`//"github.com/fsgo/cache" //注释`,
				},
			},
			want: `github.com/fsgo/cache`,
		},
		{
			name: "case 9",
			fields: fields{
				Path: `"fmt" //注释`,
			},
			want: `fmt`,
		},
		{
			name: "case 10",
			fields: fields{
				Path: `_ "net/http" //注释`,
			},
			want: `net/http`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decl := &importDecl{
				Comments: tt.fields.Comments,
				Path:     tt.fields.Path,
			}
			if got := decl.RealPath(); got != tt.want {
				t.Errorf("RealPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

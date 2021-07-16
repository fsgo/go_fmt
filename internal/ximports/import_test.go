// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/2

package ximports

import (
	"testing"
)

func Test_isImportPathLine(t *testing.T) {
	type args struct {
		bf []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				bf: []byte(""),
			},
			want: false,
		},
		{
			name: "case 2",
			args: args{
				bf: []byte(`// "github.com"`),
			},
			want: false,
		},
		{
			name: "case 3",
			args: args{
				bf: []byte(`// a "github.com"`),
			},
			want: false,
		},
		{
			name: "case 4-是注释",
			args: args{
				bf: []byte(`/*a "github.com"*/`),
			},
			want: false,
		},
		{
			name: "case 5-有换行符",
			args: args{
				bf: []byte("/*a \n\"github.com\"\n*/"),
			},
			want: false,
		},
		{
			name: "case 6-有换行符",
			args: args{
				bf: []byte("\"github.\ncom/a\""),
			},
			want: false,
		},
		{
			name: "case 7",
			args: args{
				bf: []byte(`"github.com/a"`),
			},
			want: true,
		},
		{
			name: "case 8",
			args: args{
				bf: []byte(`a "github.com/a" `),
			},
			want: true,
		},
		{
			name: "case 9",
			args: args{
				bf: []byte(`_ "github.com/a" `),
			},
			want: true,
		},
		{
			name: "case 10",
			args: args{
				bf: []byte(`_"github.com/a" `),
			},
			want: true,
		},
		{
			name: "case 11",
			args: args{
				bf: []byte(`git"github.com/a" `),
			},
			want: true,
		},
		{
			name: "case 12-后引号不匹配",
			args: args{
				bf: []byte(`"github.com/a"汉字`),
			},
			want: false,
		},
		{
			name: "case 10-不允许汉字",
			args: args{
				bf: []byte(`"汉字"`),
			},
			want: false,
		},
		{
			name: "case 11",
			args: args{
				bf: []byte(`a汉字 "b汉字c"`),
			},
			want: false,
		},
		{
			name: "case 12",
			args: args{
				bf: []byte(`git "../github.com/a" `),
			},
			want: true,
		},
		{
			name: "case 13",
			args: args{
				bf: []byte(`git "../../github_b/a123/a" `),
			},
			want: true,
		},
		{
			name: "case 14-少引号",
			args: args{
				bf: []byte(`git "../../github_b/a123/a `),
			},
			want: false,
		},
		{
			name: "case 15",
			args: args{
				bf: []byte(`"fmt" //`),
			},
			want: true,
		},
		{
			name: "case 16",
			args: args{
				bf: []byte(`"github.com/go-playground/locales/en"`),
			},
			want: true,
		},
		{
			name: "case 17",
			args: args{
				bf: []byte(`validator_engine "gopkg.in/go-playground/validator.v9"`),
			},
			want: true,
		},
		{
			name: "case 18",
			args: args{
				bf: []byte(`. "github.com/onsi/gomega"`),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isImportPathLine(tt.args.bf); got != tt.want {
				t.Errorf("isImportPathLine(%q) = %v, want %v", tt.args.bf, got, tt.want)
			}
		})
	}
}

// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/4

package pkgs

import (
	"testing"
)

func TestIsStd(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				path: "fmt",
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				path: "net/http",
			},
			want: true,
		},
		{
			name: "case 3",
			args: args{
				path: "share",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStd(tt.args.path); got != tt.want {
				t.Errorf("IsStd() = %v, want %v", got, tt.want)
			}
		})
	}
}

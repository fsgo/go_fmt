/*
 * Copyright(C) 2021 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2021/1/15
 */

package common

import (
	"strings"
	"testing"
)

func Test_currentDirAllGoFiles(t *testing.T) {
	files, err := currentDirAllGoFiles()
	if err != nil {
		t.Fatalf("currentDirAllGoFiles with error:%s", err)
	}

	for _, fpath := range files {
		if strings.Contains(fpath, "testdata") {
			t.Errorf("fpath=%s should ignored", fpath)
		}
	}
}

func Test_isGoFileName(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				"files.go",
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				"files_test.go",
			},
			want: true,
		},
		{
			name: "case 3",
			args: args{
				"abc/not_exists.go",
			},
			want: false,
		},
		{
			name: "case 4",
			args: args{
				"testdata/rule1/input/demo_1.go.txt",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGoFileName(tt.args.fileName); got != tt.want {
				t.Errorf("isGoFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/3/13
 */

package common

import (
	"testing"
)

func TestDoNotEdit(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{},
			want: false,
		},
		{
			name: "case 2",
			args: args{
				src: []byte(""),
			},
			want: false,
		},
		{
			name: "case 3",
			args: args{
				src: []byte("abcd"),
			},
			want: false,
		},
		{
			name: "case 4",
			args: args{
				src: []byte(`// Code generated by protoc-gen-go. DO NOT EDIT.`),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DoNotEdit(tt.args.src); got != tt.want {
				t.Errorf("DoNotEdit() = %v, want %v", got, tt.want)
			}
		})
	}
}

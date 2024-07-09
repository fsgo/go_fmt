// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/5

package common

import (
	"reflect"
	"testing"

	"github.com/fsgo/fst"
)

func Test_intSliceDelete(t *testing.T) {
	type args struct {
		lines  []int
		delete []int
	}

	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case 1",
			args: args{
				lines:  []int{1, 2, 5, 6},
				delete: []int{1, 3},
			},
			want: []int{1, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intSliceDelete(tt.args.lines, tt.args.delete...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intSliceDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_IsGoVersion(t *testing.T) {
	req := &Request{
		FileName: "request.go",
	}
	fst.True(t, req.GoVersionGEQ("1.13"))
	fst.True(t, req.GoVersionGEQ("1.19"))
	fst.True(t, req.GoVersionGEQ("1.18"))

	fst.False(t, req.GoVersionGEQ("1.80"))
	fst.False(t, req.GoVersionGEQ("1.99"))
}

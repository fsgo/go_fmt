// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/5

package gofmt

import (
	"testing"
)

func Test_multilineCommentSingle(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				txt: `/*
abc

def

*/`,
			},
			want: `// abc
// 
// def`,
		},
		{
			name: "case 2",
			args: args{
				txt: `/*
* abc
*
** def
*
*/`,
			},
			want: `// abc
// 
// def`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := multilineCommentSingle(tt.args.txt); got != tt.want {
				t.Errorf("txt=%s \nmultilineCommentSingle() =\n %v,\n want=\n %v", tt.args.txt, got, tt.want)
			}
		})
	}
}

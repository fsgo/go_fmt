// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/3/19

package common

import (
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "same",
			args: args{
				a: "hello",
				b: "hello",
			},
			want: "",
		},
		{
			name: "not same",
			args: args{
				a: "hello",
				b: "world",
			},
			want: "Trace: [0]\nline 1:\n  -: hello\n  +: world",
		},
		{
			name: "df1_del",
			args: args{
				a: "hello\nworld",
				b: "hello",
			},
			want: "Trace: [1->?]\nline 2:\n  -: world",
		},
		{
			name: "df1_add",
			args: args{
				a: "hello",
				b: "hello\nworld",
			},
			want: "Trace: [?->1]\nline 2:\n  +: world",
		},
		{
			name: "df2_add",
			args: args{
				a: strings.Repeat("a\n", 50) + "\nhello",
				b: strings.Repeat("a\n", 50) + "\nhello\nworld",
			},
			want: "Trace: [?->52]\nline 53:\n  +: world",
		},
		{
			name: "df3_add",
			args: args{
				a: "\nhello",
				b: "\nhello\n",
			},
			want: "Trace: [?->2]\nline 3:\n  +: \\n",
		},
		{
			name: "df2_del",
			args: args{
				a: strings.Repeat("a\n", 50) + "\nhello\nworld",
				b: strings.Repeat("a\n", 50) + "\nhello",
			},
			want: "Trace: [52->?]\nline 53:\n  -: world",
		},
		{
			name: "df3_del_space",
			args: args{
				a: strings.Repeat("a\n", 50) + "\nhello\n\n\n",
				b: strings.Repeat("a\n", 50) + "\nhello",
			},
			want: "Trace: [52->?]\nline 53:\n  -: \\n\n\nTrace: [53->?]\nline 54:\n  -: \\n\n\nTrace: [54->?]\nline 55:\n  -: \\n",
		},
		{
			name: "df4_del_space",
			args: args{
				a: " \n \n\n",
				b: "",
			},
			want: "Trace: [0->?]\nline 1:\n  -: _\\n\n\nTrace: [1->?]\nline 2:\n  -: _\\n\n\nTrace: [2->?]\nline 3:\n  -: \\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := Diff(tt.args.a, tt.args.b, true)
			var got string
			if diff != nil {
				got = diff.String()
			}
			if got != tt.want {
				t.Errorf("\nDiff() = %s \n %q,\n  want = %s \n %q", got, got, tt.want, tt.want)
			}
		})
	}
}

func Test_diffReporter_parserLineNo(t *testing.T) {
	tests := []struct {
		txt  string
		want int
	}{
		{
			txt:  "[28]",
			want: 28,
		},
		{
			txt:  "[31->?]",
			want: 31,
		},
		{
			txt:  "[36->32]",
			want: 36,
		},
		{
			txt:  "[?->39]",
			want: 39,
		},
	}
	for _, tt := range tests {
		t.Run(tt.txt, func(t *testing.T) {
			r := &diffReporter{}
			if got := r.parserLineNo(tt.txt); got != tt.want {
				t.Errorf("parserLineNo(%q) = %v, want %v", tt.txt, got, tt.want)
			}
		})
	}
}

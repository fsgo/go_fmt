// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/3/6

package common

import (
	"reflect"
	"testing"
)

func TestModuleByFile(t *testing.T) {
	type args struct {
		goModPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "mod_ok.txt",
			args: args{
				goModPath: "testdata/mod_ok.txt",
			},
			want: "github.com/fsgo/go_fmt",
		},
		{
			name: "mod_err.txt",
			args: args{
				goModPath: "testdata/mod_err.txt",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ModuleByFile(tt.args.goModPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModuleByFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ModuleByFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListModules(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Modules
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				dir: "testdata/list_modules/empty",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "many",
			args: args{
				dir: "testdata/list_modules",
			},
			want: []string{
				"github.com/test/hello",
				"github.com/test/hello/say",
				"github.com/test/world",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListModules(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListModules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListModules() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModules_In(t *testing.T) {
	type args struct {
		module string
	}
	tests := []struct {
		name string
		ms   Modules
		args args
		want bool
	}{
		{
			name: "in",
			ms: Modules{
				"github.com/test/abc",
			},
			args: args{
				module: "github.com/test/abc/say",
			},
			want: true,
		},
		{
			name: "not in",
			ms: Modules{
				"github.com/test/abc",
			},
			args: args{
				module: "github.com/test/a",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.PkgIn(tt.args.module); got != tt.want {
				t.Errorf("PkgIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInModule(t *testing.T) {
	type args struct {
		pkg    string
		module string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				pkg:    "abc/def",
				module: "abc",
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				pkg:    "abc",
				module: "abc/def",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InModule(tt.args.pkg, tt.args.module); got != tt.want {
				t.Errorf("InModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

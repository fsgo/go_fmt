// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/1/2

package ximports

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/fsgo/fst"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
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

func Test_parserImportSrc(t *testing.T) {
	type args struct {
		src string
	}

	tests := []struct {
		name      string
		args      args
		wantLines []*importDecl
		wantErr   bool
	}{
		{
			name: "case 1",
			args: args{
				src: `
import (
	"reflect"
	"testing"
)`,
			},
			wantLines: []*importDecl{
				{
					Path: `"reflect"`,
				},
				{
					Path: `"testing"`,
				},
			},
		},
		{
			name: "case 2",
			args: args{
				src: `
import (
    /* on reflect 1 */
    // on reflect 2
	"reflect"  /* after reflect */

	// on fmt 1
	// on fmt 2
	"fmt" // after fmt
)`,
			},
			wantLines: []*importDecl{
				{
					Path: `"reflect"  /* after reflect */`,
					Docs: []string{
						`/* on reflect 1 */`,
						`// on reflect 2`,
					},
				},
				{
					Path: `"fmt" // after fmt`,
					Docs: []string{
						`// on fmt 1`,
						`// on fmt 2`,
					},
				},
			},
		},
		{
			name: "case 3",
			args: args{
				src: `
import (
    // cm 0-0   

    /* on reflect 1 */
    // on reflect 2
	 r "reflect"  /* after reflect */

	// cm 2-0

	// on fmt 1
	// on fmt 2
	_ "fmt" // after fmt

   // "http"

   // cm 1-0
   /* cm 1-1 */
)`,
			},
			wantLines: []*importDecl{
				{
					Path: `r "reflect"  /* after reflect */`,
					Docs: []string{
						`// cm 0-0`,
						`/* on reflect 1 */`,
						`// on reflect 2`,
					},
				},
				{
					Path: `_ "fmt" // after fmt`,
					Docs: []string{
						`// cm 2-0`,
						`// on fmt 1`,
						`// on fmt 2`,
					},
				},
				{
					Path: ``,
					Docs: []string{
						`// "http"`,
					},
				},
				{
					Path: ``,
					Docs: []string{
						`// cm 1-0`,
						`/* cm 1-1 */`,
					},
				},
			},
		},
		{
			name: "case 4",
			args: args{
				src: `
import "fmt"
`,
			},
			wantLines: []*importDecl{
				{
					Path: `"fmt"`,
				},
			},
		},
		{
			name: "case 5",
			args: args{
				src: `
import "fmt" // after fmt 
`,
			},
			wantLines: []*importDecl{
				{
					Path: `"fmt" // after fmt`,
				},
			},
		},
		{
			name: "case 6",
			args: args{
				src: `
import (
	// on fmt
	"fmt"
	"log"
	// on net
	"net" // after net
	"github.com/go_fmt/app2/internal"
	"golang.org/x/mod/modfile"
	_ "net/http" // after http
)
`,
			},
			wantLines: []*importDecl{
				{
					Path: `"fmt"`,
					Docs: []string{
						`// on fmt`,
					},
				},
				{
					Path: `"log"`,
				},
				{
					Path: `"net" // after net`,
					Docs: []string{
						`// on net`,
					},
				},
				{
					Path: `"github.com/go_fmt/app2/internal"`,
				},
				{
					Path: `"golang.org/x/mod/modfile"`,
				},
				{
					Path: `_ "net/http" // after http`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := strings.TrimSpace(tt.args.src)
			gotLines, err := parserImportSrc([]byte(src))
			if (err != nil) != tt.wantErr {
				t.Errorf("parserImportSrc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("parserImportSrc() \ngot  = %v, \nwant = %v", gotLines, tt.wantLines)
			}
		})
	}
}

func TestFormatImports(t *testing.T) {
	ms, err := filepath.Glob("./testdata/*.input")
	fst.NoError(t, err)
	fst.NotEmpty(t, ms)
	for i := 0; i < len(ms); i++ {
		fp := ms[i]
		t.Run(fp, func(t *testing.T) {
			tmpPath := fp + ".got"
			_ = os.Remove(tmpPath)

			src, err := os.ReadFile(fp)
			fst.NoError(t, err)
			fs, af, err := common.ParseOneFile(fp, src)

			fst.NoError(t, err)
			req := &common.Request{
				FileName: fp,
				FSet:     fs,
				AstFile:  af,
				Opt: common.Options{
					LocalModule: "github.com/go_fmt/app2",
					TabIndent:   true,
					TabWidth:    8,
				},
			}
			got, err := FormatImports(req)
			fst.NoError(t, err)

			wantFp := fp[:len(fp)-len(".input")] + ".want"
			want, err := os.ReadFile(wantFp)
			fst.NoError(t, err)

			if !bytes.Equal(want, got) {
				_ = os.WriteFile(tmpPath, got, 0644)
				t.Logf("got file=%s", tmpPath)
			}

			fst.Equal(t, string(want), string(got))
		})
	}
}

func TestFormatImports2(t *testing.T) {
	xtest.CheckFile(t, "import.go", "", func(req *common.Request) {
		got, err := FormatImports(req)
		fst.NoError(t, err)
		fst.NotEmpty(t, got)
	})
}

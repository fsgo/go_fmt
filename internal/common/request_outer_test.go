// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/6

package common_test

import (
	"go/ast"
	"testing"

	"github.com/fsgo/fst"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func TestTokenLine(t *testing.T) {
	xtest.CheckFile(t,
		"testdata/request/case1.go.input",
		"testdata/request/case1.go.want",
		func(req *common.Request) {
			req.TokenLine().DeleteLine(0, 13)
			req.TokenLine().DeleteLine(0, 4)
			req.TokenLine().DeleteLine(0, 6)

			decl := req.AstFile.Decls[0].(*ast.GenDecl)
			fst.False(t, req.NoFormat(decl))

			ts := decl.Specs[0].(*ast.TypeSpec)
			sts := ts.Type.(*ast.StructType)
			field := sts.Fields.List[0]
			fst.Equal(t, "Name", field.Names[0].Obj.Name)
			req.TokenLine().AddLine(0, field.End())

			req.TokenLine().Execute()
		})
}

func TestNoFormat(t *testing.T) {
	// 检查文件级别的指令
	xtest.CheckFile(t,
		"testdata/request/case2.go.input",
		"testdata/request/case2.go.want",
		func(req *common.Request) {
			fst.True(t, req.NoFormat(req.AstFile))

			decl := req.AstFile.Decls[0].(*ast.GenDecl)
			fst.True(t, req.NoFormat(decl))
		})
}

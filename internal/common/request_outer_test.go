// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/6

package common_test

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xtest"
)

func TestTokenLine(t *testing.T) {
	xtest.Check(t,
		"testdata/request/case1.go.input",
		"testdata/request/case1.go.want",
		func(req *common.Request) {
			req.TokenLine().DeleteLine(0, 13)
			req.TokenLine().DeleteLine(0, 4)
			req.TokenLine().DeleteLine(0, 6)

			decl := req.AstFile.Decls[0].(*ast.GenDecl)
			ts := decl.Specs[0].(*ast.TypeSpec)
			sts := ts.Type.(*ast.StructType)
			field := sts.Fields.List[0]
			require.Equal(t, "Name", field.Names[0].Obj.Name)
			req.TokenLine().AddLine(0, field.End())

			req.TokenLine().Execute()
		})
}

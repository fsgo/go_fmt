// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/11

package xpasser

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFile_Load(t *testing.T) {
	f := &File{
		FileName: "file.go",
	}
	err := f.Load(nil)
	require.NoError(t, err)
	require.NotNil(t, f)

	var cty *types.Struct
	ast.Inspect(f.AstFile, func(node ast.Node) bool {
		sn, ok := node.(*ast.SelectorExpr)
		if !ok || sn.Sel.Name != "Config" {
			return true
		}
		ty := f.TypesInfo.Types[sn]
		tn := ty.Type.(*types.Named)
		cty = tn.Underlying().(*types.Struct)
		return true
	})

	require.NotNil(t, cty)
	require.Equal(t, "Fset", cty.Field(0).Name())
	require.Equal(t, "ParserMode", cty.Field(1).Name())
}

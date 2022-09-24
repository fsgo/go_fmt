// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/18

package xanalysis

import (
	"go/ast"
	"go/types"
	"log"

	"github.com/fsgo/go_fmt/internal/common"
	"github.com/fsgo/go_fmt/internal/xpasser"
)

func fixStruct(req *common.Request) {
	ast.Inspect(req.AstFile, func(n ast.Node) bool {
		if vt, ok := n.(*ast.CompositeLit); ok && len(vt.Elts) > 0 {
			if p, err := xpasser.FindPackage(req.FileName); err == nil {
				vtp := p.TypesInfo.TypeOf(vt)
				if vts, ok2 := vtp.(*types.Struct); ok2 {
					// for pointer
					// eg: type user struct{name string}
					// u:=&User{"abc"}
					doFixStruct(req, vt, vts)
				} else if vn, ok3 := vtp.(*types.Named); ok3 {
					// for value
					// eg: type user struct{name string}
					// u:=User{"abc"}
					if vts, ok4 := vn.Underlying().(*types.Struct); ok4 {
						doFixStruct(req, vt, vts)
					}
				}
			} else {
				if req.Opt.Trace {
					log.Println("cannot found pkg:", err)
				}
			}
		}
		return true
	})
}

func doFixStruct(req *common.Request, ac *ast.CompositeLit, at *types.Struct) {
	if len(ac.Elts) == 0 {
		return
	}
	switch ac.Elts[0].(type) {
	case *ast.KeyValueExpr:
		return
	}
	for i := 0; i < len(ac.Elts); i++ {
		raw := ac.Elts[i]
		et := at.Field(i)
		v1 := &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name:    et.Name(),
				NamePos: raw.Pos(),
			},
			Value: raw,
		}
		ac.Elts[i] = v1
	}
}

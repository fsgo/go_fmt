// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/10

package simplify

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"log"
)

func fixStructExprNoKey(fileSet *token.FileSet, f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		if vt, ok := n.(*ast.AssignStmt); ok && len(vt.Rhs) > 0 {
			doFixAssignStmt(fileSet, vt)
		}
		if vt, ok := n.(*ast.CompositeLit); ok && len(vt.Elts) > 0 {
			doFixStructKey(fileSet, vt)
		}
		return true
	})
}

// for
//
//	var tests=[]struct{
//	   name string
//	   id string
//	}{
//
//	   {
//	      "name",
//	      "id"
//	   },
//	}
func doFixAssignStmt(fileSet *token.FileSet, n *ast.AssignStmt) {
	// log.Println("before:")
	// ast.Print(fileSet, n)
	// log.Println(strings.Repeat("-", 80))

	for i := 0; i < len(n.Rhs); i++ {
		n1, ok := n.Rhs[i].(*ast.CompositeLit)
		if !ok {
			continue
		}

		var st *ast.StructType
		if n2, ok2 := n1.Type.(*ast.ArrayType); ok2 {
			if n3, ok3 := n2.Elt.(*ast.StructType); ok3 {
				st = n3
			}
		}
		if st == nil {
			continue
		}
		names := structFields(st)
		for j := 0; j < len(n1.Elts); j++ {
			if m1, ok1 := n1.Elts[j].(*ast.CompositeLit); ok1 {
				if len(names) != len(m1.Elts) {
					return
				}
				for z := 0; z < len(m1.Elts); z++ {
					val := m1.Elts[z]
					if _, ok3 := val.(*ast.KeyValueExpr); ok3 {
						break
					}
					// ast.Print(fileSet,m1.Elts[z])

					v1 := &ast.KeyValueExpr{
						Key: &ast.Ident{
							Name:    names[z].Name,
							NamePos: val.Pos(),
						},
						Value: val,
					}
					m1.Elts[z] = v1
				}
			}
		}
	}
	// ast.Print(fileSet, n)
}

func structFields(st *ast.StructType) []*ast.Ident {
	var fs []*ast.Ident
	for i := 0; i < len(st.Fields.List); i++ {
		fs = append(fs, st.Fields.List[i].Names...)
	}
	return fs
}

func doFixStructKey(fileSet *token.FileSet, n *ast.CompositeLit) {
	defer func() {
		if re := recover(); re != nil {
			bf := &bytes.Buffer{}
			ast.Fprint(bf, fileSet, n, nil)
			log.Println("panic:", re, "code=\n", nodeCode(n), "\nast=\n", bf.String())
		}
	}()

	if _, ok := n.Elts[0].(*ast.KeyValueExpr); ok {
		return
	}

	if n.Type == nil {
		return
	}

	var st *ast.StructType
	if n1, ok := n.Type.(*ast.Ident); ok {
		if n1.Obj != nil && n1.Obj.Decl != nil {
			if n2, ok2 := n1.Obj.Decl.(*ast.TypeSpec); ok2 {
				if n2.Type != nil {
					if n3, ok3 := n2.Type.(*ast.StructType); ok3 {
						st = n3
					}
				}
			}
		}
	}

	if st == nil {
		if n1, ok := n.Type.(*ast.ArrayType); ok {
			if n2, ok2 := n1.Elt.(*ast.StructType); ok2 {
				st = n2
			}
		}
	}

	if st == nil {
		return
	}
	names := structFields(st)
	if len(names) != len(n.Elts) {
		return
	}
	for i := 0; i < len(n.Elts); i++ {
		item := n.Elts[i]
		if _, ok := item.(*ast.CompositeLit); ok {
			break
		}
		if _, ok := item.(*ast.KeyValueExpr); ok {
			break
		}
		v1 := &ast.KeyValueExpr{
			Key: &ast.Ident{
				Name:    names[i].Name,
				NamePos: item.Pos(),
			},
			Value: item,
		}
		n.Elts[i] = v1
	}
}

func nodeCode(n ast.Node) string {
	bf := &bytes.Buffer{}
	fset := token.NewFileSet()
	_ = format.Node(bf, fset, n)
	return bf.String()
}

// fixStructBlankLine 给 struct 的字段定义添加空行
func fixStructBlankLine(fileSet *token.FileSet, f *ast.File) {
	ast.Inspect(f, func(node ast.Node) bool {
		if nv, ok := node.(*ast.StructType); ok {
			doFixStructBlankLine(fileSet, nv)
		}
		return true
	})
}

func doFixStructBlankLine(fileSet *token.FileSet, st *ast.StructType) {
	if len(st.Fields.List) == 0 {
		return
	}
	// todo
}

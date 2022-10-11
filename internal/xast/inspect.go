// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/3

package xast

import (
	"go/ast"
	"go/token"

	"github.com/fsgo/go_fmt/internal/common"
)

// dealWithInspect 删除多余的空行
// 如 struct 的 fields 定义的前后，多余的空行
// func 内部前后多余的空行
func dealWithInspect(req *common.Request) {
	ast.Inspect(req.AstFile, func(node ast.Node) bool {
		switch vt := node.(type) {
		case *ast.TypeSpec:
			switch vt.Type.(type) {
			case *ast.StructType,
				*ast.InterfaceType:
				fixStructInterfaceType(req, vt)
			}
		case *ast.BlockStmt:
			fixBlockStmt(req, vt)
		case *ast.GenDecl:
			fixGenDecl(req, vt)
		case *ast.CompositeLit:
			fixCompositeLit(req, vt)
		}
		return true
	})
}

// fixGenDecl 处理全局的和函数内部的定义
func fixGenDecl(req *common.Request, st *ast.GenDecl) {
	if common.Debug {
		common.DebugPrintln(0, "fixGenDecl",
			"start:", req.FSet.Position(st.Pos()),
			"end:", req.FSet.Position(st.End()),
		)
	}

	// 处理头部多余的空行
	{
		firstFieldPos := firstSpecPos(st.Specs, st.End())
		trimHeadEmptyLine(req, st.Pos(), firstFieldPos)
	}

	// 处理中间的部分，添加适当的空行
	{
		for j := 0; j < len(st.Specs); j++ {
			spec := st.Specs[j]
			switch vt := spec.(type) {
			case *ast.TypeSpec:
				if j > 0 && vt.Doc != nil {
					req.TokenLine().AddLine(0, vt.Doc.Pos())
				}
			case *ast.ValueSpec:
				if j > 0 && vt.Doc != nil {
					req.TokenLine().AddLine(0, vt.Doc.Pos())
				}
			}
		}
	}

	// 处理尾部多余的空行
	{
		lastFieldPos := lastSpecPos(st.Specs, st.Pos())
		trimTailEmptyLine(req, lastFieldPos, st.End())
	}
}

func firstSpecPos(fl []ast.Spec, def token.Pos) token.Pos {
	if len(fl) == 0 {
		return def
	}
	first := fl[0]
	switch fv := first.(type) {
	case *ast.ValueSpec:
		if fv.Doc != nil {
			return fv.Doc.Pos()
		}
	case *ast.TypeSpec:
		if fv.Doc != nil {
			return fv.Doc.Pos()
		}
	}
	return first.Pos()
}

func lastSpecPos(fl []ast.Spec, def token.Pos) token.Pos {
	if len(fl) == 0 {
		return def
	}
	last := fl[len(fl)-1]
	switch fv := last.(type) {
	case *ast.ValueSpec:
		if fv.Comment != nil {
			return fv.Comment.End()
		}
	case *ast.TypeSpec:
		if fv.Comment != nil {
			return fv.Comment.End()
		}
	}
	return last.End()
}

func fixStructInterfaceType(req *common.Request, ts *ast.TypeSpec) {
	tp := ts.Type
	if common.Debug {
		common.DebugPrintln(0, "fixStructType", ts.Name.Name,
			"start:", req.FSet.Position(tp.Pos()),
			"end:", req.FSet.Position(tp.End()),
		)
	}

	var fields *ast.FieldList
	switch vt := tp.(type) {
	case *ast.StructType:
		fields = vt.Fields
	case *ast.InterfaceType:
		fields = vt.Methods
	default:
		return
	}
	// 处理头部多余的空行
	{
		firstFieldPos := fieldsFirstPos(fields, tp.End())
		trimHeadEmptyLine(req, tp.Pos(), firstFieldPos)
	}

	// 给字段定义之间，若有文档，则添加空行
	{
		for i := 0; i < len(fields.List); i++ {
			fd := fields.List[i]
			if fd.Doc != nil {
				if i > 0 {
					req.TokenLine().AddLine(0, fd.Doc.Pos())
				}
				req.TokenLine().AddLine(0, nodeEndPos(req, fd))
			}
		}
	}

	// 处理尾部多余的空行
	{
		lastFieldPos := fieldsLastPos(fields, tp.Pos())
		trimTailEmptyLine(req, lastFieldPos, tp.End())
	}
}

func fieldsFirstPos(fl *ast.FieldList, def token.Pos) token.Pos {
	if fl == nil || len(fl.List) == 0 {
		return def
	}
	first := fl.List[0]
	if first.Doc != nil {
		return first.Doc.Pos()
	}
	return first.Pos()
}

func fieldsLastPos(fl *ast.FieldList, def token.Pos) token.Pos {
	if fl == nil || len(fl.List) == 0 {
		return def
	}
	last := fl.List[len(fl.List)-1]
	return endPos(last.End(), last.Comment)
}

func fixBlockStmt(req *common.Request, tts *ast.BlockStmt) {
	if common.Debug {
		common.DebugPrintln(0, "fixBlockStmt",
			"start:", req.FSet.Position(tts.Pos()),
			"end:", req.FSet.Position(tts.End()),
		)
	}
	headPos := tts.End() // 正序查找，首个元素的 pos
	tailPos := tts.Pos() // 倒序查找，最后一个元素的 pos
	if len(tts.List) > 0 {
		headPos = tts.List[0].Pos()
		tailPos = tts.List[len(tts.List)-1].End()
	}
	// 处理头部的空行
	trimHeadEmptyLine(req, tts.Pos(), headPos)

	// 处理尾部的空行
	trimTailEmptyLine(req, tailPos, tts.End())

	fixEmptyFunc(req, tts)
}

func fixCompositeLit(req *common.Request, cl *ast.CompositeLit) {
	switch tp := cl.Type.(type) {
	case *ast.MapType:
		doMapType(req, cl, tp)
	case *ast.ArrayType:
		doArrayType(req, cl, tp)
	}
}

// 处理 map 定义的空行
func doMapType(req *common.Request, cl *ast.CompositeLit, tp *ast.MapType) {
	headEltPos := cl.End()
	if len(cl.Elts) > 0 {
		headEltPos = cl.Elts[0].Pos()
	}
	// 处理头部的空行
	trimHeadEmptyLine(req, tp.End(), headEltPos)
}

// 处理 slice/array 的格式
func doArrayType(req *common.Request, cl *ast.CompositeLit, tp *ast.ArrayType) {
	headEltPos := cl.End()
	lastEltPos := cl.Lbrace
	if len(cl.Elts) > 0 {
		headEltPos = cl.Elts[0].Pos()
		lastEltPos = cl.Elts[len(cl.Elts)-1].End()
	}
	// 处理头部的空行
	trimHeadEmptyLine(req, tp.End(), headEltPos)

	// 处理尾部的空行
	// 普通原始后的空行 go/printer 会自动移除
	// 但是最后一行数注释的时候，不会自动移除
	trimTailEmptyLine(req, lastEltPos, cl.Rbrace)

	doFormatArrayPretty(req, cl, tp)
}

// 将 slice 定义格式化为多行
func doFormatArrayPretty(req *common.Request, cl *ast.CompositeLit, tp *ast.ArrayType) {
	if len(cl.Elts) < 2 {
		return
	}
	cmts := cmtsBetween(req, tp.End(), cl.End()-1)
	if len(cmts) > 0 {
		return
	}
	//
	// var _=[]int{ // l1 是这个 { 所在的行
	//  1,2,3,      // l2 是首个元素 1 所在的行
	//  4,5,6,7}    // l3 是 } 所在的行
	//
	l1 := req.FSet.Position(tp.Lbrack).Line
	l2 := req.FSet.Position(cl.Elts[0].Pos()).Line
	if l1+1 != l2 {
		return
	}
	l3 := req.FSet.Position(cl.Rbrace).Line
	if l2+1 != l3 {
		return
	}
	var l2ElCount int
	var z int
	for i := 0; i < len(cl.Elts); i++ {
		e := cl.Elts[i]
		if line := req.FSet.Position(e.Pos()).Line; line == l2 {
			l2ElCount++
		} else {
			z++
			if z == l2ElCount || i == len(cl.Elts)-1 {
				req.TokenLine().AddLine(0, e.End())
				z = 0
			}
		}
	}
}

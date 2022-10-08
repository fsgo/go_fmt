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
			if tts, ok2 := vt.Type.(*ast.StructType); ok2 {
				fixStructType(req, vt, tts)
			}
		case *ast.BlockStmt:
			fixBlockStmt(req, vt)
		case *ast.GenDecl:
			fixGenDecl(req, vt)
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

func fixStructType(req *common.Request, ts *ast.TypeSpec, st *ast.StructType) {
	if common.Debug {
		common.DebugPrintln(0, "fixStructType", ts.Name.Name,
			"start:", req.FSet.Position(st.Pos()),
			"end:", req.FSet.Position(st.End()),
		)
	}

	// 处理头部多余的空行
	{
		firstFieldPos := fieldsFirstPos(st.Fields, st.End())
		trimHeadEmptyLine(req, st.Pos(), firstFieldPos)
	}

	// 给字段定义之间，若有文档，则添加空行
	{
		for i := 0; i < len(st.Fields.List); i++ {
			fd := st.Fields.List[i]
			if fd.Doc != nil && i > 0 {
				req.TokenLine().AddLine(0, fd.Doc.Pos())
			}
		}
	}

	// 处理尾部多余的空行
	{
		lastFieldPos := fieldsLastPos(st.Fields, st.Pos())
		trimTailEmptyLine(req, lastFieldPos, st.End())
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
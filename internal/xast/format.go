// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/2

package xast

import (
	"go/ast"
	"go/token"

	"github.com/fsgo/go_fmt/internal/common"
)

// Format 其他一下高级规则
func Format(req *common.Request) {
	req.MustReParse()
	if common.Debug {
		common.DebugPrintln(0, "xast.Format",
			req.FileName, req.AstFile.Pos(), req.FSet.Base(),
		)
	}
	dealFileDecls(req)
	dealWithInspect(req)
	req.TokenLine().Execute()
}

// dealFileDecls 在适当的位置添加空行
//  1. 多个同类型定义之间，保持一个空行间隔
//     如 2 个 func 之间
//  2. 不同类型的定义之间，保持一个空行
//  3. 有文档的类型、值定义，前后各有一个空行
func dealFileDecls(req *common.Request) {
	var last ast.Decl
	total := len(req.AstFile.Decls)
	for i := 0; i < total; i++ {
		cur := req.AstFile.Decls[i]
		if gd, ok := cur.(*ast.GenDecl); ok {
			gf := &genDeclFixer{
				req:  req,
				decl: gd,
				last: last,
			}
			gf.fix()
		} else if fd, ok2 := cur.(*ast.FuncDecl); ok2 {
			fixFuncDecl(req, fd)
		}
		last = cur
	}
}

func fixFuncDecl(req *common.Request, fd *ast.FuncDecl) {
	if common.Debug {
		common.DebugPrintln(0,
			"fixFuncDecl", fd.Name.Name,
			"start:", req.FSet.Position(fd.Pos()),
			"end:", req.FSet.Position(fd.End()),
		)
	}

	// func 尾部添加空行
	req.TokenLine().AddLine(0, nodeEndPos(req, fd))

	fixEmptyFunc(req, fd.Body)
}

// 空函数定义，将{} 放在同一行
func fixEmptyFunc(req *common.Request, fd *ast.BlockStmt) {
	if fd == nil {
		// 只有函数定义，没有函数体：
		// func delete(m string)
		// see testdata/case4.go.input
		return
	}
	cmts := cmtsBetween(req, fd.Lbrace, fd.Rbrace)
	// 在函数定义区间无评论内容，同时没有表达式
	if len(cmts) == 0 && len(fd.List) == 0 {
		fd.Rbrace = fd.Lbrace + 1
	}
}

type genDeclFixer struct {
	req  *common.Request
	decl *ast.GenDecl
	last ast.Decl
}

func (gf *genDeclFixer) notSameType() bool {
	if gf.last == nil {
		return false
	}
	gd1, ok1 := gf.last.(*ast.GenDecl)
	if !ok1 {
		return true
	}

	return gd1.Tok != gf.decl.Tok
}

func (gf *genDeclFixer) fix() {
	// 前一个类型和当前类型不一样
	if gf.last != nil && gf.notSameType() {
		if common.Debug {
			common.DebugPrintln(0, "notSameType")
		}
		gf.addNewLineBefore()
	}

	// 分组定义，eg：
	// var (xxx)
	if gf.decl.Lparen.IsValid() {
		gf.addNewLineBefore()
		gf.addNewLineAfter()
	}

	if gf.decl.Tok == token.TYPE {
		for j := 0; j < len(gf.decl.Specs); j++ {
			spec := gf.decl.Specs[j]
			switch spec.(type) {
			case *ast.TypeSpec:
				// 对于全局的 type 定义分组，多个定义之间添加空行分割
				gf.req.TokenLine().AddLine(0, nodeEndPos(gf.req, spec))
			case *ast.ValueSpec:
			}
		}
	}
}

func (gf *genDeclFixer) addNewLineBefore() {
	gf.req.TokenLine().AddLine(1, nodeStartPos(gf.req, gf.decl)-1)
}

func (gf *genDeclFixer) addNewLineAfter() {
	gf.req.TokenLine().AddLine(1, nodeEndPos(gf.req, gf.decl))
}

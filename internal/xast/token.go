// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/7

package xast

import (
	"go/ast"
	"go/token"

	"github.com/fsgo/go_fmt/internal/common"
)

// nodeEndPos 查找当前节点的结尾的 Pos
func nodeEndPos(req *common.Request, node ast.Node) token.Pos {
	pos := node.End()
	switch vt := node.(type) {
	case *ast.ValueSpec:
		if vt.Comment != nil {
			pos = vt.Comment.End()
		}
	case *ast.TypeSpec:
		if vt.Comment != nil {
			pos = vt.Comment.End()
		}
	}

	line := req.FSet.Position(pos).Line
	for i := 0; i < len(req.AstFile.Comments); i++ {
		cmt := req.AstFile.Comments[i]
		p1 := req.FSet.Position(cmt.Pos())
		if cmt.Pos() >= pos && (line == p1.Line) {
			pos = cmt.End()
			line = req.FSet.Position(cmt.End()).Line
		}
	}
	return pos
}

func nodeStartPos(req *common.Request, node ast.Node) token.Pos {
	pos := node.Pos()
	switch vt := node.(type) {
	case *ast.ValueSpec:
		if vt.Doc != nil {
			pos = vt.Doc.Pos()
		}
	case *ast.TypeSpec:
		if vt.Doc != nil {
			pos = vt.Doc.Pos()
		}
	}

	line := req.FSet.Position(pos).Line
	for i := 0; i < len(req.AstFile.Comments); i++ {
		cmt := req.AstFile.Comments[i]
		p1 := req.FSet.Position(cmt.End())
		if cmt.Pos() < pos && (line == p1.Line) {
			pos = cmt.Pos()
			line = req.FSet.Position(cmt.Pos()).Line
		}
	}
	return pos
}

// trimHeadEmptyLine 处理头部的空行
func trimHeadEmptyLine(req *common.Request, start, end token.Pos) {
	lineBegin := req.FSet.Position(start).Line + 1
	lineEnd := req.FSet.Position(end).Line
	cmts := cmtsBetween(req, start, end)
	for j := 0; j < len(cmts); j++ {
		cmt := cmts[j]
		lineNo := req.FSet.Position(cmt.Pos()).Line
		// 紧接着下一行就是评论内容
		// eg: type user struct{
		//   // 这是评论
		// }
		if lineBegin == lineNo {
			return
		}
		if lineNo > lineBegin {
			lineEnd = lineNo
			break
		}
	}
	for ; lineBegin < lineEnd; lineBegin++ {
		req.TokenLine().DeleteLine(1, lineBegin)
	}
}

// trimTailEmptyLine 处理尾部的空行
func trimTailEmptyLine(req *common.Request, start, end token.Pos) {
	lineBegin := req.FSet.Position(start).Line + 1
	lineEnd := req.FSet.Position(end).Line
	cmts := cmtsBetween(req, start, end)
	for j := len(cmts) - 1; j >= 0; j-- {
		cmt := cmts[j]
		no := req.FSet.Position(cmt.End()).Line
		if no <= lineEnd {
			lineBegin = no + 1
			break
		}
	}
	for ; lineBegin < lineEnd; lineBegin++ {
		req.TokenLine().DeleteLine(1, lineBegin)
	}
}

func cmtsBetween(req *common.Request, start, end token.Pos) []*ast.CommentGroup {
	if len(req.AstFile.Comments) == 0 {
		return nil
	}
	var list []*ast.CommentGroup
	for i := 0; i < len(req.AstFile.Comments); i++ {
		cmt := req.AstFile.Comments[i]
		if cmt.Pos() > start && cmt.End() < end {
			list = append(list, cmt)
		}
	}
	return list
}

func endPos(end token.Pos, cmt *ast.CommentGroup) token.Pos {
	if cmt == nil {
		return end
	}
	return cmt.End()
}

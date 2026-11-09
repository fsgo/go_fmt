package gofmt

import (
	"go/ast"
	"strconv"
	"strings"

	"github.com/fsgo/go_fmt/internal/common"
)

func formatStructTags(req *common.Request) {
	ast.Inspect(req.AstFile, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range st.Fields.List {
			if field.Tag == nil {
				continue
			}

			raw, err := strconv.Unquote(field.Tag.Value)
			if err != nil {
				continue
			}
			// "db:\"id\""   ---> `db:"id"`
			if !strings.Contains(raw, "`") {
				field.Tag.Value = "`" + raw + "`"
			}
		}

		return true
	})
}

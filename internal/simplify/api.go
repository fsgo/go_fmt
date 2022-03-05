// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu (duv123@baidu.com)
// Date: 2022/3/5

package simplify

import (
	"go/ast"
)

// Format call simplify
func Format(f *ast.File) {
	simplify(f)
}

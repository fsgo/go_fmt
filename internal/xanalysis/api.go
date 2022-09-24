// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/18

package xanalysis

import (
	"github.com/fsgo/go_fmt/internal/common"
)

// Format 基于类型检查、语法检查等优化
func Format(req *common.Request) {
	fixStruct(req)
}

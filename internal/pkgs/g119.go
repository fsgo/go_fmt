// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/8/19

package pkgs

import (
	"sync/atomic"
)

// require go1.19
//
// https://go.dev/doc/go1.19#go-doc
var _ = atomic.Uint32{}

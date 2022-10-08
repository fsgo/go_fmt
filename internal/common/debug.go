// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/6

package common

import (
	"fmt"
	"log"
	"os"
)

// Debug 程序内部调试
var Debug = os.Getenv("go_fmt_debug") == "1"

var debugLogger = log.New(os.Stderr, "[Debug] ", log.Lshortfile)

// DebugPrintln 打印调试日志
func DebugPrintln(depth int, v ...any) {
	if !Debug {
		return
	}
	_ = debugLogger.Output(2+depth, fmt.Sprintln(v...))
}

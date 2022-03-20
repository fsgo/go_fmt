// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/3/19

package common

import (
	"fmt"
)

const consoleColorTag = 0x1B

// ConsoleRed 控制台红色字符
func ConsoleRed(txt string) string {
	return fmt.Sprintf("%c[31m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// ConsoleGreen 控制台绿色字符
func ConsoleGreen(txt string) string {
	return fmt.Sprintf("%c[32m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

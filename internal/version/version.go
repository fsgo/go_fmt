// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/7/22

package version

const versionID = "v0.5.3"

const versionDate = "2023-08-20"

// Version 版本信息
func Version() string {
	return versionID + " " + versionDate
}

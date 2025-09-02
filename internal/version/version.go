// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/7/22

package version

const versionID = "v0.7.0"

const versionDate = "2025-09-02"

// Version 版本信息
func Version() string {
	return versionID + " " + versionDate
}

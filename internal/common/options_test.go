// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/5/11

package common

import (
	"testing"

	"github.com/fsgo/fst"
)

func TestOptions_GetImportGroup(t *testing.T) {
	o := &Options{}
	fst.Equal(t, 0, o.GetImportGroup(ImportGroupGoStandard))
	fst.Equal(t, 1, o.GetImportGroup(ImportGroupThirdParty))
	fst.Equal(t, 2, o.GetImportGroup(ImportGroupCurrentModule))
}

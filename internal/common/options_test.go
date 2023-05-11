// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/5/11

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions_GetImportGroup(t *testing.T) {
	o := &Options{}
	require.Equal(t, 0, o.GetImportGroup(ImportGroupGoStandard))
	require.Equal(t, 1, o.GetImportGroup(ImportGroupThirdParty))
	require.Equal(t, 2, o.GetImportGroup(ImportGroupCurrentModule))
}

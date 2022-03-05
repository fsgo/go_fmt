// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/3/5

//go:build go1.18
// +build go1.18

package simplify

import (
	"reflect"
)

func init() {
	substKind = func(m map[string]reflect.Value, p reflect.Value, pos reflect.Value) *reflect.Value {
		if p.Kind() != reflect.Pointer {
			return nil
		}
		v := reflect.New(p.Type()).Elem()
		if elem := p.Elem(); elem.IsValid() {
			v.Set(subst(m, elem, pos).Addr())
		}
		return &v
	}
}

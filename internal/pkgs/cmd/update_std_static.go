/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/1/4
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"

	"github.com/fsgo/go_fmt/internal/pkgs/std"
)

var out = flag.String("out", "", "输出文件")

func main() {
	flag.Parse()
	pkgs, err := std.PKGs()
	if err != nil {
		log.Fatalln(err)
	}
	if *out == "" {
		fmt.Println("pkgs:", pkgs)
		return
	}

	absPath, errAbs := filepath.Abs(*out)
	if errAbs != nil {
		log.Fatalf("filepath.Abs(%q) with error:%s\n", *out, errAbs)
	}

	var buf bytes.Buffer
	buf.WriteString("// AUTO GENERATE, DOT EDIT IT\n\n")

	buf.WriteString("// GO VersionFile: ")
	buf.WriteString(runtime.Version())
	buf.WriteString("\n\n")

	buf.WriteString("package ")

	pkgName := filepath.Base(filepath.Dir(absPath))

	buf.WriteString(pkgName)
	buf.WriteString("\n\n")
	buf.WriteString("var stdPKGs=[]string{\n")
	for _, pkg := range pkgs {
		buf.WriteString(fmt.Sprintf("\t%q", pkg))
		buf.WriteString(",\n")
	}
	buf.WriteString("}\n")

	code, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("format.Source with error:%s\n", err)
	}

	if err := ioutil.WriteFile(*out, code, 0644); err != nil {
		log.Fatalf("write %s with error:%v\n", *out, err)
	}
}

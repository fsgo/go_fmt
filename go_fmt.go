/*
 * Copyright(C) 2019 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2019/12/16
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/scanner"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsgo/go_fmt/internal/gofmt"
)

var version = "v0.1 20191217"

var (
	LocalPrefix string
	write       bool

	multipleFiles bool
)

func init() {
	flag.BoolVar(&write, "w", false, "write result to (source) file instead of stdout")
	flag.StringVar(&LocalPrefix, "local", "auto", "put imports beginning with this string after 3rd-party packages; comma-separated list")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go_fmt [flags] [path ...]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nsite :    github.com/fsgo/go_fmt\n")
	fmt.Fprintf(os.Stderr, "version:  %s\n", version)
	os.Exit(2)
}

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	paths := flag.Args()

	if len(paths) == 0 {
		paths = []string{"git_change"}
	}

	if len(paths) == 1 {
		multipleFiles = false
	}

	for _, filePath := range paths {
		formatOneFile(filePath)
	}
}

func formatOneFile(fileName string) {
	var goFiles []string

	var checkGoFile = func(fileName string) {
		info, err := os.Stat(fileName)
		if err != nil {
			report(err)
		}

		if isGoFile(info) {
			goFiles = append(goFiles, fileName)
		}
	}

	if fileName == "./..." {
		multipleFiles = true

		filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
			if err == nil && isGoFile(info) {
				goFiles = append(goFiles, path)
			}
			if err != nil {
				report(err)
			}
			return nil
		})
	} else if fileName == "git_change" {
		multipleFiles = true

		files, err := gofmt.GitChangeFiles()
		if err != nil {
			report(err)
		}
		for _, filePath := range files {
			checkGoFile(filePath)
		}
	} else {
		checkGoFile(fileName)
	}

	for _, goFile := range goFiles {
		if err := formatFileByName(goFile); err != nil {
			report(err)
		}
	}
}

var consoleColorTag = 0x1B

func formatFileByName(fileName string) error {
	src, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	localPrefix, err := gofmt.DetectLocal(LocalPrefix, fileName)
	if err != nil {
		return err
	}
	out, err := gofmt.Format(fileName, src, localPrefix)
	if err != nil {
		return err
	}

	if write {
		if bytes.Equal(src, out) {
			fmt.Fprintf(os.Stderr, "%c[1;40;34m unchange : %s%c[0m\n", consoleColorTag, fileName, consoleColorTag)
		} else {
			fmt.Fprintf(os.Stdout, "%c[1;40;32m rewrited : %s%c[32m\n", consoleColorTag, fileName, consoleColorTag)
		}
		return ioutil.WriteFile(fileName, out, 0)
	} else {
		fmt.Println(string(out))
	}
	return nil
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/12

package xpasser

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsgo/go_fmt/internal/common"
)

// TryGoModTidy 更新 go.sum 文件
// 若不执行 go mod tidy 可能由于 go.sum 文件未更新，导致 go list 命令失败
// 进而导致 pkg 不能正常的 load
func TryGoModTidy(opt common.Options, fs []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	if len(fs) > 0 && len(fs[0]) > 0 {
		cmd.Dir = filepath.Dir(fs[0])
	}
	if opt.Trace {
		log.Println("exec_start:", cmd.String())
	}
	out, err := cmd.Output()
	if err != nil {
		log.Println("exec:", cmd.String(), ", failed:\n", stderr.String())
	}
	if opt.Trace {
		log.Println("exec_done:", cmd.String(), "out:", string(out), ", err:", err, "cost:", time.Since(start).String())
	}
}

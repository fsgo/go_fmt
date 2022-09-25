// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/3/19

package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// DiffResult Diff 的结果
type DiffResult interface {
	Detail() any
	String() string
}

// Diff 比较文本内容的不同
func Diff(a, b string, trace bool) DiffResult {
	linesA := strings.Split(a, "\n")
	linesB := strings.Split(b, "\n")
	r := &diffReporter{
		trace: trace,
	}
	same := cmp.Equal(linesA, linesB, cmp.Reporter(r))
	if same {
		return nil
	}
	return r
}

// DiffType diff 的类型
type DiffType string

const (
	// DiffTypeChange 格式有变化
	DiffTypeChange DiffType = "change"

	// DiffTypeAdd 是新增内容
	DiffTypeAdd = "add"

	// DoffTypeDelete 内容被删除
	DoffTypeDelete = "delete"
)

type diffDetail struct {
	Trace  string
	Delete string
	Add    string
	Type   DiffType
	LineNo int
}

func (dd *diffDetail) String() string {
	return dd.Output(true)
}

func (dd *diffDetail) Output(trace bool) string {
	var b strings.Builder
	if trace {
		b.WriteString("Trace: " + dd.Trace + "\n")
	}
	b.WriteString(fmt.Sprintf("line %d:\n", dd.LineNo))

	if len(dd.Delete) > 0 || dd.Type == DoffTypeDelete {
		b.WriteString(fmt.Sprintf("  -: %s\n", dd.quote(dd.Delete)))
	}
	if len(dd.Add) > 0 || dd.Type == DiffTypeAdd {
		b.WriteString(fmt.Sprintf("  +: %s\n", dd.quote(dd.Add)))
	}
	return strings.TrimSpace(b.String())
}

func (dd *diffDetail) quote(txt string) string {
	// 针对全部是不可见字符的情况做展现优化
	if len(strings.TrimSpace(txt)) == 0 {
		arr := []byte(txt)
		for i := 0; i < len(arr); i++ {
			switch arr[i] {
			case ' ':
				arr[i] = '_'
			case '\t':
				arr[i] = '-'
			}
		}
		return string(arr) + "\\n"
	}
	return txt
}

type diffReporter struct {
	path  cmp.Path
	diffs []*diffDetail
	trace bool
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

//	解析文本获取行数
//
// [31->?] [26] [36->32] [?->39]
func (r *diffReporter) parserLineNo(txt string) int {
	for _, b := range []byte("[]?>") {
		txt = strings.ReplaceAll(txt, string(b), "")
	}
	if strings.HasPrefix(txt, "-") || strings.HasSuffix(txt, "-") {
		txt = strings.Trim(txt, "-")
	} else if idx := strings.Index(txt, "-"); idx > 0 {
		// [36->32]
		txt = txt[0:idx]
	}
	num, _ := strconv.Atoi(txt)
	return num
}

func (r *diffReporter) Report(rs cmp.Result) {
	if rs.Equal() {
		return
	}
	vx, vy := r.path.Last().Values()
	lineNo := r.parserLineNo(r.path.Last().String())

	detail := &diffDetail{
		LineNo: lineNo + 1,
		Type:   DiffTypeChange,
		Trace:  r.path.Last().String(),
	}

	if vx.IsValid() {
		detail.Delete = r.formatTxt(vx)
	} else {
		detail.Type = DiffTypeAdd
	}

	if vy.IsValid() {
		detail.Add = r.formatTxt(vy)
	} else {
		detail.Type = DoffTypeDelete
	}

	r.diffs = append(r.diffs, detail)
}

func (r *diffReporter) formatTxt(rv reflect.Value) string {
	txt := fmt.Sprintf("%+v", rv)
	return txt
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) String() string {
	lines := make([]string, len(r.diffs))
	for i := 0; i < len(r.diffs); i++ {
		lines[i] = r.diffs[i].Output(r.trace)
	}
	return strings.Join(lines, "\n\n")
}

func (r *diffReporter) Detail() any {
	return r.diffs
}

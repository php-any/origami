package php

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SscanfFunction 实现 PHP sscanf。
// 两参数时返回解析结果数组；参数不足时返回 null。
// 额外 by-ref 变量时返回赋值个数（简化实现）。
type SscanfFunction struct{}

func NewSscanfFunction() data.FuncStmt { return &SscanfFunction{} }

func (f *SscanfFunction) GetName() string { return "sscanf" }

// 匹配 %d / %02x / %f / %[1-5] 等转换说明
var sscanfSpecRe = regexp.MustCompile(`%(\d*)(?:\.(\d+))?([diufFeEgGcsxXon]|\[\^?[^\]]*\])`)

func (f *SscanfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, ok := ctx.GetIndexValue(0)
	if !ok || strVal == nil {
		return data.NewNullValue(), nil
	}
	fmtVal, ok := ctx.GetIndexValue(1)
	if !ok || fmtVal == nil {
		return data.NewNullValue(), nil
	}
	input := strVal.AsString()
	format := fmtVal.AsString()

	matches := sscanfSpecRe.FindAllStringSubmatch(format, -1)
	if len(matches) == 0 {
		return data.NewArrayValue(nil), nil
	}

	dests := make([]any, len(matches))
	ptrs := make([]any, len(matches))
	for i, m := range matches {
		verb := m[3]
		if strings.HasPrefix(verb, "[") {
			s := ""
			dests[i] = &s
			ptrs[i] = &s
			continue
		}
		switch verb[0] {
		case 'd', 'i', 'u', 'o', 'x', 'X':
			n := 0
			dests[i] = &n
			ptrs[i] = &n
		case 'f', 'F', 'e', 'E', 'g', 'G':
			fl := 0.0
			dests[i] = &fl
			ptrs[i] = &fl
		case 'c':
			b := byte(0)
			dests[i] = &b
			ptrs[i] = &b
		default:
			s := ""
			dests[i] = &s
			ptrs[i] = &s
		}
	}

	goFormat := phpSscanfFormatToGo(format)
	n, _ := fmt.Sscanf(input, goFormat, ptrs...)

	// 仅实现两参数返回数组（Filament Color / 常用路径）。
	// by-ref 赋值后续再补；多传参时仍先返回数组，避免误入计数分支。
	if n < len(matches) {
		return data.NewNullValue(), nil
	}
	out := make([]data.Value, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, sscanfDestToValue(dests[i]))
	}
	return data.NewArrayValue(out), nil
}

func sscanfDestToValue(dest any) data.Value {
	switch p := dest.(type) {
	case *int:
		return data.NewIntValue(*p)
	case *float64:
		return data.NewFloatValue(*p)
	case *string:
		return data.NewStringValue(*p)
	case *byte:
		return data.NewStringValue(string([]byte{*p}))
	default:
		return data.NewNullValue()
	}
}

func phpSscanfFormatToGo(format string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range format {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	return b.String()
}

var sscanfFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
	node.NewParameter(nil, "format", 1, nil, nil),
}

func (f *SscanfFunction) GetParams() []data.GetValue {
	return sscanfFunctionGetParams
}

var sscanfFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "format", 1, data.NewBaseType("string")),
}

func (f *SscanfFunction) GetVariables() []data.Variable {
	return sscanfFunctionGetVariables
}

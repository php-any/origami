package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrIreplaceFunction 实现 str_ireplace 函数
// str_ireplace(mixed $search, mixed $replace, mixed $subject, int &$count = null): string|array
type StrIreplaceFunction struct{}

func NewStrIreplaceFunction() data.FuncStmt {
	return &StrIreplaceFunction{}
}

func (f *StrIreplaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	searchValue, _ := ctx.GetIndexValue(0)
	replaceValue, _ := ctx.GetIndexValue(1)
	subjectValue, _ := ctx.GetIndexValue(2)

	if searchValue == nil || replaceValue == nil || subjectValue == nil {
		return subjectValue, nil
	}

	search := searchValue.AsString()
	replace := replaceValue.AsString()
	subject := subjectValue.AsString()

	result := strings.ReplaceAll(strings.ToLower(subject), strings.ToLower(search), replace)
	// 由于使用了 ToLower 进行大小写不敏感替换，需要更精确的处理
	// 简单实现：使用 case-insensitive replace
	result = caseInsensitiveReplace(subject, search, replace)

	return data.NewStringValue(result), nil
}

func caseInsensitiveReplace(s, old, new string) string {
	if old == "" {
		return s
	}
	lower := strings.ToLower(s)
	oldLower := strings.ToLower(old)
	var result strings.Builder
	for {
		idx := strings.Index(lower, oldLower)
		if idx == -1 {
			result.WriteString(s)
			break
		}
		result.WriteString(s[:idx])
		result.WriteString(new)
		s = s[idx+len(old):]
		lower = lower[idx+len(old):]
	}
	return result.String()
}

func (f *StrIreplaceFunction) GetName() string {
	return "str_ireplace"
}

func (f *StrIreplaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "search", 0, nil, nil),
		node.NewParameter(nil, "replace", 1, nil, nil),
		node.NewParameter(nil, "subject", 2, nil, nil),
		node.NewParameterReference(nil, "count", 3, data.NewBaseType("int")),
	}
}

func (f *StrIreplaceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "search", 0, nil),
		node.NewVariable(nil, "replace", 1, nil),
		node.NewVariable(nil, "subject", 2, nil),
		node.NewVariable(nil, "count", 3, data.NewBaseType("int")),
	}
}

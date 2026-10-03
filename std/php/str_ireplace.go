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

	result, count := caseInsensitiveReplace(subject, search, replace)
	if slot := ctx.GetIndexZVal(3); slot != nil {
		slot.StoreRaw(data.NewIntValue(count))
	}

	return data.NewStringValue(result), nil
}

func caseInsensitiveReplace(s, old, new string) (string, int) {
	if old == "" {
		return s, 0
	}
	lower := strings.ToLower(s)
	oldLower := strings.ToLower(old)
	var result strings.Builder
	count := 0
	for {
		idx := strings.Index(lower, oldLower)
		if idx == -1 {
			result.WriteString(s)
			break
		}
		result.WriteString(s[:idx])
		result.WriteString(new)
		count++
		s = s[idx+len(old):]
		lower = lower[idx+len(old):]
	}
	return result.String(), count
}

func (f *StrIreplaceFunction) GetName() string {
	return "str_ireplace"
}

var strIreplaceFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "search", 0, nil, nil),
	node.NewParameter(nil, "replace", 1, nil, nil),
	node.NewParameter(nil, "subject", 2, nil, nil),
	node.NewOutputParameterReference(nil, "count", 3, node.NewNullLiteral(nil), data.NewBaseType("int")),
}

func (f *StrIreplaceFunction) GetParams() []data.GetValue {
	return strIreplaceFunctionGetParams
}

var strIreplaceFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "search", 0, nil),
	node.NewVariable(nil, "replace", 1, nil),
	node.NewVariable(nil, "subject", 2, nil),
	node.NewVariable(nil, "count", 3, data.NewBaseType("int")),
}

func (f *StrIreplaceFunction) GetVariables() []data.Variable {
	return strIreplaceFunctionGetVariables
}

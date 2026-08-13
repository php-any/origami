package array

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	caseLower = 0
	caseUpper = 1
)

// ArrayChangeKeyCaseFunction 实现 array_change_key_case
// array_change_key_case(array $array, int $case = CASE_LOWER): array
type ArrayChangeKeyCaseFunction struct{}

func NewArrayChangeKeyCaseFunction() data.FuncStmt {
	return &ArrayChangeKeyCaseFunction{}
}

func (f *ArrayChangeKeyCaseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	caseVal, _ := ctx.GetIndexValue(1)

	mode := caseLower
	if iv, ok := caseVal.(data.AsInt); ok {
		if v, err := iv.AsInt(); err == nil && v == caseUpper {
			mode = caseUpper
		}
	}

	entries := toKVEntries(arrayVal)
	if len(entries) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	transform := strings.ToLower
	if mode == caseUpper {
		transform = strings.ToUpper
	}

	result := data.NewObjectValue()
	for _, e := range entries {
		result.SetProperty(transform(e.keyStr), e.value)
	}
	return result, nil
}

func (f *ArrayChangeKeyCaseFunction) GetName() string { return "array_change_key_case" }

func (f *ArrayChangeKeyCaseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "case", 1, data.NewIntValue(caseLower), nil),
	}
}

func (f *ArrayChangeKeyCaseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "case", 1, data.NewBaseType("int")),
	}
}

package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayIsListFunction 实现 array_is_list 函数 (PHP 8.1+)
// array_is_list(array $array): bool
type ArrayIsListFunction struct{}

func NewArrayIsListFunction() data.FuncStmt {
	return &ArrayIsListFunction{}
}

func (f *ArrayIsListFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrValue, _ := ctx.GetIndexValue(0)
	if arrValue == nil {
		return data.NewBoolValue(false), nil
	}

	arr, ok := arrValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	// 空数组是 list
	if len(arr.List) == 0 {
		return data.NewBoolValue(true), nil
	}

	// 检查是否有字符串键（非空 Name 表示字符串键）
	for _, z := range arr.List {
		if z.Name != "" {
			return data.NewBoolValue(false), nil
		}
	}

	return data.NewBoolValue(true), nil
}

func (f *ArrayIsListFunction) GetName() string {
	return "array_is_list"
}

func (f *ArrayIsListFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
	}
}

func (f *ArrayIsListFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	}
}

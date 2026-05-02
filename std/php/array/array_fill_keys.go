package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayFillKeysFunction 实现 array_fill_keys 函数
// array_fill_keys(array $keys, mixed $value): array
type ArrayFillKeysFunction struct{}

func NewArrayFillKeysFunction() data.FuncStmt {
	return &ArrayFillKeysFunction{}
}

func (f *ArrayFillKeysFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	keysValue, _ := ctx.GetIndexValue(0)
	valueValue, _ := ctx.GetIndexValue(1)

	if keysValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	arr, ok := keysValue.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	var val data.Value
	if valueValue != nil {
		val = valueValue
	} else {
		val = data.NewNullValue()
	}

	seen := make(map[string]bool)
	list := make([]*data.ZVal, 0)
	for _, z := range arr.List {
		k := z.Value.AsString()
		if !seen[k] {
			list = append(list, &data.ZVal{Name: k, Value: val})
			seen[k] = true
		}
	}

	return &data.ArrayValue{List: list}, nil
}

func (f *ArrayFillKeysFunction) GetName() string {
	return "array_fill_keys"
}

func (f *ArrayFillKeysFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "keys", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}

func (f *ArrayFillKeysFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "keys", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "value", 1, nil),
	}
}

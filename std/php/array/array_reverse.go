package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayReverseFunction 实现 array_reverse
func NewArrayReverseFunction() data.FuncStmt {
	return &ArrayReverseFunction{}
}

type ArrayReverseFunction struct{}

func (f *ArrayReverseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	preserveKeysVal, _ := ctx.GetIndexValue(1)
	if arrayValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	preserve := false
	if boolean, ok := preserveKeysVal.(data.AsBool); ok {
		preserve, _ = boolean.AsBool()
	}
	result := data.NewArrayValue(nil).(*data.ArrayValue)
	appendEntry := func(key, value data.Value) {
		if _, integer := key.(*data.IntValue); integer && !preserve {
			result.AppendValue(value)
		} else {
			result.SetKey(key, value)
		}
	}
	switch array := arrayValue.(type) {
	case *data.ArrayValue:
		for position := array.Len() - 1; position >= 0; position-- {
			slot := array.At(position)
			appendEntry(slot.PHPArrayKey(position), slot.Value)
		}
	case *data.ObjectValue:
		entries := toKVEntries(array)
		for position := len(entries) - 1; position >= 0; position-- {
			key := entries[position].key
			if integer, ok := data.ParseIntArrayKeyName(key.AsString()); ok {
				key = data.NewIntValue(integer)
			}
			appendEntry(key, entries[position].value)
		}
	}
	return result, nil
}

func (f *ArrayReverseFunction) GetName() string {
	return "array_reverse"
}

var arrayReverseFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, nil),
	node.NewParameter(nil, "preserve_keys", 1, nil, nil),
}

func (f *ArrayReverseFunction) GetParams() []data.GetValue {
	return arrayReverseFunctionGetParams
}

var arrayReverseFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	node.NewVariable(nil, "preserve_keys", 1, data.NewBaseType("bool")),
}

func (f *ArrayReverseFunction) GetVariables() []data.Variable {
	return arrayReverseFunctionGetVariables
}

package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RsortFunction 实现 rsort 函数
// rsort(array &$array, int $flags = SORT_REGULAR): bool
// 对数组进行降序排序，会重新索引数组的键
type RsortFunction struct{}

func NewRsortFunction() data.FuncStmt {
	return &RsortFunction{}
}

func (f *RsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)
	flagsValue, _ := ctx.GetIndexValue(1)

	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if arrayRef.Len() == 0 {
		return data.NewBoolValue(true), nil
	}

	flags := 0
	if flagsValue != nil {
		if intVal, ok := flagsValue.(*data.IntValue); ok {
			flags, _ = intVal.AsInt()
		}
	}

	arrayRef.EditReindexing(func(slots []*data.ZVal) {
		sort.SliceStable(slots, func(i, j int) bool {
			// Reverse: compare j < i instead of i < j
			return compareValues(slots[j].Value, slots[i].Value, flags)
		})
	})

	return data.NewBoolValue(true), nil
}

func (f *RsortFunction) GetName() string {
	return "rsort"
}

var rsortFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
}

func (f *RsortFunction) GetParams() []data.GetValue {
	return rsortFunctionGetParams
}

var rsortFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
	node.NewVariable(nil, "flags", 1, data.Int{}),
}

func (f *RsortFunction) GetVariables() []data.Variable {
	return rsortFunctionGetVariables
}

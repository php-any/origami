package array

import (

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// KrsortFunction 实现 krsort 函数
// krsort(array &$array, int $flags = SORT_REGULAR): bool
// 按键名降序排序，保持键到值的关联。
type KrsortFunction struct{}

func NewKrsortFunction() data.FuncStmt {
	return &KrsortFunction{}
}

func (f *KrsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)
	flagsValue, _ := ctx.GetIndexValue(1)

	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		flags := 0
		if flag, ok := flagsValue.(data.AsInt); ok {
			flags, _ = flag.AsInt()
		}
		sortArrayKeys(v, flags, true)
		return data.NewBoolValue(true), nil

	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *KrsortFunction) GetName() string {
	return "krsort"
}

var krsortFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
}

func (f *KrsortFunction) GetParams() []data.GetValue {
	return krsortFunctionGetParams
}

var krsortFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
	node.NewVariable(nil, "flags", 1, data.Int{}),
}

func (f *KrsortFunction) GetVariables() []data.Variable {
	return krsortFunctionGetVariables
}

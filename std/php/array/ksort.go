package array

import (

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// KsortFunction 实现 ksort 函数
// ksort(array &$array, int $flags = SORT_REGULAR): bool
// 按键名升序排序，保持键到值的关联，不重新索引键。
type KsortFunction struct{}

func NewKsortFunction() data.FuncStmt {
	return &KsortFunction{}
}

func (f *KsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)
	flagsValue, _ := ctx.GetIndexValue(1) // 可选 flags

	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 支持两种内部表示：
	// - ArrayValue: 索引数组（int 键）
	// - ObjectValue: 关联数组（string 键）
	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		flags := 0
		if flag, ok := flagsValue.(data.AsInt); ok {
			flags, _ = flag.AsInt()
		}
		sortArrayKeys(v, flags, false)
		return data.NewBoolValue(true), nil

	default:
		return data.NewBoolValue(false), nil
	}
}

func (f *KsortFunction) GetName() string {
	return "ksort"
}

var ksortFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
}

func (f *KsortFunction) GetParams() []data.GetValue {
	return ksortFunctionGetParams
}

var ksortFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
	node.NewVariable(nil, "flags", 1, data.Int{}),
}

func (f *KsortFunction) GetVariables() []data.Variable {
	return ksortFunctionGetVariables
}

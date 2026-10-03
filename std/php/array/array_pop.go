package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayPopFunction 实现 array_pop。
// PHP：按引用弹出最后一个元素。关联数组在 Origami 里可能是 ObjectValue，
// 必须与 array_shift 一样就地删除，否则 Blade renderComponent 拿到 null、Livewire 报 missing root tag。

func NewArrayPopFunction() data.FuncStmt {
	return &ArrayPopFunction{}
}

type ArrayPopFunction struct{}

func (f *ArrayPopFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)

	if arr, ok := arrayValue.(*data.ArrayValue); ok {
		slot := arr.PopSlot()
		if slot == nil || slot.ReadValue() == nil {
			return data.NewNullValue(), nil
		}
		return slot.ReadValue(), nil
	}

	return data.NewNullValue(), nil
}

func (f *ArrayPopFunction) GetName() string {
	return "array_pop"
}

var arrayPopFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
}

func (f *ArrayPopFunction) GetParams() []data.GetValue {
	return arrayPopFunctionGetParams
}

var arrayPopFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
}

func (f *ArrayPopFunction) GetVariables() []data.Variable {
	return arrayPopFunctionGetVariables
}

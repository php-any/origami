package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ResetFunction 实现 reset 函数
// 将数组的内部指针移动到第一个元素，并返回该元素的值
type ResetFunction struct{}

func NewResetFunction() data.FuncStmt {
	return &ResetFunction{}
}

func (f *ResetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		return val.ResetPointer(), nil

	case *data.ClassValue:
		if ctl := data.EmitPHPError(ctx, 8192, "reset(): Calling reset() on an object is deprecated", nil); ctl != nil {
			return nil, ctl
		}
		return val.ObjectValue.ResetPointer(), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *ResetFunction) GetName() string {
	return "reset"
}

var resetFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
}

func (f *ResetFunction) GetParams() []data.GetValue {
	return resetFunctionGetParams
}

var resetFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
}

func (f *ResetFunction) GetVariables() []data.Variable {
	return resetFunctionGetVariables
}

package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PrevFunction 实现 prev 函数
// 将数组的内部指针向后移动一位，并返回该元素的值
type PrevFunction struct{}

func NewPrevFunction() data.FuncStmt {
	return &PrevFunction{}
}

func (f *PrevFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		return val.MovePointer(-1), nil

	case *data.ClassValue:
		if ctl := data.EmitPHPError(ctx, 8192, "prev(): Calling prev() on an object is deprecated", nil); ctl != nil {
			return nil, ctl
		}
		return val.ObjectValue.MovePointer(-1), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *PrevFunction) GetName() string {
	return "prev"
}

var prevFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
}

func (f *PrevFunction) GetParams() []data.GetValue {
	return prevFunctionGetParams
}

var prevFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
}

func (f *PrevFunction) GetVariables() []data.Variable {
	return prevFunctionGetVariables
}

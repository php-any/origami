package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// EndFunction 实现 end 函数
// 将数组的内部指针移动到最后一个元素，并返回该元素的值
type EndFunction struct{}

func NewEndFunction() data.FuncStmt {
	return &EndFunction{}
}

func (f *EndFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		return val.EndPointer(), nil

	case *data.ClassValue:
		if ctl := data.EmitPHPError(ctx, 8192, "end(): Calling end() on an object is deprecated", nil); ctl != nil {
			return nil, ctl
		}
		return val.ObjectValue.EndPointer(), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *EndFunction) GetName() string {
	return "end"
}

var endFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
}

func (f *EndFunction) GetParams() []data.GetValue {
	return endFunctionGetParams
}

var endFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
}

func (f *EndFunction) GetVariables() []data.Variable {
	return endFunctionGetVariables
}

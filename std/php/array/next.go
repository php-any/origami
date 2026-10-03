package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NextFunction 实现 next 函数
// 将数组的内部指针向前移动一位，并返回该元素的值
type NextFunction struct{}

func NewNextFunction() data.FuncStmt {
	return &NextFunction{}
}

func (f *NextFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		return val.MovePointer(1), nil

	case *data.ClassValue:
		if ctl := data.EmitPHPError(ctx, 8192, "next(): Calling next() on an object is deprecated", nil); ctl != nil {
			return nil, ctl
		}
		return val.ObjectValue.MovePointer(1), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *NextFunction) GetName() string {
	return "next"
}

var nextFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
}

func (f *NextFunction) GetParams() []data.GetValue {
	return nextFunctionGetParams
}

var nextFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
}

func (f *NextFunction) GetVariables() []data.Variable {
	return nextFunctionGetVariables
}

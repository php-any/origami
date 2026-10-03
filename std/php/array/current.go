package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CurrentFunction 实现 current 函数
// 返回数组中的当前元素
type CurrentFunction struct{}

func NewCurrentFunction() data.FuncStmt {
	return &CurrentFunction{}
}

func (f *CurrentFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		return val.PointerValue(), nil

	case *data.ClassValue:
		if ctl := data.EmitPHPError(ctx, 8192, "current(): Calling current() on an object is deprecated", nil); ctl != nil {
			return nil, ctl
		}
		return val.ObjectValue.PointerValue(), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *CurrentFunction) GetName() string {
	return "current"
}

var currentFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.Mixed{}),
}

func (f *CurrentFunction) GetParams() []data.GetValue {
	// PHP 8.0+：current() 按值接收，允许 current(array_slice(...)) 等临时值
	return currentFunctionGetParams
}

var currentFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
}

func (f *CurrentFunction) GetVariables() []data.Variable {
	return currentFunctionGetVariables
}

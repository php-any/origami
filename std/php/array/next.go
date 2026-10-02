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
	case *data.ObjectValue:
		return val.MovePointer(1), nil
	case *data.ClassValue:
		// 处理 Iterator 对象
		// 检查是否实现了 Iterator 接口
		if targetInterface, ok := ctx.GetVM().GetInterface("Iterator"); ok {
			if checkInterfaceStructure(val.Class, targetInterface) {
				// 移动到下一个位置
				if ctl := callVoidMethod(val, "next"); ctl != nil {
					return data.NewNullValue(), nil
				}

				// 检查当前位置是否有效
				valid, ctl := callBoolMethod(val, "valid")
				if ctl != nil {
					return data.NewNullValue(), nil
				}
				if !valid {
					return data.NewNullValue(), nil
				}

				// 获取当前元素
				currentVal, ctl := callValueMethod(val, "current")
				if ctl != nil {
					return data.NewNullValue(), nil
				}
				return currentVal, nil
			}
		}
		// 不是 Iterator 接口，返回 null
		return data.NewNullValue(), nil

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

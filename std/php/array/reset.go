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
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		// 处理数组
		if len(val.List) == 0 {
			return data.NewNullValue(), nil
		}
		// 返回第一个元素
		return val.List[0].Value, nil

	case *data.ObjectValue:
		// 处理对象（关联数组）
		var firstValue data.Value
		var hasValue bool

		// 使用 RangeProperties 按插入顺序遍历，获取第一个元素
		val.RangeProperties(func(key string, value data.Value) bool {
			if !hasValue {
				firstValue = value
				hasValue = true
			}
			return false // 只获取第一个元素
		})

		if !hasValue {
			return data.NewNullValue(), nil
		}
		return firstValue, nil

	case *data.ClassValue:
		// 处理 Iterator 对象
		// 检查是否实现了 Iterator 接口
		if targetInterface, ok := ctx.GetVM().GetInterface("Iterator"); ok {
			if checkInterfaceStructure(val.Class, targetInterface) {
				// 重置迭代器到开始位置
				if ctl := callVoidMethod(val, "rewind"); ctl != nil {
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

				// 获取当前元素（第一个元素）
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

func (f *ResetFunction) GetName() string {
	return "reset"
}

func (f *ResetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *ResetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

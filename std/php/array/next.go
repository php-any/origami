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
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		// 处理数组：对于 ArrayValue，next() 移动到第二个元素
		// 在 PHP 中，如果没有调用过指针函数，next() 会先移动到第二个元素
		if len(val.List) < 2 {
			return data.NewNullValue(), nil
		}
		// 返回第二个元素（模拟指针移动到第二个位置）
		return val.List[1].Value, nil

	case *data.ObjectValue:
		// 处理对象（关联数组）：移动到第二个元素
		var secondValue data.Value
		var count int
		var hasValue bool

		// 使用 RangeProperties 按插入顺序遍历，获取第二个元素
		val.RangeProperties(func(key string, value data.Value) bool {
			count++
			if count == 2 {
				secondValue = value
				hasValue = true
				return false // 找到第二个元素后停止
			}
			return true
		})

		if !hasValue {
			return data.NewNullValue(), nil
		}
		return secondValue, nil

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

func (f *NextFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *NextFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

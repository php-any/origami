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
		// 处理数组：对于 ArrayValue，当前元素是第一个元素（因为没有内部指针）
		// 在 PHP 中，如果没有调用过指针函数，current() 返回第一个元素
		if len(val.Value) == 0 {
			return data.NewNullValue(), nil
		}
		// 返回第一个元素（模拟当前指针在第一个位置）
		return val.Value[0], nil

	case *data.ObjectValue:
		// 处理对象（关联数组）：返回第一个元素
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
				// 获取当前元素（不移动指针）
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

func (f *CurrentFunction) GetName() string {
	return "current"
}

func (f *CurrentFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *CurrentFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

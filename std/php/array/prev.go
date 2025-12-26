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
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		// 处理数组：对于 ArrayValue，prev() 移动到倒数第二个元素
		// 在 PHP 中，如果没有调用过指针函数，prev() 会移动到倒数第二个元素
		if len(val.Value) < 2 {
			return data.NewNullValue(), nil
		}
		// 返回倒数第二个元素（模拟指针向后移动）
		return val.Value[len(val.Value)-2], nil

	case *data.ObjectValue:
		// 处理对象（关联数组）：移动到倒数第二个元素
		var prevValue data.Value
		var count int
		var totalCount int

		// 先计算总数
		val.RangeProperties(func(key string, value data.Value) bool {
			totalCount++
			return true
		})

		if totalCount < 2 {
			return data.NewNullValue(), nil
		}

		// 使用 RangeProperties 按插入顺序遍历，获取倒数第二个元素
		val.RangeProperties(func(key string, value data.Value) bool {
			count++
			if count == totalCount-1 {
				prevValue = value
				return false // 找到倒数第二个元素后停止
			}
			return true
		})

		if prevValue == nil {
			return data.NewNullValue(), nil
		}
		return prevValue, nil

	case *data.ClassValue:
		// 处理 Iterator 对象
		// 注意：标准 Iterator 接口没有 prev() 方法，所以对于 Iterator 对象，我们返回 null
		// 或者可以尝试调用 prev() 方法（如果存在）
		// 检查是否实现了 Iterator 接口
		if targetInterface, ok := ctx.GetVM().GetInterface("Iterator"); ok {
			if checkInterfaceStructure(val.Class, targetInterface) {
				// Iterator 接口没有 prev() 方法，返回 null
				// 如果需要支持 prev()，需要扩展 Iterator 接口或使用其他方式
				return data.NewNullValue(), nil
			}
		}
		// 不是 Iterator 接口，返回 null
		return data.NewNullValue(), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *PrevFunction) GetName() string {
	return "prev"
}

func (f *PrevFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *PrevFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

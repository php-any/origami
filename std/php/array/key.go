package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// KeyFunction 实现 key 函数
// 返回数组中当前元素的键名
type KeyFunction struct{}

func NewKeyFunction() data.FuncStmt {
	return &KeyFunction{}
}

func (f *KeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		// 处理数组：对于 ArrayValue，当前键是第一个索引（0）
		if len(val.List) == 0 {
			return data.NewNullValue(), nil
		}
		// 返回第一个索引（0）
		return data.NewIntValue(0), nil

	case *data.ObjectValue:
		// 处理对象（关联数组）：返回第一个键
		var firstKey string
		var hasValue bool

		// 使用 RangeProperties 按插入顺序遍历，获取第一个键
		val.RangeProperties(func(key string, value data.Value) bool {
			if !hasValue {
				firstKey = key
				hasValue = true
			}
			return false // 只获取第一个键
		})

		if !hasValue {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(firstKey), nil

	case *data.ClassValue:
		// 处理 Iterator 对象
		// 检查是否实现了 Iterator 接口
		if targetInterface, ok := ctx.GetVM().GetInterface("Iterator"); ok {
			if checkInterfaceStructure(val.Class, targetInterface) {
				// 获取当前键（不移动指针）
				keyVal, ctl := callValueMethod(val, "key")
				if ctl != nil {
					return data.NewNullValue(), nil
				}
				return keyVal, nil
			}
		}
		// 不是 Iterator 接口，返回 null
		return data.NewNullValue(), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *KeyFunction) GetName() string {
	return "key"
}

func (f *KeyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *KeyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

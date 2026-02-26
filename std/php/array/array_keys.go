package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayKeysFunction 实现 array_keys 函数
// 返回数组中所有的键（对普通数组为 0..n-1，对关联数组为属性名）
type ArrayKeysFunction struct{}

func NewArrayKeysFunction() data.FuncStmt {
	return &ArrayKeysFunction{}
}

func (f *ArrayKeysFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：数组
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		// 处理普通数组（数值索引）
		length := len(v.List)
		keys := make([]data.Value, 0, length)
		for i := 0; i < length; i++ {
			keys = append(keys, data.NewIntValue(i))
		}
		return data.NewArrayValue(keys), nil

	case *data.ObjectValue:
		// 处理对象（关联数组），按插入顺序（OrderedMap 顺序）返回键
		keys := make([]data.Value, 0)
		v.RangeProperties(func(key string, _ data.Value) bool {
			keys = append(keys, data.NewStringValue(key))
			return true
		})
		return data.NewArrayValue(keys), nil

	default:
		// 不是数组类型，返回空数组
		return data.NewArrayValue([]data.Value{}), nil
	}
}

func (f *ArrayKeysFunction) GetName() string {
	return "array_keys"
}

func (f *ArrayKeysFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
	}
}

func (f *ArrayKeysFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	}
}

package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayKeyFirstFunction 实现 PHP 内置函数 array_key_first
// array_key_first(array $array): int|string|null
//
// 语义（简化版，与 Symfony Console 用法兼容即可）：
// - 若参数不是数组/对象，返回 null
// - 若数组/对象为空，返回 null
// - 对普通索引数组，返回第一个索引（0）
// - 对关联数组/对象（使用 ObjectValue 存属性），返回按插入顺序的第一个键
type ArrayKeyFirstFunction struct{}

func NewArrayKeyFirstFunction() data.FuncStmt {
	return &ArrayKeyFirstFunction{}
}

func (f *ArrayKeyFirstFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewNullValue(), nil
	}

	switch v := val.(type) {
	case *data.ArrayValue:
		if len(v.List) == 0 {
			return data.NewNullValue(), nil
		}
		// PHP 中 array_key_first 对纯索引数组返回第一个索引 0
		return data.NewIntValue(0), nil
	case *data.ObjectValue:
		// 对象（在 Origami 中常用于表示关联数组）：
		// 使用 RangeProperties 按插入顺序遍历，取第一个键。
		var firstKey string
		found := false
		v.RangeProperties(func(key string, _ data.Value) bool {
			firstKey = key
			found = true
			return false // 只取第一个
		})
		if !found {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(firstKey), nil
	default:
		// 非数组/对象，返回 null（足以满足 Symfony 对 array_key_first 的使用场景）
		return data.NewNullValue(), nil
	}
}

func (f *ArrayKeyFirstFunction) GetName() string {
	return "array_key_first"
}

func (f *ArrayKeyFirstFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
	}
}

func (f *ArrayKeyFirstFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	}
}

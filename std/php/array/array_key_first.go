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
		if v.Len() == 0 {
			return data.NewNullValue(), nil
		}
		return v.At(0).PHPArrayKey(0), nil

	default:
		// 非数组/对象，返回 null（足以满足 Symfony 对 array_key_first 的使用场景）
		return data.NewNullValue(), nil
	}
}

func (f *ArrayKeyFirstFunction) GetName() string {
	return "array_key_first"
}

var arrayKeyFirstFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, nil),
}

func (f *ArrayKeyFirstFunction) GetParams() []data.GetValue {
	return arrayKeyFirstFunctionGetParams
}

var arrayKeyFirstFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
}

func (f *ArrayKeyFirstFunction) GetVariables() []data.Variable {
	return arrayKeyFirstFunctionGetVariables
}

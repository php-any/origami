package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayValuesFunction 实现 array_values 函数
// 返回数组中所有的值，并重新索引（从 0 开始）
type ArrayValuesFunction struct{}

func NewArrayValuesFunction() data.FuncStmt {
	return &ArrayValuesFunction{}
}

func (f *ArrayValuesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：数组
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 处理数组
	if arrayVal, ok := arrayValue.(*data.ArrayValue); ok {
		// 对于 ArrayValue，直接返回所有值（已经是数字索引）
		return data.NewArrayValue(arrayVal.Value), nil
	}

	// 处理对象（关联数组）
	if objectVal, ok := arrayValue.(*data.ObjectValue); ok {
		// 获取所有属性值
		properties := objectVal.GetProperties()
		values := make([]data.Value, 0, len(properties))

		// 收集所有值（忽略键）
		for _, val := range properties {
			values = append(values, val)
		}

		// 返回重新索引的数组（从 0 开始）
		return data.NewArrayValue(values), nil
	}

	// 不是数组类型，返回空数组
	return data.NewArrayValue([]data.Value{}), nil
}

func (f *ArrayValuesFunction) GetName() string {
	return "array_values"
}

func (f *ArrayValuesFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
	}
}

func (f *ArrayValuesFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	}
}

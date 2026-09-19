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
		return nil, throwMustBeArray("array_values", nil)
	}

	if arrayVal, ok := arrayValue.(*data.ArrayValue); ok {
		return data.NewArrayValue(arrayVal.ToValueList()), nil
	}

	if objectVal, ok := arrayValue.(*data.ObjectValue); ok {
		values := make([]data.Value, 0)
		objectVal.RangeProperties(func(_ string, val data.Value) bool {
			values = append(values, val)
			return true
		})
		return data.NewArrayValue(values), nil
	}

	return nil, throwMustBeArray("array_values", arrayValue)
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

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

	return nil, throwMustBeArray("array_values", arrayValue)
}

func (f *ArrayValuesFunction) GetName() string {
	return "array_values"
}

var arrayValuesFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, nil),
}

func (f *ArrayValuesFunction) GetParams() []data.GetValue {
	return arrayValuesFunctionGetParams
}

var arrayValuesFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
}

func (f *ArrayValuesFunction) GetVariables() []data.Variable {
	return arrayValuesFunctionGetVariables
}

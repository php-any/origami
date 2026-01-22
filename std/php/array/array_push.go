package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayPushFunction() data.FuncStmt {
	return &ArrayPushFunction{}
}

type ArrayPushFunction struct{}

func (f *ArrayPushFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewIntValue(0), nil
	}

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		// 获取 Parameters 参数（包含所有传入的值）
		paramsValue, _ := ctx.GetIndexValue(1)
		if paramsValue != nil {
			if paramsArray, ok := paramsValue.(*data.ArrayValue); ok {
				// Parameters 返回的是 ArrayValue，包含所有参数
				paramsList := paramsArray.ToValueList()
				for _, val := range paramsList {
					v.List = append(v.List, data.NewZVal(val))
				}
			}
		}

		return data.NewIntValue(len(v.List)), nil
	default:
		return data.NewIntValue(0), nil
	}
}

func (f *ArrayPushFunction) GetName() string {
	return "array_push"
}

func (f *ArrayPushFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
		node.NewParameters(nil, "values", 1, nil, nil),
	}
}

func (f *ArrayPushFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "values", 1, data.Mixed{}),
	}
}

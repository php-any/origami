package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayMergeFunction() data.FuncStmt {
	return &ArrayMergeFunction{}
}

type ArrayMergeFunction struct{}

func (f *ArrayMergeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 Parameters 参数（包含所有传入的数组）
	paramsValue, _ := ctx.GetIndexValue(0)
	if paramsValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// Parameters 返回的是 ArrayValue，包含所有参数
	paramsArray, ok := paramsValue.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 收集所有数组参数的值
	var allValues []data.Value

	// 遍历 Parameters 中的每个数组参数
	for _, paramValue := range paramsArray.Value {
		// 处理数组
		if arrayVal, ok := paramValue.(*data.ArrayValue); ok {
			allValues = append(allValues, arrayVal.Value...)
		} else if objectVal, ok := paramValue.(*data.ObjectValue); ok {
			// 处理对象（关联数组）
			properties := objectVal.GetProperties()
			for _, val := range properties {
				allValues = append(allValues, val)
			}
		} else {
			// 非数组类型，直接添加
			allValues = append(allValues, paramValue)
		}
	}

	return data.NewArrayValue(allValues), nil
}

func (f *ArrayMergeFunction) GetName() string {
	return "array_merge"
}

func (f *ArrayMergeFunction) GetParams() []data.GetValue {
	// 使用可变参数
	return []data.GetValue{
		node.NewParameters(nil, "arrays", 0, nil, nil),
	}
}

func (f *ArrayMergeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
	}
}

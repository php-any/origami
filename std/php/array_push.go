package php

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

	// 检查是否为数组引用
	arrayRef, ok := arrayValue.(*data.ReferenceValue)
	if !ok {
		return data.NewIntValue(0), nil
	}

	// 获取数组值
	parentCtx := arrayRef.Ctx
	varRef := arrayRef.Val

	v, acl := varRef.GetValue(parentCtx)
	if acl != nil {
		return data.NewIntValue(0), nil
	}

	if v == nil {
		return data.NewIntValue(0), nil
	}

	internalV, intervalCtl := v.GetValue(parentCtx)
	if intervalCtl != nil {
		return data.NewIntValue(0), nil
	}

	arrayVal, ok := internalV.(*data.ArrayValue)
	if !ok {
		return data.NewIntValue(0), nil
	}

	// 收集所有要添加的值
	var newValues []data.Value

	// 获取 Parameters 参数（包含所有传入的值）
	paramsValue, _ := ctx.GetIndexValue(1)
	if paramsValue != nil {
		if paramsArray, ok := paramsValue.(*data.ArrayValue); ok {
			// Parameters 返回的是 ArrayValue，包含所有参数
			newValues = append(newValues, paramsArray.Value...)
		}
	}

	// 添加新值到数组
	arrayVal.Value = append(arrayVal.Value, newValues...)

	// 更新数组值
	parentCtx.SetVariableValue(varRef, arrayVal)

	return data.NewIntValue(len(arrayVal.Value)), nil
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

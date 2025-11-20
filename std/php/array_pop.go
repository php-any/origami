package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayPopFunction() data.FuncStmt {
	return &ArrayPopFunction{}
}

type ArrayPopFunction struct{}

func (f *ArrayPopFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 检查是否为数组引用
	arrayRef, ok := arrayValue.(*data.ReferenceValue)
	if !ok {
		return data.NewNullValue(), nil
	}

	// 获取数组值
	parentCtx := arrayRef.Ctx
	varRef := arrayRef.Val

	v, acl := varRef.GetValue(parentCtx)
	if acl != nil {
		return data.NewNullValue(), nil
	}

	if v == nil {
		return data.NewNullValue(), nil
	}

	internalV, intervalCtl := v.GetValue(parentCtx)
	if intervalCtl != nil {
		return data.NewNullValue(), nil
	}

	arrayVal, ok := internalV.(*data.ArrayValue)
	if !ok {
		return data.NewNullValue(), nil
	}

	// 如果数组为空，返回 null
	if len(arrayVal.Value) == 0 {
		return data.NewNullValue(), nil
	}

	// 弹出最后一个元素
	lastElement := arrayVal.Value[len(arrayVal.Value)-1]
	arrayVal.Value = arrayVal.Value[:len(arrayVal.Value)-1]

	// 更新数组值
	parentCtx.SetVariableValue(varRef, arrayVal)

	return lastElement, nil
}

func (f *ArrayPopFunction) GetName() string {
	return "array_pop"
}

func (f *ArrayPopFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *ArrayPopFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

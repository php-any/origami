package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIssetFunction() data.FuncStmt {
	return &IssetFunction{}
}

type IssetFunction struct{}

func (f *IssetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	params := f.GetParams()
	if len(params) == 0 {
		return data.NewBoolValue(false), nil
	}

	// 获取第一个参数的值
	varParam := params[0]
	varValue, ctl := varParam.GetValue(ctx)
	if ctl != nil {
		return data.NewBoolValue(false), nil
	}

	// 检查参数值是否是 ReferenceValue（变量引用）
	// 如果是，说明参数是通过引用传递的，需要检查父级上下文中的变量
	if refValue, ok := varValue.(*data.ReferenceValue); ok {
		// 使用父级上下文和变量引用检查变量是否存在
		parentCtx := refValue.Ctx
		varRef := refValue.Val

		// 使用 GetIndexValue 安全地检查父级上下文中的变量（避免索引越界 panic）
		v, ok := parentCtx.GetIndexValue(varRef.GetIndex())
		if !ok {
			// 变量在父级上下文中不存在
			return data.NewBoolValue(false), nil
		}

		// 如果变量值为 nil，返回 false
		if v == nil {
			return data.NewBoolValue(false), nil
		}

		// 获取实际值
		internalV, intervalCtl := v.GetValue(parentCtx)
		if intervalCtl != nil {
			return data.NewBoolValue(false), nil
		}

		// 检查值是否为 null
		if internalV == nil {
			return data.NewBoolValue(false), nil
		}

		// 检查是否为 NullValue 类型
		if _, ok := internalV.(*data.NullValue); ok {
			return data.NewBoolValue(false), nil
		}

		// 值已设置且不为 null
		return data.NewBoolValue(true), nil
	}

	// 如果不是引用类型，说明参数是值传递，直接检查值是否为 null
	if varValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取实际值
	internalV, intervalCtl := varValue.GetValue(ctx)
	if intervalCtl != nil {
		return data.NewBoolValue(false), nil
	}

	// 检查值是否为 null
	if internalV == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查是否为 NullValue 类型
	if _, ok := internalV.(*data.NullValue); ok {
		return data.NewBoolValue(false), nil
	}

	// 值已设置且不为 null
	return data.NewBoolValue(true), nil
}

func (f *IssetFunction) GetName() string {
	return "isset"
}

func (f *IssetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "var", 0, data.Mixed{}),
	}
}

func (f *IssetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "var", 0, data.Mixed{}),
	}
}

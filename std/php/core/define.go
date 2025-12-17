package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DefineFunction 实现 define 函数
type DefineFunction struct{}

func NewDefineFunction() data.FuncStmt {
	return &DefineFunction{}
}

func (f *DefineFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：常量名（必须是字符串）
	constNameValue, _ := ctx.GetIndexValue(0)
	if constNameValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 将常量名转换为字符串
	var constName string
	if str, ok := constNameValue.(data.AsString); ok {
		constName = str.AsString()
	} else {
		constName = constNameValue.AsString()
	}

	// 获取第二个参数：常量值
	valueExpr, _ := ctx.GetIndexValue(1)
	if valueExpr == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取实际值
	actualValue, intervalCtl := valueExpr.GetValue(ctx)
	if intervalCtl != nil {
		return data.NewBoolValue(false), intervalCtl
	}

	if actualValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 将值转换为 Value 类型
	value, ok := actualValue.(data.Value)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	// 在 VM 中设置全局常量
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewBoolValue(false), nil
	}

	// 设置常量（如果已存在会返回错误）
	ctl := vm.SetConstant(constName, value)
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}

	// define 函数返回 true 表示成功
	return data.NewBoolValue(true), nil
}

func (f *DefineFunction) GetName() string {
	return "define"
}

func (f *DefineFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "constant_name", 0, nil, data.String{}),
		node.NewParameter(nil, "value", 1, nil, data.Mixed{}),
		node.NewParameter(nil, "case_insensitive", 2, node.NewBooleanLiteral(nil, false), data.Bool{}),
	}
}

func (f *DefineFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "constant_name", 0, data.String{}),
		node.NewVariable(nil, "value", 1, data.Mixed{}),
		node.NewVariable(nil, "case_insensitive", 2, data.Bool{}),
	}
}

package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DefinedFunction 实现 defined 函数
type DefinedFunction struct{}

func NewDefinedFunction() data.FuncStmt {
	return &DefinedFunction{}
}

func (f *DefinedFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取常量名参数
	constNameValue, ok := ctx.GetIndexValue(0)
	if !ok || constNameValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 将常量名转换为字符串
	var constName string
	if str, ok := constNameValue.(data.AsString); ok {
		constName = str.AsString()
	} else {
		constName = constNameValue.AsString()
	}

	// 从 VM 中获取全局常量
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查常量是否存在
	_, exists := vm.GetConstant(constName)
	return data.NewBoolValue(exists), nil
}

func (f *DefinedFunction) GetName() string {
	return "defined"
}

func (f *DefinedFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "constant_name", 0, nil, data.String{}),
	}
}

func (f *DefinedFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "constant_name", 0, data.String{}),
	}
}

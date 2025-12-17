package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ExtensionLoadedFunction 实现 extension_loaded 函数
// 当前实现固定返回 true，用于兼容性检测。
type ExtensionLoadedFunction struct{}

func NewExtensionLoadedFunction() data.FuncStmt {
	return &ExtensionLoadedFunction{}
}

func (f *ExtensionLoadedFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 忽略具体扩展名，固定返回 true
	return data.NewBoolValue(true), nil
}

func (f *ExtensionLoadedFunction) GetName() string {
	return "extension_loaded"
}

func (f *ExtensionLoadedFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "extension", 0, nil, data.String{}),
	}
}

func (f *ExtensionLoadedFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "extension", 0, data.String{}),
	}
}

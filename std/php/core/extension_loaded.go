package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ExtensionLoadedFunction 实现 extension_loaded 函数
type ExtensionLoadedFunction struct{}

func NewExtensionLoadedFunction() data.FuncStmt {
	return &ExtensionLoadedFunction{}
}

func (f *ExtensionLoadedFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	extension, _ := ctx.GetIndexValue(0)
	name := ""
	if extension != nil {
		name = strings.ToLower(extension.AsString())
	}

	// 不可把未实现的原生扩展报告为已加载，否则依赖会选择不可用的快路径。
	switch name {
	case "gmp", "bcmath":
		return data.NewBoolValue(false), nil
	}
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

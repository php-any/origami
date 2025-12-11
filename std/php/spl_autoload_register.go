package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SplAutoloadRegisterFunction 实现 spl_autoload_register
// 当前实现为占位，记录回调的功能可后续扩展
type SplAutoloadRegisterFunction struct{}

func NewSplAutoloadRegisterFunction() data.FuncStmt { return &SplAutoloadRegisterFunction{} }

func (f *SplAutoloadRegisterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 目前不执行实际的自动加载逻辑，只返回 true 以兼容调用
	return data.NewBoolValue(true), nil
}

func (f *SplAutoloadRegisterFunction) GetName() string { return "spl_autoload_register" }

func (f *SplAutoloadRegisterFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}

func (f *SplAutoloadRegisterFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
	}
}

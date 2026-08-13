package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewGetDefinedVarsFunction() data.FuncStmt {
	return &GetDefinedVarsFunction{}
}

// GetDefinedVarsFunction 实现 PHP 内置函数 get_defined_vars。
// 使用调用者上下文，返回当前作用域符号表的关联数组快照。
type GetDefinedVarsFunction struct{}

func (f *GetDefinedVarsFunction) GetName() string { return "get_defined_vars" }

func (f *GetDefinedVarsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	result := data.NewObjectValue()
	for name, value := range ctx.GetDefinedVariables() {
		result.SetProperty(name, value)
	}
	return result, nil
}

func (f *GetDefinedVarsFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewCallerContextParameter(nil)}
}

func (f *GetDefinedVarsFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

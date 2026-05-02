package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HeaderFunction 实现 header 函数（与 Net/Http 服务器集成）
type HeaderFunction struct{}

func NewHeaderFunction() data.FuncStmt { return &HeaderFunction{} }

func (f *HeaderFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// header() 通常用于设置 HTTP 响应头
	// 在当前 origami HTTP 服务器中，响应头已由 Response 对象管理
	// 这里仅返回 null 以保持兼容性
	return data.NewNullValue(), nil
}

func (f *HeaderFunction) GetName() string { return "header" }
func (f *HeaderFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "header", 0, nil, nil),
		node.NewParameter(nil, "replace", 1, node.NewBooleanLiteral(nil, true), nil),
		node.NewParameter(nil, "response_code", 2, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (f *HeaderFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "header", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "replace", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "response_code", 2, data.NewBaseType("int")),
	}
}

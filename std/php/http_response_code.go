package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HttpResponseCodeFunction 实现 http_response_code 函数
type HttpResponseCodeFunction struct{}

func NewHttpResponseCodeFunction() data.FuncStmt { return &HttpResponseCodeFunction{} }

func (f *HttpResponseCodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v != nil {
		return v, nil
	}
	return data.NewIntValue(200), nil
}

func (f *HttpResponseCodeFunction) GetName() string { return "http_response_code" }
func (f *HttpResponseCodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "response_code", 0, node.NewIntLiteral(nil, "200"), nil)}
}
func (f *HttpResponseCodeFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "response_code", 0, data.NewBaseType("int"))}
}

package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ParseStrFunction 实现 parse_str。
type ParseStrFunction struct{}

func NewParseStrFunction() data.FuncStmt { return &ParseStrFunction{} }

func (f *ParseStrFunction) GetName() string { return "parse_str" }

var parseStrFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, data.NewBaseType("string")),
	node.NewOutputParameterReference(nil, "result", 1, nil, data.NewBaseType("array")),
}

func (f *ParseStrFunction) GetParams() []data.GetValue {
	return parseStrFunctionGetParams
}

var parseStrFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "result", 1, data.NewBaseType("array")),
}

func (f *ParseStrFunction) GetVariables() []data.Variable {
	return parseStrFunctionGetVariables
}

func (f *ParseStrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	s := ""
	if strVal != nil {
		s = strVal.AsString()
	}
	parsed := data.ParseFormFields(s)

	slot := ctx.GetIndexZVal(1)
	prepared, ctl := slot.PrepareWrite(parsed, ctx)
	if ctl != nil {
		return nil, ctl
	}
	data.CowAssign(slot, prepared)
	return data.NewNullValue(), nil
}

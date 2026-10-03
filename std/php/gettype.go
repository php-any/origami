package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewGettypeFunction() data.FuncStmt {
	return &GettypeFunction{}
}

type GettypeFunction struct{}

func (f *GettypeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "value", 0, data.Mixed{}))
	if ctl != nil {
		return nil, ctl
	}

	tp := "unknown"
	switch data.ValueKindOf(v) {
	case data.ValueArray:
		tp = "array"
	case data.ValueBool:
		tp = "boolean"
	case data.ValueResource:
		tp = "resource"
		if resource, ok := v.(interface{ IsPHPResourceOpen() bool }); ok && !resource.IsPHPResourceOpen() {
			tp = "resource (closed)"
		}
	case data.ValueObject:
		tp = "object"
	case data.ValueFloat:
		tp = "double"
	case data.ValueInt:
		tp = "integer"
	case data.ValueString:
		tp = "string"
	case data.ValueNull:
		tp = "NULL"
	}
	return data.NewStringValue(tp), nil
}

func (f *GettypeFunction) GetName() string {
	return "gettype"
}

var gettypeFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, data.Mixed{}),
}

func (f *GettypeFunction) GetParams() []data.GetValue {
	return gettypeFunctionGetParams
}

var gettypeFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.Mixed{}),
}

func (f *GettypeFunction) GetVariables() []data.Variable {
	return gettypeFunctionGetVariables
}

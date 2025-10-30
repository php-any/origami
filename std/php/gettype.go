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

	internalV, intervalCtl := v.GetValue(ctx)
	if intervalCtl != nil {
		return nil, intervalCtl
	}

	tp := "unknown"
	switch internalV.(type) {
	case *data.ArrayValue:
		tp = "array"
	case *data.BoolValue:
		tp = "bool"
	case *data.ClassValue:
		tp = "class"
	case *data.FloatValue:
		tp = "float"
	case *data.IntValue:
		tp = "int"
	case *data.ObjectValue:
		tp = "object"
	case *data.StringValue:
		tp = "string"
	case *data.NullValue:
		tp = "null"
	case *data.AnyValue:
		tp = "any"
	}
	return data.NewStringValue(tp), intervalCtl
}

func (f *GettypeFunction) GetName() string {
	return "gettype"
}

func (f *GettypeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.Mixed{}),
	}
}

func (f *GettypeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
	}
}

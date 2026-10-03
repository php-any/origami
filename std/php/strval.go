package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrvalFunction() data.FuncStmt {
	return &StrvalFunction{}
}

type StrvalFunction struct{}

func (f *StrvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return strvalValue(ctx, v)
}

func strvalValue(ctx data.Context, v data.Value) (data.GetValue, data.Control) {
	if v == nil {
		return data.NewStringValue(""), nil
	}
	switch val := v.(type) {
	case *data.NullValue:
		return data.NewStringValue(""), nil
	case *data.BoolValue:
		if val.Value {
			return data.NewStringValue("1"), nil
		}
		return data.NewStringValue(""), nil
	case *data.StringValue:
		return val, nil
	case *data.IntValue, *data.FloatValue:
		return data.NewStringValue(val.AsString()), nil
	case *data.ArrayValue:
		return data.NewStringValue("Array"), nil
	case *data.ClassValue:
		result, accepted, ctl := data.ObjectToStringValue(val, ctx)
		if ctl != nil {
			return nil, ctl
		}
		if accepted {
			return result, nil
		}
		return data.NewStringValue("Object"), nil
	case *data.ThisValue:
		return strvalValue(ctx, val.ClassValue)

	default:
		return data.NewStringValue(v.AsString()), nil
	}
}

func (f *StrvalFunction) GetName() string {
	return "strval"
}

var strvalFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, data.Mixed{}),
}

func (f *StrvalFunction) GetParams() []data.GetValue {
	return strvalFunctionGetParams
}

var strvalFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.Mixed{}),
}

func (f *StrvalFunction) GetVariables() []data.Variable {
	return strvalFunctionGetVariables
}

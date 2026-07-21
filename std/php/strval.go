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
	return strvalValue(v), nil
}

func strvalValue(v data.Value) data.Value {
	if v == nil {
		return data.NewStringValue("")
	}
	switch val := v.(type) {
	case *data.NullValue:
		return data.NewStringValue("")
	case *data.BoolValue:
		if val.Value {
			return data.NewStringValue("1")
		}
		return data.NewStringValue("")
	case *data.StringValue:
		return val
	case *data.IntValue, *data.FloatValue:
		return data.NewStringValue(val.AsString())
	case *data.ArrayValue:
		return data.NewStringValue("Array")
	case *data.ClassValue:
		if method, ok := val.GetMethod("__toString"); ok {
			ret, ctl := method.Call(val.Context)
			if ctl == nil {
				if sv, ok := ret.(data.Value); ok {
					return data.NewStringValue(sv.AsString())
				}
			}
		}
		return data.NewStringValue("Object")
	case *data.ObjectValue:
		return data.NewStringValue("Object")
	default:
		return data.NewStringValue(v.AsString())
	}
}

func (f *StrvalFunction) GetName() string {
	return "strval"
}

func (f *StrvalFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.Mixed{}),
	}
}

func (f *StrvalFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
	}
}

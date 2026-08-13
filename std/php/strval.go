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
	return strvalValue(v)
}

func strvalValue(v data.Value) (data.GetValue, data.Control) {
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
		if method, ok := val.GetMethod("__toString"); ok {
			fnCtx := val.CreateContext(method.GetVariables())
			fnCtx.SetCallArgs([]data.GetValue{})
			ret, ctl := method.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			if sv, ok := ret.(data.Value); ok {
				return data.NewStringValue(sv.AsString()), nil
			}
			if gv, ok := ret.(data.GetValue); ok {
				v2, ctl2 := gv.GetValue(fnCtx)
				if ctl2 != nil {
					return nil, ctl2
				}
				if sv, ok := v2.(data.Value); ok {
					return data.NewStringValue(sv.AsString()), nil
				}
			}
		}
		return data.NewStringValue("Object"), nil
	case *data.ObjectValue:
		return data.NewStringValue("Object"), nil
	default:
		return data.NewStringValue(v.AsString()), nil
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

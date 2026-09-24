package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// BoolvalFunction 实现 PHP boolval(mixed $value): bool
type BoolvalFunction struct{}

func NewBoolvalFunction() data.FuncStmt {
	return &BoolvalFunction{}
}

func (f *BoolvalFunction) GetName() string { return "boolval" }

var boolvalFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
}

func (f *BoolvalFunction) GetParams() []data.GetValue {
	return boolvalFunctionGetParams
}

var boolvalFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, nil),
}

func (f *BoolvalFunction) GetVariables() []data.Variable {
	return boolvalFunctionGetVariables
}

func (f *BoolvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	if b, ok := v.(data.AsBool); ok {
		bv, err := b.AsBool()
		if err == nil {
			return data.NewBoolValue(bv), nil
		}
	}
	// 回退：非空字符串/非零数字等由 AsString/AsInt 语义覆盖不足时按 PHP truthiness
	switch t := v.(type) {
	case *data.BoolValue:
		return data.NewBoolValue(t.Value), nil
	case *data.NullValue:
		return data.NewBoolValue(false), nil
	case *data.IntValue:
		return data.NewBoolValue(t.Value != 0), nil
	case *data.FloatValue:
		f64, _ := t.AsFloat()
		return data.NewBoolValue(f64 != 0), nil
	case *data.StringValue:
		return data.NewBoolValue(t.Value != "" && t.Value != "0"), nil
	case *data.ArrayValue:
		return data.NewBoolValue(len(t.List) > 0), nil
	default:
		return data.NewBoolValue(true), nil
	}
}

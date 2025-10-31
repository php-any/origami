package std

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStringFunction() data.FuncStmt { return &StringFunction{} }

type StringFunction struct{}

func (f *StringFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}

	if s, ok := v.(data.AsString); ok {
		return data.NewStringValue(s.AsString()), nil
	}

	// Fallback: Value implements AsString in our Value interface
	return data.NewStringValue(v.AsString()), nil
}

func (f *StringFunction) GetName() string { return "string" }

func (f *StringFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *StringFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}

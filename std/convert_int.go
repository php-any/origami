package std

import (
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIntFunction() data.FuncStmt { return &IntFunction{} }

type IntFunction struct{}

func (f *IntFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}

	// Try common conversions in order
	switch tv := v.(type) {
	case data.AsInt:
		if i, err := tv.AsInt(); err == nil {
			return data.NewIntValue(i), nil
		}
	case data.AsFloat:
		if f64, err := tv.AsFloat(); err == nil {
			return data.NewIntValue(int(f64)), nil
		}
	case data.AsBool:
		if b, err := tv.AsBool(); err == nil {
			if b {
				return data.NewIntValue(1), nil
			}
			return data.NewIntValue(0), nil
		}
	case data.AsString:
		if i, err := strconv.Atoi(tv.AsString()); err == nil {
			return data.NewIntValue(i), nil
		}
	}

	// Fallback to string parse if possible
	if s, ok := v.(data.AsString); ok {
		if i, err := strconv.Atoi(s.AsString()); err == nil {
			return data.NewIntValue(i), nil
		}
	}
	return data.NewIntValue(0), nil
}

func (f *IntFunction) GetName() string { return "int" }

func (f *IntFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IntFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}

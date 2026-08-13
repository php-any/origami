package std

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewBoolFunction() data.FuncStmt { return &BoolFunction{} }

type BoolFunction struct{}

func (f *BoolFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}

	switch tv := v.(type) {
	case data.AsBool:
		if b, err := tv.AsBool(); err == nil {
			return data.NewBoolValue(b), nil
		}
	case data.AsString:
		s := strings.TrimSpace(strings.ToLower(tv.AsString()))
		switch s {
		case "", "0", "false", "no", "off", "null", "nil":
			return data.NewBoolValue(false), nil
		default:
			return data.NewBoolValue(true), nil
		}
	}

	// Fallback: non-nil defaults to true
	return data.NewBoolValue(true), nil
}

func (f *BoolFunction) GetName() string { return "bool" }

func (f *BoolFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *BoolFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}

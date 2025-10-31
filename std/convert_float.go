package std

import (
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewFloatFunction() data.FuncStmt { return &FloatFunction{} }

type FloatFunction struct{}

func (f *FloatFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}

	switch tv := v.(type) {
	case data.AsFloat:
		if f64, err := tv.AsFloat(); err == nil {
			return data.NewFloatValue(f64), nil
		}
	case data.AsInt:
		if i, err := tv.AsInt(); err == nil {
			return data.NewFloatValue(float64(i)), nil
		}
	case data.AsString:
		if f64, err := strconv.ParseFloat(tv.AsString(), 64); err == nil {
			return data.NewFloatValue(f64), nil
		}
	}

	if s, ok := v.(data.AsString); ok {
		if f64, err := strconv.ParseFloat(s.AsString(), 64); err == nil {
			return data.NewFloatValue(f64), nil
		}
	}

	return data.NewFloatValue(0), nil
}

func (f *FloatFunction) GetName() string { return "float" }

func (f *FloatFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *FloatFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}

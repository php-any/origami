package php

import (
	"math"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func asFloatForNanCheck(v data.Value) (float64, bool) {
	if v == nil {
		return 0, false
	}
	if fv, ok := v.(*data.FloatValue); ok {
		return fv.Value, true
	}
	return 0, false
}

// IsNanFunction 实现 is_nan
type IsNanFunction struct{}

func NewIsNanFunction() data.FuncStmt { return &IsNanFunction{} }

func (f *IsNanFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	fv, ok := asFloatForNanCheck(v)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(math.IsNaN(fv)), nil
}

func (f *IsNanFunction) GetName() string { return "is_nan" }
func (f *IsNanFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *IsNanFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, data.NewBaseType("float"))}
}

// IsInfiniteFunction 实现 is_infinite
type IsInfiniteFunction struct{}

func NewIsInfiniteFunction() data.FuncStmt { return &IsInfiniteFunction{} }

func (f *IsInfiniteFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	fv, ok := asFloatForNanCheck(v)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(math.IsInf(fv, 0)), nil
}

func (f *IsInfiniteFunction) GetName() string { return "is_infinite" }
func (f *IsInfiniteFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *IsInfiniteFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, data.NewBaseType("float"))}
}

// IsFiniteFunction 实现 is_finite
type IsFiniteFunction struct{}

func NewIsFiniteFunction() data.FuncStmt { return &IsFiniteFunction{} }

func (f *IsFiniteFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	fv, ok := asFloatForNanCheck(v)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(!math.IsNaN(fv) && !math.IsInf(fv, 0)), nil
}

func (f *IsFiniteFunction) GetName() string { return "is_finite" }
func (f *IsFiniteFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *IsFiniteFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, data.NewBaseType("float"))}
}

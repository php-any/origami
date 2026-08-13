package php

import (
	"math"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PowFunction 实现 pow 函数
type PowFunction struct{}

func NewPowFunction() data.FuncStmt { return &PowFunction{} }

func (f *PowFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	baseV, _ := ctx.GetIndexValue(0)
	expV, _ := ctx.GetIndexValue(1)

	var base, exp float64
	if af, ok := baseV.(data.AsFloat); ok {
		base, _ = af.AsFloat()
	} else if ai, ok := baseV.(*data.IntValue); ok {
		base = float64(ai.Value)
	}
	if af, ok := expV.(data.AsFloat); ok {
		exp, _ = af.AsFloat()
	} else if ai, ok := expV.(*data.IntValue); ok {
		exp = float64(ai.Value)
	}

	return data.NewFloatValue(math.Pow(base, exp)), nil
}

func (f *PowFunction) GetName() string { return "pow" }
func (f *PowFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "base", 0, nil, nil),
		node.NewParameter(nil, "exp", 1, nil, nil),
	}
}
func (f *PowFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "base", 0, nil),
		node.NewVariable(nil, "exp", 1, nil),
	}
}

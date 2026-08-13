package php

import (
	"math"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RoundFunction 实现 round 函数
type RoundFunction struct{}

func NewRoundFunction() data.FuncStmt { return &RoundFunction{} }

func (f *RoundFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	precisionV, _ := ctx.GetIndexValue(1)

	if v == nil {
		return data.NewFloatValue(0), nil
	}

	var precision int
	if precisionV != nil {
		if iv, ok := precisionV.(*data.IntValue); ok {
			precision = iv.Value
		}
	}

	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		pow := math.Pow10(precision)
		return data.NewFloatValue(math.Round(fv*pow) / pow), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *RoundFunction) GetName() string { return "round" }
func (f *RoundFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "num", 0, nil, nil),
		node.NewParameter(nil, "precision", 1, node.NewIntLiteral(nil, "0"), nil),
	}
}
func (f *RoundFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "num", 0, nil),
		node.NewVariable(nil, "precision", 1, data.NewBaseType("int")),
	}
}

package php

import (
	gmath "math"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// AbsFunction 实现 abs(int|float $num): int|float
type AbsFunction struct{}

func NewAbsFunction() data.FuncStmt { return &AbsFunction{} }

func (f *AbsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	// 优先保持整数类型（PHP：abs(int) → int）
	if _, isFloat := v.(*data.FloatValue); !isFloat {
		if ai, ok := v.(data.AsInt); ok {
			iv, err := ai.AsInt()
			if err == nil {
				if iv < 0 {
					iv = -iv
				}
				return data.NewIntValue(iv), nil
			}
		}
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(gmath.Abs(fv)), nil
	}
	return data.NewIntValue(0), nil
}

func (f *AbsFunction) GetName() string { return "abs" }
func (f *AbsFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *AbsFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, nil)}
}

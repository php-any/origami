package php

import (
	"math"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CeilFunction 实现 ceil 函数
type CeilFunction struct{}

func NewCeilFunction() data.FuncStmt { return &CeilFunction{} }

func (f *CeilFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Ceil(fv)), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *CeilFunction) GetName() string { return "ceil" }
func (f *CeilFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *CeilFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, nil)}
}

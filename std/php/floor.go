package php

import (
	"math"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FloorFunction 实现 floor 函数
type FloorFunction struct{}

func NewFloorFunction() data.FuncStmt { return &FloorFunction{} }

func (f *FloorFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewFloatValue(0), nil
	}
	if af, ok := v.(data.AsFloat); ok {
		fv, _ := af.AsFloat()
		return data.NewFloatValue(math.Floor(fv)), nil
	}
	return data.NewFloatValue(0), nil
}

func (f *FloorFunction) GetName() string { return "floor" }
func (f *FloorFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "num", 0, nil, nil)}
}
func (f *FloorFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "num", 0, nil)}
}

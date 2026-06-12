package core

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type PrintfFunction struct{}

func NewPrintfFunction() data.FuncStmt {
	return &PrintfFunction{}
}

func (f *PrintfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatVal, _ := ctx.GetIndexValue(0)
	if formatVal == nil {
		return data.NewNullValue(), nil
	}
	format := formatVal.AsString()
	args := collectPrintfArgs(ctx)
	n, _ := fmt.Printf(format, args...)
	return data.NewIntValue(n), nil
}

func (f *PrintfFunction) GetName() string {
	return "printf"
}

func (f *PrintfFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, data.String{}),
	}
}

func (f *PrintfFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.String{}),
	}
}

func collectPrintfArgs(ctx data.Context) []interface{} {
	var args []interface{}
	for i := 1; ; i++ {
		val, _ := ctx.GetIndexValue(i)
		if val == nil {
			break
		}
		if v, ok := val.(*data.IntValue); ok {
			args = append(args, v.Value)
		} else if v, ok := val.(*data.FloatValue); ok {
			args = append(args, v.Value)
		} else if v, ok := val.(*data.BoolValue); ok {
			if v.Value {
				args = append(args, 1)
			} else {
				args = append(args, 0)
			}
		} else {
			args = append(args, val.AsString())
		}
	}
	return args
}

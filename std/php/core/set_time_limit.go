package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type SetTimeLimitFunction struct{}

func NewSetTimeLimitFunction() data.FuncStmt {
	return &SetTimeLimitFunction{}
}

func (f *SetTimeLimitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	secVal, _ := ctx.GetIndexValue(0)
	seconds := 0
	if secVal != nil {
		if iv, ok := secVal.(data.AsInt); ok {
			seconds, _ = iv.AsInt()
		}
	}
	SetExecutionDeadline(seconds)
	return data.NewBoolValue(true), nil
}

func (f *SetTimeLimitFunction) GetName() string { return "set_time_limit" }

func (f *SetTimeLimitFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "seconds", 0, nil, data.Int{}),
	}
}

func (f *SetTimeLimitFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "seconds", 0, data.Int{}),
	}
}

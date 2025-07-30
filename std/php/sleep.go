package php

import (
	"github.com/php-any/origami/node"
	"time"
)
import "github.com/php-any/origami/data"

func NewSleepFunction() data.FuncStmt {
	return &SleepFunction{}
}

type SleepFunction struct{}

func (f *SleepFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "seconds", 0, data.Int{}))
	if ctl != nil {
		return nil, ctl
	}
	i, _ := v.(data.AsInt).AsInt()
	time.Sleep(time.Duration(i) * time.Second)
	return nil, nil
}
func (f *SleepFunction) GetName() string {
	return "sleep"
}

func (f *SleepFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "seconds", 0, nil, nil),
	}
}

func (f *SleepFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "seconds", 0, data.Int{}),
	}
}

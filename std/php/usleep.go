package php

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type UsleepFunction struct{}

func NewUsleepFunction() data.FuncStmt {
	return &UsleepFunction{}
}

func (f *UsleepFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	microsecondsVal, _ := ctx.GetIndexValue(0)
	if microsecondsVal == nil {
		return data.NewNullValue(), nil
	}
	raw := strings.TrimSpace(microsecondsVal.AsString())
	if raw == "" {
		return data.NewNullValue(), nil
	}
	fval, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return data.NewNullValue(), nil
	}
	if fval > 0 {
		time.Sleep(time.Duration(int64(math.Ceil(fval))) * time.Microsecond)
	}
	return data.NewNullValue(), nil
}

func (f *UsleepFunction) GetName() string {
	return "usleep"
}

func (f *UsleepFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "microseconds", 0, nil, nil),
	}
}

func (f *UsleepFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "microseconds", 0, data.NewBaseType("int")),
	}
}

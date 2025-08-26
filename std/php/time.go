package php

import (
	"time"

	"github.com/php-any/origami/data"
)

func NewTimeFunction() data.FuncStmt {
	return &TimeFunction{}
}

type TimeFunction struct {
	data.Function
}

func (f *TimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(int(time.Now().Unix())), nil
}
func (f *TimeFunction) GetName() string {
	return "time"
}

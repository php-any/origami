package php

import (
	"fmt"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewMicrotimeFunction() data.FuncStmt {
	return &MicrotimeFunction{}
}

type MicrotimeFunction struct {
	data.Function
}

func (f *MicrotimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	now := time.Now()

	// PHP microtime(false): sprintf("%.8F %ld", usec/1e6, sec)  → "0.12345600 1699999999"
	// 旧实现写成 "%d %d"（微秒整数 秒），CombGenerator 的 substr($time[0], 2, 5)
	// 会拿到空串，同一秒内时间戳部分完全相同。
	getAsFloat := false
	temp, ok := ctx.GetIndexValue(0)
	if ok {
		if ab, ok := temp.(data.AsBool); ok {
			getAsFloat, _ = ab.AsBool()
		}
	}
	if getAsFloat {
		microseconds := float64(now.Unix()) + float64(now.Nanosecond())/1e9
		return data.NewFloatValue(microseconds), nil
	}
	sec := now.Unix()
	frac := float64(now.Nanosecond()) / 1e9
	result := fmt.Sprintf("%.8f %d", frac, sec)
	return data.NewStringValue(result), nil
}

func (f *MicrotimeFunction) GetName() string {
	return "microtime"
}

func (f *MicrotimeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "get_as_float", 0, data.NewBoolValue(false), nil),
	}
}

func (f *MicrotimeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "get_as_float", 0, nil),
	}
}

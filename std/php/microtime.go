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

	// 默认返回字符串格式
	getAsFloat := false
	temp, ok := ctx.GetIndexValue(0)
	if ok {
		getAsFloat, _ = temp.(data.AsBool).AsBool()
	}
	if getAsFloat {
		// 返回浮点数格式（微秒精度）
		now := time.Now()
		// 秒部分 + 微秒部分（转换为秒）
		microseconds := float64(now.Unix()) + float64(now.Nanosecond())/1e9
		return data.NewFloatValue(microseconds), nil
	} else {
		// 返回字符串格式 "微秒 秒"
		seconds := now.Unix()
		microseconds := now.Nanosecond() / 1000
		result := fmt.Sprintf("%d %d", microseconds, seconds)
		return data.NewStringValue(result), nil
	}
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

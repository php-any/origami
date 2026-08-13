package php

import (
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GettimeofdayFunction 实现 PHP 内置函数 gettimeofday
//
//	gettimeofday(?bool $as_float = false): array|float
//
// 默认返回含 sec/usec/minuteswest/dsttime 的关联数组；as_float=true 时返回浮点秒。
func NewGettimeofdayFunction() data.FuncStmt {
	return &GettimeofdayFunction{}
}

type GettimeofdayFunction struct{}

func (f *GettimeofdayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	now := time.Now()
	asFloat := false
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			asFloat, _ = b.AsBool()
		}
	}
	if asFloat {
		return data.NewFloatValue(float64(now.Unix()) + float64(now.Nanosecond())/1e9), nil
	}
	_, offset := now.Zone()
	return &data.ArrayValue{
		List: []*data.ZVal{
			data.NewNamedZVal("sec", data.NewIntValue(int(now.Unix()))),
			data.NewNamedZVal("usec", data.NewIntValue(now.Nanosecond()/1000)),
			data.NewNamedZVal("minuteswest", data.NewIntValue(-offset/60)),
			data.NewNamedZVal("dsttime", data.NewIntValue(0)),
		},
	}, nil
}

func (f *GettimeofdayFunction) GetName() string { return "gettimeofday" }
func (f *GettimeofdayFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "as_float", 0, data.NewBoolValue(false), nil),
	}
}
func (f *GettimeofdayFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "as_float", 0, nil),
	}
}

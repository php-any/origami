package php

import (
	"time"

	"github.com/php-any/origami/data"
)

// NewDateDefaultTimezoneSetFunction 创建 date_default_timezone_set 函数。
// PHP 语义（简化版）：
//
//	date_default_timezone_set(timezone: string): bool
//
// 在 Origami 中：
// - 尝试加载指定的时区；
// - 成功则设置为 Go 运行时的 time.Local，返回 true；
// - 失败则返回 false。
func NewDateDefaultTimezoneSetFunction() data.FuncStmt {
	return &DateDefaultTimezoneSetFunction{}
}

type DateDefaultTimezoneSetFunction struct {
	data.Function
}

func (f *DateDefaultTimezoneSetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	tzVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	// 转换为字符串
	var timezone string
	switch v := tzVal.(type) {
	case *data.StringValue:
		timezone = v.Value
	case *data.NullValue:
		return data.NewBoolValue(false), nil
	default:
		timezone = v.AsString()
	}

	// 尝试加载时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// 设置 Go 运行时的本地时区
	time.Local = loc
	return data.NewBoolValue(true), nil
}

func (f *DateDefaultTimezoneSetFunction) GetName() string {
	return "date_default_timezone_set"
}

func (f *DateDefaultTimezoneSetFunction) GetParams() []data.GetValue {
	return nil
}

func (f *DateDefaultTimezoneSetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("timezone", 0, data.NewBaseType("string")),
	}
}

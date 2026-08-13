package php

import (
	"time"

	"github.com/php-any/origami/data"
)

// NewDateDefaultTimezoneGetFunction 创建 date_default_timezone_get 函数。
// PHP 语义（简化版）：
//
//	date_default_timezone_get(): string
//
// 在 Origami 中：
// - 优先从 Go 运行时的 time.Local 取 IANA 名称；
// - 若不可用，则退回 "UTC"。
func NewDateDefaultTimezoneGetFunction() data.FuncStmt {
	return &DateDefaultTimezoneGetFunction{}
}

type DateDefaultTimezoneGetFunction struct {
	data.Function
}

func (f *DateDefaultTimezoneGetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	loc := time.Local
	if loc == nil {
		return data.NewStringValue("UTC"), nil
	}
	name := loc.String()
	if name == "" || name == "Local" {
		// 无法获取具体本地时区名时，保持行为稳定，统一返回 UTC
		return data.NewStringValue("UTC"), nil
	}
	return data.NewStringValue(name), nil
}

func (f *DateDefaultTimezoneGetFunction) GetName() string {
	return "date_default_timezone_get"
}

func (f *DateDefaultTimezoneGetFunction) GetParams() []data.GetValue {
	// 无参数
	return nil
}

func (f *DateDefaultTimezoneGetFunction) GetVariables() []data.Variable {
	return nil
}

package php

import (
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewTimezoneOpenFunction 创建 timezone_open 函数
// PHP 语义：
//
//	timezone_open(string $timezoneId): DateTimeZone|false
//
// 在 Origami 中我们完整返回一个 DateTimeZone 对象：
//   - 若 time.LoadLocation($timezoneId) 成功，则构造一个 DateTimeZone 实例，并将时区 ID
//     挂载到实例属性 "name" 上（供 timezone_name_get / instanceof 使用）
//   - 若失败，则返回 false
func NewTimezoneOpenFunction() data.FuncStmt {
	return &TimezoneOpenFunction{}
}

type TimezoneOpenFunction struct {
	data.Function
}

func (f *TimezoneOpenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 读取第一个参数 $timezone
	v, ctl := ctx.GetVariableValue(node.NewVariable(nil, "timezone", 0, data.String{}))
	if ctl != nil {
		return data.NewBoolValue(false), ctl
	}
	tzID := v.(data.AsString).AsString()

	if tzID == "" {
		return data.NewBoolValue(false), nil
	}

	// 使用 Go 标准库 time.LoadLocation 检查时区是否合法
	if _, err := time.LoadLocation(tzID); err != nil {
		return data.NewBoolValue(false), nil
	}

	// 查找已注册的 DateTimeZone 类（由 std/system.Load 注册）
	vm := ctx.GetVM()
	classStmt, ok := vm.GetClass("DateTimeZone")
	if !ok || classStmt == nil {
		// 未找到类定义时按 PHP 语义返回 false，而不是抛错
		return data.NewBoolValue(false), nil
	}

	// 创建 DateTimeZone 实例，并在实例属性上挂载时区 ID
	tz := data.NewClassValue(classStmt, ctx.CreateBaseContext())
	_ = tz.SetProperty("name", data.NewStringValue(tzID))

	return tz, nil
}

func (f *TimezoneOpenFunction) GetName() string {
	return "timezone_open"
}

func (f *TimezoneOpenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "timezone", 0, nil, data.String{}),
	}
}

func (f *TimezoneOpenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "timezone", 0, data.String{}),
	}
}

package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewTimezoneNameGetFunction 创建 timezone_name_get 函数
// PHP 语义（子集）：
//
//	timezone_name_get(DateTimeZone $object): string
//
// 在 Origami 中：
// - 接受由 timezone_open 返回的 DateTimeZone 实例
// - 从实例属性 "name" 中读取时区 ID 并返回
// - 若不是 DateTimeZone 实例或缺少 name 属性，则返回 false
func NewTimezoneNameGetFunction() data.FuncStmt {
	return &TimezoneNameGetFunction{}
}

type TimezoneNameGetFunction struct {
	data.Function
}

func (f *TimezoneNameGetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 直接按位置参数读取，避免与变量表绑定不一致
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}

	cv, ok := value.(*data.ClassValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	// 确保是 DateTimeZone 实例
	if cv.Class == nil || cv.Class.GetName() != "DateTimeZone" {
		return data.NewBoolValue(false), nil
	}

	// 由于 DateTimeZoneClass 不声明固定属性，时区 ID 由 timezone_open 通过 SetProperty 注入，
	// 实际上是存放在内部的 ObjectValue.property 中。优先尝试通过 GetVariableValue 获取同名变量。
	nameVar := node.NewVariable(nil, "name", 0, nil)
	if nameVal, acl := cv.GetVariableValue(nameVar); acl == nil && nameVal != nil {
		return data.NewStringValue(nameVal.AsString()), nil
	}

	// 回退：直接从 ObjectValue 的属性表中读（与 serialize 中的行为保持一致）
	if ov := cv.ObjectValue; ov != nil {
		if raw, _ := ov.GetProperty("name"); raw != nil {
			return data.NewStringValue(raw.AsString()), nil
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *TimezoneNameGetFunction) GetName() string {
	return "timezone_name_get"
}

func (f *TimezoneNameGetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "timezone", 0, nil, nil),
	}
}

func (f *TimezoneNameGetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "timezone", 0, nil),
	}
}

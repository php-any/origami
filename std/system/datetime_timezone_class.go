package system

import (
	"errors"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DateTimeZoneClass 提供 PHP 内置类 DateTimeZone 的最小运行时模型：
// - 类名为 "DateTimeZone"（全局命名空间）
// - 支持 __construct(string $timezoneId)，内部使用 time.LoadLocation 校验并保存时区 ID
// - 提供 getName(): string，返回构造时或 timezone_open 注入的时区 ID
// - 主要用于支持 CarbonTimeZone 之类对 DateTimeZone 的继承与调用
type DateTimeZoneClass struct {
	node.Node
	construct data.Method
	getName   data.Method
}

func (c *DateTimeZoneClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 运行时实例：属性和值挂载在 ClassValue.ObjectValue 上
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *DateTimeZoneClass) GetName() string {
	return "DateTimeZone"
}

func (c *DateTimeZoneClass) GetExtend() *string {
	return nil
}

func (c *DateTimeZoneClass) GetImplements() []string {
	return nil
}

func (c *DateTimeZoneClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *DateTimeZoneClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *DateTimeZoneClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		if c.construct == nil {
			c.construct = &DateTimeZoneConstructMethod{}
		}
		return c.construct, true
	case "getName":
		if c.getName == nil {
			c.getName = &DateTimeZoneGetNameMethod{}
		}
		return c.getName, true
	}
	return nil, false
}

func (c *DateTimeZoneClass) GetMethods() []data.Method {
	// 懒加载方法集合
	methods := []data.Method{}
	if m, ok := c.GetMethod("__construct"); ok {
		methods = append(methods, m)
	}
	if m, ok := c.GetMethod("getName"); ok {
		methods = append(methods, m)
	}
	return methods
}

func (c *DateTimeZoneClass) GetConstruct() data.Method {
	if m, ok := c.GetMethod("__construct"); ok {
		return m
	}
	return nil
}

// DateTimeZoneConstructMethod 实现 DateTimeZone::__construct(string $timezoneId)
type DateTimeZoneConstructMethod struct{}

func (m *DateTimeZoneConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 位置 0：$timezone
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return nil, data.NewErrorThrow(nil, errors.New("Invalid timezone (empty)"))
	}
	tzID := v.AsString()
	if tzID == "" {
		return nil, data.NewErrorThrow(nil, errors.New("Invalid timezone (empty)"))
	}

	// 校验时区合法性，保持与 timezone_open 一致
	if _, err := time.LoadLocation(tzID); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	// 将时区 ID 挂在当前实例的 "name" 属性上
	if thisCtx, ok := ctx.(*data.ClassMethodContext); ok && thisCtx.ClassValue != nil {
		_ = thisCtx.SetProperty("name", data.NewStringValue(tzID))
	}
	return nil, nil
}

func (m *DateTimeZoneConstructMethod) GetName() string {
	return "__construct"
}

func (m *DateTimeZoneConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *DateTimeZoneConstructMethod) GetIsStatic() bool {
	return false
}

func (m *DateTimeZoneConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "timezone", 0, nil, data.String{}),
	}
}

func (m *DateTimeZoneConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "timezone", 0, data.String{}),
	}
}

func (m *DateTimeZoneConstructMethod) GetReturnType() data.Types {
	return nil
}

// DateTimeZoneGetNameMethod 实现 DateTimeZone::getName(): string
type DateTimeZoneGetNameMethod struct{}

func (m *DateTimeZoneGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if thisCtx, ok := ctx.(*data.ClassMethodContext); ok && thisCtx.ClassValue != nil {
		if nameVal, _ := thisCtx.GetProperty("name"); nameVal != nil {
			return data.NewStringValue(nameVal.AsString()), nil
		}
	}
	return data.NewStringValue(""), nil
}

func (m *DateTimeZoneGetNameMethod) GetName() string {
	return "getName"
}

func (m *DateTimeZoneGetNameMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *DateTimeZoneGetNameMethod) GetIsStatic() bool {
	return false
}

func (m *DateTimeZoneGetNameMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *DateTimeZoneGetNameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *DateTimeZoneGetNameMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

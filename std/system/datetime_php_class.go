package system

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PhpDateTimeClass 提供 PHP 全局类 DateTime 的最小实现。
// 目前仅支持：
// - format(string $format): string
// - getTimestamp(): int
// - setTimezone(DateTimeZone $timezone): static
// 并实现 DateTimeInterface，便于与现有的 System\DateTime 一起通过 instanceof / 类型提示检查。
type PhpDateTimeClass struct {
	node.Node
	format       data.Method
	getTimestamp data.Method
	setTimezone  data.Method
}

func (s *PhpDateTimeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newDateTime()

	return data.NewClassValue(&PhpDateTimeClass{
		format:       &DateTimeFormatMethod{source},
		getTimestamp: &DateTimeGetTimestampMethod{source},
	}, ctx.CreateBaseContext()), nil
}

// GetName 返回全局类名 DateTime（与 PHP 内置一致）
func (s *PhpDateTimeClass) GetName() string {
	return "DateTime"
}

func (s *PhpDateTimeClass) GetExtend() *string {
	return nil
}

func (s *PhpDateTimeClass) GetImplements() []string {
	// 与 System\DateTime 一样，实现 PHP 顶层接口 DateTimeInterface
	return []string{"DateTimeInterface"}
}

func (s *PhpDateTimeClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *PhpDateTimeClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *PhpDateTimeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "format":
		return s.format, true
	case "getTimestamp":
		return s.getTimestamp, true
	case "setTimezone":
		if s.setTimezone == nil {
			s.setTimezone = &DateTimeSetTimezoneMethod{}
		}
		return s.setTimezone, true
	}
	return nil, false
}

func (s *PhpDateTimeClass) GetMethods() []data.Method {
	return []data.Method{
		s.format,
		s.getTimestamp,
	}
}

func (t *PhpDateTimeClass) GetConstruct() data.Method {
	// 当前先不实现复杂构造逻辑（如解析字符串），行为与 System\DateTime 一致，
	// 调用 new DateTime() 后的 format/getTimestamp 使用当前时间。
	return nil
}

// DateTimeSetTimezoneMethod 实现 DateTime::setTimezone(DateTimeZone $timezone): static
// 简单实现：返回 $this 自身（忽略时区切换），避免 parent::setTimezone 调用崩溃
type DateTimeSetTimezoneMethod struct{}

func (m *DateTimeSetTimezoneMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 返回 $this 对象本身（链式调用语义）
	if thisCtx, ok := ctx.(*data.ClassMethodContext); ok && thisCtx.ClassValue != nil {
		return thisCtx.ClassValue, nil
	}
	return data.NewNullValue(), nil
}

func (m *DateTimeSetTimezoneMethod) GetName() string {
	return "setTimezone"
}

func (m *DateTimeSetTimezoneMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *DateTimeSetTimezoneMethod) GetIsStatic() bool {
	return false
}

func (m *DateTimeSetTimezoneMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "timezone", 0, nil, nil),
	}
}

func (m *DateTimeSetTimezoneMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "timezone", 0, nil),
	}
}

func (m *DateTimeSetTimezoneMethod) GetReturnType() data.Types {
	return nil
}

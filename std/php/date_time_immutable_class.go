package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DateTimeImmutableClass 实现 PHP DateTimeImmutable（复用 DateTime 的时间存储与方法）。
// 需实现 DateTimeInterface::{format,getTimestamp}，否则子类（如 Monolog\JsonSerializableDateTimeImmutable）
// 会在抽象方法校验时误报未实现。
type DateTimeImmutableClass struct {
	node.Node
}

func NewDateTimeImmutableClass() *DateTimeImmutableClass {
	return &DateTimeImmutableClass{}
}

func (c *DateTimeImmutableClass) GetName() string                               { return "DateTimeImmutable" }
func (c *DateTimeImmutableClass) GetExtend() *string                            { return nil }
func (c *DateTimeImmutableClass) GetImplements() []string                       { return []string{"DateTimeInterface"} }
func (c *DateTimeImmutableClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DateTimeImmutableClass) GetPropertyList() []data.Property              { return nil }
func (c *DateTimeImmutableClass) GetConstruct() data.Method {
	return &DateTimeConstructMethod{}
}

func (c *DateTimeImmutableClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "createFromFormat":
		return &DateTimeCreateFromFormatMethod{}, true
	case "getLastErrors":
		return &DateTimeGetLastErrorsMethod{}, true
	}
	return nil, false
}

func (c *DateTimeImmutableClass) GetMethods() []data.Method {
	return []data.Method{
		&DateTimeConstructMethod{},
		&DateTimeGetTimestampMethod{},
		&DateTimeFormatMethod{},
		&DateTimeModifyMethod{},
		&DateTimeSetTimestampMethod{},
		&DateTimeSetTimezoneMethod{},
		&DateTimeGetTimezoneMethod{},
		&DateTimeAddMethod{},
		&DateTimeSubMethod{},
		&DateTimeDiffMethod{},
		&DateTimeToStringMethod{},
	}
}

func (c *DateTimeImmutableClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &DateTimeConstructMethod{}, true
	case "getTimestamp":
		return &DateTimeGetTimestampMethod{}, true
	case "setTimestamp":
		return &DateTimeSetTimestampMethod{}, true
	case "setTimezone":
		return &DateTimeSetTimezoneMethod{}, true
	case "getTimezone":
		return &DateTimeGetTimezoneMethod{}, true
	case "format":
		return &DateTimeFormatMethod{}, true
	case "add":
		return &DateTimeAddMethod{}, true
	case "sub":
		return &DateTimeSubMethod{}, true
	case "diff":
		return &DateTimeDiffMethod{}, true
	case "modify":
		return &DateTimeModifyMethod{}, true
	case "createFromFormat":
		return &DateTimeCreateFromFormatMethod{}, true
	case "__toString":
		return &DateTimeToStringMethod{}, true
	}
	return nil, false
}

func (c *DateTimeImmutableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

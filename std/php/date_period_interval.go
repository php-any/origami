package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DatePeriodClass 最小存根，满足 CarbonPeriod 对原生 DatePeriod 的继承需求。
type DatePeriodClass struct {
	node.Node
}

func NewDatePeriodClass() *DatePeriodClass { return &DatePeriodClass{} }

func (c *DatePeriodClass) GetName() string                               { return "DatePeriod" }
func (c *DatePeriodClass) GetExtend() *string                            { return nil }
func (c *DatePeriodClass) GetImplements() []string                       { return []string{"IteratorAggregate"} }
func (c *DatePeriodClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DatePeriodClass) GetPropertyList() []data.Property              { return nil }
func (c *DatePeriodClass) GetConstruct() data.Method                     { return &DatePeriodConstructMethod{} }
func (c *DatePeriodClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DatePeriodClass) GetMethods() []data.Method {
	return []data.Method{&DatePeriodConstructMethod{}}
}
func (c *DatePeriodClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return &DatePeriodConstructMethod{}, true
	}
	return nil, false
}
func (c *DatePeriodClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

type DatePeriodConstructMethod struct{}

func (m *DatePeriodConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
func (m *DatePeriodConstructMethod) GetName() string            { return "__construct" }
func (m *DatePeriodConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DatePeriodConstructMethod) GetIsStatic() bool          { return false }
func (m *DatePeriodConstructMethod) GetReturnType() data.Types  { return nil }
func (m *DatePeriodConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "start", 0, nil, nil),
		node.NewParameter(nil, "interval", 1, nil, nil),
		node.NewParameter(nil, "end", 2, nil, nil),
		node.NewParameter(nil, "options", 3, data.NewIntValue(0), nil),
	}
}
func (m *DatePeriodConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "start", 0, nil),
		node.NewVariable(nil, "interval", 1, nil),
		node.NewVariable(nil, "end", 2, nil),
		node.NewVariable(nil, "options", 3, nil),
	}
}

// DateIntervalClass 最小存根。
type DateIntervalClass struct {
	node.Node
}

func NewDateIntervalClass() *DateIntervalClass { return &DateIntervalClass{} }

func (c *DateIntervalClass) GetName() string                               { return "DateInterval" }
func (c *DateIntervalClass) GetExtend() *string                            { return nil }
func (c *DateIntervalClass) GetImplements() []string                       { return nil }
func (c *DateIntervalClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DateIntervalClass) GetPropertyList() []data.Property              { return nil }
func (c *DateIntervalClass) GetConstruct() data.Method                     { return &DateIntervalConstructMethod{} }
func (c *DateIntervalClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DateIntervalClass) GetMethods() []data.Method {
	return []data.Method{&DateIntervalConstructMethod{}}
}
func (c *DateIntervalClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return &DateIntervalConstructMethod{}, true
	}
	return nil, false
}
func (c *DateIntervalClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

type DateIntervalConstructMethod struct{}

func (m *DateIntervalConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
func (m *DateIntervalConstructMethod) GetName() string            { return "__construct" }
func (m *DateIntervalConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *DateIntervalConstructMethod) GetIsStatic() bool          { return false }
func (m *DateIntervalConstructMethod) GetReturnType() data.Types  { return nil }
func (m *DateIntervalConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "duration", 0, nil, data.String{})}
}
func (m *DateIntervalConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "duration", 0, data.String{})}
}

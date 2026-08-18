package php

import (
	"github.com/php-any/origami/data"
)

// GcEnabledFunction 对应 PHP gc_enabled()：origami 基于 Go runtime，GC 始终启用。
type GcEnabledFunction struct {
	data.Function
}

func NewGcEnabledFunction() data.FuncStmt {
	return &GcEnabledFunction{}
}

func (f *GcEnabledFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}

func (f *GcEnabledFunction) GetName() string {
	return "gc_enabled"
}

func (f *GcEnabledFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *GcEnabledFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GcEnableFunction 对应 PHP gc_enable()：Go runtime 无法禁用 GC，无实际作用。
type GcEnableFunction struct {
	data.Function
}

func NewGcEnableFunction() data.FuncStmt {
	return &GcEnableFunction{}
}

func (f *GcEnableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func (f *GcEnableFunction) GetName() string {
	return "gc_enable"
}

func (f *GcEnableFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *GcEnableFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GcDisableFunction 对应 PHP gc_disable()：Go runtime 无法禁用 GC，无实际作用。
type GcDisableFunction struct {
	data.Function
}

func NewGcDisableFunction() data.FuncStmt {
	return &GcDisableFunction{}
}

func (f *GcDisableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func (f *GcDisableFunction) GetName() string {
	return "gc_disable"
}

func (f *GcDisableFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *GcDisableFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

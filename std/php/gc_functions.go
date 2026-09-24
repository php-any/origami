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

var gcEnabledFunctionGetParams = []data.GetValue{}

func (f *GcEnabledFunction) GetParams() []data.GetValue {
	return gcEnabledFunctionGetParams
}

var gcEnabledFunctionGetVariables = []data.Variable{}

func (f *GcEnabledFunction) GetVariables() []data.Variable {
	return gcEnabledFunctionGetVariables
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

var gcEnableFunctionGetParams = []data.GetValue{}

func (f *GcEnableFunction) GetParams() []data.GetValue {
	return gcEnableFunctionGetParams
}

var gcEnableFunctionGetVariables = []data.Variable{}

func (f *GcEnableFunction) GetVariables() []data.Variable {
	return gcEnableFunctionGetVariables
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

var gcDisableFunctionGetParams = []data.GetValue{}

func (f *GcDisableFunction) GetParams() []data.GetValue {
	return gcDisableFunctionGetParams
}

var gcDisableFunctionGetVariables = []data.Variable{}

func (f *GcDisableFunction) GetVariables() []data.Variable {
	return gcDisableFunctionGetVariables
}

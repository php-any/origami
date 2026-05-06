package php

import (
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CtypeSpaceFunction 实现 ctype_space 函数
type CtypeSpaceFunction struct{}

func NewCtypeSpaceFunction() data.FuncStmt { return &CtypeSpaceFunction{} }

func (f *CtypeSpaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeSpaceFunction) GetName() string { return "ctype_space" }
func (f *CtypeSpaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, nil),
	}
}
func (f *CtypeSpaceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
	}
}

// CtypeDigitFunction 实现 ctype_digit 函数
type CtypeDigitFunction struct{}

func NewCtypeDigitFunction() data.FuncStmt { return &CtypeDigitFunction{} }

func (f *CtypeDigitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeDigitFunction) GetName() string { return "ctype_digit" }
func (f *CtypeDigitFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, nil),
	}
}
func (f *CtypeDigitFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
	}
}

// CtypeAlphaFunction 实现 ctype_alpha 函数
type CtypeAlphaFunction struct{}

func NewCtypeAlphaFunction() data.FuncStmt { return &CtypeAlphaFunction{} }

func (f *CtypeAlphaFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeAlphaFunction) GetName() string { return "ctype_alpha" }
func (f *CtypeAlphaFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, nil),
	}
}
func (f *CtypeAlphaFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
	}
}

// CtypeAlnumFunction 实现 ctype_alnum 函数
type CtypeAlnumFunction struct{}

func NewCtypeAlnumFunction() data.FuncStmt { return &CtypeAlnumFunction{} }

func (f *CtypeAlnumFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeAlnumFunction) GetName() string { return "ctype_alnum" }
func (f *CtypeAlnumFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, nil),
	}
}
func (f *CtypeAlnumFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
	}
}

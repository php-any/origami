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

// CtypeLowerFunction 实现 ctype_lower 函数
type CtypeLowerFunction struct{}

func NewCtypeLowerFunction() data.FuncStmt { return &CtypeLowerFunction{} }

func (f *CtypeLowerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !unicode.IsLower(r) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeLowerFunction) GetName() string { return "ctype_lower" }
func (f *CtypeLowerFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}
}
func (f *CtypeLowerFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}
}

// CtypeUpperFunction 实现 ctype_upper 函数
type CtypeUpperFunction struct{}

func NewCtypeUpperFunction() data.FuncStmt { return &CtypeUpperFunction{} }

func (f *CtypeUpperFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeUpperFunction) GetName() string { return "ctype_upper" }
func (f *CtypeUpperFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}
}
func (f *CtypeUpperFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}
}

// CtypePrintFunction 实现 ctype_print 函数
type CtypePrintFunction struct{}

func NewCtypePrintFunction() data.FuncStmt { return &CtypePrintFunction{} }

func (f *CtypePrintFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if r < 0x20 || r > 0x7e {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypePrintFunction) GetName() string { return "ctype_print" }
func (f *CtypePrintFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}
}
func (f *CtypePrintFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}
}

// CtypePunctFunction 实现 ctype_punct 函数
type CtypePunctFunction struct{}

func NewCtypePunctFunction() data.FuncStmt { return &CtypePunctFunction{} }

func (f *CtypePunctFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !unicode.IsPunct(r) && !unicode.IsSymbol(r) {
			return data.NewBoolValue(false), nil
		}
		if r > 0x7e {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypePunctFunction) GetName() string { return "ctype_punct" }
func (f *CtypePunctFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}
}
func (f *CtypePunctFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}
}

// CtypeXdigitFunction 实现 ctype_xdigit 函数
type CtypeXdigitFunction struct{}

func NewCtypeXdigitFunction() data.FuncStmt { return &CtypeXdigitFunction{} }

func (f *CtypeXdigitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeXdigitFunction) GetName() string { return "ctype_xdigit" }
func (f *CtypeXdigitFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}
}
func (f *CtypeXdigitFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}
}

// CtypeGraphFunction 实现 ctype_graph 函数
type CtypeGraphFunction struct{}

func NewCtypeGraphFunction() data.FuncStmt { return &CtypeGraphFunction{} }

func (f *CtypeGraphFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	if s == "" {
		return data.NewBoolValue(false), nil
	}
	for _, r := range s {
		if r < 0x21 || r > 0x7e {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *CtypeGraphFunction) GetName() string { return "ctype_graph" }
func (f *CtypeGraphFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}
}
func (f *CtypeGraphFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}
}

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
var ctypeSpaceFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "text", 0, nil, nil),
}

func (f *CtypeSpaceFunction) GetParams() []data.GetValue {
	return ctypeSpaceFunctionGetParams
}
var ctypeSpaceFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
}

func (f *CtypeSpaceFunction) GetVariables() []data.Variable {
	return ctypeSpaceFunctionGetVariables
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
var ctypeDigitFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "text", 0, nil, nil),
}

func (f *CtypeDigitFunction) GetParams() []data.GetValue {
	return ctypeDigitFunctionGetParams
}
var ctypeDigitFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
}

func (f *CtypeDigitFunction) GetVariables() []data.Variable {
	return ctypeDigitFunctionGetVariables
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
var ctypeAlphaFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "text", 0, nil, nil),
}

func (f *CtypeAlphaFunction) GetParams() []data.GetValue {
	return ctypeAlphaFunctionGetParams
}
var ctypeAlphaFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
}

func (f *CtypeAlphaFunction) GetVariables() []data.Variable {
	return ctypeAlphaFunctionGetVariables
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
var ctypeAlnumFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "text", 0, nil, nil),
}

func (f *CtypeAlnumFunction) GetParams() []data.GetValue {
	return ctypeAlnumFunctionGetParams
}
var ctypeAlnumFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "text", 0, data.NewBaseType("mixed")),
}

func (f *CtypeAlnumFunction) GetVariables() []data.Variable {
	return ctypeAlnumFunctionGetVariables
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
var ctypeLowerFunctionGetParams = []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}

func (f *CtypeLowerFunction) GetParams() []data.GetValue {
	return ctypeLowerFunctionGetParams
}
var ctypeLowerFunctionGetVariables = []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}

func (f *CtypeLowerFunction) GetVariables() []data.Variable {
	return ctypeLowerFunctionGetVariables
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
var ctypeUpperFunctionGetParams = []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}

func (f *CtypeUpperFunction) GetParams() []data.GetValue {
	return ctypeUpperFunctionGetParams
}
var ctypeUpperFunctionGetVariables = []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}

func (f *CtypeUpperFunction) GetVariables() []data.Variable {
	return ctypeUpperFunctionGetVariables
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
var ctypePrintFunctionGetParams = []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}

func (f *CtypePrintFunction) GetParams() []data.GetValue {
	return ctypePrintFunctionGetParams
}
var ctypePrintFunctionGetVariables = []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}

func (f *CtypePrintFunction) GetVariables() []data.Variable {
	return ctypePrintFunctionGetVariables
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
var ctypePunctFunctionGetParams = []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}

func (f *CtypePunctFunction) GetParams() []data.GetValue {
	return ctypePunctFunctionGetParams
}
var ctypePunctFunctionGetVariables = []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}

func (f *CtypePunctFunction) GetVariables() []data.Variable {
	return ctypePunctFunctionGetVariables
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
var ctypeXdigitFunctionGetParams = []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}

func (f *CtypeXdigitFunction) GetParams() []data.GetValue {
	return ctypeXdigitFunctionGetParams
}
var ctypeXdigitFunctionGetVariables = []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}

func (f *CtypeXdigitFunction) GetVariables() []data.Variable {
	return ctypeXdigitFunctionGetVariables
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
var ctypeGraphFunctionGetParams = []data.GetValue{node.NewParameter(nil, "text", 0, nil, nil)}

func (f *CtypeGraphFunction) GetParams() []data.GetValue {
	return ctypeGraphFunctionGetParams
}
var ctypeGraphFunctionGetVariables = []data.Variable{node.NewVariable(nil, "text", 0, data.NewBaseType("mixed"))}

func (f *CtypeGraphFunction) GetVariables() []data.Variable {
	return ctypeGraphFunctionGetVariables
}

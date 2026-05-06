package php

import (
	"strings"
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MbStrtoupperFunction 实现 mb_strtoupper 函数
type MbStrtoupperFunction struct{}

func NewMbStrtoupperFunction() data.FuncStmt { return &MbStrtoupperFunction{} }

func (f *MbStrtoupperFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(strings.ToUpper(v.AsString())), nil
}

func (f *MbStrtoupperFunction) GetName() string { return "mb_strtoupper" }
func (f *MbStrtoupperFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "encoding", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *MbStrtoupperFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "encoding", 1, data.NewNullableType(data.NewBaseType("string"))),
	}
}

// MbStrtolowerFunction 实现 mb_strtolower 函数
type MbStrtolowerFunction struct{}

func NewMbStrtolowerFunction() data.FuncStmt { return &MbStrtolowerFunction{} }

func (f *MbStrtolowerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(strings.ToLower(v.AsString())), nil
}

func (f *MbStrtolowerFunction) GetName() string { return "mb_strtolower" }
func (f *MbStrtolowerFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "encoding", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *MbStrtolowerFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "encoding", 1, data.NewNullableType(data.NewBaseType("string"))),
	}
}

// MbStrlenFunction 实现 mb_strlen 函数
type MbStrlenFunction struct{}

func NewMbStrlenFunction() data.FuncStmt { return &MbStrlenFunction{} }

func (f *MbStrlenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	s := v.AsString()
	return data.NewIntValue(utf8.RuneCountInString(s)), nil
}

func (f *MbStrlenFunction) GetName() string { return "mb_strlen" }
func (f *MbStrlenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "encoding", 1, node.NewNullLiteral(nil), nil),
	}
}
func (f *MbStrlenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "encoding", 1, data.NewNullableType(data.NewBaseType("string"))),
	}
}

// MbStrposFunction 实现 mb_strpos 函数
type MbStrposFunction struct{}

func NewMbStrposFunction() data.FuncStmt { return &MbStrposFunction{} }

func (f *MbStrposFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystack, _ := ctx.GetIndexValue(0)
	needle, _ := ctx.GetIndexValue(1)
	if haystack == nil || needle == nil {
		return data.NewBoolValue(false), nil
	}
	h := haystack.AsString()
	n := needle.AsString()
	if n == "" {
		return data.NewBoolValue(false), nil
	}

	offset := 0
	if off, ok := ctx.GetIndexValue(2); ok && off != nil {
		if asInt, ok := off.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				offset = v
			}
		}
	}

	// Convert rune offset to byte offset
	if offset > 0 {
		runeIdx := 0
		byteIdx := 0
		for byteIdx < len(h) && runeIdx < offset {
			_, size := utf8.DecodeRuneInString(h[byteIdx:])
			byteIdx += size
			runeIdx++
		}
		if byteIdx >= len(h) {
			return data.NewBoolValue(false), nil
		}
		idx := strings.Index(h[byteIdx:], n)
		if idx < 0 {
			return data.NewBoolValue(false), nil
		}
		return data.NewIntValue(utf8.RuneCountInString(h[:byteIdx+idx])), nil
	}

	if offset < 0 {
		offset = utf8.RuneCountInString(h) + offset
		if offset < 0 {
			offset = 0
		}
		runeIdx := 0
		byteIdx := 0
		for byteIdx < len(h) && runeIdx < offset {
			_, size := utf8.DecodeRuneInString(h[byteIdx:])
			byteIdx += size
			runeIdx++
		}
		idx := strings.Index(h[byteIdx:], n)
		if idx < 0 {
			return data.NewBoolValue(false), nil
		}
		return data.NewIntValue(utf8.RuneCountInString(h[:byteIdx+idx])), nil
	}

	idx := strings.Index(h, n)
	if idx < 0 {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(utf8.RuneCountInString(h[:idx])), nil
}

func (f *MbStrposFunction) GetName() string { return "mb_strpos" }
func (f *MbStrposFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, nil),
		node.NewParameter(nil, "needle", 1, nil, nil),
		node.NewParameter(nil, "offset", 2, data.NewIntValue(0), nil),
		node.NewParameter(nil, "encoding", 3, node.NewNullLiteral(nil), nil),
	}
}
func (f *MbStrposFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "encoding", 3, data.NewNullableType(data.NewBaseType("string"))),
	}
}

// MbSubstrFunction 实现 mb_substr 函数
type MbSubstrFunction struct{}

func NewMbSubstrFunction() data.FuncStmt { return &MbSubstrFunction{} }

func runeIndexToByte(s string, runeIdx int) int {
	if runeIdx <= 0 {
		return 0
	}
	r := 0
	b := 0
	for b < len(s) && r < runeIdx {
		_, size := utf8.DecodeRuneInString(s[b:])
		b += size
		r++
	}
	return b
}

func (f *MbSubstrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	startVal, _ := ctx.GetIndexValue(1)
	if strVal == nil || startVal == nil {
		return data.NewStringValue(""), nil
	}

	s := strVal.AsString()
	start := 0
	if asInt, ok := startVal.(data.AsInt); ok {
		if v, err := asInt.AsInt(); err == nil {
			start = v
		}
	}
	runeLen := utf8.RuneCountInString(s)

	if start < 0 {
		start = runeLen + start
		if start < 0 {
			start = 0
		}
	}
	if start >= runeLen {
		return data.NewStringValue(""), nil
	}

	length := runeLen
	hasLength := false
	if lenVal, ok := ctx.GetIndexValue(2); ok && lenVal != nil {
		if asInt, ok := lenVal.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				length = v
				hasLength = true
			}
		}
	}

	byteStart := runeIndexToByte(s, start)

	if !hasLength {
		return data.NewStringValue(s[byteStart:]), nil
	}

	if length < 0 {
		end := runeLen + length
		if end <= start {
			return data.NewStringValue(""), nil
		}
		byteEnd := runeIndexToByte(s, end)
		return data.NewStringValue(s[byteStart:byteEnd]), nil
	}

	end := start + length
	if end > runeLen {
		end = runeLen
	}
	byteEnd := runeIndexToByte(s, end)
	return data.NewStringValue(s[byteStart:byteEnd]), nil
}

func (f *MbSubstrFunction) GetName() string { return "mb_substr" }
func (f *MbSubstrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "start", 1, nil, nil),
		node.NewParameter(nil, "length", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "encoding", 3, node.NewNullLiteral(nil), nil),
	}
}
func (f *MbSubstrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "start", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "length", 2, data.NewNullableType(data.NewBaseType("int"))),
		node.NewVariable(nil, "encoding", 3, data.NewNullableType(data.NewBaseType("string"))),
	}
}

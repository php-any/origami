package php

import (
	"strings"

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

// MbStrlenFunction 实现 mb_strlen 函数
type MbStrlenFunction struct{}

func NewMbStrlenFunction() data.FuncStmt { return &MbStrlenFunction{} }

func (f *MbStrlenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}
	s := v.AsString()
	return data.NewIntValue(len(s)), nil
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

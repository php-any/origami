package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrrevFunction 实现 PHP strrev() 的字节反转语义。
type StrrevFunction struct{}

func NewStrrevFunction() data.FuncStmt { return &StrrevFunction{} }

func (f *StrrevFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewStringValue(""), nil
	}
	bytes := []byte(value.AsString())
	for left, right := 0, len(bytes)-1; left < right; left, right = left+1, right-1 {
		bytes[left], bytes[right] = bytes[right], bytes[left]
	}
	return data.NewStringValue(string(bytes)), nil
}

func (f *StrrevFunction) GetName() string { return "strrev" }

func (f *StrrevFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.String{}),
	}
}

func (f *StrrevFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.String{}),
	}
}

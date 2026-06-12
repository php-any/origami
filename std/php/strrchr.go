package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type StrrchrFunction struct{}

func NewStrrchrFunction() data.FuncStmt { return &StrrchrFunction{} }

func (f *StrrchrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystack, _ := ctx.GetIndexValue(0)
	needle, _ := ctx.GetIndexValue(1)

	if haystack == nil || needle == nil {
		return data.NewBoolValue(false), nil
	}

	h := haystack.AsString()
	n := needle.AsString()

	pos := strings.LastIndex(h, n)
	if pos < 0 {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(h[pos:]), nil
}

func (f *StrrchrFunction) GetName() string            { return "strrchr" }
func (f *StrrchrFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *StrrchrFunction) GetIsStatic() bool          { return false }
func (f *StrrchrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, data.String{}),
		node.NewParameter(nil, "needle", 1, nil, data.String{}),
	}
}
func (f *StrrchrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
	}
}
func (f *StrrchrFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrStrFunction 实现 strstr (以及 strchr 别名)
type StrStrFunction struct {
	name string // "strstr" 或 "strchr"
}

func NewStrStrFunction() data.FuncStmt           { return &StrStrFunction{name: "strstr"} }
func NewStrChrFunction() data.FuncStmt           { return &StrStrFunction{name: "strchr"} }

func (f *StrStrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystack, _ := ctx.GetIndexValue(0)
	needle, _ := ctx.GetIndexValue(1)
	beforeVal, _ := ctx.GetIndexValue(2)

	if haystack == nil || needle == nil {
		return data.NewBoolValue(false), nil
	}

	h := haystack.AsString()
	n := needle.AsString()

	beforeNeedle := false
	if beforeVal != nil {
		if b, ok := beforeVal.(*data.BoolValue); ok {
			beforeNeedle = b.Value
		}
	}

	pos := strings.Index(h, n)
	if pos < 0 {
		return data.NewBoolValue(false), nil
	}

	if beforeNeedle {
		return data.NewStringValue(h[:pos]), nil
	}
	return data.NewStringValue(h[pos:]), nil
}

func (f *StrStrFunction) GetName() string { return f.name }
func (f *StrStrFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *StrStrFunction) GetIsStatic() bool { return false }
func (f *StrStrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, data.String{}),
		node.NewParameter(nil, "needle", 1, nil, data.String{}),
		node.NewParameter(nil, "before_needle", 2, data.NewBoolValue(false), data.Bool{}),
	}
}
func (f *StrStrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "before_needle", 2, data.NewBaseType("bool")),
	}
}
func (f *StrStrFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

// StrIStrFunction 实现 stristr (大小写不敏感 strstr)
type StrIStrFunction struct{}

func NewStrIStrFunction() data.FuncStmt { return &StrIStrFunction{} }

func (f *StrIStrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	haystack, _ := ctx.GetIndexValue(0)
	needle, _ := ctx.GetIndexValue(1)
	beforeVal, _ := ctx.GetIndexValue(2)

	if haystack == nil || needle == nil {
		return data.NewBoolValue(false), nil
	}

	h := strings.ToLower(haystack.AsString())
	n := strings.ToLower(needle.AsString())
	origH := haystack.AsString()

	beforeNeedle := false
	if beforeVal != nil {
		if b, ok := beforeVal.(*data.BoolValue); ok {
			beforeNeedle = b.Value
		}
	}

	pos := strings.Index(h, n)
	if pos < 0 {
		return data.NewBoolValue(false), nil
	}

	if beforeNeedle {
		return data.NewStringValue(origH[:pos]), nil
	}
	return data.NewStringValue(origH[pos:]), nil
}

func (f *StrIStrFunction) GetName() string { return "stristr" }
func (f *StrIStrFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *StrIStrFunction) GetIsStatic() bool { return false }
func (f *StrIStrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "haystack", 0, nil, data.String{}),
		node.NewParameter(nil, "needle", 1, nil, data.String{}),
		node.NewParameter(nil, "before_needle", 2, data.NewBoolValue(false), data.Bool{}),
	}
}
func (f *StrIStrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "haystack", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "needle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "before_needle", 2, data.NewBaseType("bool")),
	}
}
func (f *StrIStrFunction) GetReturnType() data.Types { return data.NewBaseType("string") }

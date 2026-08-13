package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrspnFunction() data.FuncStmt {
	return &StrspnFunction{}
}

type StrspnFunction struct{}

func (fn *StrspnFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	maskVal, _ := ctx.GetIndexValue(1)

	str := ""
	if s, ok := strVal.(data.AsString); ok {
		str = s.AsString()
	}
	mask := ""
	if s, ok := maskVal.(data.AsString); ok {
		mask = s.AsString()
	}

	for i, c := range str {
		if !strings.ContainsRune(mask, c) {
			return data.NewIntValue(i), nil
		}
	}
	return data.NewIntValue(len(str)), nil
}

func (fn *StrspnFunction) GetName() string { return "strspn" }
func (fn *StrspnFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str", 0, nil, nil),
		node.NewParameter(nil, "mask", 1, nil, nil),
	}
}
func (fn *StrspnFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "mask", 1, data.NewBaseType("string")),
	}
}

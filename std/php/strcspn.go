package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewStrcspnFunction() data.FuncStmt {
	return &StrcspnFunction{}
}

type StrcspnFunction struct{}

func (fn *StrcspnFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	charListVal, _ := ctx.GetIndexValue(1)

	str := ""
	if s, ok := strVal.(data.AsString); ok {
		str = s.AsString()
	}
	charList := ""
	if s, ok := charListVal.(data.AsString); ok {
		charList = s.AsString()
	}

	// strcspn returns the length of the initial segment not containing any chars from charList
	for i, c := range str {
		if strings.ContainsRune(charList, c) {
			return data.NewIntValue(i), nil
		}
	}
	return data.NewIntValue(len(str)), nil
}

func (fn *StrcspnFunction) GetName() string { return "strcspn" }
func (fn *StrcspnFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str", 0, nil, nil),
		node.NewParameter(nil, "char_list", 1, nil, nil),
	}
}
func (fn *StrcspnFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "char_list", 1, data.NewBaseType("string")),
	}
}

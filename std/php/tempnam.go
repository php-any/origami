package php

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewTempnamFunction() data.FuncStmt {
	return &TempnamFunction{}
}

type TempnamFunction struct{}

func (fn *TempnamFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	dir, _ := ctx.GetIndexValue(0)
	prefix, _ := ctx.GetIndexValue(1)

	dirStr := ""
	if s, ok := dir.(data.AsString); ok {
		dirStr = s.AsString()
	}
	if dirStr == "" {
		dirStr = os.TempDir()
	}

	prefixStr := ""
	if s, ok := prefix.(data.AsString); ok {
		prefixStr = s.AsString()
	}
	if prefixStr == "" {
		prefixStr = "tmp"
	}

	tmpFile, err := os.CreateTemp(dirStr, prefixStr+"_*")
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	name := tmpFile.Name()
	tmpFile.Close()

	return data.NewStringValue(name), nil
}

func (fn *TempnamFunction) GetName() string {
	return "tempnam"
}

func (fn *TempnamFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, nil),
		node.NewParameter(nil, "prefix", 1, nil, nil),
	}
}

func (fn *TempnamFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "prefix", 1, data.NewBaseType("string")),
	}
}

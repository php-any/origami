package php

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewChmodFunction() data.FuncStmt {
	return &ChmodFunction{}
}

type ChmodFunction struct{}

func (fn *ChmodFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathVal, _ := ctx.GetIndexValue(0)
	modeVal, _ := ctx.GetIndexValue(1)

	path := ""
	if s, ok := pathVal.(data.AsString); ok {
		path = s.AsString()
	}

	if path == "" {
		return data.NewBoolValue(false), nil
	}

	mode := os.FileMode(0644)
	if m, ok := modeVal.(data.AsInt); ok {
		if v, err := m.AsInt(); err == nil {
			mode = os.FileMode(v)
		}
	}

	// chmod 在 Windows 上无效，但不报错
	_ = os.Chmod(path, mode)

	return data.NewBoolValue(true), nil
}

func (fn *ChmodFunction) GetName() string {
	return "chmod"
}

func (fn *ChmodFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
		node.NewParameter(nil, "permissions", 1, nil, nil),
	}
}

func (fn *ChmodFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "permissions", 1, data.NewBaseType("int")),
	}
}

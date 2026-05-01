package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewUmaskFunction() data.FuncStmt {
	return &UmaskFunction{}
}

type UmaskFunction struct{}

func (fn *UmaskFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// Windows 不支持 umask，返回 0
	return data.NewIntValue(0), nil
}

func (fn *UmaskFunction) GetName() string {
	return "umask"
}

func (fn *UmaskFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "mask", 0, nil, nil),
	}
}

func (fn *UmaskFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "mask", 0, data.NewBaseType("int")),
	}
}

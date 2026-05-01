package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewGetCfgVarFunction() data.FuncStmt {
	return &GetCfgVarFunction{}
}

type GetCfgVarFunction struct{}

func (fn *GetCfgVarFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

func (fn *GetCfgVarFunction) GetName() string { return "get_cfg_var" }

func (fn *GetCfgVarFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "option", 0, nil, nil),
	}
}

func (fn *GetCfgVarFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "option", 0, data.NewBaseType("string")),
	}
}

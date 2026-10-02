package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

func NewIsSubclassOfFunction() data.FuncStmt {
	return &IsSubclassOfFunction{}
}

type IsSubclassOfFunction struct{}

func (fn *IsSubclassOfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	source, _ := ctx.GetIndexValue(0)
	target, err := utils.ConvertFromIndex[string](ctx, 1)
	if err != nil || target == "" {
		return data.NewBoolValue(false), nil
	}
	allowString := true
	if value, ok := ctx.GetIndexValue(2); ok && value != nil {
		if b, ok := value.(data.AsBool); ok {
			allowString, _ = b.AsBool()
		}
	}
	result, ctl := checkNominalRelation(ctx, source, target, allowString, true)
	return data.NewBoolValue(result), ctl
}

func (fn *IsSubclassOfFunction) GetName() string { return "is_subclass_of" }

var isSubclassOfFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "object_or_class", 0, nil, nil),
	node.NewParameter(nil, "class", 1, nil, data.String{}),
	node.NewParameter(nil, "allow_string", 2, data.NewBoolValue(true), data.Bool{}),
}

func (fn *IsSubclassOfFunction) GetParams() []data.GetValue {
	return isSubclassOfFunctionGetParams
}

var isSubclassOfFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "object_or_class", 0, data.Mixed{}),
	node.NewVariable(nil, "class", 1, data.String{}),
	node.NewVariable(nil, "allow_string", 2, data.Bool{}),
}

func (fn *IsSubclassOfFunction) GetVariables() []data.Variable {
	return isSubclassOfFunctionGetVariables
}

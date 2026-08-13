package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewClassImplementsFunction() data.FuncStmt {
	return &ClassImplementsFunction{}
}

type ClassImplementsFunction struct{}

func (fn *ClassImplementsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	classVal, _ := ctx.GetIndexValue(0)

	className := ""
	switch c := classVal.(type) {
	case *data.ClassValue:
		className = c.Class.GetName()
	case *data.StringValue:
		className = c.Value
	case data.AsString:
		className = c.AsString()
	}

	if className == "" {
		return data.NewArrayValue([]data.Value{}), nil
	}

	vm := ctx.GetVM()
	cls, acl := vm.GetOrLoadClass(className)
	if acl != nil || cls == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	implements := cls.GetImplements()
	result := make([]data.Value, 0, len(implements))
	for _, iface := range implements {
		result = append(result, data.NewStringValue(iface))
	}

	return data.NewArrayValue(result), nil
}

func (fn *ClassImplementsFunction) GetName() string {
	return "class_implements"
}

func (fn *ClassImplementsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, nil),
	}
}

func (fn *ClassImplementsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, data.Mixed{}),
	}
}

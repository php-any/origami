package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewClassParentsFunction() data.FuncStmt {
	return &ClassParentsFunction{}
}

type ClassParentsFunction struct{}

func (fn *ClassParentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	objVal, _ := ctx.GetIndexValue(0)

	className := ""
	switch v := objVal.(type) {
	case *data.ClassValue:
		className = v.Class.GetName()
	case *data.StringValue:
		className = v.Value
	case data.AsString:
		className = v.AsString()
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}

	if className == "" {
		return data.NewArrayValue([]data.Value{}), nil
	}

	vm := ctx.GetVM()
	cls, acl := vm.GetOrLoadClass(className)
	if acl != nil || cls == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	result := make([]data.Value, 0)
	extend := cls.GetExtend()
	for extend != nil {
		result = append(result, data.NewStringValue(*extend))
		parent, acl := vm.GetOrLoadClass(*extend)
		if acl != nil || parent == nil {
			break
		}
		extend = parent.GetExtend()
	}

	return data.NewArrayValue(result), nil
}

func (fn *ClassParentsFunction) GetName() string { return "class_parents" }
func (fn *ClassParentsFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object_or_class", 0, nil, nil)}
}
func (fn *ClassParentsFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object_or_class", 0, data.Mixed{})}
}

package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIsSubclassOfFunction() data.FuncStmt {
	return &IsSubclassOfFunction{}
}

type IsSubclassOfFunction struct{}

func (fn *IsSubclassOfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	objVal, _ := ctx.GetIndexValue(0)
	classVal, _ := ctx.GetIndexValue(1)

	objClassName := ""
	switch v := objVal.(type) {
	case *data.ClassValue:
		objClassName = v.Class.GetName()
	case *data.StringValue:
		objClassName = v.Value
	case data.AsString:
		objClassName = v.AsString()
	default:
		return data.NewBoolValue(false), nil
	}

	targetClass := ""
	switch v := classVal.(type) {
	case *data.StringValue:
		targetClass = v.Value
	case data.AsString:
		targetClass = v.AsString()
	default:
		return data.NewBoolValue(false), nil
	}

	if objClassName == "" || targetClass == "" {
		return data.NewBoolValue(false), nil
	}

	vm := ctx.GetVM()
	cls, acl := vm.GetOrLoadClass(objClassName)
	if acl != nil || cls == nil {
		return data.NewBoolValue(false), nil
	}

	extend := cls.GetExtend()
	for extend != nil {
		if *extend == targetClass {
			return data.NewBoolValue(true), nil
		}
		parent, acl := vm.GetOrLoadClass(*extend)
		if acl != nil || parent == nil {
			break
		}
		extend = parent.GetExtend()
	}

	return data.NewBoolValue(false), nil
}

func (fn *IsSubclassOfFunction) GetName() string { return "is_subclass_of" }
func (fn *IsSubclassOfFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, nil),
		node.NewParameter(nil, "class", 1, nil, nil),
	}
}
func (fn *IsSubclassOfFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, data.Mixed{}),
		node.NewVariable(nil, "class", 1, data.Mixed{}),
	}
}

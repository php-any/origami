package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewClassUsesFunction() data.FuncStmt {
	return &ClassUsesFunction{}
}

type ClassUsesFunction struct{}

func (fn *ClassUsesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 返回类直接使用的 trait 名，键与值均为 trait 全限定名（对齐 PHP class_uses()）
	classStmt, ok := cls.(*node.ClassStatement)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	list := make([]*data.ZVal, 0, len(classStmt.Traits))
	for _, traitName := range classStmt.Traits {
		list = append(list, data.NewNamedZVal(traitName, data.NewStringValue(traitName)))
	}
	return &data.ArrayValue{List: list}, nil
}

func (fn *ClassUsesFunction) GetName() string { return "class_uses" }
func (fn *ClassUsesFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object_or_class", 0, nil, nil)}
}
func (fn *ClassUsesFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object_or_class", 0, data.Mixed{})}
}

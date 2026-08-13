package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// GetParentClassFunction 实现 get_parent_class 函数
// get_parent_class(object|string $object_or_class = ?): string|false
type GetParentClassFunction struct{}

func NewGetParentClassFunction() data.FuncStmt {
	return &GetParentClassFunction{}
}

func (f *GetParentClassFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, ok := parentClassLookupName(ctx)
	if !ok || className == "" {
		return data.NewBoolValue(false), nil
	}

	cls, acl := ctx.GetVM().GetOrLoadClass(className)
	if acl != nil || cls == nil {
		// PHP：类不存在时返回 false，不抛错
		return data.NewBoolValue(false), nil
	}
	if cls.GetExtend() == nil || *cls.GetExtend() == "" {
		return data.NewBoolValue(false), nil
	}

	return data.NewStringValue(*cls.GetExtend()), nil
}

func parentClassLookupName(ctx data.Context) (string, bool) {
	objectValue, hasObject := ctx.GetIndexValue(0)
	if !hasObject || objectValue == nil {
		return "", false
	}
	if cv, ok := objectValue.(*data.ClassValue); ok && cv.Class != nil {
		return cv.Class.GetName(), true
	}
	if s, ok := objectValue.(data.AsString); ok {
		name := strings.TrimSpace(s.AsString())
		if name == "" {
			return "", false
		}
		return name, true
	}
	return "", false
}

func (f *GetParentClassFunction) GetName() string {
	return "get_parent_class"
}

func (f *GetParentClassFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, nil),
	}
}

func (f *GetParentClassFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, nil),
	}
}

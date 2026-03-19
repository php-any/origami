package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MethodExistsFunction 实现 method_exists 函数
// method_exists(object|string $object_or_class, string $method): bool
// 检查对象或类是否具有该方法（包含继承链）
type MethodExistsFunction struct{}

func NewMethodExistsFunction() data.FuncStmt {
	return &MethodExistsFunction{}
}

func (f *MethodExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：对象或类名
	objectOrClassValue, ok := ctx.GetIndexValue(0)
	if !ok || objectOrClassValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取第二个参数：方法名
	methodValue, ok := ctx.GetIndexValue(1)
	if !ok || methodValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取方法名字符串
	methodName := methodValue.AsString()
	if methodName == "" {
		return data.NewBoolValue(false), nil
	}

	vm := ctx.GetVM()

	// 如果第一个参数是 ClassValue（对象实例），直接通过 ClassValue.GetMethod 查询（含继承链）
	if classValue, ok := objectOrClassValue.(*data.ClassValue); ok {
		_, found := classValue.GetMethod(methodName)
		return data.NewBoolValue(found), nil
	}

	// 如果是 ThisValue，展开为 ClassValue
	if thisValue, ok := objectOrClassValue.(*data.ThisValue); ok {
		_, found := thisValue.ClassValue.GetMethod(methodName)
		return data.NewBoolValue(found), nil
	}

	// 如果是字符串，视为类名，加载类后查询
	var className string
	switch o := objectOrClassValue.(type) {
	case *data.StringValue:
		className = o.Value
	case data.GetName:
		className = o.GetName()
	default:
		className = objectOrClassValue.AsString()
	}

	if className == "" {
		return data.NewBoolValue(false), nil
	}

	classStmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查当前类方法
	if _, found := classStmt.GetMethod(methodName); found {
		return data.NewBoolValue(true), nil
	}

	// 检查继承链中的方法
	last := classStmt
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		parentClass, acl := vm.GetClass(*ext)
		if !acl {
			break
		}
		if _, found := parentClass.GetMethod(methodName); found {
			return data.NewBoolValue(true), nil
		}
		last = parentClass
	}

	return data.NewBoolValue(false), nil
}

func (f *MethodExistsFunction) GetName() string {
	return "method_exists"
}

func (f *MethodExistsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, nil),
		node.NewParameter(nil, "method", 1, nil, nil),
	}
}

func (f *MethodExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, data.NewBaseType("object|string")),
		node.NewVariable(nil, "method", 1, data.NewBaseType("string")),
	}
}

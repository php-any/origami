package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PropertyExistsFunction 实现 property_exists 函数
// property_exists(object|string $object_or_class, string $property): bool
// 检查对象或类是否具有该属性
type PropertyExistsFunction struct{}

func NewPropertyExistsFunction() data.FuncStmt {
	return &PropertyExistsFunction{}
}

func (f *PropertyExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：对象或类名
	objectOrClassValue, _ := ctx.GetIndexValue(0)
	if objectOrClassValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取第二个参数：属性名
	propertyValue, _ := ctx.GetIndexValue(1)
	if propertyValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取属性名字符串
	var propertyName string
	if strValue, ok := propertyValue.(data.AsString); ok {
		propertyName = strValue.AsString()
	} else {
		propertyName = propertyValue.AsString()
	}

	// 获取类名
	var className string
	if classValue, ok := objectOrClassValue.(*data.ClassValue); ok {
		// 第一个参数是 ClassValue 对象，获取其类名
		className = classValue.Class.GetName()
	} else if strValue, ok := objectOrClassValue.(*data.StringValue); ok {
		// 第一个参数是字符串字面量，视为类名
		className = strValue.AsString()
	} else {
		// 不是 ClassValue 也不是 StringValue，返回 false
		return data.NewBoolValue(false), nil
	}

	vm := ctx.GetVM()

	// 获取类
	classStmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return data.NewBoolValue(false), nil
	}
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查当前类的属性列表
	properties := classStmt.GetPropertyList()
	for _, prop := range properties {
		if prop.GetName() == propertyName {
			return data.NewBoolValue(true), nil
		}
	}

	// 检查继承的属性
	last := classStmt
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		parentClass, ok := vm.GetClass(*ext)
		if !ok {
			break
		}
		parentProperties := parentClass.GetPropertyList()
		for _, prop := range parentProperties {
			if prop.GetName() == propertyName {
				// 找到属性，返回 true（不考虑访问修饰符，property_exists 会检查所有可见性）
				return data.NewBoolValue(true), nil
			}
		}
		last = parentClass
	}

	return data.NewBoolValue(false), nil
}

func (f *PropertyExistsFunction) GetName() string {
	return "property_exists"
}

func (f *PropertyExistsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object_or_class", 0, nil, nil),
		node.NewParameter(nil, "property", 1, nil, nil),
	}
}

func (f *PropertyExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object_or_class", 0, data.NewBaseType("object|string")),
		node.NewVariable(nil, "property", 1, data.NewBaseType("string")),
	}
}

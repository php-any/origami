// Package reflection 提供 PHP 反射功能的实现
package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionAttributeClass 提供 PHP ReflectionAttribute 类定义
// ReflectionAttribute 用于获取属性（attributes/annotations）的信息
type ReflectionAttributeClass struct {
	node.Node
}

// GetValue 创建 ReflectionAttribute 的实例
func (c *ReflectionAttributeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionAttribute"
func (c *ReflectionAttributeClass) GetName() string { return "ReflectionAttribute" }

// GetExtend 返回父类名，ReflectionAttribute 没有父类
func (c *ReflectionAttributeClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表，ReflectionAttribute 不实现任何接口
func (c *ReflectionAttributeClass) GetImplements() []string { return nil }

// GetProperty 获取属性，ReflectionAttribute 没有属性
func (c *ReflectionAttributeClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，ReflectionAttribute 没有属性
func (c *ReflectionAttributeClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
func (c *ReflectionAttributeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ReflectionAttributeConstructMethod{}, true
	case "getName":
		return &ReflectionAttributeGetNameMethod{}, true
	case "getArguments":
		return &ReflectionAttributeGetArgumentsMethod{}, true
	case "newInstance":
		return &ReflectionAttributeNewInstanceMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *ReflectionAttributeClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionAttributeConstructMethod{},
		&ReflectionAttributeGetNameMethod{},
		&ReflectionAttributeGetArgumentsMethod{},
		&ReflectionAttributeNewInstanceMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *ReflectionAttributeClass) GetConstruct() data.Method {
	return &ReflectionAttributeConstructMethod{}
}

// getReflectionAttributeInfo 从上下文中获取 ReflectionAttribute 的属性信息
func getReflectionAttributeInfo(ctx data.Context) data.GetValue {
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 从 ObjectValue 的 property 中获取注解对象
		if objCtx.ObjectValue != nil {
			props := objCtx.ObjectValue.GetProperties()
			annotationVal, hasAnnotation := props["_annotation"]
			if hasAnnotation {
				return annotationVal
			}
		}
	}
	return nil
}

// Package reflection 提供 PHP 反射功能的实现
package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassClass 提供 PHP ReflectionClass 类定义
// ReflectionClass 用于获取类的信息，包括方法、属性、继承关系等
type ReflectionClassClass struct {
	node.Node
}

// GetValue 创建 ReflectionClass 的实例
func (c *ReflectionClassClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionClass"
func (c *ReflectionClassClass) GetName() string { return "ReflectionClass" }

// GetExtend 返回父类名，ReflectionClass 没有父类
func (c *ReflectionClassClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表，ReflectionClass 不实现任何接口
func (c *ReflectionClassClass) GetImplements() []string { return nil }

// GetProperty 获取属性，ReflectionClass 没有属性
func (c *ReflectionClassClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，ReflectionClass 没有属性
func (c *ReflectionClassClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
// name: 方法名
// 返回: 方法对象和是否存在
func (c *ReflectionClassClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ReflectionClassConstructMethod{}, true
	case "getName":
		return &ReflectionClassGetNameMethod{}, true
	case "getMethods":
		return &ReflectionClassGetMethodsMethod{}, true
	case "getMethod":
		return &ReflectionClassGetMethodMethod{}, true
	case "getProperties":
		return &ReflectionClassGetPropertiesMethod{}, true
	case "getProperty":
		return &ReflectionClassGetPropertyMethod{}, true
	case "hasMethod":
		return &ReflectionClassHasMethodMethod{}, true
	case "hasProperty":
		return &ReflectionClassHasPropertyMethod{}, true
	case "isSubclassOf":
		return &ReflectionClassIsSubclassOfMethod{}, true
	case "getParentClass":
		return &ReflectionClassGetParentClassMethod{}, true
	case "isInstance":
		return &ReflectionClassIsInstanceMethod{}, true
	case "newInstance":
		return &ReflectionClassNewInstanceMethod{}, true
	case "newInstanceWithoutConstructor":
		return &ReflectionClassNewInstanceWithoutConstructorMethod{}, true
	case "isInstantiable":
		return &ReflectionClassIsInstantiableMethod{}, true
	case "getConstructor":
		return &ReflectionClassGetConstructorMethod{}, true
	case "newInstanceArgs":
		return &ReflectionClassNewInstanceArgsMethod{}, true
	case "getAttributes":
		return &ReflectionClassGetAttributesMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *ReflectionClassClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionClassConstructMethod{},
		&ReflectionClassGetNameMethod{},
		&ReflectionClassGetMethodsMethod{},
		&ReflectionClassGetMethodMethod{},
		&ReflectionClassGetPropertiesMethod{},
		&ReflectionClassGetPropertyMethod{},
		&ReflectionClassHasMethodMethod{},
		&ReflectionClassHasPropertyMethod{},
		&ReflectionClassIsSubclassOfMethod{},
		&ReflectionClassGetParentClassMethod{},
		&ReflectionClassIsInstanceMethod{},
		&ReflectionClassNewInstanceMethod{},
		&ReflectionClassNewInstanceWithoutConstructorMethod{},
		&ReflectionClassNewInstanceArgsMethod{},
		&ReflectionClassIsInstantiableMethod{},
		&ReflectionClassGetConstructorMethod{},
		&ReflectionClassGetAttributesMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *ReflectionClassClass) GetConstruct() data.Method {
	return &ReflectionClassConstructMethod{}
}

// ReflectionClassValue 表示 ReflectionClass 的实例
// 用于存储被反射的类信息
type ReflectionClassValue struct {
	*data.ClassValue
	className string         // 被反射的类名
	classStmt data.ClassStmt // 被反射的类语句
}

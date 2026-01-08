// Package reflection 提供 PHP 反射功能的实现
package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionNamedTypeClass 提供 PHP ReflectionNamedType 类定义
// ReflectionNamedType 继承自 ReflectionType，用于表示命名类型（如 string, int, MyClass 等）
type ReflectionNamedTypeClass struct {
	node.Node
}

// GetValue 创建 ReflectionNamedType 的实例
func (c *ReflectionNamedTypeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionNamedType"
func (c *ReflectionNamedTypeClass) GetName() string { return "ReflectionNamedType" }

// GetExtend 返回父类名，ReflectionNamedType 继承自 ReflectionType
func (c *ReflectionNamedTypeClass) GetExtend() *string {
	parent := "ReflectionType"
	return &parent
}

// GetImplements 返回实现的接口列表，ReflectionNamedType 不实现任何接口
func (c *ReflectionNamedTypeClass) GetImplements() []string { return nil }

// GetProperty 获取属性，ReflectionNamedType 没有属性
func (c *ReflectionNamedTypeClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，ReflectionNamedType 没有属性
func (c *ReflectionNamedTypeClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
// ReflectionNamedType 继承 ReflectionType 的所有方法
func (c *ReflectionNamedTypeClass) GetMethod(name string) (data.Method, bool) {
	// ReflectionNamedType 继承 ReflectionType 的所有方法
	// 使用 ReflectionType 的方法实现
	typeClass := &ReflectionTypeClass{}
	return typeClass.GetMethod(name)
}

// GetMethods 返回所有方法列表
// ReflectionNamedType 继承 ReflectionType 的所有方法
func (c *ReflectionNamedTypeClass) GetMethods() []data.Method {
	typeClass := &ReflectionTypeClass{}
	return typeClass.GetMethods()
}

// GetConstruct 返回构造函数
func (c *ReflectionNamedTypeClass) GetConstruct() data.Method {
	return &ReflectionNamedTypeConstructMethod{}
}

// newReflectionNamedType 创建一个新的 ReflectionNamedType 实例
// 这是一个辅助函数，用于创建 ReflectionNamedType 对象
func newReflectionNamedType(ctx data.Context, typeInfo data.Types) *data.ClassValue {
	typeClass := &ReflectionNamedTypeClass{}
	typeValue := data.NewClassValue(typeClass, ctx.CreateBaseContext())

	// 存储类型信息到实例属性中
	typeName := ""
	if typeInfo != nil {
		typeName = typeInfo.String()
	}
	typeValue.ObjectValue.SetProperty("_typeName", data.NewStringValue(typeName))

	return typeValue
}

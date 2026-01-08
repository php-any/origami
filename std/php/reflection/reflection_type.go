// Package reflection 提供 PHP 反射功能的实现
package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionTypeClass 提供 PHP ReflectionType 类定义
// ReflectionType 用于获取类型的信息，包括类型名、是否为内置类型等
type ReflectionTypeClass struct {
	node.Node
}

// GetValue 创建 ReflectionType 的实例
func (c *ReflectionTypeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionType"
func (c *ReflectionTypeClass) GetName() string { return "ReflectionType" }

// GetExtend 返回父类名，ReflectionType 没有父类
func (c *ReflectionTypeClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表，ReflectionType 不实现任何接口
func (c *ReflectionTypeClass) GetImplements() []string { return nil }

// GetProperty 获取属性，ReflectionType 没有属性
func (c *ReflectionTypeClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，ReflectionType 没有属性
func (c *ReflectionTypeClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
func (c *ReflectionTypeClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ReflectionTypeConstructMethod{}, true
	case "getName":
		return &ReflectionTypeGetNameMethod{}, true
	case "__toString":
		return &ReflectionTypeToStringMethod{}, true
	case "isBuiltin":
		return &ReflectionTypeIsBuiltinMethod{}, true
	case "allowsNull":
		return &ReflectionTypeAllowsNullMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *ReflectionTypeClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionTypeConstructMethod{},
		&ReflectionTypeGetNameMethod{},
		&ReflectionTypeToStringMethod{},
		&ReflectionTypeIsBuiltinMethod{},
		&ReflectionTypeAllowsNullMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *ReflectionTypeClass) GetConstruct() data.Method {
	return &ReflectionTypeConstructMethod{}
}

// newReflectionType 创建一个新的 ReflectionType 实例
// 这是一个辅助函数，用于创建 ReflectionType 对象
func newReflectionType(ctx data.Context, typeInfo data.Types) *data.ClassValue {
	typeClass := &ReflectionTypeClass{}
	typeValue := data.NewClassValue(typeClass, ctx.CreateBaseContext())

	// 存储类型信息到实例属性中
	// 我们需要存储类型的字符串表示，以便后续使用
	typeName := ""
	if typeInfo != nil {
		typeName = typeInfo.String()
	}
	typeValue.ObjectValue.SetProperty("_typeName", data.NewStringValue(typeName))

	// 存储原始类型对象（如果可能）
	// 注意：由于 data.Types 是接口，我们不能直接存储，所以存储字符串表示

	return typeValue
}

// getReflectionTypeInfo 从上下文中获取 ReflectionType 的类型信息
func getReflectionTypeInfo(ctx data.Context) (string, data.Types) {
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 从 ObjectValue 的 property 中获取类型名
		if objCtx.ObjectValue != nil {
			props := objCtx.ObjectValue.GetProperties()
			typeNameVal, hasTypeName := props["_typeName"]

			if hasTypeName {
				var typeName string
				if strVal, ok := typeNameVal.(*data.StringValue); ok {
					typeName = strVal.AsString()
				}

				if typeName != "" {
					// 尝试从类型名重建类型对象
					typeInfo := data.NewBaseType(typeName)
					return typeName, typeInfo
				}
			}
		}
	}
	return "", nil
}

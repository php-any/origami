// Package reflection 提供 PHP 反射功能的实现
package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionParameterClass 提供 PHP ReflectionParameter 类定义
// ReflectionParameter 用于获取方法参数的信息，包括参数名、位置、默认值等
type ReflectionParameterClass struct {
	node.Node
}

// GetValue 创建 ReflectionParameter 的实例
func (c *ReflectionParameterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionParameter"
func (c *ReflectionParameterClass) GetName() string { return "ReflectionParameter" }

// GetExtend 返回父类名，ReflectionParameter 没有父类
func (c *ReflectionParameterClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表，ReflectionParameter 不实现任何接口
func (c *ReflectionParameterClass) GetImplements() []string { return nil }

// GetProperty 获取属性，ReflectionParameter 没有属性
func (c *ReflectionParameterClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，ReflectionParameter 没有属性
func (c *ReflectionParameterClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
func (c *ReflectionParameterClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ReflectionParameterConstructMethod{}, true
	case "getName":
		return &ReflectionParameterGetNameMethod{}, true
	case "getPosition":
		return &ReflectionParameterGetPositionMethod{}, true
	case "isOptional":
		return &ReflectionParameterIsOptionalMethod{}, true
	case "isDefaultValueAvailable":
		return &ReflectionParameterIsDefaultValueAvailableMethod{}, true
	case "getDefaultValue":
		return &ReflectionParameterGetDefaultValueMethod{}, true
	case "getType":
		return &ReflectionParameterGetTypeMethod{}, true
	case "getDeclaringClass":
		return &ReflectionParameterGetDeclaringClassMethod{}, true
	case "isVariadic":
		return &ReflectionParameterIsVariadicMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *ReflectionParameterClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionParameterConstructMethod{},
		&ReflectionParameterGetNameMethod{},
		&ReflectionParameterGetPositionMethod{},
		&ReflectionParameterIsOptionalMethod{},
		&ReflectionParameterIsDefaultValueAvailableMethod{},
		&ReflectionParameterGetDefaultValueMethod{},
		&ReflectionParameterGetTypeMethod{},
		&ReflectionParameterGetDeclaringClassMethod{},
		&ReflectionParameterIsVariadicMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *ReflectionParameterClass) GetConstruct() data.Method {
	return &ReflectionParameterConstructMethod{}
}

// newReflectionParameter 创建一个新的 ReflectionParameter 实例
// 这是一个辅助函数，用于创建 ReflectionParameter 对象
func newReflectionParameter(ctx data.Context, className string, methodName string, paramIndex int, param data.GetValue) *data.ClassValue {
	paramClass := &ReflectionParameterClass{}
	paramValue := data.NewClassValue(paramClass, ctx.CreateBaseContext())

	// 存储参数信息到实例属性中
	paramValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	paramValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))
	paramValue.ObjectValue.SetProperty("_paramIndex", data.NewIntValue(paramIndex))

	// 存储参数对象本身（用于后续获取参数信息）
	// 注意：这里我们需要存储参数的索引，因为 param 是 GetValue 类型，不能直接存储
	// 实际使用时，我们会从方法中重新获取参数

	return paramValue
}

// getReflectionParameterInfo 从上下文中获取 ReflectionParameter 的参数信息
func getReflectionParameterInfo(ctx data.Context) (string, string, int, data.GetValue) {
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 从 ObjectValue 的 property 中获取类名、方法名和参数索引
		// 使用 GetProperties 方法获取所有属性，确保能获取到动态设置的属性
		if objCtx.ObjectValue != nil {
			props := objCtx.ObjectValue.GetProperties()
			classNameVal, hasClassName := props["_className"]
			methodNameVal, hasMethodName := props["_methodName"]
			paramIndexVal, hasParamIndex := props["_paramIndex"]

			if hasClassName && hasMethodName && hasParamIndex {
				var className, methodName string
				var paramIndex int

				if strVal, ok := classNameVal.(*data.StringValue); ok {
					className = strVal.AsString()
				}
				if strVal, ok := methodNameVal.(*data.StringValue); ok {
					methodName = strVal.AsString()
				}
				if intVal, ok := paramIndexVal.(*data.IntValue); ok {
					paramIndex, _ = intVal.AsInt()
				}

				if className != "" && methodName != "" {
					vm := ctx.GetVM()
					stmt, _ := vm.GetOrLoadClass(className)
					if stmt != nil {
						method, exists := stmt.GetMethod(methodName)
						if exists {
							params := method.GetParams()
							if paramIndex >= 0 && paramIndex < len(params) {
								return className, methodName, paramIndex, params[paramIndex]
							}
						}
					}
				}
			}
		}
	}
	return "", "", -1, nil
}

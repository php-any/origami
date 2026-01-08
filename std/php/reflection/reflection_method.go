// Package reflection 提供 PHP 反射功能的实现
package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ReflectionMethodClass 提供 PHP ReflectionMethod 类定义
// ReflectionMethod 用于获取方法的信息，包括方法名、参数、修饰符等
type ReflectionMethodClass struct {
	node.Node
}

// GetValue 创建 ReflectionMethod 的实例
func (c *ReflectionMethodClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// GetName 返回类名 "ReflectionMethod"
func (c *ReflectionMethodClass) GetName() string { return "ReflectionMethod" }

// GetExtend 返回父类名，ReflectionMethod 没有父类
func (c *ReflectionMethodClass) GetExtend() *string { return nil }

// GetImplements 返回实现的接口列表，ReflectionMethod 不实现任何接口
func (c *ReflectionMethodClass) GetImplements() []string { return nil }

// GetProperty 获取属性，ReflectionMethod 没有属性
func (c *ReflectionMethodClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取属性列表，ReflectionMethod 没有属性
func (c *ReflectionMethodClass) GetPropertyList() []data.Property {
	return nil
}

// GetMethod 根据方法名获取方法
func (c *ReflectionMethodClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ReflectionMethodConstructMethod{}, true
	case "getName":
		return &ReflectionMethodGetNameMethod{}, true
	case "getModifiers":
		return &ReflectionMethodGetModifiersMethod{}, true
	case "isStatic":
		return &ReflectionMethodIsStaticMethod{}, true
	case "isPublic":
		return &ReflectionMethodIsPublicMethod{}, true
	case "isProtected":
		return &ReflectionMethodIsProtectedMethod{}, true
	case "isPrivate":
		return &ReflectionMethodIsPrivateMethod{}, true
	case "getParameters":
		return &ReflectionMethodGetParametersMethod{}, true
	case "getNumberOfParameters":
		return &ReflectionMethodGetNumberOfParametersMethod{}, true
	}
	return nil, false
}

// GetMethods 返回所有方法列表
func (c *ReflectionMethodClass) GetMethods() []data.Method {
	return []data.Method{
		&ReflectionMethodConstructMethod{},
		&ReflectionMethodGetNameMethod{},
		&ReflectionMethodGetModifiersMethod{},
		&ReflectionMethodIsStaticMethod{},
		&ReflectionMethodIsPublicMethod{},
		&ReflectionMethodIsProtectedMethod{},
		&ReflectionMethodIsPrivateMethod{},
		&ReflectionMethodGetParametersMethod{},
		&ReflectionMethodGetNumberOfParametersMethod{},
	}
}

// GetConstruct 返回构造函数
func (c *ReflectionMethodClass) GetConstruct() data.Method {
	return &ReflectionMethodConstructMethod{}
}

// ReflectionMethodValue 表示 ReflectionMethod 的实例
// 用于存储被反射的方法信息
type ReflectionMethodValue struct {
	*data.ClassValue
	className  string      // 被反射的类名
	methodName string      // 被反射的方法名
	method     data.Method // 被反射的方法对象
}

// newReflectionMethod 创建一个新的 ReflectionMethod 实例
// 这是一个辅助函数，用于创建 ReflectionMethod 对象
func newReflectionMethod(ctx data.Context, className string, methodName string) *data.ClassValue {
	methodClass := &ReflectionMethodClass{}
	methodValue := data.NewClassValue(methodClass, ctx.CreateBaseContext())

	// 存储方法信息到实例属性中
	methodValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	methodValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))

	return methodValue
}

// getReflectionMethodInfo 从上下文中获取 ReflectionMethod 的方法信息
func getReflectionMethodInfo(ctx data.Context) (string, string, data.Method) {
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 从 ObjectValue 的 property 中直接获取类名和方法名
		// ClassMethodContext 嵌入了 *ClassValue，所以 objCtx 本身就是 ClassValue
		if objCtx.ObjectValue != nil {
			// 使用 GetProperties 方法获取所有属性，然后查找 _className 和 _methodName
			props := objCtx.ObjectValue.GetProperties()
			classNameVal, hasClassName := props["_className"]
			methodNameVal, hasMethodName := props["_methodName"]

			if hasClassName && hasMethodName {
				var className, methodName string
				if strVal, ok := classNameVal.(*data.StringValue); ok {
					className = strVal.AsString()
				}
				if strVal, ok := methodNameVal.(*data.StringValue); ok {
					methodName = strVal.AsString()
				}

				if className != "" && methodName != "" {
					vm := ctx.GetVM()
					stmt, _ := vm.GetOrLoadClass(className)
					if stmt != nil {
						method, exists := stmt.GetMethod(methodName)
						if exists {
							return className, methodName, method
						} else if methodName == token.ConstructName && stmt.GetConstruct() != nil {
							return className, methodName, stmt.GetConstruct()
						}
					}
				}
			}
		}
	}
	return "", "", nil
}

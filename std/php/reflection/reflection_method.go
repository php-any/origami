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
	case "getDeclaringClass":
		return &ReflectionMethodGetDeclaringClassMethod{}, true
	case "setAccessible":
		return &ReflectionMethodSetAccessibleMethod{}, true
	case "invoke":
		return &ReflectionMethodInvokeMethod{}, true
	case "hasReturnType":
		return &ReflectionMethodHasReturnTypeMethod{}, true
	case "getAttributes":
		return &ReflectionMethodGetAttributesMethod{}, true
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
		&ReflectionMethodGetDeclaringClassMethod{},
		&ReflectionMethodSetAccessibleMethod{},
		&ReflectionMethodInvokeMethod{},
		&ReflectionMethodHasReturnTypeMethod{},
		&ReflectionMethodGetAttributesMethod{},
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

func newReflectionMethod(ctx data.Context, className, methodName string) *data.ClassValue {
	methodClass := &ReflectionMethodClass{}
	methodValue := data.NewClassValue(methodClass, ctx.CreateBaseContext())
	methodValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	methodValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))
	return methodValue
}

// getReflectionMethodInfo 从上下文中获取 ReflectionMethod 的方法信息
// 支持查找继承的方法
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

					// 优先按接口解析（ReflectionMethod 也可作用于 interface 方法）
					if iface, acl := vm.GetOrLoadInterface(className); acl == nil && iface != nil {
						if method, exists := iface.GetMethod(methodName); exists {
							return className, methodName, method
						}
						for _, ext := range iface.GetExtends() {
							parent, pacl := vm.GetOrLoadInterface(ext)
							if pacl != nil || parent == nil {
								continue
							}
							if method, exists := parent.GetMethod(methodName); exists {
								return ext, methodName, method
							}
						}
					}

					v, acl := vm.LoadPkg(className)
					if acl != nil {
						return "", "", nil
					}
					if v != nil {
						stmt, ok := v.(data.ClassStmt)
						if !ok {
							return "", "", nil
						}
						// 首先在当前类中查找方法
						method, exists := stmt.GetMethod(methodName)
						if exists {
							return className, methodName, method
						}
						if staticMethods, ok := stmt.(data.GetStaticMethod); ok {
							if method, exists := staticMethods.GetStaticMethod(methodName); exists {
								return className, methodName, method
							}
						}

						// 如果是构造函数，检查当前类的构造函数
						if methodName == token.ConstructName && stmt.GetConstruct() != nil {
							return className, methodName, stmt.GetConstruct()
						}

						// 如果当前类没有找到，查找继承的方法
						declaringClassName := className
						last := stmt
						for last.GetExtend() != nil {
							ext := last.GetExtend()
							pkgVal, acl := vm.LoadPkg(*ext)
							if acl != nil {
								return "", "", nil
							}
							parentStmt, ok := pkgVal.(data.ClassStmt)
							if !ok || parentStmt == nil {
								break
							}

							// 在父类中查找方法
							method, exists := parentStmt.GetMethod(methodName)
							if exists {
								// 检查访问权限，私有方法不能继承
								if method.GetModifier() != data.ModifierPrivate {
									declaringClassName = *ext
									return declaringClassName, methodName, method
								}
							}
							if staticMethods, ok := parentStmt.(data.GetStaticMethod); ok {
								if method, exists := staticMethods.GetStaticMethod(methodName); exists &&
									method.GetModifier() != data.ModifierPrivate {
									declaringClassName = *ext
									return declaringClassName, methodName, method
								}
							}

							// 如果是构造函数，检查父类的构造函数
							if methodName == token.ConstructName && parentStmt.GetConstruct() != nil {
								declaringClassName = *ext
								return declaringClassName, methodName, parentStmt.GetConstruct()
							}

							last = parentStmt
						}
					}
				}
			}
		}
	}
	return "", "", nil
}

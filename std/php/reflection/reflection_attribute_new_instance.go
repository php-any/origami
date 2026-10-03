package reflection

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionAttributeNewInstanceMethod 实现 ReflectionAttribute::newInstance
// 创建属性（注解）的新实例
type ReflectionAttributeNewInstanceMethod struct{}

// GetName 返回方法名 "newInstance"
func (m *ReflectionAttributeNewInstanceMethod) GetName() string { return "newInstance" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionAttributeNewInstanceMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionAttributeNewInstanceMethod) GetIsStatic() bool { return false }

var reflectionAttributeNewInstanceMethodGetParams = []data.GetValue{}

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionAttributeNewInstanceMethod) GetParams() []data.GetValue {
	return reflectionAttributeNewInstanceMethodGetParams
}

var reflectionAttributeNewInstanceMethodGetVariables = []data.Variable{}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionAttributeNewInstanceMethod) GetVariables() []data.Variable {
	return reflectionAttributeNewInstanceMethodGetVariables
}

// GetReturnType 返回返回类型，返回对象类型
func (m *ReflectionAttributeNewInstanceMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// Call 执行 newInstance 方法
// 创建属性（注解）的新实例
func (m *ReflectionAttributeNewInstanceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if annotation := reflectionAttributeDeclaration(ctx); annotation != nil {
		class, ctl := ctx.GetVM().GetOrLoadClass(annotation.Name)
		if ctl != nil {
			return nil, ctl
		}
		if object, ok := ctx.(*data.ClassMethodContext); ok {
			if metadata, ok := object.Class.(*ReflectionAttributeClass); ok && metadata.target != 0 {
				flags := -1
				if declared, ok := class.(*node.ClassStatement); ok {
					for _, marker := range declared.Annotations {
						if data.TypeNameEqual(marker.Class.GetName(), "Attribute") {
							if value, ctl := marker.GetProperty("flags"); ctl == nil {
								if number, ok := value.(data.AsInt); ok {
									flags, _ = number.AsInt()
								}
							}
						}
					}
				}
				if flags < 0 {
					return nil, data.NewErrorThrowByName(nil, errors.New("Attempting to use non-attribute class "+annotation.Name+" as attribute"), "Error")
				}
				if flags&metadata.target == 0 {
					return nil, data.NewErrorThrowByName(nil, errors.New("Attribute "+annotation.Name+" cannot target function"), "Error")
				}
				if metadata.repeated && flags&64 == 0 {
					return nil, data.NewErrorThrowByName(nil, errors.New("Attribute "+annotation.Name+" must not be repeated"), "Error")
				}
			}
		}
		return node.CreateInstanceFromClass(class, annotation.Arguments, ctx)
	}
	annotationValue := getReflectionAttributeInfo(ctx)
	if annotationValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionAttribute::newInstance() failed"))
	}

	// 如果注解对象是 ClassValue，创建一个新的实例
	if classVal, ok := annotationValue.(*data.ClassValue); ok {
		// 创建新的 ClassValue 实例
		newInstance := data.NewClassValue(classVal.Class, ctx.CreateBaseContext())
		for s, value := range classVal.GetProperties() {
			newInstance.SetProperty(s, value)
		}
		return newInstance, nil
	}

	return nil, data.NewErrorThrow(nil, errors.New("ReflectionAttribute::newInstance() failed"))
}

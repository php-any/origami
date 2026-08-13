package reflection

import (
	"errors"

	"github.com/php-any/origami/data"
)

// ReflectionClassNewInstanceWithoutConstructorMethod 实现 ReflectionClass::newInstanceWithoutConstructor
// 创建被反射类的新实例，但不调用构造函数
type ReflectionClassNewInstanceWithoutConstructorMethod struct{}

// GetName 返回方法名 "newInstanceWithoutConstructor"
func (m *ReflectionClassNewInstanceWithoutConstructorMethod) GetName() string {
	return "newInstanceWithoutConstructor"
}

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassNewInstanceWithoutConstructorMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassNewInstanceWithoutConstructorMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionClassNewInstanceWithoutConstructorMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionClassNewInstanceWithoutConstructorMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回对象类型
func (m *ReflectionClassNewInstanceWithoutConstructorMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// Call 执行 newInstanceWithoutConstructor 方法
// 创建被反射类的新实例，但不调用构造函数
// 如果类不存在，抛出异常
func (m *ReflectionClassNewInstanceWithoutConstructorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return nil, data.NewErrorThrow(nil, errors.New("Class does not exist"))
	}

	// 创建实例但不调用构造函数
	object := data.NewClassValue(classStmt, ctx.CreateBaseContext())
	return object, nil
}

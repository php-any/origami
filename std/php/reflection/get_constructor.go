package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionClassGetConstructorMethod 实现 ReflectionClass::getConstructor
// 返回被反射类的构造函数
type ReflectionClassGetConstructorMethod struct{}

// GetName 返回方法名 "getConstructor"
func (m *ReflectionClassGetConstructorMethod) GetName() string { return "getConstructor" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetConstructorMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetConstructorMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionClassGetConstructorMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionClassGetConstructorMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回混合类型（ReflectionMethod 对象或 null）
func (m *ReflectionClassGetConstructorMethod) GetReturnType() data.Types {
	return data.Mixed{}
}

// Call 执行 getConstructor 方法
// 返回被反射类的构造函数
// 如果类没有构造函数，返回 null
func (m *ReflectionClassGetConstructorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewNullValue(), nil
	}

	// 获取构造函数
	constructor := classStmt.GetConstruct()
	if constructor == nil {
		return data.NewNullValue(), nil
	}

	// 创建 ReflectionMethod 实例
	methodClass := &ReflectionMethodClass{}
	methodValue := data.NewClassValue(methodClass, ctx.CreateBaseContext())

	// 存储方法信息到实例属性中
	methodValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	methodValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(constructor.GetName()))

	return methodValue, nil
}

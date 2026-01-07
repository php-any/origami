package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassIsInstantiableMethod 实现 ReflectionClass::isInstantiable
// 检查被反射的类是否可以被实例化
type ReflectionClassIsInstantiableMethod struct{}

// GetName 返回方法名 "isInstantiable"
func (m *ReflectionClassIsInstantiableMethod) GetName() string { return "isInstantiable" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassIsInstantiableMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassIsInstantiableMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionClassIsInstantiableMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionClassIsInstantiableMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionClassIsInstantiableMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 isInstantiable 方法
// 检查被反射的类是否可以被实例化
// 接口和抽象类不能被实例化，其他类可以被实例化
// 返回 true 表示可以实例化，false 表示不能实例化
func (m *ReflectionClassIsInstantiableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查是否是接口
	vm := ctx.GetVM()
	if _, isInterface := vm.GetInterface(className); isInterface {
		return data.NewBoolValue(false), nil
	}

	// 检查是否是抽象类
	// 通过类型断言检查是否是 AbstractClassStatement
	if _, isAbstract := classStmt.(*node.AbstractClassStatement); isAbstract {
		return data.NewBoolValue(false), nil
	}

	// 其他情况都可以实例化
	return data.NewBoolValue(true), nil
}

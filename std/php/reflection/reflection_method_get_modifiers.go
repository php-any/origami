package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodGetModifiersMethod 实现 ReflectionMethod::getModifiers
// 返回被反射方法的修饰符标志
type ReflectionMethodGetModifiersMethod struct{}

// GetName 返回方法名 "getModifiers"
func (m *ReflectionMethodGetModifiersMethod) GetName() string { return "getModifiers" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodGetModifiersMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodGetModifiersMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionMethodGetModifiersMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodGetModifiersMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回整数类型
func (m *ReflectionMethodGetModifiersMethod) GetReturnType() data.Types {
	return data.Int{}
}

// Call 执行 getModifiers 方法
// 返回被反射方法的修饰符标志（位掩码）
func (m *ReflectionMethodGetModifiersMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, _, method := getReflectionMethodInfo(ctx)
	if method == nil {
		return data.NewIntValue(0), nil
	}

	modifiers := 0
	modifier := method.GetModifier()
	isStatic := method.GetIsStatic()

	// PHP 反射修饰符常量
	// IS_PUBLIC = 1
	// IS_PROTECTED = 2
	// IS_PRIVATE = 4
	// IS_STATIC = 16

	if modifier == data.ModifierPublic {
		modifiers |= 1 // IS_PUBLIC
	} else if modifier == data.ModifierProtected {
		modifiers |= 2 // IS_PROTECTED
	} else if modifier == data.ModifierPrivate {
		modifiers |= 4 // IS_PRIVATE
	}

	if isStatic {
		modifiers |= 16 // IS_STATIC
	}

	return data.NewIntValue(modifiers), nil
}

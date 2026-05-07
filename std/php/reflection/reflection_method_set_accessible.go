package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionMethodSetAccessibleMethod 实现 ReflectionMethod::setAccessible
// 在 PHP 8.1+ 中此方法已弃用，设置为无操作
type ReflectionMethodSetAccessibleMethod struct{}

// GetName 返回方法名 "setAccessible"
func (m *ReflectionMethodSetAccessibleMethod) GetName() string { return "setAccessible" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionMethodSetAccessibleMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionMethodSetAccessibleMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
func (m *ReflectionMethodSetAccessibleMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionMethodSetAccessibleMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型
func (m *ReflectionMethodSetAccessibleMethod) GetReturnType() data.Types { return nil }

// Call 执行 setAccessible 方法
// PHP 8.1+ 中此方法为无操作
func (m *ReflectionMethodSetAccessibleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

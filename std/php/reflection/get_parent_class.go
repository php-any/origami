package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionClassGetParentClassMethod 实现 ReflectionClass::getParentClass
// 返回被反射类的父类名称
type ReflectionClassGetParentClassMethod struct{}

// GetName 返回方法名 "getParentClass"
func (m *ReflectionClassGetParentClassMethod) GetName() string { return "getParentClass" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetParentClassMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetParentClassMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionClassGetParentClassMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionClassGetParentClassMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回混合类型（当前实现返回字符串或 false）
func (m *ReflectionClassGetParentClassMethod) GetReturnType() data.Types {
	return data.Mixed{}
}

// Call 执行 getParentClass 方法
// 返回被反射类的父类的 ReflectionClass 实例
// 如果类没有父类，返回 false
func (m *ReflectionClassGetParentClassMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	if classStmt.GetExtend() == nil {
		return data.NewBoolValue(false), nil
	}

	parentClassName := *classStmt.GetExtend()

	// 创建 ReflectionClass 实例用于父类
	classClass := &ReflectionClassClass{}
	classValue := data.NewClassValue(classClass, ctx.CreateBaseContext())

	// 存储父类名到实例属性中
	classValue.ObjectValue.SetProperty("_className", data.NewStringValue(parentClassName))

	return classValue, nil
}

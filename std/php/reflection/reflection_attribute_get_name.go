package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionAttributeGetNameMethod 实现 ReflectionAttribute::getName
// 返回属性（注解）的名称
type ReflectionAttributeGetNameMethod struct{}

// GetName 返回方法名 "getName"
func (m *ReflectionAttributeGetNameMethod) GetName() string { return "getName" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionAttributeGetNameMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionAttributeGetNameMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionAttributeGetNameMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionAttributeGetNameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回返回类型，返回字符串类型
func (m *ReflectionAttributeGetNameMethod) GetReturnType() data.Types {
	return data.String{}
}

// Call 执行 getName 方法
// 返回属性（注解）的名称
func (m *ReflectionAttributeGetNameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	annotationValue := getReflectionAttributeInfo(ctx)
	if annotationValue == nil {
		return data.NewStringValue(""), nil
	}

	// 获取注解对象的类名
	if classVal, ok := annotationValue.(*data.ClassValue); ok {
		if classStmt := classVal.Class; classStmt != nil {
			return data.NewStringValue(classStmt.GetName()), nil
		}
	}

	return data.NewStringValue(""), nil
}

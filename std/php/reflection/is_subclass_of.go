package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassIsSubclassOfMethod 实现 ReflectionClass::isSubclassOf
// 检查被反射的类是否是指定类的子类
type ReflectionClassIsSubclassOfMethod struct{}

// GetName 返回方法名 "isSubclassOf"
func (m *ReflectionClassIsSubclassOfMethod) GetName() string { return "isSubclassOf" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassIsSubclassOfMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassIsSubclassOfMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - class: 父类名（字符串）或类对象，类型为 Mixed
func (m *ReflectionClassIsSubclassOfMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class", 0, nil, data.Mixed{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassIsSubclassOfMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.Mixed{}),
	}
}

// GetReturnType 返回返回类型，返回布尔类型
func (m *ReflectionClassIsSubclassOfMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 isSubclassOf 方法
// 检查被反射的类是否是指定类的子类
// 通过遍历继承链来检查继承关系
// 返回 true 表示是子类，false 表示不是
func (m *ReflectionClassIsSubclassOfMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	classValue, _ := ctx.GetIndexValue(0)
	if classValue == nil {
		return data.NewBoolValue(false), nil
	}

	var targetClassName string
	if classVal, ok := classValue.(*data.ClassValue); ok {
		targetClassName = classVal.Class.GetName()
	} else if strValue, ok := classValue.(*data.StringValue); ok {
		targetClassName = strValue.AsString()
	} else {
		targetClassName = classValue.AsString()
	}

	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查继承关系
	vm := ctx.GetVM()
	last := classStmt
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		if *ext == targetClassName {
			return data.NewBoolValue(true), nil
		}
		next, ok := vm.GetClass(*ext)
		if !ok {
			break
		}
		last = next
	}

	return data.NewBoolValue(false), nil
}

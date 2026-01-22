package reflection

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassNewInstanceMethod 实现 ReflectionClass::newInstance
// 创建被反射类的新实例，并调用构造函数
type ReflectionClassNewInstanceMethod struct{}

// GetName 返回方法名 "newInstance"
func (m *ReflectionClassNewInstanceMethod) GetName() string { return "newInstance" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassNewInstanceMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassNewInstanceMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 使用可变参数接收构造函数的所有参数
func (m *ReflectionClassNewInstanceMethod) GetParams() []data.GetValue {
	// 可变参数，使用 Parameters 接收所有参数
	return []data.GetValue{
		node.NewParameters(nil, "args", 0, nil, nil),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassNewInstanceMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, nil),
	}
}

// GetReturnType 返回返回类型，返回对象类型
func (m *ReflectionClassNewInstanceMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// Call 执行 newInstance 方法
// 创建被反射类的新实例，并将传入的参数传递给构造函数
// 如果类不存在或参数数量超出限制，抛出异常
func (m *ReflectionClassNewInstanceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Class %s does not exist", className))
	}

	// 获取可变参数数组
	argsVar := node.NewVariable(nil, "args", 0, nil)
	argsValue, acl := ctx.GetVariableValue(argsVar)
	if acl != nil {
		return nil, acl
	}

	// 将参数数组转换为 GetValue 列表
	args := make([]data.GetValue, 0)
	if arrayValue, ok := argsValue.(*data.ArrayValue); ok {
		valueList := arrayValue.ToValueList()
		for _, v := range valueList {
			args = append(args, v)
		}
	}

	// 使用 createInstanceAndCallConstructor 创建实例
	// 需要从 node 包导入，但由于是内部函数，我们直接实现类似逻辑
	object, acl := classStmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	if object, ok := object.(*data.ClassValue); ok {
		if method := object.Class.GetConstruct(); method != nil {
			varies := method.GetVariables()
			fnCtx := object.CreateContext(varies)
			// 入参的值设置到上下文中
			for index, arg := range args {
				tempV, acl := arg.GetValue(ctx)
				if acl != nil {
					return nil, acl
				}
				if index >= len(varies) {
					return nil, data.NewErrorThrow(nil, fmt.Errorf("对象(%v)构造函数参数数量超出限制: %d", object.Class.GetName(), index))
				}
				fnCtx.SetVariableValue(varies[index], tempV.(data.Value))
			}
			_, acl = method.Call(fnCtx)
			if acl != nil {
				return nil, acl
			}
		}
	}

	return object, nil
}

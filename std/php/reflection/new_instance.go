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

var reflectionClassNewInstanceMethodGetParams = []data.GetValue{
	node.NewParameters(nil, "args", 0, nil, nil),
}

// GetParams 返回参数列表
// 使用可变参数接收构造函数的所有参数
func (m *ReflectionClassNewInstanceMethod) GetParams() []data.GetValue {
	// 可变参数，使用 Parameters 接收所有参数
	return reflectionClassNewInstanceMethodGetParams
}

var reflectionClassNewInstanceMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "args", 0, nil),
}

// GetVariables 返回变量列表
func (m *ReflectionClassNewInstanceMethod) GetVariables() []data.Variable {
	return reflectionClassNewInstanceMethodGetVariables
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

	var args []data.GetValue
	if argValue, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := argValue.(*data.ArrayValue); ok {
			args = reflectionConstructorArguments(arr)
		} else {
			args = append(args, argValue)
		}
	}
	return reflectionNewInstance(classStmt, args, ctx)
}

func reflectionNewInstance(class data.ClassStmt, arguments []data.GetValue, ctx data.Context) (data.GetValue, data.Control) {
	if node.IsAbstractClassStmt(class) {
		return nil, createReflectionException("Cannot instantiate abstract class "+class.GetName(), ctx, nil)
	}
	if constructor := class.GetConstruct(); constructor != nil && constructor.GetModifier() != data.ModifierPublic {
		return nil, createReflectionException("Access to non-public constructor of class "+class.GetName(), ctx, nil)
	}
	return node.CreateInstanceFromClass(class, arguments, ctx)
}

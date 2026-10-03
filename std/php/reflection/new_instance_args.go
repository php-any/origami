package reflection

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassNewInstanceArgsMethod 实现 ReflectionClass::newInstanceArgs
// 创建被反射类的新实例，并使用数组中的参数调用构造函数
type ReflectionClassNewInstanceArgsMethod struct{}

// GetName 返回方法名 "newInstanceArgs"
func (m *ReflectionClassNewInstanceArgsMethod) GetName() string { return "newInstanceArgs" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassNewInstanceArgsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassNewInstanceArgsMethod) GetIsStatic() bool { return false }

var reflectionClassNewInstanceArgsMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "args", 0, nil, data.Arrays{}),
}

// GetParams 返回参数列表
// 参数:
//   - args: 构造函数参数数组，类型为 array
func (m *ReflectionClassNewInstanceArgsMethod) GetParams() []data.GetValue {
	return reflectionClassNewInstanceArgsMethodGetParams
}

var reflectionClassNewInstanceArgsMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "args", 0, data.Arrays{}),
}

// GetVariables 返回变量列表
func (m *ReflectionClassNewInstanceArgsMethod) GetVariables() []data.Variable {
	return reflectionClassNewInstanceArgsMethodGetVariables
}

// GetReturnType 返回返回类型，返回对象类型
func (m *ReflectionClassNewInstanceArgsMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// Call 执行 newInstanceArgs 方法
// 创建被反射类的新实例，并使用数组中的参数传递给构造函数
// 如果类不存在或参数数量超出限制，抛出异常
func (m *ReflectionClassNewInstanceArgsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Class %s does not exist", className))
	}

	// 获取参数数组
	argsValue, _ := ctx.GetIndexValue(0)
	if argsValue == nil {
		// 如果没有提供参数，使用空数组
		argsValue = data.NewArrayValue([]data.Value{})
	}

	array, ok := argsValue.(*data.ArrayValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionClass::newInstanceArgs() expects parameter 1 to be array"))
	}
	return reflectionNewInstance(classStmt, reflectionConstructorArguments(array), ctx)
}

func reflectionConstructorArguments(array *data.ArrayValue) []data.GetValue {
	args := make([]data.GetValue, 0, array.Len())
	slots := array.View()
	for i := 0; i < slots.Len(); i++ {
		slot := slots.At(i)
		var arg data.GetValue = slot.ReadValue()
		if slot.RefCount() > 0 {
			arg = data.NewZValValue(slot)
		}
		if _, named := slot.PHPArrayKey(i).(*data.StringValue); named {
			arg = node.NewNamedArgument(nil, slot.Name, arg)
		}
		args = append(args, arg)
	}
	return args
}

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

// GetParams 返回参数列表
// 参数:
//   - args: 构造函数参数数组，类型为 array
func (m *ReflectionClassNewInstanceArgsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "args", 0, nil, data.Arrays{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassNewInstanceArgsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "args", 0, data.Arrays{}),
	}
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

	// 将参数数组转换为 GetValue 列表
	args := make([]data.GetValue, 0)
	if arrayValue, ok := argsValue.(*data.ArrayValue); ok {
		valueList := arrayValue.ToValueList()
		for _, v := range valueList {
			args = append(args, v)
		}
	} else {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionClass::newInstanceArgs() expects parameter 1 to be array"))
	}

	// 创建实例并调用构造函数
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

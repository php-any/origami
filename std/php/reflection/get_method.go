package reflection

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetMethodMethod 实现 ReflectionClass::getMethod
// 根据方法名获取被反射类的指定方法
type ReflectionClassGetMethodMethod struct{}

// GetName 返回方法名 "getMethod"
func (m *ReflectionClassGetMethodMethod) GetName() string { return "getMethod" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetMethodMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetMethodMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - name: 方法名（字符串）
func (m *ReflectionClassGetMethodMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassGetMethodMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}

// GetReturnType 返回返回类型，返回混合类型（当前实现返回字符串）
func (m *ReflectionClassGetMethodMethod) GetReturnType() data.Types {
	return data.Mixed{}
}

// Call 执行 getMethod 方法
// 根据方法名查找并返回对应的方法
// 如果方法不存在，抛出异常
func (m *ReflectionClassGetMethodMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取方法名参数
	methodNameValue, _ := ctx.GetIndexValue(0)
	if methodNameValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionClass::getMethod() expects parameter 1 to be string"))
	}

	methodName := methodNameValue.AsString()
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s does not exist", methodName))
	}

	// 查找方法
	_, exists := classStmt.GetMethod(methodName)
	if !exists {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s does not exist", methodName))
	}

	// 创建 ReflectionMethod 实例
	methodClass := &ReflectionMethodClass{}
	methodValue := data.NewClassValue(methodClass, ctx.CreateBaseContext())

	// 存储方法信息到实例属性中
	methodValue.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	methodValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))

	return methodValue, nil
}

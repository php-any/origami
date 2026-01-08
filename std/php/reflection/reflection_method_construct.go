package reflection

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionMethodConstructMethod 实现 ReflectionMethod::__construct
// 构造函数用于初始化 ReflectionMethod 实例，接收类名和方法名作为参数
type ReflectionMethodConstructMethod struct{}

// GetName 返回方法名 "__construct"
func (m *ReflectionMethodConstructMethod) GetName() string { return "__construct" }

// GetModifier 返回方法修饰符，构造函数是公开的
func (m *ReflectionMethodConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，构造函数不是静态方法
func (m *ReflectionMethodConstructMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - class: 类名（字符串）或对象实例，类型为 Mixed
//   - method: 方法名（字符串），类型为 String
func (m *ReflectionMethodConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "method", 1, nil, data.String{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionMethodConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.Mixed{}),
		node.NewVariable(nil, "method", 1, data.String{}),
	}
}

// GetReturnType 返回返回类型，构造函数无返回值
func (m *ReflectionMethodConstructMethod) GetReturnType() data.Types { return nil }

// Call 执行构造函数
// 从参数中获取类名和方法名，加载对应的方法，并将信息存储到实例的属性中
func (m *ReflectionMethodConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：类名或对象
	classValue, _ := ctx.GetIndexValue(0)
	if classValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionMethod::__construct() expects parameter 1 to be string or object"))
	}

	var className string

	// 检查参数类型
	if classVal, ok := classValue.(*data.ClassValue); ok {
		// 参数是对象，获取其类名
		className = classVal.Class.GetName()
	} else if strValue, ok := classValue.(*data.StringValue); ok {
		// 参数是字符串，视为类名
		className = strValue.AsString()
	} else {
		// 尝试转换为字符串
		className = classValue.AsString()
	}

	// 获取第二个参数：方法名
	methodNameValue, _ := ctx.GetIndexValue(1)
	if methodNameValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionMethod::__construct() expects parameter 2 to be string"))
	}
	methodName := methodNameValue.AsString()

	// 加载类
	vm := ctx.GetVM()
	stmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Class %s does not exist", className))
	}

	// 查找方法
	_, exists := stmt.GetMethod(methodName)
	if !exists {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s::%s() does not exist", className, methodName))
	}

	// 将方法信息存储到当前对象的属性中
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 存储类名和方法名到 ObjectValue 的实例属性中
		objCtx.ObjectValue.SetProperty("_className", data.NewStringValue(className))
		objCtx.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))
	}

	return nil, nil
}

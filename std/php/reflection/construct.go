package reflection

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassConstructMethod 实现 ReflectionClass::__construct
// 构造函数用于初始化 ReflectionClass 实例，接收一个类名或对象作为参数
type ReflectionClassConstructMethod struct{}

// GetName 返回方法名 "__construct"
func (m *ReflectionClassConstructMethod) GetName() string { return "__construct" }

// GetModifier 返回方法修饰符，构造函数是公开的
func (m *ReflectionClassConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，构造函数不是静态方法
func (m *ReflectionClassConstructMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - class: 类名（字符串）或对象实例，类型为 Mixed
func (m *ReflectionClassConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class", 0, nil, data.Mixed{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class", 0, data.Mixed{}),
	}
}

// GetReturnType 返回返回类型，构造函数无返回值
func (m *ReflectionClassConstructMethod) GetReturnType() data.Types { return nil }

// Call 执行构造函数
// 从参数中获取类名或对象，加载对应的类，并将类名存储到实例的 _className 属性中
func (m *ReflectionClassConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：类名或对象
	classValue, _ := ctx.GetIndexValue(0)
	if classValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionClass::__construct() expects parameter 1 to be string or object"))
	}

	var className string

	// 检查参数类型
	if classVal, ok := classValue.(data.GetName); ok {
		// 参数是对象，获取其类名
		className = classVal.GetName()
	} else if classVal, ok := classValue.(*data.StringValue); ok {
		className = classVal.AsString()
	} else {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionClass::__construct() expects parameter 1 to be string or object"))
	}

	// 加载类
	vm := ctx.GetVM()
	stmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Class %s does not exist", className))
	}

	// 将类信息存储到当前对象的属性中
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 存储类名到 ObjectValue 的实例属性中
		objCtx.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	}

	return nil, nil
}

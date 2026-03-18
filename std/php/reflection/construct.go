package reflection

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// createReflectionException 创建 ReflectionException 异常
func createReflectionException(message string, ctx data.Context, from data.From) data.Control {
	// 使用 NewExpression 创建实例并调用构造函数
	messageExpr := node.NewStringLiteral(nil, message)
	newExpr := node.NewNewExpression(nil, "ReflectionException", []data.GetValue{messageExpr})

	object, acl := newExpr.GetValue(ctx)
	if acl != nil {
		return acl
	}

	classValue, ok := object.(*data.ClassValue)
	if !ok {
		return data.NewErrorThrow(from, fmt.Errorf("ReflectionClass error: failed to create instance"))
	}

	// 抛出异常，使用传入的 from 作为位置信息
	return data.NewErrorThrowFromClassValue(from, classValue)
}

// TODO: 修复 try-catch 异常传播问题
// 目前的问题是：当 catch 块抛出新的异常时，这个新异常没有被正确传播
// 导致 try-catch 块后面的代码继续执行
// 临时解决方案：直接抛出致命错误，终止程序

// ReflectionClassConstructMethod 实现 ReflectionClass::__construct
// 构造函数用于初始化 ReflectionClass 实例，接收一个类名或对象作为参数
type ReflectionClassConstructMethod struct {
	node.Node
}

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
		return nil, createReflectionException("Argument #1 ($class) must be of type object|string, null given", ctx, m.GetFrom())
	}

	var className string

	switch classVal := classValue.(type) {
	case data.GetName:
		// 参数是对象，获取其类名
		className = classVal.GetName()
	case *data.StringValue:
		className = classVal.AsString()
	case *data.BoolValue:
		return nil, createReflectionException(fmt.Sprintf("Argument #1 ($class) must be of type object|string, false given (value: %v). This error typically occurs when Laravel's container tries to resolve an invalid dependency. Check your service provider bindings.", classVal.Value), ctx, m.GetFrom())
	case *data.NullValue:
		return nil, createReflectionException("Argument #1 ($class) must not be null", ctx, m.GetFrom())
	default:
		// 对于非字符串和非对象的类型，抛出 ReflectionException
		// TODO: 修复 try-catch 异常传播后，改回使用 createReflectionException
		typeName := fmt.Sprintf("%T", classValue)
		// 临时方案：直接抛出致命错误
		return nil, createReflectionException(fmt.Sprintf("ReflectionClass::__construct(): Argument #1 ($class) must be of type object|string, %s given", typeName), ctx, m.GetFrom())
	}

	// 加载类
	vm := ctx.GetVM()
	stmt, acl := vm.LoadPkg(className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		// 类不存在，创建带有位置信息的异常
		var from data.From
		if objCtx, ok := ctx.(*data.ClassMethodContext); ok && objCtx.Class != nil {
			from = objCtx.Class.GetFrom()
		}
		throw := data.NewErrorThrow(from, fmt.Errorf("class %s does not exist", className))
		if tv, ok := throw.(*data.ThrowValue); ok {
			// 添加当前位置作为堆栈帧
			tv.AddStackWithInfo(from, "ReflectionClass", "__construct")
		}
		return nil, throw
	}

	// 将类信息存储到当前对象的属性中
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 存储类名到 ObjectValue 的实例属性中
		objCtx.ObjectValue.SetProperty("_className", data.NewStringValue(className))
	}

	return nil, nil
}

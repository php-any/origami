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

// reflectionConstructArg 读取 Reflection*::__construct 实参。
// 优先 FlatCallArgs（调用方求值快照）：并发请求下池化槽位可能变成 *IntValue。
func reflectionConstructArg(ctx data.Context, i int) data.Value {
	if ctx == nil {
		return nil
	}
	if flat := ctx.GetFlatCallArgs(); i >= 0 && i < len(flat) && flat[i] != nil {
		v := unwrapReflectionValue(flat[i])
		if i == 0 {
			if _, ok := classNameFromObjectOrString(v); ok {
				return v
			}
		} else if _, isNull := v.(*data.NullValue); !isNull {
			return v
		}
	}
	v, ok := ctx.GetIndexValue(i)
	if !ok {
		return nil
	}
	return v
}

// classNameFromObjectOrString 对齐 PHP ReflectionClass/ReflectionMethod 第一个参数：
// object|string。GetIndexValue 可能包一层 ZValValue，必须先解开，否则 $this / 实例
// 会落到 “expects parameter 1 to be string or object”，Filament schema 懒加载失败。
func classNameFromObjectOrString(v data.Value) (string, bool) {
	v = unwrapReflectionValue(v)
	if v == nil {
		return "", false
	}
	if av, ok := v.(*data.AnyValue); ok && av != nil {
		if inner, ok := av.Value.(data.Value); ok {
			v = unwrapReflectionValue(inner)
		}
	}
	switch t := v.(type) {
	case *data.ThisValue:
		if t != nil && t.ClassValue != nil && t.Class != nil {
			return t.Class.GetName(), true
		}
	case *data.ClassValue:
		if t != nil && t.Class != nil {
			return t.Class.GetName(), true
		}
	case *data.ClassMethodContext:
		if t != nil && t.ClassValue != nil && t.Class != nil {
			return t.Class.GetName(), true
		}
	case *data.StringValue:
		return t.AsString(), true
	}
	if named, ok := v.(data.GetName); ok {
		if n := named.GetName(); n != "" {
			return n, true
		}
	}
	return "", false
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

var reflectionClassConstructMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "class", 0, nil, data.Mixed{}),
}

// GetParams 返回参数列表
// 参数:
//   - class: 类名（字符串）或对象实例，类型为 Mixed
func (m *ReflectionClassConstructMethod) GetParams() []data.GetValue {
	return reflectionClassConstructMethodGetParams
}

var reflectionClassConstructMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "class", 0, data.Mixed{}),
}

// GetVariables 返回变量列表
func (m *ReflectionClassConstructMethod) GetVariables() []data.Variable {
	return reflectionClassConstructMethodGetVariables
}

// GetReturnType 返回返回类型，构造函数无返回值
func (m *ReflectionClassConstructMethod) GetReturnType() data.Types { return nil }

// Call 执行构造函数
// 从参数中获取类名或对象，加载对应的类，并将类名存储到实例的 _className 属性中
func (m *ReflectionClassConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：类名或对象
	classValue := reflectionConstructArg(ctx, 0)
	if classValue == nil {
		return nil, createReflectionException("Argument #1 ($class) must be of type object|string, null given", ctx, m.GetFrom())
	}

	var className string
	if name, ok := classNameFromObjectOrString(classValue); ok && name != "" {
		className = name
	} else {
		switch classVal := unwrapReflectionValue(classValue).(type) {
		case *data.BoolValue:
			return nil, createReflectionException(fmt.Sprintf("Argument #1 ($class) must be of type object|string, bool given (value: %v)", classVal.Value), ctx, m.GetFrom())
		case *data.NullValue:
			return nil, createReflectionException("Argument #1 ($class) must not be null", ctx, m.GetFrom())
		default:
			typeName := fmt.Sprintf("%T", unwrapReflectionValue(classValue))
			return nil, createReflectionException(fmt.Sprintf("ReflectionClass::__construct(): Argument #1 ($class) must be of type object|string, %s given", typeName), ctx, m.GetFrom())
		}
	}

	// 加载类；失败须抛 ReflectionException，供调用方按 PHP 语义捕获。
	vm := ctx.GetVM()
	stmt, acl := vm.LoadPkg(className)
	if acl != nil {
		if tv, ok := acl.(*data.ThrowValue); ok && tv.Object != nil {
			return nil, acl
		}
		return nil, createReflectionException(
			fmt.Sprintf(`Class "%s" does not exist`, className), ctx, m.GetFrom())
	}
	if stmt == nil {
		return nil, createReflectionException(
			fmt.Sprintf(`Class "%s" does not exist`, className), ctx, m.GetFrom())
	}

	// 将类信息存储到当前对象的属性中
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 存储类名到 ObjectValue 的实例属性中
		setReflectionClassIdentity(objCtx.ClassValue, className)
	}

	return nil, nil
}

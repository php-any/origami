package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ClosureClass 提供 PHP Closure 相关静态方法
type ClosureClass struct {
	node.Node
	bind data.Method
}

func (c *ClosureClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if c.bind == nil {
		c.bind = &ClosureBindMethod{}
	}
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ClosureClass) GetName() string { return "Closure" }

func (c *ClosureClass) GetExtend() *string { return nil }

func (c *ClosureClass) GetImplements() []string { return nil }

func (c *ClosureClass) GetProperty(name string) (data.Property, bool) { return nil, false }

func (c *ClosureClass) GetPropertyList() []data.Property { return nil }

func (c *ClosureClass) GetMethod(name string) (data.Method, bool) { return nil, false }

func (c *ClosureClass) GetMethods() []data.Method { return nil }

// 静态方法
func (c *ClosureClass) GetStaticMethod(name string) (data.Method, bool) {
	if c.bind == nil {
		c.bind = &ClosureBindMethod{}
	}
	switch name {
	case "bind":
		return c.bind, true
	case "fromCallable", "fromcallable":
		return &ClosureFromCallableMethod{}, true
	}
	return nil, false
}

// GetConstruct 无构造函数
func (c *ClosureClass) GetConstruct() data.Method { return nil }

// ClosureBindMethod 实现 Closure::bind
// 当前简化实现：返回原始闭包，不改变绑定对象/作用域
type ClosureBindMethod struct{}

func (m *ClosureBindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	closureVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: closure"))
	}

	newThisVal, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: newThis"))
	}

	newScope, _ := ctx.GetIndexValue(2)

	// 解析绑定的 $this 对象
	var boundThis *data.ClassValue
	switch tv := newThisVal.(type) {
	case *data.ClassValue:
		boundThis = tv
	case *data.ThisValue:
		boundThis = tv.ClassValue
	case *data.StringValue:
		// 静态调用 Closure::bind(..., 'ClassName', ...)
		// $this 为 null（静态绑定），仅改变作用域
	}

	// 仅接受可调用类型
	switch fv := closureVal.(type) {
	case *data.FuncValue:
		scopeClass := ""
		if scopeStr, ok := newScope.(*data.StringValue); ok && scopeStr.Value != "" {
			scopeClass = scopeStr.Value
		}
		// 即使 scopeClass 为空，也保留 $this 绑定
		if boundThis != nil || scopeClass != "" {
			return data.NewBoundFuncValue(fv.Value, scopeClass, boundThis), nil
		}
		return closureVal, nil
	case *data.BoundFuncValue:
		scopeClass := fv.ScopeClass
		if scopeStr, ok := newScope.(*data.StringValue); ok && scopeStr.Value != "" {
			scopeClass = scopeStr.Value
		}
		if boundThis == nil {
			boundThis = fv.BoundObject
		}
		if boundThis != nil || scopeClass != "" {
			return data.NewBoundFuncValue(fv.Value, scopeClass, boundThis), nil
		}
		return closureVal, nil
	default:
		return nil, utils.NewThrow(errors.New("Closure::bind 需要传入闭包/可调用类型"))
	}
}

func (m *ClosureBindMethod) GetName() string { return "bind" }

func (m *ClosureBindMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ClosureBindMethod) GetIsStatic() bool { return true }

func (m *ClosureBindMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "closure", 0, nil, nil),
		node.NewParameter(nil, "newThis", 1, nil, nil),  // 占位，当前未使用
		node.NewParameter(nil, "newScope", 2, nil, nil), // 占位，当前未使用
	}
}

func (m *ClosureBindMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "closure", 0, data.Mixed{}),
		node.NewVariable(nil, "newThis", 1, data.Mixed{}),
		node.NewVariable(nil, "newScope", 2, data.Mixed{}),
	}
}

func (m *ClosureBindMethod) GetReturnType() data.Types { return nil }

// ClosureFromCallableMethod 实现 Closure::fromCallable()。
// 将各类可调用形式（函数名、[对象/类, 方法] 数组、已有闭包等）包装为 Closure。
type ClosureFromCallableMethod struct{}

func (m *ClosureFromCallableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	callable, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("fromCallable() expects parameter 1 to be callable"))
	}
	switch cv := callable.(type) {
	case *data.FuncValue:
		return cv, nil
	case *data.BoundFuncValue:
		return cv, nil
	case *data.ArrayValue:
		list := cv.ToValueList()
		if len(list) < 2 {
			return nil, utils.NewThrow(errors.New("fromCallable(): array callback must have exactly 2 members"))
		}
		methodName := list[1].AsString()
		// 对象实例方法
		if obj, ok := list[0].(*data.ClassValue); ok {
			return data.NewFuncValue(node.NewObjectMethodCallable(obj, methodName)), nil
		}
		if tv, ok := list[0].(*data.ThisValue); ok && tv.ClassValue != nil {
			return data.NewFuncValue(node.NewObjectMethodCallable(tv.ClassValue, methodName)), nil
		}
		// 静态/类名方法
		className := list[0].AsString()
		stmt, acl := ctx.GetVM().GetOrLoadClass(className)
		if acl != nil {
			return nil, acl
		}
		method, ok := stmt.GetMethod(methodName)
		if !ok {
			if sm, ok2 := stmt.(data.GetStaticMethod); ok2 {
				method, ok = sm.GetStaticMethod(methodName)
			}
		}
		if !ok {
			return nil, utils.NewThrow(errors.New("fromCallable(): 未找到方法 " + className + "::" + methodName))
		}
		fn, acl := node.NewStaticMethodFuncValue(stmt, method).GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		return fn, nil
	default:
		// 字符串函数名
		if str, ok := callable.(data.AsString); ok {
			name := str.AsString()
			fn, ok := ctx.GetVM().GetFunc(name)
			if !ok {
				return nil, utils.NewThrow(errors.New("fromCallable(): function " + name + " does not exist"))
			}
			return data.NewFuncValue(fn), nil
		}
		return nil, utils.NewThrow(errors.New("fromCallable(): 不可调用类型"))
	}
}

func (m *ClosureFromCallableMethod) GetName() string { return "fromCallable" }

func (m *ClosureFromCallableMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *ClosureFromCallableMethod) GetIsStatic() bool { return true }

func (m *ClosureFromCallableMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "callable", 0, nil, nil)}
}

func (m *ClosureFromCallableMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "callable", 0, data.Mixed{})}
}

func (m *ClosureFromCallableMethod) GetReturnType() data.Types { return nil }

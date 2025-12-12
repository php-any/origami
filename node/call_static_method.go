package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CallStaticMethod 表示对象属性访问表达式
type CallStaticMethod struct {
	*Node  `pp:"-"`
	stmt   data.GetValue // 类名称 Class::fn() or Class::test::one
	Method string        // 函数名
}

func NewCallStaticMethod(from *TokenFrom, path data.GetValue, method string) *CallStaticMethod {
	return &CallStaticMethod{
		Node:   NewNode(from),
		stmt:   path,
		Method: method,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var method data.Method
	var classStmt data.ClassStmt
	var has bool

	switch expr := pe.stmt.(type) {
	case data.GetStaticMethod:
		method, has = expr.GetStaticMethod(pe.Method)
		if has {
			if cls, ok := expr.(data.ClassStmt); ok {
				classStmt = cls
			}
		} else {
			return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法调用函数(%s)。", pe.Method))
		}
	default:
		c, acl := pe.stmt.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		switch expr := c.(type) {
		case data.GetStaticMethod:
			method, has = expr.GetStaticMethod(pe.Method)
			if has {
				if cls, ok := expr.(data.ClassStmt); ok {
					classStmt = cls
				}
			}
		case data.GetMethod:
			method, has = expr.GetMethod(pe.Method)
			if has {
				// 实例方法，直接返回 FuncValue
				return data.NewFuncValue(method), nil
			}
		}
	}

	if !has {
		name := ""
		if getName, ok := pe.stmt.(data.ClassStmt); ok {
			name = getName.GetName()
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("(%v)无法调用函数(%s)。", name, pe.Method))
	}

	// 静态方法需要 ClassMethodContext，返回包装器让 CallMethod 正确处理
	if classStmt != nil {
		return NewStaticMethodFuncValue(classStmt, method), nil
	}

	// 如果没有类信息，直接返回 FuncValue（向后兼容）
	return data.NewFuncValue(method), nil
}

// CallStaticMethodLater 延迟的静态方法调用（类未加载时）
type CallStaticMethodLater struct {
	*Node
	className string // 类名（字符串形式）
	method    string // 方法名
	namespace string // 命名空间
}

// NewCallStaticMethodLater 创建延迟的静态方法调用
func NewCallStaticMethodLater(from *TokenFrom, className, method, namespace string) *CallStaticMethodLater {
	return &CallStaticMethodLater{
		Node:      NewNode(from),
		className: className,
		method:    method,
		namespace: namespace,
	}
}

// GetValue 获取延迟静态方法调用的值
func (pe *CallStaticMethodLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 尝试加载类
	stmt, acl := ctx.GetVM().GetOrLoadClass(pe.className)
	if acl != nil {
		return nil, acl
	}
	if stmt == nil {
		// 如果还是找不到，尝试使用命名空间
		fullClassName := pe.className
		if pe.namespace != "" {
			fullClassName = pe.namespace + "\\" + pe.className
		}
		stmt, acl = ctx.GetVM().GetOrLoadClass(fullClassName)
		if acl != nil {
			return nil, acl
		}
		if stmt == nil {
			return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法调用静态方法(%s::%s), 未找到类", pe.className, pe.method))
		}
	}

	// 创建实际的静态方法调用
	tokenFrom, ok := pe.GetFrom().(*TokenFrom)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("无法获取TokenFrom信息"))
	}
	callStaticMethod := NewCallStaticMethod(tokenFrom, stmt, pe.method)

	return callStaticMethod.GetValue(ctx)
}

func NewStaticMethodFuncValue(class data.ClassStmt, method data.Method) *StaticMethodFuncValue {
	return &StaticMethodFuncValue{
		class:  class,
		method: method,
	}
}

// StaticMethodFuncValue 静态方法函数值包装器，确保调用时使用 ClassMethodContext
type StaticMethodFuncValue struct {
	class  data.ClassStmt
	method data.Method
}

// staticMethodFunc 适配器：将 data.Method 包装为 data.FuncStmt，并在调用时切换到 ClassMethodContext
type staticMethodFunc struct {
	class  data.ClassStmt
	method data.Method
}

func (s *staticMethodFunc) GetName() string               { return s.method.GetName() }
func (s *staticMethodFunc) GetParams() []data.GetValue    { return s.method.GetParams() }
func (s *staticMethodFunc) GetVariables() []data.Variable { return s.method.GetVariables() }
func (s *staticMethodFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	// 创建类方法上下文，使用传入的 callCtx（包含已设置的参数），绑定当前类，保证 self:: 可用
	classValue := data.NewClassValue(s.class, callCtx)
	methodCtx := &data.ClassMethodContext{ClassValue: classValue}
	return s.method.Call(methodCtx)
}

func (s *StaticMethodFuncValue) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 返回 FuncValue，但内部使用 staticMethodFunc 包装，确保调用时使用 ClassMethodContext
	return data.NewFuncValue(&staticMethodFunc{class: s.class, method: s.method}), nil
}

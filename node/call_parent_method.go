package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallParentMethod 表示父类方法调用表达式
type CallParentMethod struct {
	*Node  `pp:"-"`
	Method string // 方法名
}

func NewCallParentMethod(from data.From, method string) *CallParentMethod {
	return &CallParentMethod{
		Node:   NewNode(from),
		Method: method,
	}
}

// GetValue 获取父类方法调用表达式的值
func (pe *CallParentMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("parent:: 只能在类方法中使用"))
	}

	// 获取父类
	if classCtx.Class.GetExtend() == nil {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("当前类没有父类"))
	}

	parentClassName := *classCtx.Class.GetExtend()
	vm := ctx.GetVM()
	parentClass, acl := vm.GetOrLoadClass(parentClassName)
	if acl != nil {
		return nil, acl
	}

	// 获取父类方法
	method, has := parentClass.GetMethod(pe.Method)
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("父类 %s 没有方法 %s", parentClassName, pe.Method))
	}

	// 检查方法访问权限
	if method.GetModifier() == data.ModifierPrivate {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("父类方法 %s 是私有的，无法访问", pe.Method))
	}

	return &ChangeCtxAndCallFuncValue{ctx: classCtx, fun: method}, nil
}

// ChangeCtxAndCallFuncValue parent::A 后内部再调用 parent::B 上下文没有切换导致死循环
type ChangeCtxAndCallFuncValue struct {
	ctx *data.ClassMethodContext
	fun data.Method
}

// parentMethodFunc 适配器：将 data.Method 包装为 data.FuncStmt，并在调用时切换到父类上下文
type parentMethodFunc struct {
	baseCtx *data.ClassMethodContext
	method  data.Method
}

func (p *parentMethodFunc) GetName() string               { return p.method.GetName() }
func (p *parentMethodFunc) GetParams() []data.GetValue    { return p.method.GetParams() }
func (p *parentMethodFunc) GetVariables() []data.Variable { return p.method.GetVariables() }
func (p *parentMethodFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	// 取父类
	if p.baseCtx == nil || p.baseCtx.Class == nil || p.baseCtx.Class.GetExtend() == nil {
		return nil, data.NewErrorThrow(nil, errors.New("当前类没有父类"))
	}

	parentName := *p.baseCtx.Class.GetExtend()
	vm := callCtx.GetVM()
	parentClass, acl := vm.GetOrLoadClass(parentName)
	if acl != nil {
		return nil, acl
	}

	// 使用调用时创建好的变量上下文(callCtx)，但切换到父类的方法上下文，保持同一个对象实例
	newCtx := &data.ClassMethodContext{ClassValue: &data.ClassValue{
		ObjectValue: p.baseCtx.ObjectValue,
		Class:       parentClass,
		Context:     callCtx,
	}}

	return p.method.Call(newCtx)
}

func (c *ChangeCtxAndCallFuncValue) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewFuncValue(&parentMethodFunc{baseCtx: c.ctx, method: c.fun}).Call(ctx)
}

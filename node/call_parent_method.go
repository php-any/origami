package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallParentMethod 表示父类方法调用表达式
type CallParentMethod struct {
	*Node        `pp:"-"`
	Method       string // 方法名
	CurrentClass string
	Arguments    []data.GetValue
}

func NewCallParentMethod(from data.From, currentClass, method string, args []data.GetValue) *CallParentMethod {
	return &CallParentMethod{
		Node:         NewNode(from),
		CurrentClass: currentClass,
		Method:       method,
		Arguments:    args,
	}
}

// GetValue 获取父类方法调用表达式的值
func (pe *CallParentMethod) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("parent:: 只能在类方法中使用"))
	}
	class, ok := ctx.GetVM().GetClass(pe.CurrentClass)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("parent:: 只能在类方法中使用"))
	}

	// 获取父类
	if class.GetExtend() == nil {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New("当前类没有父类"))
	}

	parentClassName := *class.GetExtend()
	vm := ctx.GetVM()
	parentClass, acl := vm.GetOrLoadClass(parentClassName)
	if acl != nil {
		return nil, acl
	}

	// 获取父类方法：需要沿继承链向上查找（父类本身或其父类中定义的方法）
	var (
		method data.Method
		has    bool
	)
	current := parentClass
	for current != nil {
		method, has = current.GetMethod(pe.Method)
		if has {
			break
		}
		if current.GetExtend() == nil {
			break
		}
		nextName := *current.GetExtend()
		next, acl := vm.GetOrLoadClass(nextName)
		if acl != nil {
			return nil, acl
		}
		current = next
	}
	if !has {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("父类 %s 及其继承链中都没有方法 %s", parentClassName, pe.Method))
	}

	// 检查方法访问权限
	if method.GetModifier() == data.ModifierPrivate {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("父类方法 %s 是私有的，无法访问", pe.Method))
	}

	temp := &CallObjectMethod{
		Node:   pe.Node,
		Object: classCtx.ClassValue,
		Args:   pe.Arguments,
	}

	newCtx, acl := temp.callMethodParams(classCtx.ClassValue, ctx, method)
	if acl != nil {
		return nil, acl
	}

	return method.Call(newCtx)
}

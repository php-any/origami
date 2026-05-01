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
	// 检查是否在类上下文中（类方法或类级初始化器）
	var class data.ClassStmt
	var object *data.ClassValue
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		object = classCtx.ClassValue
		// SelfClass 优先：由外层 parent:: 调用设置，表示代码定义所在类
		if classCtx.SelfClass != nil {
			class = classCtx.SelfClass
		} else if pe.CurrentClass != "" {
			if cls, has := ctx.GetVM().GetClass(pe.CurrentClass); has && cls.GetExtend() != nil {
				class = cls
			} else {
				class = classCtx.Class
			}
		} else {
			class = classCtx.Class
		}
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		object = classVal
		if pe.CurrentClass != "" {
			if cls, has := ctx.GetVM().GetClass(pe.CurrentClass); has && cls.GetExtend() != nil {
				class = cls
			} else {
				class = classVal.Class
			}
		} else {
			class = classVal.Class
		}
	} else {
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
		method     data.Method
		foundClass data.ClassStmt // 实际定义方法的类
		has        bool
	)
	current := parentClass
	for current != nil {
		method, has = current.GetMethod(pe.Method)
		if !has {
			if gsm, ok := current.(data.GetStaticMethod); ok {
				method, has = gsm.GetStaticMethod(pe.Method)
			}
		}
		if has {
			foundClass = current
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

	if method.GetModifier() == data.ModifierPrivate {
		return nil, data.NewErrorThrow(pe.GetFrom(), fmt.Errorf("父类方法 %s 是私有的，无法访问", pe.Method))
	}

	temp := &CallObjectMethod{
		Node:   pe.Node,
		Object: object,
		Args:   pe.Arguments,
	}

	newCtx, acl := temp.callMethodParams(object, ctx, method)
	if acl != nil {
		return nil, acl
	}

	// 通过 SelfClass 告知被调用方法其代码定义所在的类。
	// 这样方法内的 self:: / parent:: 会从 foundClass 开始解析，
	// 而非从运行时子类开始，从而正确处理 trait 合并场景。
	if cmc, ok := newCtx.(*data.ClassMethodContext); ok {
		cmc.SelfClass = foundClass
		if cmc.StaticClass == nil {
			cmc.StaticClass = cmc.ClassValue.Class
		}
	}

	return method.Call(newCtx)
}

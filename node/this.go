package node

import (
	"errors"
	"github.com/php-any/origami/data"
)

type This struct {
	*Node `pp:"-"`
}

func NewThis(from data.From) *This {
	return &This{
		Node: NewNode(from),
	}
}

func (u *This) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// PHP / Livewire：
	// - 类方法体内 $this 永远是方法接收者（ClassMethodContext），即使 Context 链上
	//   挂着外层 Closure::bind（否则 Collection 方法会错成 Login）。
	// - 视图闭包由 Closure::bind 注入时，顶层是 BoundContext，应用 BoundThis
	//   （嵌套 Livewire 组件视图）。
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok && classCtx.ObjectValue != nil {
		return data.NewThisValue(classCtx.ClassValue), nil
	}
	if bc, ok := ctx.(*data.BoundContext); ok && bc.BoundThis != nil {
		return data.NewThisValue(bc.BoundThis), nil
	}
	// 沿 BoundContext 链查找；一旦遇到 ClassMethodContext 即用其接收者，不继续往里挖 BoundThis
	for c := ctx; c != nil; {
		switch v := c.(type) {
		case *data.BoundContext:
			if v.BoundThis != nil {
				return data.NewThisValue(v.BoundThis), nil
			}
			c = v.Context
		case *data.ClassMethodContext:
			if v.ObjectValue != nil {
				return data.NewThisValue(v.ClassValue), nil
			}
			return nil, data.NewErrorThrow(u.from, errors.New("this关键字只能在类中使用"))
		case *data.ClassValue:
			c = v.Context
		default:
			return nil, data.NewErrorThrow(u.from, errors.New("this关键字只能在类中使用"))
		}
	}
	return nil, data.NewErrorThrow(u.from, errors.New("this关键字只能在类中使用"))
}

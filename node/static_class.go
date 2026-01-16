package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// StaticClass 表示 static::class 表达式
// 在 PHP 中，static::class 使用 late static binding，返回实际调用时的类名
// 当前实现暂时返回当前类的类名，后续可以增强为真正的 late static binding
type StaticClass struct {
	*Node `pp:"-"`
}

// NewStaticClass 创建一个新的 static::class 表达式节点
func NewStaticClass(from data.From) *StaticClass {
	return &StaticClass{
		Node: NewNode(from),
	}
}

// GetValue 获取 static::class 的值（当前类的类名，或实际调用时的类名）
func (s *StaticClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(s.from, errors.New("static::class 只能在类方法中使用"))
	}

	// 获取当前类的类名
	parent, ok := classCtx.Context.(*data.ClassMethodContext)
	for ok {
		classCtx = parent
		parent, ok = classCtx.Context.(*data.ClassMethodContext)
	}
	currentClass := classCtx.Class
	className := currentClass.GetName()

	return data.NewStringValue(className), nil
}

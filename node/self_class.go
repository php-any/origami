package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// SelfClass 表示 self::class 表达式
type SelfClass struct {
	*Node `pp:"-"`
}

// NewSelfClass 创建一个新的 self::class 表达式节点
func NewSelfClass(from data.From) *SelfClass {
	return &SelfClass{
		Node: NewNode(from),
	}
}

// GetValue 获取 self::class 的值（当前类的类名）
// 注意：在 PHP 中，self::class 应该返回定义该方法的类名
// 当前实现返回调用时的类名，这在大多数情况下是正确的
// 只有在子类调用父类方法时，才会出现差异（应该返回父类名，但当前返回子类名）
func (s *SelfClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类方法上下文中
	classCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok {
		return nil, data.NewErrorThrow(s.from, errors.New("self::class 只能在类方法中使用"))
	}

	// 获取当前类的类名
	// TODO: 实现真正的语义，返回定义该方法的类名（需要从调用栈获取方法名）
	currentClass := classCtx.Class
	className := currentClass.GetName()

	return data.NewStringValue(className), nil
}

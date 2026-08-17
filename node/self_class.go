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

// GetValue 获取 self::class 的值（代码定义所在类的类名）
// PHP 语义：self::class 是编译期（词法）绑定，应返回"定义该方法的类"，
// 而非调用时的（后期绑定）类。例如父类/ trait 中定义的方法被子类实例调用时，
// self::class 仍应返回父类/trait 宿主类，而不是子类。
// 类方法上下文通过 SelfClass 字段记录代码定义所在类，优先使用它；
// 若未设置（如类级初始化器），回退到当前上下文类。
func (s *SelfClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否在类上下文中（类方法或类级初始化器）
	var currentClass data.ClassStmt
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// SelfClass 为代码定义所在类（trait/父类），优先使用，符合 self::class 词法绑定语义
		if classCtx.SelfClass != nil {
			currentClass = classCtx.SelfClass
		} else {
			currentClass = classCtx.Class
		}
	} else if classVal, ok := ctx.(*data.ClassValue); ok {
		currentClass = classVal.Class
	} else {
		return nil, data.NewErrorThrow(s.from, errors.New("self::class 只能在类方法中使用"))
	}

	className := currentClass.GetName()

	return data.NewStringValue(className), nil
}

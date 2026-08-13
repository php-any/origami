package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CloneExpression 表示 clone 表达式
// 语法：clone <expression>
type CloneExpression struct {
	*Node  `pp:"-"`
	Target data.GetValue
}

// NewCloneExpression 创建一个新的 clone 表达式节点
func NewCloneExpression(from *TokenFrom, target data.GetValue) *CloneExpression {
	return &CloneExpression{
		Node:   NewNode(from),
		Target: target,
	}
}

// GetValue 实现 Value 接口
// 行为与 PHP 类似：
// - 只允许对对象进行 clone
// - 生成一个新对象，复制实例属性
// - 如果存在 __clone 方法，则在新对象上调用该方法
func (n *CloneExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 先计算被克隆的目标值
	value, acl := n.Target.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	if value == nil {
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("clone 关键字后面的表达式不能为 null"))
	}

	// $this 在运行时为 ThisValue，需解包为 ClassValue 再克隆
	if thisVal, ok := value.(*data.ThisValue); ok {
		value = thisVal.ClassValue
	}

	// 仅支持对象克隆
	obj, ok := value.(*data.ClassValue)
	if !ok {
		if v, ok2 := value.(data.Value); ok2 {
			return nil, data.NewErrorThrow(n.from, fmt.Errorf("clone 关键字只能用于对象, 当前类型: %T", v))
		}
		return nil, data.NewErrorThrow(n.from, fmt.Errorf("clone 关键字只能用于对象"))
	}

	// 创建新实例，保持同一个类与上下文
	cloned := data.NewClassValue(obj.Class, obj.Context)

	// 复制实例属性（浅拷贝属性值，符合 PHP 克隆语义）
	obj.RangeProperties(func(key string, v data.Value) bool {
		cloned.SetProperty(key, v)
		return true
	})

	// 如果类定义了 __clone 方法，则在新对象上调用它
	if method, ok := cloned.GetMethod("__clone"); ok && method != nil {
		varies := method.GetVariables()
		fnCtx := cloned.CreateContext(varies)
		// 记录调用参数（无参）
		fnCtx.SetCallArgs([]data.GetValue{})

		_, acl = method.Call(fnCtx)
		if acl != nil {
			if throwValue, ok := acl.(*data.ThrowValue); ok {
				throwValue.AddStackWithInfo(n.from, cloned.Class.GetName(), "__clone")
			}
			return nil, acl
		}
	}

	return cloned, nil
}

package node

import (
	"github.com/php-any/origami/data"
)

// InstanceOfExpression 表示 instanceof 表达式
type InstanceOfExpression struct {
	*Node
	Object    data.GetValue // 对象表达式
	ClassName string        // 类名
}

// NewInstanceOfExpression 创建一个新的 instanceof 表达式
func NewInstanceOfExpression(from data.From, object data.GetValue, className string) *InstanceOfExpression {
	return &InstanceOfExpression{
		Node:      NewNode(from),
		Object:    object,
		ClassName: className,
	}
}

// GetValue 获取 instanceof 表达式的值
func (i *InstanceOfExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算对象表达式的值
	objectValue, c := i.Object.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 检查对象值是否为类实例
	if classValue, ok := objectValue.(*data.ClassValue); ok {
		// 使用 checkClassIs 函数检查类关系
		checkC, ok := ctx.GetVM().GetClass(i.ClassName)
		if ok {
			result := checkClassIs(ctx, classValue.Class, checkC.GetName())
			return data.NewBoolValue(result), nil
		}
		checkI, ok := ctx.GetVM().GetInterface(i.ClassName)
		if ok {
			result := checkClassIs(ctx, classValue.Class, checkI.GetName())
			return data.NewBoolValue(result), nil
		}
	}

	// 如果不是类实例，返回 false
	return data.NewBoolValue(false), nil
}

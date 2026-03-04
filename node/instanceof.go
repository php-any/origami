package node

import (
	"github.com/php-any/origami/data"
)

// InstanceOfExpression 表示 instanceof 表达式
type InstanceOfExpression struct {
	*Node
	Object    data.GetValue // 对象表达式
	ClassName data.GetValue // 类名
}

// NewInstanceOfExpression 创建一个新的 instanceof 表达式
func NewInstanceOfExpression(from data.From, object data.GetValue, className data.GetValue) *InstanceOfExpression {
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

	switch right := i.ClassName.(type) {
	case *StringLiteral:
		class := right.Value
		return instanceof(ctx, class, objectValue)
	case *VariableExpression:
		r, acl := right.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		class := r.(data.AsString).AsString()
		return instanceof(ctx, class, objectValue)
	}

	// 如果不是类实例，返回 false
	return data.NewBoolValue(false), nil
}

func instanceof(ctx data.Context, class string, objectValue data.GetValue) (data.GetValue, data.Control) {
	// 检查对象值是否为类实例
	if classValue, ok := objectValue.(*data.ClassValue); ok {
		c, acl := ctx.GetVM().LoadPkg(class)
		if acl != nil {
			return nil, acl
		}
		if c != nil {
			switch checkC := c.(type) {
			case data.ClassStmt:
				result, acl := checkClassIs(ctx, classValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			case data.InterfaceStmt:
				result, acl := checkClassIs(ctx, classValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			}
		}
	} else if thisValue, ok := objectValue.(*data.ThisValue); ok {
		// 处理 ThisValue（$this）
		c, acl := ctx.GetVM().LoadPkg(class)
		if acl != nil {
			return nil, acl
		}
		if c != nil {
			switch checkC := c.(type) {
			case data.ClassStmt:
				result, acl := checkClassIs(ctx, thisValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			case data.InterfaceStmt:
				result, acl := checkClassIs(ctx, thisValue.Class, checkC.GetName())
				return data.NewBoolValue(result), acl
			}
		}
	}

	switch class {
	case "object":
		switch objectValue.(type) {
		case *data.ClassValue:
			return data.NewBoolValue(true), nil
		case *data.ObjectValue:
			return data.NewBoolValue(true), nil
		}
	case "Closure", "closure":
		switch objectValue.(type) {
		case *data.FuncValue:
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

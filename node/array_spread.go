package node

import "github.com/php-any/origami/data"

// ArraySpread 表示数组展开运算符 ...$array
type ArraySpread struct {
	Expr data.GetValue // 要展开的表达式
}

// NewArraySpread 创建一个数组展开节点
func NewArraySpread(expr data.GetValue) data.GetValue {
	return &ArraySpread{
		Expr: expr,
	}
}

// GetValue 获取展开后的数组元素
func (a *ArraySpread) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取要展开的表达式值
	exprValue, ctl := a.Expr.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	if exprValue == nil {
		return nil, data.NewErrorThrow(nil, data.NewError(nil, "展开运算符的操作数不能为 null", nil))
	}

	// 检查是否是数组
	if arrayValue, ok := exprValue.(*data.ArrayValue); ok {
		// 返回数组的所有元素
		return arrayValue, nil
	}

	// 如果不是数组，返回错误
	return nil, data.NewErrorThrow(nil, data.NewError(nil, "展开运算符只能用于数组", nil))
}

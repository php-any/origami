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

// GetValue 获取展开后的数组元素。
// 关联数组在 Origami 中可能是 ObjectValue，展开时转换为带键名的 ArrayValue，以保留键。
func (a *ArraySpread) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	exprValue, ctl := a.Expr.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	if exprValue == nil {
		return nil, data.NewErrorThrow(nil, data.NewError(nil, "展开运算符的操作数不能为 null", nil))
	}

	switch v := exprValue.(type) {
	case *data.ArrayValue:
		return v, nil
	case *data.ObjectValue:
		list := make([]*data.ZVal, 0)
		v.RangeProperties(func(key string, value data.Value) bool {
			list = append(list, data.NewNamedZVal(key, value))
			return true
		})
		return &data.ArrayValue{List: list}, nil
	default:
		return nil, data.NewErrorThrow(nil, data.NewError(nil, "展开运算符只能用于数组", nil))
	}
}

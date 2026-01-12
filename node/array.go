package node

import (
	"github.com/php-any/origami/data"
)

type Array struct {
	*Node `pp:"-"`
	V     []data.GetValue
}

func NewArray(token *TokenFrom, arr []data.GetValue) data.GetValue {
	return &Array{
		Node: NewNode(token),
		V:    arr,
	}
}

// GetValue 获取数字字面量的值
func (n *Array) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	arr := []data.Value{}
	for _, statement := range n.V {
		v, acl := statement.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}

		// 检查是否是展开运算符
		if spread, ok := statement.(*ArraySpread); ok {
			// 展开数组元素
			spreadValue, acl := spread.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if arrayValue, ok := spreadValue.(*data.ArrayValue); ok {
				// 将展开的数组元素添加到结果数组中
				arr = append(arr, arrayValue.Value...)
			} else {
				return nil, data.NewErrorThrow(n.from, data.NewError(n.from, "展开运算符只能用于数组", nil))
			}
		} else {
			// 普通元素
			arr = append(arr, v.(data.Value))
		}
	}
	return data.NewArrayValue(arr), nil
}

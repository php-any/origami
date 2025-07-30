package node

import (
	"github.com/php-any/origami/data"
	"strconv"
)

type IntLiteral struct {
	*Node `pp:"-"`
	V     data.Value
}

func NewIntLiteral(token *TokenFrom, str string) data.GetValue {
	// 将字符串转换为 float64
	i, err := strconv.Atoi(str)
	if err != nil {
		// 如果转换失败，返回 0
		return &IntLiteral{
			Node: NewNode(token),
			V:    data.NewIntValue(0),
		}
	}

	return &IntLiteral{
		Node: NewNode(token),
		V:    data.NewIntValue(i),
	}
}

// GetValue 获取数字字面量的值
func (n *IntLiteral) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return n.V, nil
}

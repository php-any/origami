package node

import (
	"github.com/php-any/origami/data"
	"strconv"
)

type FloatLiteral struct {
	*Node `pp:"-"`
	V     data.Value
}

func NewFloatLiteral(token *TokenFrom, str string) data.GetValue {
	// 将字符串转换为 float64
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		// 如果转换失败，返回 0
		return &FloatLiteral{
			Node: NewNode(token),
			V:    data.NewFloatValue(0),
		}
	}

	return &FloatLiteral{
		Node: NewNode(token),
		V:    data.NewFloatValue(i),
	}
}

// GetValue 获取数字字面量的值
func (n *FloatLiteral) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return n.V, nil
}

package node

import (
	"strconv"

	"github.com/php-any/origami/data"
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

// NewNumberLiteral 解析复杂数字面量（十六进制、二进制、八进制、科学计数法）
func NewNumberLiteral(token *TokenFrom, str string) data.GetValue {
	// 检查是否是科学计数法（包含 e 或 E）
	for i := 0; i < len(str); i++ {
		if str[i] == 'e' || str[i] == 'E' {
			// 科学计数法，解析为浮点数
			f, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return &FloatLiteral{
					Node: NewNode(token),
					V:    data.NewFloatValue(0),
				}
			}
			return &FloatLiteral{
				Node: NewNode(token),
				V:    data.NewFloatValue(f),
			}
		}
	}

	// 检查是否是十六进制（0x 或 0X 开头）
	if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
		// 去掉 0x 前缀
		hexStr := str[2:]
		i, err := strconv.ParseInt(hexStr, 16, 64)
		if err != nil {
			return &IntLiteral{
				Node: NewNode(token),
				V:    data.NewIntValue(0),
			}
		}
		return &IntLiteral{
			Node: NewNode(token),
			V:    data.NewIntValue(int(i)),
		}
	}

	// 检查是否是二进制（0b 或 0B 开头）
	if len(str) > 2 && str[0] == '0' && (str[1] == 'b' || str[1] == 'B') {
		// 去掉 0b 前缀
		binStr := str[2:]
		i, err := strconv.ParseInt(binStr, 2, 64)
		if err != nil {
			return &IntLiteral{
				Node: NewNode(token),
				V:    data.NewIntValue(0),
			}
		}
		return &IntLiteral{
			Node: NewNode(token),
			V:    data.NewIntValue(int(i)),
		}
	}

	// 检查是否是八进制（0 开头，但不是 0x 或 0b）
	if len(str) > 1 && str[0] == '0' {
		// 去掉前导 0
		octStr := str[1:]
		i, err := strconv.ParseInt(octStr, 8, 64)
		if err != nil {
			return &IntLiteral{
				Node: NewNode(token),
				V:    data.NewIntValue(0),
			}
		}
		return &IntLiteral{
			Node: NewNode(token),
			V:    data.NewIntValue(int(i)),
		}
	}

	// 默认尝试解析为整数
	i, err := strconv.Atoi(str)
	if err != nil {
		// 如果解析整数失败，尝试解析为浮点数
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return &IntLiteral{
				Node: NewNode(token),
				V:    data.NewIntValue(0),
			}
		}
		return &FloatLiteral{
			Node: NewNode(token),
			V:    data.NewFloatValue(f),
		}
	}

	return &IntLiteral{
		Node: NewNode(token),
		V:    data.NewIntValue(i),
	}
}

package node

import "github.com/php-any/origami/data"

// SpreadArgument 表示调用实参中的展开参数 ...expr
// 例如: t(...$arr)
type SpreadArgument struct {
	*Node `pp:"-"`
	Expr  data.GetValue
}

// NewSpreadArgument 创建一个新的展开参数节点
func NewSpreadArgument(from data.From, expr data.GetValue) *SpreadArgument {
	return &SpreadArgument{
		Node: NewNode(from),
		Expr: expr,
	}
}

// GetValue 默认直接转发内部表达式的值
// 具体“展开”语义由调用处按需处理
func (s *SpreadArgument) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return s.Expr.GetValue(ctx)
}


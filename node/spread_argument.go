package node

import (
	"github.com/php-any/origami/data"
)

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
	if s.Expr == nil {
		// PHP 8.1 first-class callable：method(...) 的占位展开，不是可抛出的错误
		return nil, ToClosure{}
	}
	return s.Expr.GetValue(ctx)
}

// ToClosure 标记「一等可调用」语法，必须实现完整 Control/Value，
// 不可嵌入 nil 的 data.Control（否则 ShowControl/AsString 会空指针 panic）。
type ToClosure struct{}

func (ToClosure) AsString() string { return "ToClosure" }

func (ToClosure) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return nil, ToClosure{}
}

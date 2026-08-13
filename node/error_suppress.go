package node

import "github.com/php-any/origami/data"

// ErrorSuppress 表示 PHP 的 @ 错误抑制：求值内层表达式，若抛出则抑制并返回 null
type ErrorSuppress struct {
	*Node `pp:"-"`
	Inner data.GetValue
}

// NewErrorSuppress 创建 @expr 节点
func NewErrorSuppress(from data.From, inner data.GetValue) *ErrorSuppress {
	return &ErrorSuppress{
		Node:  NewNode(from),
		Inner: inner,
	}
}

// GetValue 执行内层表达式；若内层返回 Control（如 ErrorThrow）则抑制并返回 null
func (e *ErrorSuppress) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	val, ctrl := e.Inner.GetValue(ctx)
	if ctrl != nil {
		return data.NewNullValue(), nil
	}
	return val, nil
}

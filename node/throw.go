package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// ThrowStatement 表示throw语句
type ThrowStatement struct {
	*Node `pp:"-"`
	Value data.GetValue
}

// NewThrowStatement 创建一个新的throw语句
func NewThrowStatement(from *TokenFrom, value data.GetValue) *ThrowStatement {
	return &ThrowStatement{
		Node:  NewNode(from),
		Value: value,
	}
}

// GetValue 获取throw语句的值
func (t *ThrowStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取要抛出的值
	v, ctl := t.Value.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	if obj, ok := v.(*data.ClassValue); ok {
		return nil, data.NewErrorThrowFromClassValue(t.from, obj)
	}

	return nil, data.NewErrorThrow(t.from, errors.New(v.(data.Value).AsString()))
}

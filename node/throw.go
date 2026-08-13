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

	switch obj := v.(type) {
	case *data.ClassValue:
		return nil, data.NewErrorThrowFromClassValue(t.from, obj)
	case *data.ThisValue:
		// return $this / 方法链式返回 $this 后再 throw
		return nil, data.NewErrorThrowFromClassValue(t.from, obj.ClassValue)
	case *data.ThrowValue:
		return nil, obj
	}

	return nil, data.NewErrorThrow(t.from, errors.New(v.(data.Value).AsString()))
}

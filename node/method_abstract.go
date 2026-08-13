package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// AbstractMethod 表示抽象方法
type AbstractMethod struct {
	*ClassMethod
}

// NewAbstractMethod 创建一个新的抽象方法
func NewAbstractMethod(method *ClassMethod) *AbstractMethod {
	return &AbstractMethod{
		ClassMethod: method,
	}
}

// Call 调用抽象方法
// 抽象方法不能被直接调用
func (m *AbstractMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, data.NewErrorThrow(m.GetFrom(), errors.New("抽象方法不能被直接调用"))
}

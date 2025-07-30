package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

type Parent struct {
	*Node `pp:"-"`
}

func NewParent(from data.From) *Parent {
	return &Parent{
		Node: NewNode(from),
	}
}

func (p *Parent) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return classCtx.ClassValue, nil
	}
	return nil, data.NewErrorThrow(p.from, errors.New("parent 只能在类作用域中使用"))
}

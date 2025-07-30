package node

import (
	"errors"
	"github.com/php-any/origami/data"
)

type This struct {
	*Node `pp:"-"`
}

func NewThis(from data.From) *This {
	return &This{
		Node: NewNode(from),
	}
}

func (u *This) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return data.NewThisValue(classCtx.ClassValue), nil
	}
	return nil, data.NewErrorThrow(u.from, errors.New("this关键字只能在类中使用"))
}

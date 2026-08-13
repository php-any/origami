package node

import "github.com/php-any/origami/data"

func (u *Todo) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

type Todo struct {
	*Node `pp:"-"`
	Value data.GetValue
}

func NewTodo() *Todo {
	return &Todo{}
}

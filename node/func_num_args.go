package node

import (
	"github.com/php-any/origami/data"
)

// FuncNumArgs 表示 func_num_args 关键字表达式
type FuncNumArgs struct {
	*Node `pp:"-"`
}

func NewFuncNumArgs(from data.From) data.GetValue {
	return &FuncNumArgs{
		Node: NewNode(from),
	}
}

func (f *FuncNumArgs) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	values, acl := expandCallArgsValues(ctx)
	if acl != nil {
		return nil, acl
	}
	return data.NewIntValue(len(values)), nil
}

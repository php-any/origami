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
	args := ctx.GetCallArgs()
	if args == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(args)), nil
}

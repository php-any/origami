package node

import (
	"github.com/php-any/origami/data"
)

// FuncNumArgs 表示 func_num_args 关键字表达式
type FuncNumArgs struct {
	*Node `pp:"-"`
}

// NewFuncNumArgs 创建一个新的 func_num_args 表达式
func NewFuncNumArgs(from data.From) data.GetValue {
	return &FuncNumArgs{
		Node: NewNode(from),
	}
}

// GetValue 获取 func_num_args 的值（返回函数参数数量）
func (f *FuncNumArgs) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	args := ctx.GetCallArgs()
	if args == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(args)), nil
}

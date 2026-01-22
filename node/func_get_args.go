package node

import (
	"github.com/php-any/origami/data"
)

// FuncGetArgs 表示 func_get_args 关键字表达式
type FuncGetArgs struct {
	*Node `pp:"-"`
}

// NewFuncGetArgs 创建一个新的 func_get_args 表达式
func NewFuncGetArgs(from data.From) data.GetValue {
	return &FuncGetArgs{
		Node: NewNode(from),
	}
}

// GetValue 获取 func_get_args 的值（返回所有函数参数）
func (f *FuncGetArgs) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	args := ctx.GetCallArgs()
	if args == nil || len(args) == 0 {
		return data.NewArrayValue(nil), nil
	}

	// 对每个参数表达式求值
	values := make([]data.Value, 0, len(args))
	for _, arg := range args {
		v, acl := arg.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		if val, ok := v.(data.Value); ok {
			values = append(values, val)
		} else {
			values = append(values, data.NewNullValue())
		}
	}

	return data.NewArrayValue(values), nil
}

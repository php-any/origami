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
// 使用 GetIndexValue 按索引读取已求值的参数，避免在错误上下文中重新求值表达式
func (f *FuncGetArgs) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	args := ctx.GetCallArgs()
	if args == nil || len(args) == 0 {
		return data.NewArrayValue(nil), nil
	}

	values := make([]data.Value, 0, len(args))
	for i := 0; i < len(args); i++ {
		v, ok := ctx.GetIndexValue(i)
		if ok && v != nil {
			values = append(values, v)
		} else {
			values = append(values, data.NewNullValue())
		}
	}

	return data.NewArrayValue(values), nil
}

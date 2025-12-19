package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FuncGetArgsFunction 实现 func_get_args 函数
// 在调用者上下文中，按顺序返回所有入参。
type FuncGetArgsFunction struct{}

func NewFuncGetArgsFunction() data.FuncStmt {
	return &FuncGetArgsFunction{}
}

func (f *FuncGetArgsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

func (f *FuncGetArgsFunction) GetName() string {
	return "func_get_args"
}

func (f *FuncGetArgsFunction) GetParams() []data.GetValue {
	// 使用 CallerContextParameter 标记：在调用时不创建新的函数上下文，
	// 而是直接在调用者的 Context 中执行 Call，这样就能读到上级入参。
	return []data.GetValue{
		node.NewCallerContextParameter(nil),
	}
}

func (f *FuncGetArgsFunction) GetVariables() []data.Variable {
	// 无显式局部变量
	return []data.Variable{}
}

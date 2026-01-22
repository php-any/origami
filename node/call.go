package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// CallExpression 表示函数调用表达式
type CallExpression struct {
	*Node   `pp:"-"`
	FunName string // 被调用的表达式
	Fun     data.FuncStmt
	Args    []data.GetValue
}

// NewCallExpression 创建一个新的函数调用表达式
func NewCallExpression(token *TokenFrom, fn string, arguments []data.GetValue, fun data.FuncStmt) *CallExpression {
	if fn[0:1] == "\\" {
		fn = fn[1:]
	}

	return &CallExpression{
		Node:    NewNode(token),
		FunName: fn,
		Fun:     fun,
		Args:    arguments,
	}
}

// GetValue 获取函数调用表达式的值
func (pe *CallExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	fn := pe.Fun
	varies := fn.GetVariables()
	params := fn.GetParams()
	arguments := pe.Args
	fnCtx := ctx.CreateContext(varies)
	var acl data.Control
	// 入参的值设置到上下文中
	for index, param := range params {
		if len(arguments) > index {
			arg := arguments[index]
			switch argTV := arg.(type) {
			case *NamedArgument:
				param, err := findParams(params, argTV.Name)
				if err != nil {
					return nil, data.NewErrorThrow(pe.from, err)
				}
				acl = paramSetValue(fnCtx, ctx, nil, param, argTV, varies, index, arguments)
			default:
				acl = paramSetValue(fnCtx, ctx, nil, param, argTV, varies, index, arguments)
			}
		} else {
			switch param := param.(type) {
			case *CallerContextParameter:
				fnCtx = ctx
			default:
				_, acl = param.GetValue(fnCtx)
			}
		}
	}

	if acl != nil {
		return nil, acl
	}

	// 将本次调用的参数表达式列表记录到函数上下文中
	fnCtx.SetCallArgs(pe.Args)

	return fn.Call(fnCtx)
}

func NewCallTodo(call *CallExpression, namespace string) *CallLater {
	return &CallLater{
		CallExpression: call,
		namespace:      namespace,
	}
}

// CallLater 未确认的函数调用
type CallLater struct {
	*CallExpression
	namespace string
}

func (pe *CallLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if pe.Fun == nil {
		fn, ok := ctx.GetVM().GetFunc(pe.FunName)
		if !ok {
			fn, ok = ctx.GetVM().GetFunc(pe.namespace + "\\" + pe.FunName)
			if !ok {
				namespace := ""
				if pe.namespace != "" {
					namespace = pe.namespace + "\\"
				}

				fn, ok = ctx.GetVM().GetFunc(namespace + pe.FunName)
				if !ok {
					return nil, data.NewErrorThrow(pe.from, errors.New(fmt.Sprintf("无法调用函数(%s), 未找到函数", pe.FunName)))
				}
			}
		}

		pe.FunName = fn.GetName()
		pe.Fun = fn
	}
	return pe.CallExpression.GetValue(ctx)
}

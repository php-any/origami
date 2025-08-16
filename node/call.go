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
	fnCtx := ctx.CreateContext(varies)
	// 入参的值设置到上下文中
	for index, arg := range fn.GetParams() {
		switch argObj := arg.(type) {
		case *Parameter:
			if index < len(pe.Args) {
				param := pe.Args[index]
				switch paramTV := param.(type) {
				case *NamedArgument:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					vari, err := findVariable(varies, paramTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(pe.from, err)
					}
					fnCtx.SetVariableValue(vari, tempV.(data.Value))
				default:
					tempV, acl := paramTV.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					acl = argObj.SetValue(fnCtx, tempV.(data.Value))
					if acl != nil {
						return nil, acl
					}
				}
			} else if argObj.DefaultValue == nil {
				return nil, data.NewErrorThrow(pe.from, errors.New(fmt.Sprintf("无法调用函数(%s), 缺少参数", pe.FunName)))
			} else {
				argObj.GetValue(fnCtx)
			}
		case *Parameters:
			args, _ := fnCtx.GetVariableValue(argObj)
			var ares *data.ArrayValue
			var ok bool
			if ares, ok = args.(*data.ArrayValue); !ok {
				ares = data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
			}

			for i := index; i < len(pe.Args); i++ {
				param := pe.Args[i]
				tempV, acl := param.GetValue(ctx)
				if acl != nil {
					return nil, acl
				}
				ares.Value = append(ares.Value, tempV.(data.Value))
				fnCtx.SetVariableValue(argObj, ares)
			}
		default:
			return nil, data.NewErrorThrow(pe.from, errors.New("未识别的参数类型"))
		}
	}

	return fn.Call(fnCtx)
}

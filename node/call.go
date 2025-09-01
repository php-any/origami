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

					if val, ok := tempV.(data.Value); ok {
						acl = argObj.SetValue(fnCtx, val)
					} else {
						acl = argObj.SetValue(fnCtx, data.NewNullValue())
					}

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
		case *ParameterReference:
			if index < len(pe.Args) {
				param := pe.Args[index]
				switch paramTV := param.(type) {
				case *NamedArgument:
					vari, err := findVariable(varies, paramTV.Name)
					if err != nil {
						return nil, data.NewErrorThrow(pe.from, err)
					}
					if val, ok := paramTV.Value.(data.Variable); ok {
						acl := vari.SetValue(fnCtx, data.NewReferenceValue(val, ctx))
						if acl != nil {
							return nil, acl
						}
					} else {
						return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能传入变量, fn: %s", pe.FunName))
					}
				default:
					if val, ok := paramTV.(data.Variable); ok {
						acl := argObj.SetValue(fnCtx, data.NewReferenceValue(val, ctx))
						if acl != nil {
							return nil, acl
						}
					} else {
						return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能传入变量, fn: %s", pe.FunName))
					}
				}
			} else {
				return nil, data.NewErrorThrow(pe.from, fmt.Errorf("引用参数只能是必传参数, fn: %s", pe.FunName))
			}
		default:
			return nil, data.NewErrorThrow(pe.from, errors.New("未识别的参数类型"))
		}
	}

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
		fn, ok := ctx.GetVM().GetFunc(pe.namespace + "\\" + pe.FunName)
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
		pe.FunName = fn.GetName()
		pe.Fun = fn
	}
	return pe.CallExpression.GetValue(ctx)
}

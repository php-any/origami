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
	if fn == nil {
		return nil, data.NewErrorThrow(pe.from, fmt.Errorf("无法调用函数(%s), 未找到函数", pe.FunName))
	}
	varies := fn.GetVariables()
	params := fn.GetParams()
	arguments := pe.Args
	fnCtx := ctx.CreateContext(varies)

	// 简单函数 + 展开参数的快速通道：所有形参都是普通 Parameter，且存在 SpreadArgument
	simpleParams := true
	for _, p := range params {
		if _, ok := p.(*Parameter); !ok {
			simpleParams = false
			break
		}
	}
	if simpleParams {
		hasSpread := false
		for _, a := range arguments {
			if _, ok := a.(*SpreadArgument); ok {
				hasSpread = true
				break
			}
		}
		if hasSpread {
			// 按调用实参顺序，将普通参数与 ...expr 展平成一维数组，然后依次绑定到形参
			var flat []data.Value
			for _, arg := range arguments {
				if spread, ok := arg.(*SpreadArgument); ok {
					// first-class callable 场景（Expr==nil）退回通用路径
					if spread.Expr == nil {
						hasSpread = false
						break
					}
					spreadVal, acl := spread.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					if spreadVal == nil {
						continue
					}
					switch v := spreadVal.(type) {
					case *data.ArrayValue:
						for _, z := range v.List {
							flat = append(flat, z.Value)
						}
					case *data.ObjectValue:
						v.RangeProperties(func(_ string, val data.Value) bool {
							flat = append(flat, val)
							return true
						})
					default:
						if val, ok := spreadVal.(data.Value); ok {
							flat = append(flat, val)
						}
					}
				} else {
					v, acl := arg.GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
					if v == nil {
						flat = append(flat, data.NewNullValue())
					} else if val, ok := v.(data.Value); ok {
						flat = append(flat, val)
					}
				}
			}

			if hasSpread {
				for i := 0; i < len(params) && i < len(flat) && i < len(varies); i++ {
					fnCtx.SetVariableValue(varies[i], flat[i])
				}
				fnCtx.SetCallArgs(pe.Args)
				return fn.Call(fnCtx)
			}
		}
	}

	var acl data.Control
	// 入参的值设置到上下文中
	for index, param := range params {
		// CallerContextParameter 无论有无传参都需切换为调用者上下文
		if _, ok := param.(*CallerContextParameter); ok {
			fnCtx = ctx
			continue
		}
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
			_, acl = param.GetValue(fnCtx)
		}
		if acl != nil {
			if addStack, ok := acl.(data.AddStack); ok {
				addStack.AddStackWithInfo(pe.from, "", pe.FunName+fmt.Sprintf("(%d:%s)", index, TryGetCallClassName(param)))
			}
			// 实参求值中的 throw 须向上冒泡，由调用方的 try/catch 处理，不能在此 fatal
			if _, ok := acl.(data.ThrowControl); ok {
				return nil, acl
			}
			ctx.GetVM().ThrowControl(acl)
			return nil, acl
		}
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

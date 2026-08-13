package node

import (
	"errors"
	"fmt"
	"sync"

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
	// PHP 8.1 一等函数可调用语法：strlen(...) 返回 Closure，不执行函数。
	if isFirstClassCallableArgs(pe.Args) {
		return data.NewFuncValue(fn), nil
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
	// PHP 8 命名参数：先按名绑定，再按位置填未绑定形参，最后补默认值。
	// 不能按 arguments[i] 对齐 params[i]，否则 named 会占位并在后续循环被默认值覆盖。
	bound := make([]bool, len(params))
	variadicIdx := -1
	for i, p := range params {
		if _, ok := p.(*Parameters); ok {
			variadicIdx = i
			break
		}
	}
	paramIndexByName := func(name string) (int, error) {
		for i, p := range params {
			if gn, ok := p.(data.GetName); ok && gn.GetName() == name {
				return i, nil
			}
		}
		return -1, errors.New("无法找到变量: " + name)
	}
	var positional []data.GetValue
	for _, arg := range arguments {
		switch a := arg.(type) {
		case *NamedArgument:
			idx, err := paramIndexByName(a.Name)
			if err != nil {
				if variadicIdx >= 0 {
					positional = append(positional, a)
					continue
				}
				return nil, data.NewErrorThrow(pe.from, err)
			}
			acl = paramSetValue(fnCtx, ctx, nil, params[idx], a, varies, idx, arguments)
			if acl != nil {
				break
			}
			bound[idx] = true
		case *CallerContextParameter:
			continue
		default:
			positional = append(positional, arg)
		}
		if acl != nil {
			break
		}
	}
	if acl == nil {
		pos := 0
		for index, param := range params {
			if _, ok := param.(*CallerContextParameter); ok {
				fnCtx = ctx
				continue
			}
			if bound[index] {
				continue
			}
			if pos < len(positional) {
				arg := positional[pos]
				pos++
				acl = paramSetValue(fnCtx, ctx, nil, param, arg, varies, index, arguments)
				if acl != nil {
					break
				}
				bound[index] = true
				continue
			}
			_, acl = param.GetValue(fnCtx)
			if acl != nil {
				break
			}
		}
	}
	if acl != nil {
		if addStack, ok := acl.(data.AddStack); ok {
			addStack.AddStackWithInfo(pe.from, "", pe.FunName)
		}
		if _, ok := acl.(data.ThrowControl); ok {
			return nil, acl
		}
		ctx.GetVM().ThrowControl(acl)
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
	resolveMu sync.Mutex
}

func (pe *CallLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	pe.resolveMu.Lock()
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
					pe.resolveMu.Unlock()
					return nil, data.NewErrorThrow(pe.from, errors.New(fmt.Sprintf("无法调用函数(%s), 未找到函数", pe.FunName)))
				}
			}
		}

		pe.FunName = fn.GetName()
		pe.Fun = fn
	}
	pe.resolveMu.Unlock()
	return pe.CallExpression.GetValue(ctx)
}

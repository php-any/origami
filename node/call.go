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
	return callValue(pe.GetReferenceValue(ctx))
}

func (pe *CallExpression) GetReferenceValue(ctx data.Context) (data.GetValue, data.Control) {
	return pe.invoke(ctx, pe.Fun)
}

func (pe *CallExpression) invoke(ctx data.Context, fn data.FuncStmt) (data.GetValue, data.Control) {
	// PHP 8.1 一等函数可调用语法：strlen(...) 返回 Closure，不执行函数。
	if isFirstClassCallableArgs(pe.Args) {
		return data.NewFuncValue(fn), nil
	}
	varies := fn.GetVariables()
	params := fn.GetParams()
	arguments := pe.Args
	if usesCallerContextParams(params) {
		return callOnCallerContext(ctx, arguments, fn.Call)
	}
	fnCtx := ctx.CreateContext(varies)
	allocated := fnCtx

	if canFastPositionalBind(params, arguments) {
		if acl := bindPositionalParameters(fnCtx, ctx, params, arguments, varies); acl != nil {
			if _, ok := acl.(ToClosure); ok {
				finishPooledCall(fn, allocated, ctx, nil, nil)
				return nil, acl
			}
			if addStack, ok := acl.(data.AddStack); ok {
				addStack.AddStackWithInfo(pe.from, "", pe.FunName)
			}
			if _, ok := acl.(data.ThrowControl); ok {
				finishPooledCall(fn, allocated, ctx, nil, nil)
				return nil, acl
			}
			ctx.GetVM().ThrowControl(acl)
			finishPooledCall(fn, allocated, ctx, nil, nil)
			return nil, acl
		}
		ret, ctl := fn.Call(fnCtx)
		return finishPooledCall(fn, allocated, ctx, ret, ctl)
	}

	if hasReferenceParameters(params) {
		if ctl := bindReferenceCall(fnCtx, ctx, params, arguments, nil); ctl != nil {
			return finishPooledCall(fn, allocated, ctx, nil, ctl)
		}
		ret, ctl := fn.Call(fnCtx)
		return finishPooledCall(fn, allocated, ctx, ret, ctl)
	}
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
			if ctl := bindReferenceCall(fnCtx, ctx, params, arguments, nil); ctl != nil {
				return finishPooledCall(fn, allocated, ctx, nil, ctl)
			}
			ret, ctl := fn.Call(fnCtx)
			return finishPooledCall(fn, allocated, ctx, ret, ctl)
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
		// first-class callable 占位不应走 ThrowControl/错误打印
		if _, ok := acl.(ToClosure); ok {
			return nil, acl
		}
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
	fnCtx.SetFlatCallArgs(collectCallArgValues(ctx, positional, fnCtx))

	ret, ctl := fn.Call(fnCtx)
	return finishPooledCall(fn, allocated, ctx, ret, ctl)
}

func NewCallTodo(call *CallExpression, namespace string) *CallLater {
	later := &CallLater{
		CallExpression: call,
		namespace:      namespace,
		functionID:     data.Symbols.Intern(call.FunName),
	}
	if namespace != "" {
		later.namespaceID = data.Symbols.Intern(namespace + "\\" + call.FunName)
	}
	return later
}

// CallLater 未确认的函数调用
type CallLater struct {
	*CallExpression
	namespace               string
	functionID, namespaceID data.SymbolID
}

func (pe *CallLater) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return callValue(pe.GetReferenceValue(ctx))
}

func (pe *CallLater) GetReferenceValue(ctx data.Context) (data.GetValue, data.Control) {
	fn, ctl := pe.resolveFun(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return pe.invoke(ctx, fn)
}

// Call sites retain only immutable symbol IDs. Resolution always uses the
// executing VM, so cached programs cannot retain another request's function.
func (pe *CallLater) resolveFun(ctx data.Context) (data.FuncStmt, data.Control) {
	functions := ctx.GetVM()
	if pe.functionID != 0 {
		if pe.namespaceID != 0 {
			if fn, found := functions.GetFuncBySymbol(pe.namespaceID); found {
				return fn, nil
			}
		}
		if fn, found := functions.GetFuncBySymbol(pe.functionID); found {
			return fn, nil
		}
		return nil, data.NewErrorThrow(pe.from, fmt.Errorf("无法调用函数(%s), 未找到函数", pe.FunName))
	}
	if pe.namespace != "" {
		if fn, ok := functions.GetFunc(pe.namespace + "\\" + pe.FunName); ok {
			return fn, nil
		}
	}
	if fn, ok := functions.GetFunc(pe.FunName); ok {
		return fn, nil
	}
	return nil, data.NewErrorThrow(pe.from, fmt.Errorf("无法调用函数(%s), 未找到函数", pe.FunName))
}

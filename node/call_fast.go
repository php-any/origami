package node

import (
	"github.com/php-any/origami/data"
)

type pooledContext interface {
	ReleasePooled()
}

func canFastPositionalBind(params []data.GetValue, args []data.GetValue) bool {
	if isFirstClassCallableArgs(args) {
		return false
	}
	for _, p := range params {
		switch p.(type) {
		case *Parameter, *ParameterRawAST:
		default:
			return false
		}
	}
	for _, a := range args {
		switch a.(type) {
		case *NamedArgument, *SpreadArgument, *CallerContextParameter:
			return false
		}
	}
	return true
}

func bindPositionalParameters(fnCtx, ctx data.Context, params []data.GetValue, args []data.GetValue, varies []data.Variable) data.Control {
	nArgs := len(args)
	nParams := len(params)
	if nArgs == 0 {
		for i := 0; i < nParams; i++ {
			if _, acl := params[i].GetValue(fnCtx); acl != nil {
				return acl
			}
		}
		return nil
	}

	needFlat := nArgs > nParams
	var flat []data.Value
	if needFlat {
		flat = make([]data.Value, 0, nArgs)
	}
	for i, arg := range args {
		if i < nParams {
			if raw, ok := params[i].(*ParameterRawAST); ok {
				if acl := raw.BindUnevaluated(fnCtx, ctx, arg); acl != nil {
					return acl
				}
				continue
			}
		}
		val, acl := evalCallArg(ctx, arg)
		if acl != nil {
			return acl
		}
		if needFlat {
			flat = append(flat, val)
		}
		if i < nParams {
			p := params[i].(*Parameter)
			if acl := bindPositionalValue(fnCtx, p, val); acl != nil {
				return acl
			}
		}
	}
	for i := nArgs; i < nParams; i++ {
		if _, acl := params[i].GetValue(fnCtx); acl != nil {
			return acl
		}
	}
	fnCtx.SetCallArgs(args)
	if needFlat {
		fnCtx.SetFlatCallArgs(flat)
	}
	return nil
}

func evalCallArg(ctx data.Context, arg data.GetValue) (data.Value, data.Control) {
	switch a := arg.(type) {
	case *VariableExpression:
		zv := ctx.GetIndexZVal(a.Index)
		if zv == nil || zv.Value == nil {
			return data.NewNullValue(), nil
		}
		return zv.Value, nil
	case *IntLiteral:
		return a.V, nil
	case *FloatLiteral:
		return a.V, nil
	case *StringLiteral:
		if a.intern == nil {
			a.intern = data.NewStringValue(a.Value)
		}
		return a.intern, nil
	case *BooleanLiteral:
		return data.NewBoolValue(a.Value), nil
	case *NullLiteral:
		return data.NewNullValue(), nil
	default:
		v, acl := arg.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}
		val, _ := v.(data.Value)
		if val == nil {
			return data.NewNullValue(), nil
		}
		return val, nil
	}
}

func bindPositionalValue(fnCtx data.Context, p *Parameter, val data.Value) data.Control {
	if p.Type == nil {
		if data.IsScalarAssignFast(val) {
			data.AssignScalarToZVal(fnCtx.GetIndexZVal(p.Index), val)
			return nil
		}
		return fnCtx.SetVariableValue(p, val)
	}
	return p.SetValue(fnCtx, val)
}

func tryReleaseCallContext(fn any, fnCtx data.Context) {
	if e, ok := fnCtx.(data.ContextEscaper); ok && e.IsEscaped() {
		return
	}
	switch f := fn.(type) {
	case *FunctionStatement:
		if f.IsGenerator || f.ReturnsReference {
			return
		}
	case *ClassMethod:
		if f.IsGenerator {
			return
		}
	}
	if r, ok := fnCtx.(pooledContext); ok {
		r.ReleasePooled()
	}
}

// usesCallerContextParams 表示函数必须在调用者符号表上执行（extract / get_defined_vars / unset）。
func usesCallerContextParams(params []data.GetValue) bool {
	for _, p := range params {
		if _, ok := p.(*CallerContextParameter); ok {
			return true
		}
	}
	return false
}

// callOnCallerContext 在调用者 Context 上执行内置函数。
// 不可 CreateContext + tryReleaseCallContext：CallerContextParameter 路径若把 fnCtx
// 指回调用方再 ReleasePooled，会把仍在使用的闭包/函数符号表还回 pool（slots 变成 0），
// Laravel Filesystem::getRequire 的 extract 之后 require $__path 会因此失败。
func callOnCallerContext(ctx data.Context, args []data.GetValue, invoke func(data.Context) (data.GetValue, data.Control)) (data.GetValue, data.Control) {
	prevArgs := ctx.GetCallArgs()
	prevFlat := ctx.GetFlatCallArgs()
	ctx.SetCallArgs(args)
	ret, ctl := invoke(ctx)
	ctx.SetCallArgs(prevArgs)
	ctx.SetFlatCallArgs(prevFlat)
	return ret, ctl
}

// finishPooledCall 回收本次调用创建的 pooled context。
// allocated 是 CreateContext 得到的帧；若执行改用了调用方 ctx（CallerContextParameter），
// 只回收 allocated，绝不能把 caller 还回 pool。
func finishPooledCall(fn any, allocated, caller data.Context, ret data.GetValue, ctl data.Control) (data.GetValue, data.Control) {
	if allocated != caller {
		tryReleaseCallContext(fn, allocated)
	}
	return ret, ctl
}

func markContextEscaped(ctx data.Context) {
	for ctx != nil {
		if e, ok := ctx.(data.ContextEscaper); ok {
			e.MarkEscaped()
		}
		switch t := ctx.(type) {
		case *data.BoundContext:
			ctx = t.Context
		case *data.ClassMethodContext:
			ctx = t.Context
		default:
			return
		}
	}
}

func phpCallEnter(ctx data.Context, frame data.CallFrame) int {
	rec := ctx.(data.CallRecorder)
	d := rec.EnterCall()
	rec.PushCallFrame(frame)
	return d
}

func phpCallLeave(ctx data.Context) {
	rec := ctx.(data.CallRecorder)
	rec.PopCallFrame()
	rec.LeaveCall()
}

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
	flat := make([]data.Value, 0, len(args))
	for i, arg := range args {
		if i < len(params) {
			if raw, ok := params[i].(*ParameterRawAST); ok {
				if acl := raw.BindUnevaluated(fnCtx, ctx, arg); acl != nil {
					return acl
				}
				continue
			}
		}
		v, acl := arg.GetValue(ctx)
		if acl != nil {
			return acl
		}
		val, _ := v.(data.Value)
		if val == nil {
			val = data.NewNullValue()
		}
		flat = append(flat, val)
		if i < len(params) {
			if acl := params[i].(*Parameter).SetValue(fnCtx, val); acl != nil {
				return acl
			}
		}
	}
	for i := len(args); i < len(params); i++ {
		if _, acl := params[i].GetValue(fnCtx); acl != nil {
			return acl
		}
	}
	fnCtx.SetCallArgs(args)
	fnCtx.SetFlatCallArgs(flat)
	return nil
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

func phpCallEnter(ctx data.Context, frame data.CallFrame) (leave func(), depth int) {
	rec := ctx.(data.CallRecorder)
	d := rec.EnterCall()
	rec.PushCallFrame(frame)
	return func() {
		rec.PopCallFrame()
		rec.LeaveCall()
	}, d
}

package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	netdata "github.com/php-any/origami/std/net/data"
	"github.com/php-any/origami/utils"
)

// NextFunc 实现 data.FuncStmt，作为洋葱模型中 $next 回调
type NextFunc struct {
	name     string
	fn       func(request data.Value, response data.Value) (data.GetValue, data.Control)
	variable []data.Variable
}

func (f NextFunc) Call(ctx data.Context) (_ data.GetValue, acl data.Control) {
	request, err := utils.ConvertFromIndex[data.Value](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	response, err := utils.ConvertFromIndex[data.Value](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	defer func() {
		if r := recover(); r != nil {
			if acl2, ok2 := r.(data.Control); ok2 {
				acl = acl2
				return
			}
			panic(r)
		}
	}()

	return f.fn(request, response)
}

func (f NextFunc) GetName() string { return f.name }
func (f NextFunc) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
	}
}
func (f NextFunc) GetVariables() []data.Variable {
	return f.variable
}

// executeMiddlewareChain 使用洋葱模型执行中间件链和控制器
// 链结构: middleware[0].handle → middleware[1].handle → ... → controller
func executeMiddlewareChain(vm data.VM, ctx data.Context, rt netdata.Route, request data.Value, response data.Value) (data.GetValue, data.Control) {
	middlewares := rt.Middlewares

	// 如果没有中间件，直接执行控制器
	if len(middlewares) == 0 {
		return executeControllerMethod(vm, ctx, rt, request, response)
	}

	// 构建洋葱链，从内到外
	chainIdx := 0

	var buildNext func() (data.GetValue, data.Control)
	buildNext = func() (data.GetValue, data.Control) {
		if chainIdx >= len(middlewares) {
			// 链末端：执行控制器方法
			return executeControllerMethod(vm, ctx, rt, request, response)
		}

		mw := middlewares[chainIdx]
		chainIdx++

		cls, ok := vm.GetClass(mw.ClassName)
		if !ok {
			// 中间件类不存在，跳过，继续下一个
			return buildNext()
		}

		// 实例化中间件类
		inst, acl := cls.GetValue(ctx)
		if acl != nil {
			return nil, acl
		}

		cv, ok := inst.(*data.ClassValue)
		if !ok {
			return buildNext()
		}

		// 获取 handle 方法
		method, has := cv.GetMethod("handle")
		if !has {
			// 没有 handle 方法，跳过
			return buildNext()
		}

		vars := method.GetVariables()
		if len(vars) < 3 {
			return buildNext()
		}

		var lastNextReturn data.GetValue
		// 创建 $next 回调，指向链中的下一个节点
		nextFunc := data.NewFuncValue(NextFunc{
			name: "next",
			fn: func(request data.Value, response data.Value) (data.GetValue, data.Control) {
				ret, acl := buildNext()
				lastNextReturn = ret
				return ret, acl
			},
			variable: []data.Variable{
				node.NewVariable(nil, "request", 0, nil),
				node.NewVariable(nil, "response", 1, nil),
			},
		})

		// 绑定参数并调用 handle($request, $response, $next)
		fnCtx := cv.CreateContext(vars)
		fnCtx.SetVariableValue(vars[0], request)
		fnCtx.SetVariableValue(vars[1], response)
		fnCtx.SetVariableValue(vars[2], nextFunc)

		ret, acl := method.Call(fnCtx)
		if isEmptyHandlerReturn(ret) && !isEmptyHandlerReturn(lastNextReturn) {
			ret = lastNextReturn
		}
		return ret, acl
	}

	return buildNext()
}

func isEmptyHandlerReturn(v data.GetValue) bool {
	if v == nil {
		return true
	}
	if _, ok := v.(*data.NullValue); ok {
		return true
	}
	return false
}

// executeControllerMethod 执行控制器方法
func executeControllerMethod(vm data.VM, callCtx data.Context, rt netdata.Route, reqProxy data.Value, resProxy data.Value) (data.GetValue, data.Control) {
	args, acl := resolveHandlerArgs(vm, callCtx, rt.HandlerSpec, reqProxy, resProxy)
	if acl != nil {
		return nil, acl
	}
	if args == nil {
		return nil, nil
	}

	receiver := rt.Receiver
	if rt.ReceiverFactory != nil {
		var acl data.Control
		receiver, acl = rt.ReceiverFactory(callCtx)
		if acl != nil {
			return nil, acl
		}
	}
	if receiver != nil {
		return node.CallHTTPControllerMethod(receiver, rt.Target, args)
	}

	mute := callCtx.CreateContext(rt.Target.GetVariables())
	for i, arg := range args {
		if i < len(rt.Target.GetVariables()) {
			mute.SetVariableValue(rt.Target.GetVariables()[i], arg)
		}
	}
	return rt.Target.Call(mute)
}

package kit

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
)

// RegisterConditionable 挂上 when / unless（对齐 Conditionable trait）。
func RegisterConditionable(methods map[string]data.Method, newWhenProxy func(data.Context, data.Value) (data.Value, data.Control)) {
	methods["when"] = InstanceMethodOpt("when", []string{"value", "callback", "default"}, 0, func(ctx data.Context) (data.GetValue, data.Control) {
		return conditionableWhen(ctx, newWhenProxy, false)
	})
	methods["unless"] = InstanceMethodOpt("unless", []string{"value", "callback", "default"}, 0, func(ctx data.Context) (data.GetValue, data.Control) {
		return conditionableWhen(ctx, newWhenProxy, true)
	})
}

func conditionableWhen(ctx data.Context, newWhenProxy func(data.Context, data.Value) (data.Value, data.Control), unless bool) (data.GetValue, data.Control) {
	recv := Receiver(ctx)
	if recv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Conditionable when/unless missing $this"))
	}
	n := callArgCount(ctx)
	if n == 0 {
		return newWhenProxy(ctx, recv)
	}
	value := Arg(ctx, 0)
	if fv, ok := value.(*data.FuncValue); ok && fv != nil {
		ret, ctl := Call(ctx, fv, recv)
		if ctl != nil {
			return nil, ctl
		}
		if v, ok := ret.(data.Value); ok {
			value = v
		}
	}
	if n == 1 {
		proxy, ctl := newWhenProxy(ctx, recv)
		if ctl != nil {
			return nil, ctl
		}
		cv, ok := proxy.(*data.ClassValue)
		if !ok {
			return proxy, nil
		}
		cond := Truthy(value)
		if unless {
			cond = !cond
		}
		return CallInstanceMethod(ctx, cv, "condition", data.NewBoolValue(cond))
	}
	callback := Arg(ctx, 1)
	defaultCb := Arg(ctx, 2)
	truthy := Truthy(value)
	if unless {
		truthy = !truthy
	}
	if truthy && isCallableValue(callback) {
		ret, ctl := Call(ctx, callback, recv, value)
		if ctl != nil {
			return nil, ctl
		}
		if ret != nil {
			if v, ok := ret.(data.Value); ok && !IsNull(v) {
				return v, nil
			}
		}
		return recv, nil
	}
	if !truthy && isCallableValue(defaultCb) {
		ret, ctl := Call(ctx, defaultCb, recv, value)
		if ctl != nil {
			return nil, ctl
		}
		if ret != nil {
			if v, ok := ret.(data.Value); ok && !IsNull(v) {
				return v, nil
			}
		}
		return recv, nil
	}
	return recv, nil
}

func isCallableValue(v data.Value) bool {
	if v == nil || IsNull(v) {
		return false
	}
	switch v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return true
	default:
		return false
	}
}

func callArgCount(ctx data.Context) int {
	if args := ctx.GetFlatCallArgs(); args != nil {
		return len(args)
	}
	n := 0
	for i := 0; i < 8; i++ {
		if v, ok := ctx.GetIndexValue(i); ok && v != nil && !IsNull(v) {
			n = i + 1
		}
	}
	return n
}

// CallInstanceMethod 在 $this 上调用实例方法。
func CallInstanceMethod(ctx data.Context, recv *data.ClassValue, method string, args ...data.Value) (data.GetValue, data.Control) {
	if recv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("CallInstanceMethod: nil receiver"))
	}
	m, ok := recv.GetMethod(strings.ToLower(method))
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s::%s does not exist.", recv.Class.GetName(), method))
	}
	vars := m.GetVariables()
	inner := recv.CreateContext(vars)
	fnCtx := data.WrapMethodFrame(inner, recv, recv.Class, recv.Class)
	for i, v := range vars {
		if i < len(args) && args[i] != nil {
			_ = fnCtx.SetVariableValue(v, args[i])
		}
	}
	fnCtx.SetFlatCallArgs(args)
	return m.Call(fnCtx)
}

package events

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

// listenerWrapper 对齐 Dispatcher::makeListener 返回的闭包：
// - 调用约定始终为 ($event, $payload)
// - GetStaticVariables()['listener'] 供 Telescope formatListeners 使用
type listenerWrapper struct {
	owner    *data.ClassValue
	listener data.Value
	wildcard bool
	params   []data.GetValue
	vars     []data.Variable
}

func makeListenerWrapper(owner *data.ClassValue, listener data.Value, wildcard bool) *data.FuncValue {
	w := &listenerWrapper{
		owner:    owner,
		listener: listener,
		wildcard: wildcard,
	}
	w.params = []data.GetValue{
		node.NewParameter(nil, "event", 0, nil, nil),
		node.NewParameter(nil, "payload", 1, nil, nil),
	}
	w.vars = []data.Variable{
		node.NewVariable(nil, "event", 0, nil),
		node.NewVariable(nil, "payload", 1, nil),
	}
	return data.NewFuncValue(w)
}

func (w *listenerWrapper) GetName() string              { return "Closure" }
func (w *listenerWrapper) GetParams() []data.GetValue   { return w.params }
func (w *listenerWrapper) GetVariables() []data.Variable { return w.vars }

func (w *listenerWrapper) GetStaticVariables() map[string]data.Value {
	return map[string]data.Value{
		"listener": w.listener,
	}
}

func (w *listenerWrapper) Call(ctx data.Context) (data.GetValue, data.Control) {
	event := kit.Arg(ctx, 0)
	payload := kit.Arg(ctx, 1)
	eventName := ""
	if event != nil {
		eventName = event.AsString()
	}
	payload = ensureArrayPayload(payload)
	args := listenerArgs(eventName, payload, w.wildcard)
	return invokeRawListener(ctx, w.owner, w.listener, args)
}

// invokeRawListener 调用尚未包装的原始 listener（闭包 / 类名 / [obj, method]）。
func invokeRawListener(ctx data.Context, cv *data.ClassValue, listener data.Value, args []data.Value) (data.GetValue, data.Control) {
	listener = kit.Unwrap(listener)
	switch t := listener.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return kit.Call(ctx, t, args...)
	case *data.StringValue:
		return invokeStringListener(ctx, cv, t.AsString(), args)
	case *data.ArrayValue:
		entries := kit.Entries(t)
		if len(entries) >= 2 {
			target := entries[0].Value
			method := entries[1].Value.AsString()
			if scv, ok := kit.Unwrap(target).(*data.ClassValue); ok {
				if m, ok := scv.GetMethod(method); ok && m != nil {
					nctx := scv.CreateContext(m.GetVariables())
					data.BindDeclaredArgs(nctx, m, args)
					return m.Call(nctx)
				}
			}
			if s, ok := kit.Unwrap(target).(*data.StringValue); ok {
				return invokeClassMethod(ctx, cv, s.AsString(), method, args)
			}
		}
	}
	return data.NewNullValue(), nil
}

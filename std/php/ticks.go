package php

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

type TickFunction struct{ unregister bool }

func (f *TickFunction) GetName() string {
	if f.unregister {
		return "unregister_tick_function"
	}
	return "register_tick_function"
}
func (f *TickFunction) GetParams() []data.GetValue {
	if f.unregister {
		return registerShutdownFunctionFunctionGetParams[:1]
	}
	return registerShutdownFunctionFunctionGetParams
}
func (f *TickFunction) GetVariables() []data.Variable {
	return registerShutdownFunctionFunctionGetVariables
}
func (f *TickFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cb, _ := ctx.GetIndexValue(0)
	provider, ok := ctx.GetVM().(data.TickProvider)
	if !ok {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("tick state unavailable"), "Error")
	}
	state := provider.PHPTickState()
	if f.unregister {
		retained := make([]data.TickCallback, 0, len(state.Callbacks))
		for _, entry := range state.Callbacks {
			equal, ctl := node.NewBinaryEqStrict(nil, cb, entry.Callback).GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}
			if !equal.(*data.BoolValue).Value {
				retained = append(retained, entry)
			}
		}
		state.Callbacks = retained
		return data.NewNullValue(), nil
	}
	check := core.NewIsCallableFunction()
	checkCtx := ctx.CreateContext(check.GetVariables())
	checkCtx.SetIndexZVal(0, data.NewZVal(cb))
	valid, ctl := check.Call(checkCtx)
	if ctl != nil {
		return nil, ctl
	}
	if b, ok := valid.(*data.BoolValue); !ok || !b.Value {
		return nil, data.NewTypeError(nil, fmt.Errorf("register_tick_function(): Argument #1 ($callback) must be a valid callback"))
	}
	var args []data.Value
	if value, ok := ctx.GetIndexValue(1); ok {
		if array, ok := value.(*data.ArrayValue); ok {
			args = array.ToValueList()
		}
	}
	state.Callbacks = append(state.Callbacks, data.TickCallback{Callback: cb, Invocation: &shutdownInvocation{callback: cb, args: args}})
	return data.NewBoolValue(true), nil
}

package container

import (
	"fmt"

	"github.com/php-any/origami/data"
)

func callableValue(v data.GetValue) (data.FuncStmt, bool) {
	switch fn := v.(type) {
	case *data.FuncValue:
		return fn.Value, true
	case *data.BoundFuncValue:
		return fn.Value, true
	default:
		return nil, false
	}
}

func invokeFactory(ctx data.Context, factory data.GetValue, containerHost data.GetValue) (data.GetValue, data.Control) {
	stmt, ok := callableValue(factory)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Binding factory is not callable."))
	}

	callCtx := ctx.CreateContext(stmt.GetVariables())
	if len(stmt.GetParams()) > 0 && containerHost != nil {
		if val, ok := containerHost.(data.Value); ok {
			callCtx.SetIndexZVal(0, data.NewZVal(val))
		}
	}

	switch fn := factory.(type) {
	case *data.FuncValue:
		return fn.Call(callCtx)
	case *data.BoundFuncValue:
		return fn.Call(callCtx)
	default:
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Binding factory is not callable."))
	}
}

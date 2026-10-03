package runtime

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"os"
)

func (vm *VM) PHPErrorState() *data.PHPErrorState {
	if call := currentRequestCallState(); call != nil {
		if call.phpErrors == nil {
			call.phpErrors = vm.phpErrors.NewRequest()
		}
		return call.phpErrors
	}
	return &vm.phpErrors
}
func (vm *RequestVM) PHPErrorState() *data.PHPErrorState {
	if call := vm.requestCall(); call != &vm.call {
		if call.phpErrors == nil {
			call.phpErrors = vm.Base.phpErrors.NewRequest()
		}
		return call.phpErrors
	}
	if vm.phpErrors == nil {
		vm.phpErrors = vm.Base.phpErrors.NewRequest()
	}
	return vm.phpErrors
}

func (vm *VM) ReportPHPError(ctx data.Context, level int, message string, from data.From) data.Control {
	return reportPHPError(ctx, vm.PHPErrorState(), vm.GetErrorHandlerFor(level), level, message, from)
}
func (vm *RequestVM) ReportPHPError(ctx data.Context, level int, message string, from data.From) data.Control {
	return reportPHPError(ctx, vm.PHPErrorState(), vm.GetErrorHandlerFor(level), level, message, from)
}
func reportPHPError(ctx data.Context, state *data.PHPErrorState, callback data.Value, level int, message string, from data.From) data.Control {
	file, line := "", 0
	if from != nil {
		file = from.GetSource()
		line, _ = from.GetStartPosition()
		line++
	} else if recorder, ok := ctx.(data.CallStackTracker); ok {
		frames := recorder.SnapshotCallStack()
		if len(frames) != 0 {
			frame := frames[len(frames)-1]
			file, line = frame.File, frame.Line
		}
	}
	if callback != nil && state.BeginHandler() {
		defer state.EndHandler()
		fn, ctl := node.ResolveCallback(ctx, callback)
		if ctl != nil {
			return ctl
		}
		var function data.FuncStmt
		var invoke func(data.Context) (data.GetValue, data.Control)
		switch callable := fn.(type) {
		case *data.FuncValue:
			function, invoke = callable.Value, callable.Call
		case *data.BoundFuncValue:
			function, invoke = callable.Value, callable.Call
		default:
			return data.NewTypeError(nil, fmt.Errorf("invalid error handler"))
		}
		fnCtx := ctx.CreateContext(function.GetVariables())
		if ctl := data.BindDeclaredArgs(fnCtx, function, []data.Value{data.NewIntValue(level), data.NewStringValue(message), data.NewStringValue(file), data.NewIntValue(line)}); ctl != nil {
			return ctl
		}
		ret, ctl := invoke(fnCtx)
		if ctl != nil {
			return ctl
		}
		if b, ok := ret.(*data.BoolValue); !ok || b.Value {
			return nil
		}
	}
	state.SetLast(&data.PHPErrorInfo{Type: level, Message: message, File: file, Line: line})
	if level&state.Reporting() != 0 {
		label := "Warning"
		switch level {
		case 8, 1024:
			label = "Notice"
		case 8192, 16384:
			label = "Deprecated"
		case 256:
			label = "Fatal error"
		}
		fmt.Fprintf(os.Stderr, "%s: %s in %s on line %d\n", label, message, file, line)
	}
	if level == 256 {
		return data.NewExitControl(255)
	}
	return nil
}

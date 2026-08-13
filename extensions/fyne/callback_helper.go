package fyne

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/php-any/origami/data"
)

// callPHPCallback 安全调用无参数的 PHP 闭包
func callPHPCallback(callback data.FuncStmt, ctx data.Context) {
	if callback == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "[fyne] PANIC in PHP callback: %v\n%s\n", r, debug.Stack())
		}
	}()
	fnCtx := ctx.CreateContext(callback.GetVariables())
	_, ctl := callback.Call(fnCtx)
	if ctl != nil {
		fmt.Fprintf(os.Stderr, "[fyne] PHP callback error: %s\n", ctl.AsString())
	}
}

// callPHPCallbackWith 安全调用带参数的 PHP 闭包
func callPHPCallbackWith(callback data.FuncStmt, ctx data.Context, values ...data.Value) {
	if callback == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "[fyne] PANIC in PHP callback: %v\n%s\n", r, debug.Stack())
		}
	}()
	fnCtx := ctx.CreateContext(callback.GetVariables())
	for i, v := range values {
		fnCtx.SetIndexZVal(i, data.NewZVal(v))
	}
	_, ctl := callback.Call(fnCtx)
	if ctl != nil {
		fmt.Fprintf(os.Stderr, "[fyne] PHP callback error: %s\n", ctl.AsString())
	}
}

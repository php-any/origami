package support

import (
	"fmt"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
)

type macroTable struct {
	mu sync.RWMutex
	fn map[string]data.Value
}

var illuminateMacros = map[string]*macroTable{}

func macrosFor(class string) *macroTable {
	if t, ok := illuminateMacros[class]; ok {
		return t
	}
	t := &macroTable{fn: map[string]data.Value{}}
	illuminateMacros[class] = t
	return t
}

func registerMacroable(methods map[string]data.Method, class string) {
	methods["macro"] = newStaticMethod("macro", []string{"name", "macro"}, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		return macroRegister(class, ctx)
	}, false)
	methods["hasmacro"] = newStaticMethod("hasMacro", []string{"name"}, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		t := macrosFor(class)
		t.mu.RLock()
		_, ok := t.fn[strings.ToLower(strArg(ctx, 0))]
		t.mu.RUnlock()
		return data.NewBoolValue(ok), nil
	}, false)
	methods["flushmacros"] = newStaticMethod("flushMacros", nil, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		t := macrosFor(class)
		t.mu.Lock()
		t.fn = map[string]data.Value{}
		t.mu.Unlock()
		return data.NewNullValue(), nil
	}, false)
	methods["__callstatic"] = newStaticMethod("__callStatic", []string{"method", "parameters"}, -1, func(ctx data.Context) (data.GetValue, data.Control) {
		return macroInvoke(class, ctx, nil)
	}, false)
}

func macroRegister(class string, ctx data.Context) (data.GetValue, data.Control) {
	name := strings.ToLower(strArg(ctx, 0))
	fn, _ := ctx.GetIndexValue(1)
	t := macrosFor(class)
	t.mu.Lock()
	t.fn[name] = fn
	t.mu.Unlock()
	return data.NewNullValue(), nil
}

func lookupMacro(class, name string) (data.Value, bool) {
	t := macrosFor(class)
	t.mu.RLock()
	fn, ok := t.fn[strings.ToLower(name)]
	t.mu.RUnlock()
	return fn, ok && fn != nil
}

func macroInvoke(class string, ctx data.Context, recv *data.ClassValue) (data.GetValue, data.Control) {
	name := strArg(ctx, 0)
	fn, ok := lookupMacro(class, name)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s::%s does not exist.", class, name))
	}
	var args []data.Value
	if params, ok := ctx.GetIndexValue(1); ok {
		for _, e := range toEntries(params) {
			args = append(args, e.value)
		}
	}
	if recv != nil {
		if fv, ok := fn.(*data.FuncValue); ok {
			fn = data.NewBoundFuncValue(fv.Value, class, recv)
		}
	}
	return callValue(ctx, fn, args...)
}

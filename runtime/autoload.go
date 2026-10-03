package runtime

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

func AddAutoLoad(fun *data.FuncValue) {
	parser.AddAutoLoad(fun)
}

func RemoveAutoLoad(fun *data.FuncValue) {
	parser.RemoveAutoLoad(fun)
}

// ClearAutoLoad 清空全部 autoload 回调，供开发模式热重载使用。
func ClearAutoLoad() {
	parser.ClearAutoLoad()
}

func GetAutoLoad() []*data.FuncValue {
	return parser.GetAutoLoad()
}

func CallAutoLoad(name string, ctx data.Context) (bool, data.Control) {
	return parser.CallAutoLoad(name, ctx)
}

type autoloadHost interface {
	AutoloadFunctions() []*data.FuncValue
	AddAutoload(*data.FuncValue)
	RemoveAutoload(*data.FuncValue)
}

func AddAutoLoadInContext(ctx data.Context, fn *data.FuncValue) {
	if host, ok := ctx.GetVM().(autoloadHost); ok {
		host.AddAutoload(fn)
		return
	}
	AddAutoLoad(fn)
}
func RemoveAutoLoadInContext(ctx data.Context, fn *data.FuncValue) {
	if host, ok := ctx.GetVM().(autoloadHost); ok {
		host.RemoveAutoload(fn)
		return
	}
	RemoveAutoLoad(fn)
}
func GetAutoLoadInContext(ctx data.Context) []*data.FuncValue {
	if host, ok := ctx.GetVM().(autoloadHost); ok {
		return host.AutoloadFunctions()
	}
	return GetAutoLoad()
}

type autoloadRegistrationHost interface {
	RegisterAutoload(data.Value, data.Value, bool)
	UnregisterAutoload(data.Value) bool
	AutoloadOriginals() []data.Value
}

func RegisterAutoloadInContext(ctx data.Context, original, resolved data.Value, prepend bool) {
	if host, ok := ctx.GetVM().(autoloadRegistrationHost); ok {
		host.RegisterAutoload(original, resolved, prepend)
		return
	}
	if fn, ok := resolved.(*data.FuncValue); ok {
		AddAutoLoad(fn)
	}
}
func UnregisterAutoloadInContext(ctx data.Context, original data.Value) bool {
	if host, ok := ctx.GetVM().(autoloadRegistrationHost); ok {
		return host.UnregisterAutoload(original)
	}
	if fn, ok := original.(*data.FuncValue); ok {
		RemoveAutoLoad(fn)
		return true
	}
	return false
}
func AutoloadOriginalsInContext(ctx data.Context) []data.Value {
	if host, ok := ctx.GetVM().(autoloadRegistrationHost); ok {
		return host.AutoloadOriginals()
	}
	var out []data.Value
	for _, fn := range GetAutoLoad() {
		out = append(out, fn)
	}
	return out
}

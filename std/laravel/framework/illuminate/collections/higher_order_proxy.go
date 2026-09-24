package collections

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const higherOrderProxyName = "Illuminate\\Support\\HigherOrderCollectionProxy"

var collectionProxies = map[string]struct{}{
	"average": {}, "avg": {}, "contains": {}, "doesntcontain": {},
	"each": {}, "every": {}, "filter": {}, "first": {}, "flatmap": {},
	"groupby": {}, "hasmany": {}, "hassole": {}, "keyby": {}, "last": {},
	"map": {}, "max": {}, "min": {}, "partition": {}, "percentage": {},
	"reject": {}, "skipuntil": {}, "skipwhile": {}, "some": {},
	"sortby": {}, "sortbydesc": {}, "sum": {}, "takeuntil": {}, "takewhile": {},
	"unique": {}, "unless": {}, "until": {}, "when": {},
}

type HigherOrderProxyClass struct {
	node.Node
	methods map[string]data.Method
}

func NewHigherOrderProxyClass() data.ClassStmt {
	c := &HigherOrderProxyClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = kit.InstanceMethod("__construct", []string{"collection", "method"}, hopConstruct)
	c.methods["__get"] = kit.InstanceMethod("__get", []string{"key"}, hopGet)
	c.methods["__call"] = kit.InstanceMethod("__call", []string{"method", "parameters"}, hopCall)
	return c
}

func (c *HigherOrderProxyClass) GetName() string                          { return higherOrderProxyName }
func (c *HigherOrderProxyClass) GetExtend() *string                       { return nil }
func (c *HigherOrderProxyClass) GetImplements() []string                  { return nil }
func (c *HigherOrderProxyClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "collection", "method":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *HigherOrderProxyClass) GetPropertyList() []data.Property {
	p1, _ := c.GetProperty("collection")
	p2, _ := c.GetProperty("method")
	return []data.Property{p1, p2}
}
func (c *HigherOrderProxyClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *HigherOrderProxyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *HigherOrderProxyClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *HigherOrderProxyClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *HigherOrderProxyClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func newHigherOrderProxy(ctx data.Context, coll data.Value, method string) (*data.ClassValue, data.Control) {
	cv := data.NewClassValue(NewHigherOrderProxyClass(), ctx.CreateBaseContext())
	_ = cv.SetProperty("collection", coll)
	_ = cv.SetProperty("method", data.NewStringValue(method))
	return cv, nil
}

func hopConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("HigherOrderCollectionProxy missing $this"))
	}
	_ = cv.SetProperty("collection", kit.Arg(ctx, 0))
	_ = cv.SetProperty("method", kit.Arg(ctx, 1))
	return data.NewNullValue(), nil
}

func hopGet(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	key := ""
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	return hopApply(ctx, cv, func(itemCtx data.Context) (data.GetValue, data.Control) {
		item := kit.Arg(itemCtx, 0)
		if av, ok := kit.Unwrap(item).(*data.ArrayValue); ok {
			if got, ok := dataGetPath(av, key); ok {
				return got, nil
			}
			return data.NewNullValue(), nil
		}
		if icv, ok := kit.Unwrap(item).(*data.ClassValue); ok {
			v, ctl := icv.GetProperty(key)
			if ctl != nil {
				return nil, ctl
			}
			if v == nil {
				return data.NewNullValue(), nil
			}
			return v, nil
		}
		return data.NewNullValue(), nil
	})
}

func hopCall(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	method := ""
	if v := kit.Arg(ctx, 0); v != nil {
		method = v.AsString()
	}
	params := kit.Arg(ctx, 1)
	var args []data.Value
	if av, ok := kit.Unwrap(params).(*data.ArrayValue); ok && av != nil {
		args = av.ToValueList()
	}
	return hopApply(ctx, cv, func(itemCtx data.Context) (data.GetValue, data.Control) {
		item := kit.Unwrap(kit.Arg(itemCtx, 0))
		if icv, ok := item.(*data.ClassValue); ok {
			m, ok := icv.GetMethod(method)
			if !ok || m == nil {
				return nil, data.NewErrorThrow(nil, fmt.Errorf("当前值(%q)不支持调用函数, 你调用的函数(%s)", icv.Class.GetName(), method))
			}
			nctx := icv.CreateContext(m.GetVariables())
			data.BindDeclaredArgs(nctx, m, args)
			return m.Call(nctx)
		}
		switch c := item.(type) {
		case *data.FuncValue, *data.BoundFuncValue:
			return kit.Call(itemCtx, c, args...)
		}
		return nil, data.NewErrorThrow(nil, fmt.Errorf("当前值(%q)不支持调用函数, 你调用的函数(%s)", item.AsString(), method))
	})
}

func hopApply(ctx data.Context, proxy *data.ClassValue, mapFn func(data.Context) (data.GetValue, data.Control)) (data.GetValue, data.Control) {
	collV, _ := proxy.GetProperty("collection")
	methodV, _ := proxy.GetProperty("method")
	coll, ok := kit.Unwrap(collV).(*data.ClassValue)
	if !ok || coll == nil {
		return data.NewNullValue(), nil
	}
	methodName := ""
	if methodV != nil {
		methodName = methodV.AsString()
	}
	m, ok := coll.GetMethod(methodName)
	if !ok || m == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method Illuminate\\Support\\Collection::%s does not exist.", methodName))
	}
	cb := data.NewFuncValue(&hopFuncStmt{fn: mapFn})
	nctx := coll.CreateContext(m.GetVariables())
	data.BindDeclaredArgs(nctx, m, []data.Value{cb})
	return m.Call(nctx)
}

type hopFuncStmt struct {
	fn func(data.Context) (data.GetValue, data.Control)
}

func (h *hopFuncStmt) Call(ctx data.Context) (data.GetValue, data.Control) { return h.fn(ctx) }
func (h *hopFuncStmt) GetName() string                                     { return "{closure}" }
var hopFuncStmtGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
	node.NewParameter(nil, "key", 1, data.NewNullValue(), nil),
}

func (h *hopFuncStmt) GetParams() []data.GetValue {
	return hopFuncStmtGetParams
}
var hopFuncStmtGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, nil),
	node.NewVariable(nil, "key", 1, nil),
}

func (h *hopFuncStmt) GetVariables() []data.Variable {
	return hopFuncStmtGetVariables
}
func (h *hopFuncStmt) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *hopFuncStmt) GetIsStatic() bool          { return false }
func (h *hopFuncStmt) GetReturnType() data.Types  { return nil }

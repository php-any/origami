package conditionable

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const whenProxyClassName = "Illuminate\\Support\\HigherOrderWhenProxy"

type WhenProxyClass struct {
	node.Node
	methods map[string]data.Method
}

func NewWhenProxyClass() data.ClassStmt {
	c := &WhenProxyClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *WhenProxyClass) GetName() string                          { return whenProxyClassName }
func (c *WhenProxyClass) GetExtend() *string                       { return nil }
func (c *WhenProxyClass) GetImplements() []string                  { return nil }
func (c *WhenProxyClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "target", "condition", "hasCondition", "negateConditionOnCapture":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *WhenProxyClass) GetPropertyList() []data.Property {
	out := make([]data.Property, 0, 4)
	for _, n := range []string{"target", "condition", "hasCondition", "negateConditionOnCapture"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *WhenProxyClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *WhenProxyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *WhenProxyClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *WhenProxyClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *WhenProxyClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *WhenProxyClass) register() {
	c.methods["__construct"] = kit.InstanceMethod("__construct", []string{"target"}, whenProxyConstruct)
	c.methods["condition"] = kit.InstanceMethod("condition", []string{"condition"}, whenProxyCondition)
	c.methods["negateconditiononcapture"] = kit.InstanceMethod("negateConditionOnCapture", nil, whenProxyNegate)
	c.methods["__get"] = kit.InstanceMethod("__get", []string{"key"}, whenProxyGet)
	c.methods["__call"] = kit.InstanceMethod("__call", []string{"method", "parameters"}, whenProxyCall)
}

func whenProxyConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := whenProxyRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("target", kit.Arg(ctx, 0))
	_ = cv.SetProperty("hasCondition", data.NewBoolValue(false))
	_ = cv.SetProperty("negateConditionOnCapture", data.NewBoolValue(false))
	return cv, nil
}

func whenProxyCondition(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := whenProxyRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("condition", kit.Arg(ctx, 0))
	_ = cv.SetProperty("hasCondition", data.NewBoolValue(true))
	return cv, nil
}

func whenProxyNegate(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := whenProxyRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("negateConditionOnCapture", data.NewBoolValue(true))
	return cv, nil
}

func whenProxyGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := whenProxyRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	hasCondV, _ := cv.GetProperty("hasCondition")
	hasCond := kit.Truthy(hasCondV)
	keyV := kit.Arg(ctx, 0)
	key := ""
	if keyV != nil {
		key = keyV.AsString()
	}
	targetV, _ := cv.GetProperty("target")
	if !hasCond {
		val, ctl := readTargetProperty(ctx, targetV, key)
		if ctl != nil {
			return nil, ctl
		}
		negV, _ := cv.GetProperty("negateConditionOnCapture")
		cond := truthyGetValue(val)
		if bv, ok := negV.(*data.BoolValue); ok && bv.Value {
			cond = !cond
		}
		return kit.CallInstanceMethod(ctx, cv, "condition", data.NewBoolValue(cond))
	}
	condV, _ := cv.GetProperty("condition")
	if kit.Truthy(condV) {
		return readTargetProperty(ctx, targetV, key)
	}
	return targetV, nil
}

func whenProxyCall(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := whenProxyRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	hasCondV, _ := cv.GetProperty("hasCondition")
	hasCond := kit.Truthy(hasCondV)
	methodV := kit.Arg(ctx, 0)
	paramsV := kit.Arg(ctx, 1)
	method := ""
	if methodV != nil {
		method = methodV.AsString()
	}
	args := []data.Value{}
	if av, ok := paramsV.(*data.ArrayValue); ok && av != nil {
		args = av.ToValueList()
	}
	targetV, _ := cv.GetProperty("target")
	if !hasCond {
		ret, ctl := callTargetMethod(ctx, targetV, method, args)
		if ctl != nil {
			return nil, ctl
		}
		negV, _ := cv.GetProperty("negateConditionOnCapture")
		cond := truthyGetValue(ret)
		if bv, ok := negV.(*data.BoolValue); ok && bv.Value {
			cond = !cond
		}
		return kit.CallInstanceMethod(ctx, cv, "condition", data.NewBoolValue(cond))
	}
	condV, _ := cv.GetProperty("condition")
	if kit.Truthy(condV) {
		return callTargetMethod(ctx, targetV, method, args)
	}
	return targetV, nil
}

func whenProxyRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("HigherOrderWhenProxy missing $this"))
}

func readTargetProperty(ctx data.Context, target data.Value, key string) (data.GetValue, data.Control) {
	target = kit.Unwrap(target)
	if cv, ok := target.(*data.ClassValue); ok && cv != nil {
		if v, ctl := cv.GetProperty(key); ctl == nil && v != nil && !kit.IsNull(v) {
			return v, nil
		}
		return kit.CallInstanceMethod(ctx, cv, "__get", data.NewStringValue(key))
	}
	if ov, ok := target.(*data.ObjectValue); ok && ov != nil {
		v, ctl := ov.GetProperty(key)
		return v, ctl
	}
	if av, ok := target.(*data.ArrayValue); ok && av != nil {
		if z, ok := av.LookupZValByStringKey(key); ok && z != nil {
			return z.Value, nil
		}
	}
	return data.NewNullValue(), nil
}

func callTargetMethod(ctx data.Context, target data.Value, method string, args []data.Value) (data.GetValue, data.Control) {
	target = kit.Unwrap(target)
	if cv, ok := target.(*data.ClassValue); ok && cv != nil {
		return kit.CallInstanceMethod(ctx, cv, method, args...)
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("HigherOrderWhenProxy target is not an object"))
}

func truthyGetValue(v data.GetValue) bool {
	if v == nil {
		return false
	}
	if val, ok := v.(data.Value); ok {
		return kit.Truthy(val)
	}
	return true
}

// NewWhenProxy 构造 HigherOrderWhenProxy 实例。
func NewWhenProxy(ctx data.Context, target data.Value) (data.Value, data.Control) {
	stmt := NewWhenProxyClass()
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	_ = cv.SetProperty("target", target)
	_ = cv.SetProperty("hasCondition", data.NewBoolValue(false))
	_ = cv.SetProperty("negateConditionOnCapture", data.NewBoolValue(false))
	return cv, nil
}

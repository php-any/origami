package session

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const storeClassName = "Illuminate\\Session\\Store"

type StoreClass struct {
	node.Node
	methods map[string]data.Method
}

func NewStoreClass() data.ClassStmt {
	c := &StoreClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *StoreClass) GetName() string { return storeClassName }
func (c *StoreClass) GetExtend() *string {
	return nil
}
func (c *StoreClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Session\\Session"}
}
func (c *StoreClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "id", "name", "attributes", "handler", "serialization", "started":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *StoreClass) GetPropertyList() []data.Property {
	names := []string{"id", "name", "attributes", "handler", "serialization", "started"}
	out := make([]data.Property, 0, len(names))
	for _, n := range names {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *StoreClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *StoreClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	_ = cv.SetProperty("attributes", data.NewArrayValue(nil))
	_ = cv.SetProperty("started", data.NewBoolValue(false))
	return cv, nil
}
func (c *StoreClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *StoreClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *StoreClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *StoreClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"name", "handler", "id", "serialization"}, 0, sessConstruct)
	c.methods["get"] = kit.InstanceMethodOpt("get", []string{"key", "default"}, 1, sessGet)
	c.methods["put"] = kit.InstanceMethodOpt("put", []string{"key", "value"}, 1, sessPut)
	c.methods["has"] = kit.InstanceMethod("has", []string{"key"}, sessHas)
	c.methods["all"] = kit.InstanceMethod("all", nil, sessAll)
	c.methods["pull"] = kit.InstanceMethodOpt("pull", []string{"key", "default"}, 1, sessPull)
	c.methods["flash"] = kit.InstanceMethodOpt("flash", []string{"key", "value"}, 1, sessFlash)
	c.methods["forget"] = kit.InstanceMethod("forget", []string{"keys"}, sessForget)
}

func sessRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Session\\Store missing $this"))
}

func sessAttrs(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("attributes")
	if av, ok := kit.Unwrap(v).(*data.ArrayValue); ok && av != nil {
		return av
	}
	av := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("attributes", av)
	return av
}

func sessConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := sessRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if n := kit.Arg(ctx, 0); n != nil {
		_ = cv.SetProperty("name", n)
	}
	if h := kit.Arg(ctx, 1); h != nil {
		_ = cv.SetProperty("handler", h)
	}
	if id := kit.Arg(ctx, 2); id != nil {
		_ = cv.SetProperty("id", id)
	}
	if ser := kit.Arg(ctx, 3); ser != nil {
		_ = cv.SetProperty("serialization", ser)
	}
	if v, _ := cv.GetProperty("attributes"); v == nil || kit.IsNull(kit.Unwrap(v)) {
		_ = cv.SetProperty("attributes", data.NewArrayValue(nil))
	}
	return cv, nil
}

func sessKey(v data.Value) string {
	if v == nil {
		return ""
	}
	return kit.Unwrap(v).AsString()
}

func sessGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := sessRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := sessKey(kit.Arg(ctx, 0))
	def := kit.Arg(ctx, 1)
	attrs := sessAttrs(cv)
	if z, ok := attrs.LookupZValByStringKey(key); ok && z != nil && z.Value != nil {
		return z.Value, nil
	}
	if def != nil && !kit.IsNull(def) {
		return def, nil
	}
	return data.NewNullValue(), nil
}

func sessPut(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := sessRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keyArg := kit.Arg(ctx, 0)
	value := kit.Arg(ctx, 1)
	attrs := sessAttrs(cv)
	if av, ok := kit.Unwrap(keyArg).(*data.ArrayValue); ok && av != nil {
		for _, e := range av.List {
			if e == nil {
				continue
			}
			k := e.Name
			if k == "" {
				continue
			}
			attrs.SetStringKey(k, e.Value)
		}
		return data.NewNullValue(), nil
	}
	key := sessKey(keyArg)
	attrs.SetStringKey(key, value)
	return data.NewNullValue(), nil
}

func sessHas(ctx data.Context) (data.GetValue, data.Control) {
	val, ctl := sessGet(ctx)
	if ctl != nil {
		return nil, ctl
	}
	v, _ := val.(data.Value)
	return data.NewBoolValue(!kit.IsNull(v)), nil
}

func sessAll(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := sessRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return sessAttrs(cv), nil
}

func sessPull(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := sessRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := sessKey(kit.Arg(ctx, 0))
	def := kit.Arg(ctx, 1)
	attrs := sessAttrs(cv)
	var val data.Value = data.NewNullValue()
	if z, ok := attrs.LookupZValByStringKey(key); ok && z != nil && z.Value != nil {
		val = z.Value
	} else if def != nil && !kit.IsNull(def) {
		val = def
	}
	attrs.UnsetKey(data.NewStringValue(key))
	return val, nil
}

func sessFlash(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := sessRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := sessKey(kit.Arg(ctx, 0))
	value := kit.Arg(ctx, 1)
	if value == nil {
		value = data.NewBoolValue(true)
	}
	attrs := sessAttrs(cv)
	attrs.SetStringKey(key, value)
	flashNew, _ := attrs.LookupZValByStringKey("_flash.new")
	var flashArr *data.ArrayValue
	if flashNew != nil && flashNew.Value != nil {
		flashArr, _ = kit.Unwrap(flashNew.Value).(*data.ArrayValue)
	}
	if flashArr == nil {
		flashArr = data.NewArrayValue(nil).(*data.ArrayValue)
		attrs.SetStringKey("_flash.new", flashArr)
	}
	flashArr.SetIntKey(len(flashArr.List), data.NewStringValue(key))
	return data.NewNullValue(), nil
}

func sessForget(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := sessRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys := kit.Arg(ctx, 0)
	attrs := sessAttrs(cv)
	if av, ok := kit.Unwrap(keys).(*data.ArrayValue); ok && av != nil {
		for _, e := range av.List {
			if e == nil {
				continue
			}
			k := e.Name
			if k == "" && e.Value != nil {
				k = kit.Unwrap(e.Value).AsString()
			}
			if k != "" {
				attrs.UnsetKey(data.NewStringValue(k))
			}
		}
		return data.NewNullValue(), nil
	}
	attrs.UnsetKey(keys)
	return data.NewNullValue(), nil
}

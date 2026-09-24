package cache

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const repositoryClassName = "Illuminate\\Cache\\Repository"

type RepositoryClass struct {
	node.Node
	methods map[string]data.Method
}

func NewRepositoryClass() data.ClassStmt {
	c := &RepositoryClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *RepositoryClass) GetName() string { return repositoryClassName }
func (c *RepositoryClass) GetExtend() *string {
	return nil
}
func (c *RepositoryClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Cache\\Repository"}
}
func (c *RepositoryClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "store", "config", "__memory":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *RepositoryClass) GetPropertyList() []data.Property {
	out := make([]data.Property, 0, 3)
	for _, n := range []string{"store", "config", "__memory"} {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *RepositoryClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *RepositoryClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	_ = cv.SetProperty("__memory", data.NewArrayValue(nil))
	return cv, nil
}
func (c *RepositoryClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *RepositoryClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *RepositoryClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *RepositoryClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"store", "config"}, 0, cacheConstruct)
	c.methods["has"] = kit.InstanceMethod("has", []string{"key"}, cacheHas)
	c.methods["get"] = kit.InstanceMethodOpt("get", []string{"key", "default"}, 1, cacheGet)
	c.methods["put"] = kit.InstanceMethodOpt("put", []string{"key", "value", "ttl"}, 2, cachePut)
	c.methods["forever"] = kit.InstanceMethod("forever", []string{"key", "value"}, cacheForever)
	c.methods["forget"] = kit.InstanceMethod("forget", []string{"key"}, cacheForget)
}

func cacheRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Cache\\Repository missing $this"))
}

func cacheMemory(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("__memory")
	if av, ok := kit.Unwrap(v).(*data.ArrayValue); ok && av != nil {
		return av
	}
	av := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("__memory", av)
	return av
}

func cacheKey(v data.Value) string {
	if v == nil {
		return ""
	}
	return kit.Unwrap(v).AsString()
}

func cacheConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := cacheRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if store := kit.Arg(ctx, 0); store != nil {
		_ = cv.SetProperty("store", store)
	}
	if cfg := kit.Arg(ctx, 1); cfg != nil {
		_ = cv.SetProperty("config", cfg)
	}
	if v, _ := cv.GetProperty("__memory"); v == nil || kit.IsNull(kit.Unwrap(v)) {
		_ = cv.SetProperty("__memory", data.NewArrayValue(nil))
	}
	return cv, nil
}

func cacheGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := cacheRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := cacheKey(kit.Arg(ctx, 0))
	mem := cacheMemory(cv)
	if z, ok := mem.LookupZValByStringKey(key); ok && z != nil && z.Value != nil {
		return z.Value, nil
	}
	def := kit.Arg(ctx, 1)
	if fv, ok := kit.Unwrap(def).(*data.FuncValue); ok && fv != nil {
		ret, ctl := kit.Call(ctx, fv)
		if ctl != nil {
			return nil, ctl
		}
		if v, ok := ret.(data.Value); ok {
			return v, nil
		}
	}
	if def != nil && !kit.IsNull(def) {
		return def, nil
	}
	return data.NewNullValue(), nil
}

func cacheHas(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := cacheRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := cacheKey(kit.Arg(ctx, 0))
	mem := cacheMemory(cv)
	if z, ok := mem.LookupZValByStringKey(key); ok && z != nil && z.Value != nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func cachePut(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := cacheRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := cacheKey(kit.Arg(ctx, 0))
	value := kit.Arg(ctx, 1)
	mem := cacheMemory(cv)
	mem.SetStringKey(key, value)
	return data.NewBoolValue(true), nil
}

func cacheForever(ctx data.Context) (data.GetValue, data.Control) {
	return cachePut(ctx)
}

func cacheForget(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := cacheRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := cacheKey(kit.Arg(ctx, 0))
	mem := cacheMemory(cv)
	mem.UnsetKey(data.NewStringValue(key))
	return data.NewBoolValue(true), nil
}

package config

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/illuminate/collections"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const repositoryClassName = "Illuminate\\Config\\Repository"

// RepositoryClass 对齐 Illuminate\Config\Repository。
type RepositoryClass struct {
	node.Node
	methods map[string]data.Method
}

func NewRepositoryClass() data.ClassStmt {
	c := &RepositoryClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *RepositoryClass) GetName() string    { return repositoryClassName }
func (c *RepositoryClass) GetExtend() *string { return nil }
func (c *RepositoryClass) GetImplements() []string {
	return []string{"ArrayAccess", "Illuminate\\Contracts\\Config\\Repository"}
}
func (c *RepositoryClass) GetProperty(name string) (data.Property, bool) {
	if name == "items" {
		return node.NewProperty(nil, "items", "protected", false, data.NewArrayValue(nil)), true
	}
	return nil, false
}
func (c *RepositoryClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "items", "protected", false, data.NewArrayValue(nil))}
}
func (c *RepositoryClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *RepositoryClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *RepositoryClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *RepositoryClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *RepositoryClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *RepositoryClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"items"}, 0, repoConstruct)
	c.methods["has"] = kit.InstanceMethod("has", []string{"key"}, repoHas)
	c.methods["get"] = kit.InstanceMethodOpt("get", []string{"key", "default"}, 1, repoGet)
	c.methods["getmany"] = kit.InstanceMethod("getMany", []string{"keys"}, repoGetMany)
	c.methods["string"] = kit.InstanceMethodOpt("string", []string{"key", "default"}, 1, repoString)
	c.methods["integer"] = kit.InstanceMethodOpt("integer", []string{"key", "default"}, 1, repoInteger)
	c.methods["float"] = kit.InstanceMethodOpt("float", []string{"key", "default"}, 1, repoFloat)
	c.methods["boolean"] = kit.InstanceMethodOpt("boolean", []string{"key", "default"}, 1, repoBoolean)
	c.methods["array"] = kit.InstanceMethodOpt("array", []string{"key", "default"}, 1, repoArray)
	c.methods["collection"] = kit.InstanceMethodOpt("collection", []string{"key", "default"}, 1, repoCollection)
	c.methods["set"] = kit.InstanceMethodOpt("set", []string{"key", "value"}, 1, repoSet)
	c.methods["prepend"] = kit.InstanceMethod("prepend", []string{"key", "value"}, repoPrepend)
	c.methods["push"] = kit.InstanceMethod("push", []string{"key", "value"}, repoPush)
	c.methods["all"] = kit.InstanceMethod("all", nil, repoAll)
	c.methods["offsetexists"] = kit.InstanceMethod("offsetExists", []string{"offset"}, repoHas)
	c.methods["offsetget"] = kit.InstanceMethod("offsetGet", []string{"offset"}, repoGet)
	c.methods["offsetset"] = kit.InstanceMethod("offsetSet", []string{"offset", "value"}, repoSet)
	c.methods["offsetunset"] = kit.InstanceMethod("offsetUnset", []string{"offset"}, repoOffsetUnset)
	kit.RegisterMacroable(c.methods, repositoryClassName)
}

func repoReceiver(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Config\\Repository method missing $this"))
}

func repoItems(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("items")
	if av, ok := v.(*data.ArrayValue); ok && av != nil {
		return av
	}
	empty := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("items", empty)
	return empty
}

func repoConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := toAssocArray(kit.Arg(ctx, 0))
	_ = cv.SetProperty("items", items)
	return cv, nil
}

// toAssocArray 把构造参数统一成 ArrayValue（PHP 数组字面量有时是 ObjectValue）。
func toAssocArray(v data.Value) *data.ArrayValue {
	v = kit.Unwrap(v)
	if v == nil || kit.IsNull(v) {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return deepArrayCopy(av)
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(v) {
		zv := data.NewZVal(normalizeConfigValue(e.Value))
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
	}
	return out
}

func normalizeConfigValue(v data.Value) data.Value {
	v = kit.Unwrap(v)
	if v == nil {
		return data.NewNullValue()
	}
	switch t := v.(type) {
	case *data.ArrayValue:
		return deepArrayCopy(t)
	case *data.ObjectValue:
		return toAssocArray(t)
	default:
		return v
	}
}

func deepArrayCopy(src *data.ArrayValue) *data.ArrayValue {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, z := range src.List {
		if z == nil {
			continue
		}
		nz := data.NewZVal(normalizeConfigValue(z.Value))
		nz.Name = z.Name
		out.List = append(out.List, nz)
	}
	return out
}

func asArrayValue(v data.Value) *data.ArrayValue {
	v = kit.Unwrap(v)
	if v == nil || kit.IsNull(v) {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	if av, ok := v.(*data.ArrayValue); ok {
		return deepArrayCopy(av)
	}
	if _, ok := v.(*data.ObjectValue); ok {
		return toAssocArray(v)
	}
	return data.NewArrayValue(nil).(*data.ArrayValue)
}

func resolveDefault(ctx data.Context, def data.Value) (data.GetValue, data.Control) {
	if def == nil || kit.IsNull(def) {
		return data.NewNullValue(), nil
	}
	switch def.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return kit.Call(ctx, def)
	default:
		return def, nil
	}
}

func repoHas(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0)
	if key == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(collections.PathHas(repoItems(cv), key.AsString())), nil
}

func repoGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0)
	if av, ok := key.(*data.ArrayValue); ok {
		return repoGetManyKeys(ctx, cv, av)
	}
	if key == nil || kit.IsNull(key) {
		return repoItems(cv), nil
	}
	if v, ok := collections.PathGet(repoItems(cv), key.AsString()); ok {
		return v, nil
	}
	return resolveDefault(ctx, kit.Arg(ctx, 1))
}

func repoGetMany(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys, _ := kit.Arg(ctx, 0).(*data.ArrayValue)
	if keys == nil {
		return data.NewArrayValue(nil), nil
	}
	return repoGetManyKeys(ctx, cv, keys)
}

func repoGetManyKeys(ctx data.Context, cv *data.ClassValue, keys *data.ArrayValue) (data.GetValue, data.Control) {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	items := repoItems(cv)
	for _, e := range kit.Entries(keys) {
		var keyStr string
		var def data.Value
		if _, isInt := e.Key.(*data.IntValue); isInt || e.KeyStr == "" {
			keyStr = e.Value.AsString()
			def = data.NewNullValue()
		} else {
			keyStr = e.KeyStr
			def = e.Value
		}
		var v data.Value
		if got, ok := collections.PathGet(items, keyStr); ok {
			v = got
		} else {
			resolved, ctl := resolveDefault(ctx, def)
			if ctl != nil {
				return nil, ctl
			}
			v = asValue(resolved)
		}
		zv := data.NewZVal(v)
		zv.Name = keyStr
		out.List = append(out.List, zv)
	}
	return out, nil
}

func asValue(v data.GetValue) data.Value {
	if v == nil {
		return data.NewNullValue()
	}
	if val, ok := v.(data.Value); ok {
		return val
	}
	return data.NewNullValue()
}

func repoTyped(ctx data.Context, want string, check func(data.Value) bool) (data.GetValue, data.Control) {
	key := kit.Arg(ctx, 0)
	v, ctl := repoGet(ctx)
	if ctl != nil {
		return nil, ctl
	}
	val := asValue(v)
	if !check(val) {
		got := "NULL"
		if val != nil && !kit.IsNull(val) {
			got = fmt.Sprintf("%T", val)
		}
		ks := ""
		if key != nil {
			ks = key.AsString()
		}
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Configuration value for key [%s] must be a %s, %s given.", ks, want, got))
	}
	return val, nil
}

func repoString(ctx data.Context) (data.GetValue, data.Control) {
	return repoTyped(ctx, "string", func(v data.Value) bool { _, ok := v.(*data.StringValue); return ok })
}
func repoInteger(ctx data.Context) (data.GetValue, data.Control) {
	return repoTyped(ctx, "integer", func(v data.Value) bool { _, ok := v.(*data.IntValue); return ok })
}
func repoFloat(ctx data.Context) (data.GetValue, data.Control) {
	return repoTyped(ctx, "float", func(v data.Value) bool { _, ok := v.(*data.FloatValue); return ok })
}
func repoBoolean(ctx data.Context) (data.GetValue, data.Control) {
	return repoTyped(ctx, "boolean", func(v data.Value) bool { _, ok := v.(*data.BoolValue); return ok })
}
func repoArray(ctx data.Context) (data.GetValue, data.Control) {
	return repoTyped(ctx, "array", func(v data.Value) bool { _, ok := v.(*data.ArrayValue); return ok })
}

func repoCollection(ctx data.Context) (data.GetValue, data.Control) {
	arr, ctl := repoArray(ctx)
	if ctl != nil {
		return nil, ctl
	}
	vm := ctx.GetVM()
	cls, ok := vm.GetClass("Illuminate\\Support\\Collection")
	if !ok || cls == nil {
		var ctl2 data.Control
		cls, ctl2 = vm.GetOrLoadClass("Illuminate\\Support\\Collection")
		if ctl2 != nil {
			return nil, ctl2
		}
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	if ctor := cls.GetConstruct(); ctor != nil {
		nctx := cv.CreateContext(ctor.GetVariables())
		data.BindDeclaredArgs(nctx, ctor, []data.Value{asValue(arr)})
		if _, ctl := ctor.Call(nctx); ctl != nil {
			return nil, ctl
		}
	} else {
		_ = cv.SetProperty("items", asValue(arr))
	}
	return cv, nil
}

func repoSet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0)
	value := kit.Arg(ctx, 1)
	items := repoItems(cv)
	if av, ok := key.(*data.ArrayValue); ok {
		for _, e := range kit.Entries(av) {
			if e.KeyStr == "" {
				continue
			}
			collections.PathSet(items, e.KeyStr, e.Value)
		}
		return data.NewNullValue(), nil
	}
	if key != nil && !kit.IsNull(key) {
		collections.PathSet(items, key.AsString(), value)
	}
	return data.NewNullValue(), nil
}

func repoPrepend(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0)
	value := kit.Arg(ctx, 1)
	cur, _ := collections.PathGet(repoItems(cv), key.AsString())
	arr := asArrayValue(cur)
	newArr := data.NewArrayValue(nil).(*data.ArrayValue)
	newArr.List = append([]*data.ZVal{data.NewZVal(value)}, arr.List...)
	collections.PathSet(repoItems(cv), key.AsString(), newArr)
	return data.NewNullValue(), nil
}

func repoPush(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0)
	value := kit.Arg(ctx, 1)
	cur, _ := collections.PathGet(repoItems(cv), key.AsString())
	arr := asArrayValue(cur)
	newArr := data.NewArrayValue(nil).(*data.ArrayValue)
	newArr.List = append(append([]*data.ZVal{}, arr.List...), data.NewZVal(value))
	collections.PathSet(repoItems(cv), key.AsString(), newArr)
	return data.NewNullValue(), nil
}

func repoAll(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return repoItems(cv), nil
}

func repoOffsetUnset(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := repoReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0)
	if key != nil {
		collections.PathSet(repoItems(cv), key.AsString(), data.NewNullValue())
	}
	return data.NewNullValue(), nil
}

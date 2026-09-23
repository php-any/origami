package translation

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const translatorClassName = "Illuminate\\Translation\\Translator"

type TranslatorClass struct {
	node.Node
	methods map[string]data.Method
}

func NewTranslatorClass() data.ClassStmt {
	c := &TranslatorClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *TranslatorClass) GetName() string    { return translatorClassName }
func (c *TranslatorClass) GetExtend() *string { return nil }
func (c *TranslatorClass) GetImplements() []string {
	return nil
}
func (c *TranslatorClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "loader", "locale", "fallback", "loaded", "selector":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *TranslatorClass) GetPropertyList() []data.Property {
	names := []string{"loader", "locale", "fallback", "loaded", "selector"}
	out := make([]data.Property, 0, len(names))
	for _, n := range names {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *TranslatorClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *TranslatorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	_ = cv.SetProperty("locale", data.NewStringValue("en"))
	_ = cv.SetProperty("loaded", data.NewArrayValue(nil))
	return cv, nil
}
func (c *TranslatorClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *TranslatorClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *TranslatorClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *TranslatorClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"loader", "locale"}, 0, transConstruct)
	c.methods["get"] = kit.InstanceMethodOpt("get", []string{"key", "replace", "locale", "fallback"}, 1, transGet)
	c.methods["has"] = kit.InstanceMethodOpt("has", []string{"key", "locale", "fallback"}, 1, transHas)
	c.methods["setlocale"] = kit.InstanceMethod("setLocale", []string{"locale"}, transSetLocale)
}

func transRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Translation\\Translator missing $this"))
}

func transConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := transRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if l := kit.Arg(ctx, 0); l != nil {
		_ = cv.SetProperty("loader", l)
	}
	if loc := kit.Arg(ctx, 1); loc != nil {
		_ = cv.SetProperty("locale", loc)
	}
	if v, _ := cv.GetProperty("loaded"); v == nil || kit.IsNull(kit.Unwrap(v)) {
		_ = cv.SetProperty("loaded", data.NewArrayValue(nil))
	}
	return cv, nil
}

func transLines(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("loaded")
	if av, ok := kit.Unwrap(v).(*data.ArrayValue); ok && av != nil {
		return av
	}
	av := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("loaded", av)
	return av
}

func transLocale(cv *data.ClassValue) string {
	v, _ := cv.GetProperty("locale")
	if v == nil {
		return "en"
	}
	return kit.Unwrap(v).AsString()
}

func transLookup(lines *data.ArrayValue, key string) (data.Value, bool) {
	if z, ok := lines.LookupZValByStringKey(key); ok && z != nil && z.Value != nil {
		return z.Value, true
	}
	parts := strings.Split(key, ".")
	if len(parts) <= 1 {
		return nil, false
	}
	cur := data.Value(lines)
	for _, p := range parts {
		av, ok := kit.Unwrap(cur).(*data.ArrayValue)
		if !ok {
			return nil, false
		}
		z, ok := av.LookupZValByStringKey(p)
		if !ok || z == nil {
			return nil, false
		}
		cur = z.Value
	}
	return cur, true
}

func transGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := transRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Unwrap(kit.Arg(ctx, 0)).AsString()
	lines := transLines(cv)
	if v, ok := transLookup(lines, key); ok {
		return v, nil
	}
	return data.NewStringValue(key), nil
}

func transHas(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := transRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Unwrap(kit.Arg(ctx, 0)).AsString()
	lines := transLines(cv)
	_, ok := transLookup(lines, key)
	return data.NewBoolValue(ok), nil
}

func transSetLocale(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := transRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	loc := kit.Unwrap(kit.Arg(ctx, 0))
	_ = cv.SetProperty("locale", loc)
	return data.NewNullValue(), nil
}

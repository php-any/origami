package support

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const stringableClassName = "Illuminate\\Support\\Stringable"

type StringableClass struct {
	node.Node
	methods map[string]data.Method
}

func NewStringableClass() data.ClassStmt {
	c := &StringableClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *StringableClass) GetName() string                          { return stringableClassName }
func (c *StringableClass) GetExtend() *string                       { return nil }
func (c *StringableClass) GetImplements() []string                  {
	return []string{"JsonSerializable", "ArrayAccess", "Stringable"}
}
func (c *StringableClass) GetProperty(name string) (data.Property, bool) {
	if name == "value" {
		return node.NewProperty(nil, "value", "protected", false, data.NewStringValue("")), true
	}
	return nil, false
}
func (c *StringableClass) GetPropertyList() []data.Property {
	p, _ := c.GetProperty("value")
	return []data.Property{p}
}
func (c *StringableClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *StringableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *StringableClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *StringableClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *StringableClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *StringableClass) register() {
	inst := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = kit.InstanceMethod(name, params, fn)
	}
	inst("__construct", []string{"value"}, stringableConstruct)
	inst("__tostring", nil, stringableToString)
	inst("tostring", nil, stringableToString)
	inst("jsonserialize", nil, stringableToString)
	inst("after", []string{"search"}, stringableAfter)
	inst("before", []string{"search"}, stringableBefore)
	inst("camel", nil, stringableCamel)
	inst("snake", []string{"delimiter"}, stringableSnake)
	inst("studly", nil, stringableStudly)
	inst("kebab", nil, stringableKebab)
	inst("contains", []string{"needles"}, stringableContains)
	inst("lower", nil, stringableLower)
	inst("upper", nil, stringableUpper)
	inst("limit", []string{"limit", "end"}, stringableLimit)
	inst("slug", []string{"separator", "language", "dictionary"}, stringableSlug)
	inst("trim", []string{"charlist"}, stringableTrim)
	inst("length", nil, stringableLength)
	inst("substr", []string{"start", "length"}, stringableSubstr)
	inst("is", []string{"pattern"}, stringableIs)
	c.registerStringableMore()
	kit.RegisterMacroable(c.methods, stringableClassName)
	// 覆盖 Macroable::__call：先代理到 Str::{method}($value, ...$args)
	c.methods["__call"] = kit.InstanceMethod("__call", []string{"method", "parameters"}, stringableCall)
}

func newStringableValue(ctx data.Context, s string) (data.GetValue, data.Control) {
	cv := data.NewClassValue(NewStringableClass(), ctx.CreateBaseContext())
	_ = cv.SetProperty("value", data.NewStringValue(s))
	return cv, nil
}

func stringableValue(ctx data.Context) string {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return ""
	}
	v, _ := cv.GetProperty("value")
	if v == nil {
		return ""
	}
	return v.AsString()
}

func stringableConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	s := ""
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		s = v.AsString()
	}
	_ = cv.SetProperty("value", data.NewStringValue(s))
	return data.NewNullValue(), nil
}

func stringableToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(stringableValue(ctx)), nil
}

func stringableWrap(ctx data.Context, s string) (data.GetValue, data.Control) {
	return newStringableValue(ctx, s)
}

func stringableAfter(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	search := strArg(ctx, 0)
	return stringableWrap(ctx, strAfterImpl(v, search))
}

func stringableBefore(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	return stringableWrap(ctx, strBeforeImpl(v, strArg(ctx, 0)))
}

func stringableCamel(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := strCamel(withStrArg(ctx, stringableValue(ctx)))
	if ctl != nil {
		return nil, ctl
	}
	return stringableWrap(ctx, strValueString(s))
}

func stringableSnake(ctx data.Context) (data.GetValue, data.Control) {
	// delegate via snakeCached
	d := "_"
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		d = v.AsString()
	}
	return stringableWrap(ctx, snakeCached(stringableValue(ctx), d))
}

func stringableStudly(ctx data.Context) (data.GetValue, data.Control) {
	return stringableWrap(ctx, studlyCached(stringableValue(ctx), false))
}

func stringableKebab(ctx data.Context) (data.GetValue, data.Control) {
	return stringableWrap(ctx, snakeCached(stringableValue(ctx), "-"))
}

func stringableContains(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableValue(ctx)
	for _, n := range strNeedles(kit.Arg(ctx, 0)) {
		if n != "" && strings.Contains(v, n) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func stringableLower(ctx data.Context) (data.GetValue, data.Control) {
	return stringableWrap(ctx, strings.ToLower(stringableValue(ctx)))
}

func stringableUpper(ctx data.Context) (data.GetValue, data.Control) {
	return stringableWrap(ctx, strings.ToUpper(stringableValue(ctx)))
}

func stringableLimit(ctx data.Context) (data.GetValue, data.Control) {
	inner := withStrArgs(ctx, data.NewStringValue(stringableValue(ctx)), kit.Arg(ctx, 0), kit.Arg(ctx, 1))
	s, ctl := strLimit(inner)
	if ctl != nil {
		return nil, ctl
	}
	return stringableWrap(ctx, strValueString(s))
}

func stringableSlug(ctx data.Context) (data.GetValue, data.Control) {
	inner := withStrArgs(ctx, data.NewStringValue(stringableValue(ctx)), kit.Arg(ctx, 0), kit.Arg(ctx, 1), kit.Arg(ctx, 2))
	s, ctl := strSlug(inner)
	if ctl != nil {
		return nil, ctl
	}
	return stringableWrap(ctx, strValueString(s))
}

func stringableTrim(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := strTrim(withStrArg(ctx, stringableValue(ctx)))
	if ctl != nil {
		return nil, ctl
	}
	return stringableWrap(ctx, strValueString(s))
}

func stringableLength(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(mbLen(stringableValue(ctx))), nil
}

func stringableSubstr(ctx data.Context) (data.GetValue, data.Control) {
	inner := withStrArgs(ctx, data.NewStringValue(stringableValue(ctx)), kit.Arg(ctx, 0), kit.Arg(ctx, 1))
	s, ctl := strSubstr(inner)
	if ctl != nil {
		return nil, ctl
	}
	return stringableWrap(ctx, strValueString(s))
}

func stringableIs(ctx data.Context) (data.GetValue, data.Control) {
	return strIs(withStrArgs(ctx, kit.Arg(ctx, 0), data.NewStringValue(stringableValue(ctx))))
}

// stringableCall 对齐 Macroable + 未注册实例方法：优先 Str 静态，再走宏。
func stringableCall(ctx data.Context) (data.GetValue, data.Control) {
	method := ""
	if v := kit.Arg(ctx, 0); v != nil {
		method = v.AsString()
	}
	if method == "" {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Method %s:: does not exist.", stringableClassName))
	}
	var extra []data.Value
	if params := kit.Arg(ctx, 1); params != nil {
		if av, ok := kit.Unwrap(params).(*data.ArrayValue); ok && av != nil {
			extra = av.ToValueList()
		}
	}
	strStmt := strClassForStringable()
	if m, ok := strStmt.GetStaticMethod(method); ok && m != nil {
		return stringableInvokeStr(ctx, method, extra)
	}
	return kit.CallMacro(ctx, stringableClassName, false)
}

// withStrArg 构造单参数 Str 静态调用的临时帧（仅测试/Stringable 委托用，非热路径）。
type strArgFrame struct {
	data.Context
	args []data.Value
}

func (f *strArgFrame) GetIndexValue(i int) (data.Value, bool) {
	if i < 0 || i >= len(f.args) {
		return nil, false
	}
	return f.args[i], true
}

func withStrArg(ctx data.Context, a0 string) data.Context {
	return &strArgFrame{Context: ctx, args: []data.Value{data.NewStringValue(a0)}}
}

func withStrArgs(ctx data.Context, args ...data.Value) data.Context {
	return &strArgFrame{Context: ctx, args: args}
}

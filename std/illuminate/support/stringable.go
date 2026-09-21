package support

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const stringableClassName = "Illuminate\\Support\\Stringable"

type StringableClass struct {
	node.Node
	methods map[string]data.Method
}

func NewStringableClass() data.ClassStmt {
	c := &StringableClass{methods: map[string]data.Method{}}
	add := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = newInstanceMethod(name, params, fn, false)
	}
	add("__construct", []string{"value"}, stringableConstruct)
	add("__toString", nil, stringableToString)
	add("toString", nil, stringableToString)
	add("value", nil, stringableToString)
	add("after", []string{"search"}, stringableAfter)
	add("append", []string{"values"}, stringableAppend)
	add("lower", nil, stringableLower)
	add("upper", nil, stringableUpper)
	add("camel", nil, stringableCamel)
	add("snake", []string{"delimiter"}, stringableSnake)
	add("kebab", nil, stringableKebab)
	add("studly", nil, stringableStudly)
	add("finish", []string{"cap"}, stringableFinish)
	add("start", []string{"prefix"}, stringableStart)
	add("contains", []string{"needles"}, stringableContains)
	add("startsWith", []string{"needles"}, stringableStartsWith)
	add("endsWith", []string{"needles"}, stringableEndsWith)
	add("replace", []string{"search", "replace"}, stringableReplace)
	add("substr", []string{"start", "length"}, stringableSubstr)
	add("length", nil, stringableLength)
	add("isEmpty", nil, stringableIsEmpty)
	add("isNotEmpty", nil, stringableIsNotEmpty)
	add("slug", []string{"separator"}, stringableSlug)
	add("is", []string{"pattern", "ignoreCase"}, stringableIs)
	add("when", []string{"condition", "callback", "default"}, stringableWhen)
	add("__call", []string{"method", "parameters"}, stringableCall)
	registerMacroable(c.methods, stringableClassName)
	return c
}

func (c *StringableClass) GetName() string { return stringableClassName }
func (c *StringableClass) GetExtend() *string {
	return nil
}
func (c *StringableClass) GetImplements() []string {
	return []string{"JsonSerializable", "ArrayAccess", "Stringable"}
}
func (c *StringableClass) GetProperty(name string) (data.Property, bool) {
	if name == "value" {
		return node.NewProperty(nil, "value", "protected", false, data.NewStringValue("")), true
	}
	return nil, false
}
func (c *StringableClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "value", "protected", false, data.NewStringValue(""))}
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

func newStringableValue(ctx data.Context, s string) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewStringValue(s), nil
	}
	cls, ok := vm.GetClass(stringableClassName)
	if !ok {
		return data.NewStringValue(s), nil
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	_ = cv.SetProperty("value", data.NewStringValue(s))
	return cv, nil
}

func stringableStr(ctx data.Context) string {
	cv := overlayReceiver(ctx)
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
	cv := overlayReceiver(ctx)
	s := ""
	if v, ok := ctx.GetIndexValue(0); ok && v != nil && !isNull(v) {
		s = v.AsString()
	}
	_ = cv.SetProperty("value", data.NewStringValue(s))
	return cv, nil
}

func stringableToString(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(stringableStr(ctx)), nil
}

func stringableAfter(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, afterString(stringableStr(ctx), strArg(ctx, 0)))
}

func afterString(subject, search string) string {
	if search == "" {
		return subject
	}
	i := strings.Index(subject, search)
	if i < 0 {
		return subject
	}
	return subject[i+len(search):]
}

func stringableAppend(ctx data.Context) (data.GetValue, data.Control) {
	s := stringableStr(ctx)
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		s += v.AsString()
	}
	return newStringableValue(ctx, s)
}

func stringableLower(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, strings.ToLower(stringableStr(ctx)))
}
func stringableUpper(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, strings.ToUpper(stringableStr(ctx)))
}
func stringableCamel(ctx data.Context) (data.GetValue, data.Control) {
	s := studly(stringableStr(ctx))
	if s == "" {
		return newStringableValue(ctx, "")
	}
	r := []rune(s)
	r[0] = toLowerRune(r[0])
	return newStringableValue(ctx, string(r))
}

func toLowerRune(r rune) rune {
	return []rune(strings.ToLower(string(r)))[0]
}

func stringableSnake(ctx data.Context) (data.GetValue, data.Control) {
	delim := "_"
	if v, ok := ctx.GetIndexValue(0); ok && v != nil && !isNull(v) {
		delim = v.AsString()
	}
	return newStringableValue(ctx, snakeString(stringableStr(ctx), delim))
}
func stringableKebab(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, kebabString(stringableStr(ctx)))
}
func stringableStudly(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, studly(stringableStr(ctx)))
}
func stringableFinish(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableStr(ctx)
	cap := strArg(ctx, 0)
	for cap != "" && strings.HasSuffix(v, cap) {
		v = strings.TrimSuffix(v, cap)
	}
	return newStringableValue(ctx, v+cap)
}
func stringableStart(ctx data.Context) (data.GetValue, data.Control) {
	v := stringableStr(ctx)
	p := strArg(ctx, 0)
	if p != "" && !strings.HasPrefix(v, p) {
		v = p + v
	}
	return newStringableValue(ctx, v)
}
func stringableContains(ctx data.Context) (data.GetValue, data.Control) {
	h := stringableStr(ctx)
	for _, n := range strNeedles(indexVal(ctx, 0)) {
		if n != "" && strings.Contains(h, n) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}
func stringableStartsWith(ctx data.Context) (data.GetValue, data.Control) {
	h := stringableStr(ctx)
	for _, n := range strNeedles(indexVal(ctx, 0)) {
		if n != "" && strings.HasPrefix(h, n) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}
func stringableEndsWith(ctx data.Context) (data.GetValue, data.Control) {
	h := stringableStr(ctx)
	for _, n := range strNeedles(indexVal(ctx, 0)) {
		if n != "" && strings.HasSuffix(h, n) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}
func stringableReplace(ctx data.Context) (data.GetValue, data.Control) {
	return newStringableValue(ctx, strings.ReplaceAll(stringableStr(ctx), strArg(ctx, 0), strArg(ctx, 1)))
}
func stringableSubstr(ctx data.Context) (data.GetValue, data.Control) {
	s, _ := strSubstr(fakeStrCtx(ctx, stringableStr(ctx), indexVal(ctx, 0), indexVal(ctx, 1)))
	if sv, ok := s.(data.Value); ok {
		return newStringableValue(ctx, sv.AsString())
	}
	return newStringableValue(ctx, "")
}

func fakeStrCtx(ctx data.Context, s string, start, length data.Value) data.Context {
	vars := []data.Variable{
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "start", 1, nil),
		node.NewVariable(nil, "length", 2, nil),
	}
	n := ctx.GetVM().CreateContext(vars)
	_ = vars[0].SetValue(n, data.NewStringValue(s))
	if start != nil {
		_ = vars[1].SetValue(n, start)
	}
	if length != nil {
		_ = vars[2].SetValue(n, length)
	}
	return n
}

func indexVal(ctx data.Context, i int) data.Value {
	v, _ := ctx.GetIndexValue(i)
	return v
}

func stringableLength(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(len([]rune(stringableStr(ctx)))), nil
}
func stringableIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(stringableStr(ctx) == ""), nil
}
func stringableIsNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(stringableStr(ctx) != ""), nil
}
func stringableSlug(ctx data.Context) (data.GetValue, data.Control) {
	sep := "-"
	if v := strArg(ctx, 0); v != "" {
		sep = v
	}
	got, _ := strSlug(fakeSlugCtx(ctx, stringableStr(ctx), sep))
	if sv, ok := got.(data.Value); ok {
		return newStringableValue(ctx, sv.AsString())
	}
	return newStringableValue(ctx, "")
}

func fakeSlugCtx(ctx data.Context, title, sep string) data.Context {
	vars := []data.Variable{
		node.NewVariable(nil, "title", 0, nil),
		node.NewVariable(nil, "separator", 1, nil),
	}
	n := ctx.GetVM().CreateContext(vars)
	_ = vars[0].SetValue(n, data.NewStringValue(title))
	_ = vars[1].SetValue(n, data.NewStringValue(sep))
	return n
}

func stringableWhen(ctx data.Context) (data.GetValue, data.Control) {
	cond, _ := ctx.GetIndexValue(0)
	truthy := false
	if cond != nil {
		if b, ok := cond.(data.AsBool); ok {
			truthy, _ = b.AsBool()
		} else {
			truthy = cond.AsString() != "" && cond.AsString() != "0"
		}
	}
	cv := overlayReceiver(ctx)
	if truthy {
		cb, _ := ctx.GetIndexValue(1)
		if cb != nil && !isNull(cb) {
			return callValue(ctx, cb, cv)
		}
	} else {
		def, _ := ctx.GetIndexValue(2)
		if def != nil && !isNull(def) {
			return callValue(ctx, def, cv)
		}
	}
	return cv, nil
}

func stringableIs(ctx data.Context) (data.GetValue, data.Control) {
	pattern := strArg(ctx, 0)
	value := stringableStr(ctx)
	if ignoreCase(ctx, 1) {
		pattern = strings.ToLower(pattern)
		value = strings.ToLower(value)
	}
	if pattern == value {
		return data.NewBoolValue(true), nil
	}
	if strings.Contains(pattern, "*") {
		parts := strings.Split(pattern, "*")
		if len(parts) == 2 {
			return data.NewBoolValue(strings.HasPrefix(value, parts[0]) && strings.HasSuffix(value, parts[1])), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func ignoreCase(ctx data.Context, i int) bool {
	v, ok := ctx.GetIndexValue(i)
	if !ok || v == nil {
		return false
	}
	if b, ok := v.(data.AsBool); ok {
		okb, _ := b.AsBool()
		return okb
	}
	return false
}

func stringableCall(ctx data.Context) (data.GetValue, data.Control) {
	method := strArg(ctx, 0)
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewNullValue(), nil
	}
	cls, ok := vm.GetClass(strClassName)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("类(%s)不存在对应函数(%s)", stringableClassName, method))
	}
	gsm, ok := cls.(data.GetStaticMethod)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("类(%s)不存在对应函数(%s)", stringableClassName, method))
	}
	m, ok := gsm.GetStaticMethod(method)
	if !ok || m == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("类(%s)不存在对应函数(%s)", stringableClassName, method))
	}
	args := []data.Value{data.NewStringValue(stringableStr(ctx))}
	if params, ok := ctx.GetIndexValue(1); ok {
		for _, e := range toEntries(params) {
			args = append(args, e.value)
		}
	}
	nctx := vm.CreateContext(m.GetVariables())
	vars := m.GetVariables()
	for i, a := range args {
		if i >= len(vars) {
			break
		}
		if ctl := vars[i].SetValue(nctx, a); ctl != nil {
			return nil, ctl
		}
	}
	got, ctl := m.Call(nctx)
	if ctl != nil {
		return nil, ctl
	}
	if sv, ok := got.(*data.StringValue); ok {
		return newStringableValue(ctx, sv.Value)
	}
	return got, nil
}

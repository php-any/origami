package support

import (
	"html"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const attributeBagClassName = "Illuminate\\View\\ComponentAttributeBag"

type AttributeBagClass struct {
	node.Node
	methods map[string]data.Method
}

func NewAttributeBagClass() data.ClassStmt {
	c := &AttributeBagClass{methods: map[string]data.Method{}}
	add := func(name string, params []string, optionalFrom int, static bool, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[strings.ToLower(name)] = newBagMethod(name, params, optionalFrom, static, fn)
	}
	add("__construct", []string{"attributes"}, 0, false, bagConstruct)
	add("extractPropNames", []string{"keys"}, -1, true, bagExtractPropNames)
	add("all", []string{"keys"}, 0, false, bagAll)
	add("get", []string{"key", "default"}, 1, false, bagGet)
	add("getAttributes", nil, -1, false, bagGetAttributes)
	add("setAttributes", []string{"attributes"}, -1, false, bagSetAttributes)
	add("class", []string{"classList"}, -1, false, bagClass)
	add("merge", []string{"attributeDefaults", "escape"}, 0, false, bagMerge)
	add("toHtml", nil, -1, false, bagToHtml)
	add("__toString", nil, -1, false, bagToHtml)
	add("toArray", nil, -1, false, bagGetAttributes)
	add("isEmpty", nil, -1, false, bagIsEmpty)
	add("isNotEmpty", nil, -1, false, bagIsNotEmpty)
	add("offsetExists", []string{"offset"}, -1, false, bagOffsetExists)
	add("offsetGet", []string{"offset"}, -1, false, bagOffsetGet)
	add("offsetSet", []string{"offset", "value"}, -1, false, bagOffsetSet)
	add("offsetUnset", []string{"offset"}, -1, false, bagOffsetUnset)
	add("jsonSerialize", nil, -1, false, bagGetAttributes)
	add("has", []string{"key"}, -1, false, bagHas)
	add("exists", []string{"key"}, -1, false, bagHas)
	add("hasAny", []string{"keys"}, -1, false, bagHasAny)
	add("missing", []string{"key"}, -1, false, bagMissing)
	add("filled", []string{"key"}, -1, false, bagFilled)
	add("only", []string{"keys"}, 0, false, bagOnly)
	add("except", []string{"keys"}, 0, false, bagExcept)
	add("first", []string{"default"}, 0, false, bagFirst)
	add("filter", []string{"callback"}, 0, false, bagFilter)
	add("when", []string{"condition", "callback", "default"}, 1, false, bagWhen)
	add("unless", []string{"condition", "callback", "default"}, 1, false, bagUnless)
	add("getIterator", nil, -1, false, bagGetIterator)
	add("whereStartsWith", []string{"needles"}, -1, false, bagWhereStartsWith)
	add("whereDoesntStartWith", []string{"needles"}, -1, false, bagWhereDoesntStartWith)
	add("thatStartWith", []string{"needles"}, -1, false, bagWhereStartsWith)
	add("onlyProps", []string{"keys"}, -1, false, bagOnlyProps)
	add("exceptProps", []string{"keys"}, -1, false, bagExceptProps)
	add("style", []string{"styleList"}, -1, false, bagStyle)
	add("shouldEscapeAttributeValue", []string{"escape", "value"}, -1, false, bagShouldEscapeAttributeValue)
	add("prepends", []string{"value"}, -1, false, bagPrepends)
	add("resolveAppendableAttributeDefault", []string{"attributeDefaults", "key", "escape"}, -1, false, bagResolveAppendableDefault)
	registerMacroable(c.methods, attributeBagClassName)
	c.methods["__call"] = newInstanceMethod("__call", []string{"method", "parameters"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return macroInvoke(attributeBagClassName, ctx, overlayReceiver(ctx))
	}, false)
	return c
}

func newBagMethod(name string, params []string, optionalFrom int, static bool, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	ps := make([]data.GetValue, len(params))
	vs := make([]data.Variable, len(params))
	for i, p := range params {
		var def data.GetValue
		if optionalFrom >= 0 && i >= optionalFrom {
			if p == "escape" {
				def = data.NewBoolValue(true)
			} else {
				def = data.NewNullValue()
			}
		}
		ps[i] = node.NewParameter(nil, p, i, def, nil)
		vs[i] = node.NewVariable(nil, p, i, nil)
	}
	return &collMethod{name: name, params: ps, vars: vs, fn: fn, static: static}
}

func (c *AttributeBagClass) GetName() string { return attributeBagClassName }
func (c *AttributeBagClass) GetExtend() *string {
	return nil
}
func (c *AttributeBagClass) GetImplements() []string {
	return []string{
		"Illuminate\\Contracts\\Support\\Arrayable",
		"Illuminate\\Contracts\\Support\\Htmlable",
		"ArrayAccess",
		"IteratorAggregate",
		"JsonSerializable",
		"Stringable",
	}
}
func (c *AttributeBagClass) GetProperty(name string) (data.Property, bool) {
	if name == "attributes" {
		return node.NewProperty(nil, "attributes", "protected", false, data.NewArrayValue(nil)), true
	}
	return nil, false
}
func (c *AttributeBagClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "attributes", "protected", false, data.NewArrayValue(nil))}
}
func (c *AttributeBagClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *AttributeBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *AttributeBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *AttributeBagClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *AttributeBagClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func newBagValue(ctx data.Context, attrs *data.ArrayValue) *data.ClassValue {
	vm := ctx.GetVM()
	cls, _ := vm.GetClass(attributeBagClassName)
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	if attrs == nil {
		attrs = data.NewArrayValue(nil).(*data.ArrayValue)
	}
	_ = cv.SetProperty("attributes", attrs)
	return cv
}

func bagConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := overlayReceiver(ctx)
	attrs, _ := ctx.GetIndexValue(0)
	if attrs == nil || isNull(attrs) {
		attrs = data.NewArrayValue(nil)
	}
	if av, ok := unwrapValue(attrs).(*data.ArrayValue); ok {
		if nested, ok := arrayGet(av, "attributes"); ok {
			if ncv, ok := nested.(*data.ClassValue); ok && ncv.Class != nil && ncv.Class.GetName() == attributeBagClassName {
				base := cloneArrayValue(bagAttrs(ncv))
				av.UnsetKey(data.NewStringValue("attributes"))
				for _, e := range toEntries(av) {
					base.SetStringKey(e.keyStr, e.value)
				}
				attrs = base
			}
		}
	}
	_ = cv.SetProperty("attributes", unwrapValue(attrs))
	return cv, nil
}

func bagAttrs(cv *data.ClassValue) *data.ArrayValue {
	if cv == nil {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	v, _ := cv.GetProperty("attributes")
	if av, ok := unwrapValue(v).(*data.ArrayValue); ok {
		return av
	}
	empty := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("attributes", empty)
	return empty
}

func bagExtractPropNames(ctx data.Context) (data.GetValue, data.Control) {
	keys, _ := ctx.GetIndexValue(0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(keys) {
		key := e.keyStr
		if cssAssocKeyIsNumeric(key) {
			key = e.value.AsString()
		}
		out.AppendValue(data.NewStringValue(key))
		out.AppendValue(data.NewStringValue(kebabString(key)))
	}
	return out, nil
}

func bagGetAttributes(ctx data.Context) (data.GetValue, data.Control) {
	return bagAttrs(overlayReceiver(ctx)), nil
}

func bagSetAttributes(ctx data.Context) (data.GetValue, data.Control) {
	return bagConstruct(ctx)
}

func bagAll(ctx data.Context) (data.GetValue, data.Control) {
	keys, ok := ctx.GetIndexValue(0)
	if !ok || keys == nil || isNull(keys) {
		return bagGetAttributes(ctx)
	}
	return bagOnlyKeys(bagAttrs(overlayReceiver(ctx)), keys), nil
}

func bagGet(ctx data.Context) (data.GetValue, data.Control) {
	key, _ := ctx.GetIndexValue(0)
	def, _ := ctx.GetIndexValue(1)
	if v, ok := arrayGet(bagAttrs(overlayReceiver(ctx)), keyToString(key)); ok {
		return v, nil
	}
	if def == nil {
		return data.NewNullValue(), nil
	}
	return laravelValue(ctx, def)
}

func bagShouldEscapeAttributeValue(ctx data.Context) (data.GetValue, data.Control) {
	escape := true
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			escape, _ = b.AsBool()
		}
	}
	if !escape {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(shouldEscapeAttr(indexVal(ctx, 1))), nil
}

func bagPrepends(ctx data.Context) (data.GetValue, data.Control) {
	return newAppendableValue(ctx, indexVal(ctx, 0))
}

func bagResolveAppendableDefault(ctx data.Context) (data.GetValue, data.Control) {
	defs, _ := ctx.GetIndexValue(0)
	key := strArg(ctx, 1)
	escape := true
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			escape, _ = b.AsBool()
		}
	}
	av, _ := unwrapValue(defs).(*data.ArrayValue)
	val, _ := arrayGet(av, key)
	if cv, ok := val.(*data.ClassValue); ok {
		if inner, ctl := cv.GetProperty("value"); ctl == nil && inner != nil {
			val = inner
		}
	}
	if escape && shouldEscapeAttr(val) {
		return data.NewStringValue(html.EscapeString(val.AsString())), nil
	}
	if val == nil {
		return data.NewNullValue(), nil
	}
	return val, nil
}

func bagClass(ctx data.Context) (data.GetValue, data.Control) {
	list, _ := ctx.GetIndexValue(0)
	css := cssClassesFromValue(arrWrapValue(list))
	defs := data.NewArrayValue(nil).(*data.ArrayValue)
	defs.SetStringKey("class", data.NewStringValue(css))
	return bagMergeValues(ctx, overlayReceiver(ctx), defs, true)
}

func bagMerge(ctx data.Context) (data.GetValue, data.Control) {
	defs, _ := ctx.GetIndexValue(0)
	escape := true
	if v, ok := ctx.GetIndexValue(1); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			escape, _ = b.AsBool()
		}
	}
	av, _ := unwrapValue(defs).(*data.ArrayValue)
	if av == nil {
		av = data.NewArrayValue(nil).(*data.ArrayValue)
	}
	return bagMergeValues(ctx, overlayReceiver(ctx), av, escape)
}

func bagMergeValues(ctx data.Context, cv *data.ClassValue, defaults *data.ArrayValue, escape bool) (data.GetValue, data.Control) {
	cur := bagAttrs(cv)
	escaped := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(defaults) {
		val := e.value
		if escape && shouldEscapeAttr(val) {
			val = data.NewStringValue(html.EscapeString(val.AsString()))
		}
		escaped.SetStringKey(e.keyStr, val)
	}
	out := cloneArrayValue(escaped)
	for _, e := range toEntries(cur) {
		key := e.keyStr
		if key == "class" || key == "style" {
			def, _ := arrayGet(escaped, key)
			combined := joinUnique(def, e.value, key == "style")
			out.SetStringKey(key, data.NewStringValue(combined))
			continue
		}
		out.SetStringKey(key, e.value)
	}
	return newBagValue(ctx, out), nil
}

func shouldEscapeAttr(v data.Value) bool {
	if v == nil || isNull(v) {
		return false
	}
	if _, ok := v.(*data.ClassValue); ok {
		return false
	}
	if _, ok := v.(*data.BoolValue); ok {
		return false
	}
	return true
}

func joinUnique(def, cur data.Value, style bool) string {
	parts := []string{}
	add := func(s string) {
		s = strings.TrimSpace(s)
		if style && s != "" && !strings.HasSuffix(s, ";") {
			s += ";"
		}
		if s == "" {
			return
		}
		for _, p := range parts {
			if p == s {
				return
			}
		}
		parts = append(parts, s)
	}
	if def != nil && !isNull(def) {
		add(def.AsString())
	}
	if cur != nil && !isNull(cur) {
		add(cur.AsString())
	}
	return strings.Join(parts, " ")
}

func bagToHtml(ctx data.Context) (data.GetValue, data.Control) {
	var b strings.Builder
	for _, e := range toEntries(bagAttrs(overlayReceiver(ctx))) {
		if isNull(e.value) {
			continue
		}
		if bv, ok := e.value.(*data.BoolValue); ok {
			if !bv.Value {
				continue
			}
			val := e.keyStr
			if e.keyStr == "x-data" || strings.HasPrefix(e.keyStr, "wire:") {
				val = ""
			}
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(e.keyStr)
			b.WriteString(`="`)
			b.WriteString(strings.ReplaceAll(strings.TrimSpace(val), `"`, `\"`))
			b.WriteString(`"`)
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(e.keyStr)
		b.WriteString(`="`)
		b.WriteString(strings.ReplaceAll(strings.TrimSpace(e.value.AsString()), `"`, `\"`))
		b.WriteString(`"`)
	}
	return data.NewStringValue(b.String()), nil
}

func bagIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	s, ctl := bagToHtml(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(strings.TrimSpace(s.(data.Value).AsString()) == ""), nil
}

func bagIsNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	empty, ctl := bagIsEmpty(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := empty.(*data.BoolValue); ok {
		return data.NewBoolValue(!bv.Value), nil
	}
	return data.NewBoolValue(true), nil
}

func bagOffsetExists(ctx data.Context) (data.GetValue, data.Control) {
	_, ok := arrayGet(bagAttrs(overlayReceiver(ctx)), strArg(ctx, 0))
	return data.NewBoolValue(ok), nil
}
func bagOffsetGet(ctx data.Context) (data.GetValue, data.Control) { return bagGet(ctx) }
func bagOffsetSet(ctx data.Context) (data.GetValue, data.Control) {
	attrs := bagAttrs(overlayReceiver(ctx))
	attrs.SetStringKey(strArg(ctx, 0), indexVal(ctx, 1))
	return data.NewNullValue(), nil
}
func bagOffsetUnset(ctx data.Context) (data.GetValue, data.Control) {
	bagAttrs(overlayReceiver(ctx)).UnsetKey(data.NewStringValue(strArg(ctx, 0)))
	return data.NewNullValue(), nil
}

func arrWrapValue(v data.Value) data.Value {
	if v == nil || isNull(v) {
		return data.NewArrayValue(nil)
	}
	if _, ok := v.(*data.ArrayValue); ok {
		return v
	}
	return data.NewArrayValue([]data.Value{v})
}

func cssClassesFromValue(v data.Value) string {
	parts := make([]string, 0)
	for _, e := range toEntries(v) {
		if cssAssocKeyIsNumeric(e.keyStr) {
			if s := strings.TrimSpace(e.value.AsString()); s != "" {
				parts = append(parts, s)
			}
			continue
		}
		if truthy(e.value) {
			if s := strings.TrimSpace(e.keyStr); s != "" {
				parts = append(parts, s)
			}
		}
	}
	return strings.Join(parts, " ")
}

func cloneArrayValue(src *data.ArrayValue) *data.ArrayValue {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if src == nil {
		return out
	}
	for _, e := range toEntries(src) {
		if e.keyStr != "" && !cssAssocKeyIsNumeric(e.keyStr) {
			out.SetStringKey(e.keyStr, e.value)
		} else {
			out.AppendValue(e.value)
		}
	}
	return out
}

func bagGetIterator(ctx data.Context) (data.GetValue, data.Control) {
	attrs := bagAttrs(overlayReceiver(ctx))
	vm := ctx.GetVM()
	if vm == nil {
		return attrs, nil
	}
	stmt, ok := vm.GetClass("ArrayIterator")
	if !ok || stmt == nil {
		return attrs, nil
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return attrs, nil
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		return cv, nil
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	if len(method.GetVariables()) > 0 {
		_ = fnCtx.SetVariableValue(method.GetVariables()[0], attrs)
	}
	if _, acl := method.Call(fnCtx); acl != nil {
		return nil, acl
	}
	return cv, nil
}

func bagKeyHasPrefix(key string, needles []string) bool {
	for _, n := range needles {
		if n != "" && strings.HasPrefix(key, n) {
			return true
		}
	}
	return false
}

func bagWhereStartsWith(ctx data.Context) (data.GetValue, data.Control) {
	return bagFilterByKeyPrefix(ctx, true)
}

func bagWhereDoesntStartWith(ctx data.Context) (data.GetValue, data.Control) {
	return bagFilterByKeyPrefix(ctx, false)
}

func bagFilterByKeyPrefix(ctx data.Context, want bool) (data.GetValue, data.Control) {
	needles := strNeedles(indexVal(ctx, 0))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(bagAttrs(overlayReceiver(ctx))) {
		if bagKeyHasPrefix(e.keyStr, needles) == want {
			out.SetStringKey(e.keyStr, e.value)
		}
	}
	return newBagValue(ctx, out), nil
}

func bagOnlyProps(ctx data.Context) (data.GetValue, data.Control) {
	names, ctl := bagExtractPropNames(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out, _ := bagOnlyKeys(bagAttrs(overlayReceiver(ctx)), names.(data.Value)).(*data.ArrayValue)
	return newBagValue(ctx, out), nil
}

func bagExceptProps(ctx data.Context) (data.GetValue, data.Control) {
	names, ctl := bagExtractPropNames(ctx)
	if ctl != nil {
		return nil, ctl
	}
	drop := map[string]struct{}{}
	for _, e := range toEntries(names.(data.Value)) {
		drop[e.value.AsString()] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(bagAttrs(overlayReceiver(ctx))) {
		if _, skip := drop[e.keyStr]; skip {
			continue
		}
		out.SetStringKey(e.keyStr, e.value)
	}
	return newBagValue(ctx, out), nil
}

func bagStyle(ctx data.Context) (data.GetValue, data.Control) {
	list, _ := ctx.GetIndexValue(0)
	css, ctl := arrToCssStyles(fakeOneArgCtx(ctx, arrWrapValue(list)))
	if ctl != nil {
		return nil, ctl
	}
	defs := data.NewArrayValue(nil).(*data.ArrayValue)
	if sv, ok := css.(data.Value); ok {
		defs.SetStringKey("style", sv)
	}
	return bagMergeValues(ctx, overlayReceiver(ctx), defs, true)
}

func fakeOneArgCtx(ctx data.Context, v data.Value) data.Context {
	vars := []data.Variable{node.NewVariable(nil, "array", 0, nil)}
	n := ctx.CreateContext(vars)
	_ = vars[0].SetValue(n, v)
	return n
}

func bagKeyArgs(ctx data.Context) []string {
	v, ok := ctx.GetIndexValue(0)
	if !ok || v == nil || isNull(v) {
		return nil
	}
	out := strNeedles(v)
	for i := 1; ; i++ {
		x, ok := ctx.GetIndexValue(i)
		if !ok || x == nil {
			break
		}
		out = append(out, strNeedles(x)...)
	}
	return out
}

func bagHas(ctx data.Context) (data.GetValue, data.Control) {
	attrs := bagAttrs(overlayReceiver(ctx))
	for _, k := range bagKeyArgs(ctx) {
		if _, ok := dataGetPath(attrs, k); !ok {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func bagHasAny(ctx data.Context) (data.GetValue, data.Control) {
	attrs := bagAttrs(overlayReceiver(ctx))
	for _, k := range bagKeyArgs(ctx) {
		if _, ok := dataGetPath(attrs, k); ok {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func bagMissing(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := bagHas(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if bv, ok := v.(*data.BoolValue); ok {
		return data.NewBoolValue(!bv.Value), nil
	}
	return data.NewBoolValue(true), nil
}

func bagFilled(ctx data.Context) (data.GetValue, data.Control) {
	attrs := bagAttrs(overlayReceiver(ctx))
	for _, k := range bagKeyArgs(ctx) {
		v, ok := dataGetPath(attrs, k)
		if !ok || v == nil || isNull(v) || strings.TrimSpace(v.AsString()) == "" {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func bagOnly(ctx data.Context) (data.GetValue, data.Control) {
	keys, _ := ctx.GetIndexValue(0)
	attrs := bagAttrs(overlayReceiver(ctx))
	if keys == nil || isNull(keys) {
		return newBagValue(ctx, cloneArrayValue(attrs)), nil
	}
	out, _ := bagOnlyKeys(attrs, keys).(*data.ArrayValue)
	return newBagValue(ctx, out), nil
}

func bagExcept(ctx data.Context) (data.GetValue, data.Control) {
	keys, _ := ctx.GetIndexValue(0)
	attrs := bagAttrs(overlayReceiver(ctx))
	if keys == nil || isNull(keys) {
		return newBagValue(ctx, cloneArrayValue(attrs)), nil
	}
	drop := map[string]struct{}{}
	for _, e := range toEntries(keys) {
		drop[e.value.AsString()] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(attrs) {
		if _, skip := drop[e.keyStr]; skip {
			continue
		}
		if e.keyStr != "" && !cssAssocKeyIsNumeric(e.keyStr) {
			out.SetStringKey(e.keyStr, e.value)
		} else {
			out.AppendValue(e.value)
		}
	}
	return newBagValue(ctx, out), nil
}

func bagFirst(ctx data.Context) (data.GetValue, data.Control) {
	ents := toEntries(bagAttrs(overlayReceiver(ctx)))
	if len(ents) > 0 {
		return ents[0].value, nil
	}
	def, _ := ctx.GetIndexValue(0)
	return laravelValue(ctx, def)
}

func bagFilter(ctx data.Context) (data.GetValue, data.Control) {
	cb, ok := ctx.GetIndexValue(0)
	attrs := bagAttrs(overlayReceiver(ctx))
	if !ok || cb == nil || isNull(cb) {
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, e := range toEntries(attrs) {
			if !isEmptyFilterValue(e.value) {
				out.SetStringKey(e.keyStr, e.value)
			}
		}
		return newBagValue(ctx, out), nil
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(attrs) {
		got, ctl := callValue(ctx, cb, e.value, data.NewStringValue(e.keyStr))
		if ctl != nil {
			return nil, ctl
		}
		keep := false
		if got != nil {
			if b, ok := got.(data.AsBool); ok {
				keep, _ = b.AsBool()
			} else if v, ok := got.(data.Value); ok {
				keep = !isEmptyFilterValue(v)
			}
		}
		if keep {
			out.SetStringKey(e.keyStr, e.value)
		}
	}
	return newBagValue(ctx, out), nil
}

func isEmptyFilterValue(v data.Value) bool {
	if v == nil || isNull(v) {
		return true
	}
	if b, ok := v.(data.AsBool); ok {
		if bv, err := b.AsBool(); err == nil && !bv {
			return true
		}
	}
	return v.AsString() == ""
}

func bagWhen(ctx data.Context) (data.GetValue, data.Control) {
	return bagWhenUnless(ctx, false)
}

func bagUnless(ctx data.Context) (data.GetValue, data.Control) {
	return bagWhenUnless(ctx, true)
}

func bagWhenUnless(ctx data.Context, invert bool) (data.GetValue, data.Control) {
	recv := overlayReceiver(ctx)
	cond, _ := ctx.GetIndexValue(0)
	truthy := false
	if cond != nil {
		if b, ok := cond.(data.AsBool); ok {
			truthy, _ = b.AsBool()
		} else {
			truthy = cond.AsString() != "" && cond.AsString() != "0"
		}
	}
	if invert {
		truthy = !truthy
	}
	cb, _ := ctx.GetIndexValue(1)
	def, _ := ctx.GetIndexValue(2)
	if truthy && cb != nil && !isNull(cb) {
		got, ctl := callValue(ctx, cb, recv)
		if ctl != nil {
			return nil, ctl
		}
		if gv, ok := got.(data.Value); !ok || gv == nil || isNull(gv) {
			return recv, nil
		}
		return got, nil
	}
	if !truthy && def != nil && !isNull(def) {
		return callValue(ctx, def, recv)
	}
	return recv, nil
}

func bagOnlyKeys(arr *data.ArrayValue, keys data.Value) data.Value {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	want := map[string]struct{}{}
	for _, e := range toEntries(keys) {
		want[e.value.AsString()] = struct{}{}
	}
	for _, e := range toEntries(arr) {
		if _, ok := want[e.keyStr]; ok {
			out.SetStringKey(e.keyStr, e.value)
		}
	}
	return out
}

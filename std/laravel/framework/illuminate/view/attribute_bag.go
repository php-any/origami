package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/illuminate/conditionable"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const attributeBagClassName = "Illuminate\\View\\ComponentAttributeBag"
const appendableClassName = "Illuminate\\View\\AppendableAttributeValue"

// AttributeBagClass 对齐 ComponentAttributeBag 热路径。
type AttributeBagClass struct {
	node.Node
	methods map[string]data.Method
}

func NewAttributeBagClass() data.ClassStmt {
	c := &AttributeBagClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *AttributeBagClass) GetName() string    { return attributeBagClassName }
func (c *AttributeBagClass) GetExtend() *string { return nil }
func (c *AttributeBagClass) GetImplements() []string {
	return []string{"ArrayAccess", "IteratorAggregate", "Illuminate\\Contracts\\Support\\Htmlable", "Illuminate\\Contracts\\Support\\Arrayable"}
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
	m, ok := c.methods[data.MethodLookupKey(name)]
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

func (c *AttributeBagClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"attributes"}, 0, bagConstruct)
	c.methods["get"] = kit.InstanceMethodOpt("get", []string{"key", "default"}, 1, bagGet)
	c.methods["has"] = kit.InstanceMethod("has", []string{"key"}, bagHas)
	c.methods["missing"] = kit.InstanceMethod("missing", []string{"key"}, bagMissing)
	c.methods["all"] = kit.InstanceMethodOpt("all", []string{"keys"}, 0, bagAll)
	c.methods["first"] = kit.InstanceMethodOpt("first", []string{"default"}, 0, bagFirst)
	c.methods["filter"] = kit.InstanceMethod("filter", []string{"callback"}, bagFilter)
	c.methods["getattributes"] = kit.InstanceMethod("getAttributes", nil, bagAll)
	c.methods["setattributes"] = kit.InstanceMethod("setAttributes", []string{"attributes"}, bagSetAttributes)
	c.methods["merge"] = kit.InstanceMethodOpt("merge", []string{"attributeDefaults", "escape"}, 0, bagMerge)
	c.methods["class"] = kit.InstanceMethod("class", []string{"classList"}, bagClass)
	c.methods["style"] = kit.InstanceMethod("style", []string{"styleList"}, bagStyle)
	c.methods["only"] = kit.InstanceMethod("only", []string{"keys"}, bagOnly)
	c.methods["except"] = kit.InstanceMethod("except", []string{"keys"}, bagExcept)
	c.methods["exceptprops"] = kit.InstanceMethod("exceptProps", []string{"props"}, bagExceptProps)
	c.methods["onlyprops"] = kit.InstanceMethod("onlyProps", []string{"props"}, bagOnlyProps)
	c.methods["wherestartswith"] = kit.InstanceMethod("whereStartsWith", []string{"needles"}, bagWhereStartsWith)
	c.methods["wheredoesntstartswith"] = kit.InstanceMethod("whereDoesntStartWith", []string{"needles"}, bagWhereDoesntStartWith)
	c.methods["thatstartswith"] = kit.InstanceMethod("thatStartWith", []string{"needles"}, bagWhereStartsWith)
	c.methods["prepends"] = kit.InstanceMethod("prepends", []string{"value"}, bagPrepends)
	c.methods["isempty"] = kit.InstanceMethod("isEmpty", nil, bagIsEmpty)
	c.methods["isnotempty"] = kit.InstanceMethod("isNotEmpty", nil, bagIsNotEmpty)
	c.methods["jsonserialize"] = kit.InstanceMethod("jsonSerialize", nil, bagAll)
	c.methods["__tostring"] = kit.InstanceMethod("__toString", nil, bagToHtml)
	c.methods["tohtml"] = kit.InstanceMethod("toHtml", nil, bagToHtml)
	c.methods["toarray"] = kit.InstanceMethod("toArray", nil, bagAll)
	c.methods["__invoke"] = kit.InstanceMethodOpt("__invoke", []string{"attributeDefaults"}, 0, bagInvoke)
	c.methods["getiterator"] = kit.InstanceMethod("getIterator", nil, bagGetIterator)
	c.methods["offsetexists"] = kit.InstanceMethod("offsetExists", []string{"offset"}, bagHas)
	c.methods["offsetget"] = kit.InstanceMethod("offsetGet", []string{"offset"}, bagOffsetGet)
	c.methods["offsetset"] = kit.InstanceMethod("offsetSet", []string{"offset", "value"}, bagOffsetSet)
	c.methods["offsetunset"] = kit.InstanceMethod("offsetUnset", []string{"offset"}, bagOffsetUnset)
	c.methods["shouldescapeattributevalue"] = kit.InstanceMethod("shouldEscapeAttributeValue", []string{"key", "value"}, bagShouldEscape)
	c.methods["extractpropnames"] = kit.StaticMethod("extractPropNames", []string{"keys"}, -1, bagExtractPropNames)
	kit.RegisterMacroable(c.methods, attributeBagClassName)
	kit.RegisterConditionable(c.methods, bagWhenProxy)
}

func bagRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("ComponentAttributeBag missing $this"))
}

func bagAttrs(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("attributes")
	if av, ok := kit.Unwrap(v).(*data.ArrayValue); ok && av != nil {
		return av
	}
	empty := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("attributes", empty)
	return empty
}

func bagConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	attrs := kit.Arg(ctx, 0)
	if attrs == nil || kit.IsNull(attrs) {
		attrs = data.NewArrayValue(nil)
	}
	if av, ok := kit.Unwrap(attrs).(*data.ArrayValue); ok {
		_ = cv.SetProperty("attributes", av)
	} else {
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, e := range kit.Entries(attrs) {
			zv := data.NewZVal(e.Value)
			zv.Name = e.KeyStr
			out.List = append(out.List, zv)
		}
		_ = cv.SetProperty("attributes", out)
	}
	return cv, nil
}

func newBag(ctx data.Context, attrs *data.ArrayValue) *data.ClassValue {
	vm := ctx.GetVM()
	cls, _ := vm.GetClass(attributeBagClassName)
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	_ = cv.SetProperty("attributes", attrs)
	return cv
}

func bagGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if e.KeyStr == key {
			return e.Value, nil
		}
	}
	def := kit.Arg(ctx, 1)
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func bagHas(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if e.KeyStr == key {
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
	b, _ := v.(data.AsBool).AsBool()
	return data.NewBoolValue(!b), nil
}

func bagAll(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys := kit.Arg(ctx, 0)
	if keys == nil || kit.IsNull(keys) {
		return bagAttrs(cv), nil
	}
	ret, ctl := bagOnly(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if bag, ok := kit.Unwrap(ret.(data.Value)).(*data.ClassValue); ok {
		return bagAttrs(bag), nil
	}
	return ret, nil
}

func bagFirst(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	attrs := bagAttrs(cv)
	for _, e := range kit.Entries(attrs) {
		return e.Value, nil
	}
	def := kit.Arg(ctx, 0)
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func bagFilter(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(bagAttrs(cv)) {
		ret, ctl := kit.Call(ctx, cb, e.Value, data.NewStringValue(e.KeyStr))
		if ctl != nil {
			return nil, ctl
		}
		keep := false
		if ret != nil {
			if v, ok := ret.(data.Value); ok {
				keep = kit.Truthy(v)
			}
		}
		if keep {
			zv := data.NewZVal(e.Value)
			zv.Name = e.KeyStr
			out.List = append(out.List, zv)
		}
	}
	return newBag(ctx, out), nil
}

func bagPrepends(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	cls, _ := vm.GetClass(appendableClassName)
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	_ = cv.SetProperty("value", kit.Arg(ctx, 0))
	return cv, nil
}

func bagIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	htmlRet, ctl := bagToHtml(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s := strings.TrimSpace(htmlRet.(data.Value).AsString())
	return data.NewBoolValue(s == ""), nil
}

func bagIsNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := bagIsEmpty(ctx)
	if ctl != nil {
		return nil, ctl
	}
	b, _ := v.(data.AsBool).AsBool()
	return data.NewBoolValue(!b), nil
}

func bagInvoke(ctx data.Context) (data.GetValue, data.Control) {
	merged, ctl := bagMerge(ctx)
	if ctl != nil {
		return nil, ctl
	}
	htmlStr := ""
	if bag, ok := kit.Unwrap(merged.(data.Value)).(*data.ClassValue); ok {
		htmlStr = renderBagAttributes(bagAttrs(bag))
	} else if merged != nil {
		htmlStr = merged.(data.Value).AsString()
	}
	vm := ctx.GetVM()
	cls, ok := vm.GetClass("Illuminate\\Support\\HtmlString")
	if !ok {
		return data.NewStringValue(htmlStr), nil
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	_ = cv.SetProperty("html", data.NewStringValue(htmlStr))
	return cv, nil
}

func bagExtractPropNames(ctx data.Context) (data.GetValue, data.Control) {
	keys := kit.Arg(ctx, 0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(keys) {
		name := e.KeyStr
		// 对齐 PHP：is_numeric($key) ? $default : $key
		if isNumericPropKey(e) {
			v := kit.Unwrap(e.Value)
			if v != nil {
				name = v.AsString()
			}
		}
		out.List = append(out.List, data.NewZVal(data.NewStringValue(name)))
		out.List = append(out.List, data.NewZVal(data.NewStringValue(toKebab(name))))
	}
	return out, nil
}

func isNumericPropKey(e kit.KV) bool {
	if e.Key == nil && e.KeyStr == "" {
		return true
	}
	if _, ok := e.Key.(*data.IntValue); ok {
		return true
	}
	if e.KeyStr == "" {
		return true
	}
	n, err := strconv.Atoi(e.KeyStr)
	if err != nil {
		return false
	}
	return strconv.Itoa(n) == e.KeyStr
}

func toKebab(name string) string {
	kebab := strings.ReplaceAll(name, "_", "-")
	var b strings.Builder
	for i, r := range kebab {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func bagWhenProxy(ctx data.Context, target data.Value) (data.Value, data.Control) {
	return conditionable.NewWhenProxy(ctx, target)
}

func bagSetAttributes(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	attrs := kit.Arg(ctx, 0)
	if attrs == nil || kit.IsNull(attrs) {
		attrs = data.NewArrayValue(nil)
	}
	if av, ok := kit.Unwrap(attrs).(*data.ArrayValue); ok {
		_ = cv.SetProperty("attributes", av)
	} else {
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, e := range kit.Entries(attrs) {
			zv := data.NewZVal(e.Value)
			zv.Name = e.KeyStr
			out.List = append(out.List, zv)
		}
		_ = cv.SetProperty("attributes", out)
	}
	return cv, nil
}

func bagMerge(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	defaults := kit.Arg(ctx, 0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	seen := map[string]bool{}
	for _, e := range kit.Entries(bagAttrs(cv)) {
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
		seen[e.KeyStr] = true
	}
	for _, e := range kit.Entries(defaults) {
		if seen[e.KeyStr] {
			continue
		}
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
	}
	return newBag(ctx, out), nil
}

func bagClass(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	list := classListString(kit.Arg(ctx, 0))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(bagAttrs(cv)) {
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
	}
	if list != "" {
		existing := ""
		for _, e := range kit.Entries(out) {
			if e.KeyStr == "class" {
				existing = e.Value.AsString()
			}
		}
		merged := strings.TrimSpace(existing + " " + list)
		found := false
		for i, z := range out.List {
			if z != nil && z.Name == "class" {
				out.List[i] = data.NewZVal(data.NewStringValue(merged))
				out.List[i].Name = "class"
				found = true
				break
			}
		}
		if !found {
			zv := data.NewZVal(data.NewStringValue(merged))
			zv.Name = "class"
			out.List = append(out.List, zv)
		}
	}
	return newBag(ctx, out), nil
}

// classListString 对齐 Arr::toCssClasses：
//
//	foreach ($classList as $class => $constraint) {
//	    if (is_numeric($class))      { $classes[] = $constraint; }
//	    elseif ($constraint)         { $classes[] = $class; }
//	}
//
// 数字键取值（列表写法 ['a','b']），字符串键为真时取键（条件写法 ['a' => true]）。
func classListString(v data.Value) string {
	v = kit.Unwrap(v)
	if v == nil {
		return ""
	}
	if s, ok := v.(*data.StringValue); ok {
		return s.AsString()
	}
	parts := []string{}
	for _, e := range kit.Entries(v) {
		if numericEntryKey(e) {
			if s := e.Value.AsString(); s != "" {
				parts = append(parts, s)
			}
			continue
		}
		if kit.Truthy(e.Value) {
			parts = append(parts, e.KeyStr)
		}
	}
	return strings.Join(parts, " ")
}

// numericEntryKey 对齐 PHP 的 is_numeric($class)：列表下标（ArrayValue 的整数键）与
// 数字属性名（混合键数组会以 ObjectValue 承载，键一律是 StringValue）都算数字键。
// 数字键取「值」当类名，字符串键在约束为真时取「键」当类名。
func numericEntryKey(e kit.KV) bool {
	if _, isInt := e.Key.(*data.IntValue); isInt {
		return true
	}
	_, ok := data.ParseIntArrayKeyName(e.KeyStr)
	return ok
}

// styleListString 对齐 Arr::toCssStyles：数字键取值、字符串键为真时取键，每段以 ';' 收尾。
func styleListString(v data.Value) string {
	v = kit.Unwrap(v)
	if v == nil {
		return ""
	}
	if s, ok := v.(*data.StringValue); ok {
		return s.AsString()
	}
	parts := []string{}
	for _, e := range kit.Entries(v) {
		var raw string
		if numericEntryKey(e) {
			raw = e.Value.AsString()
		} else if kit.Truthy(e.Value) {
			raw = e.KeyStr
		}
		if raw == "" {
			continue
		}
		if !strings.HasSuffix(raw, ";") {
			raw += ";"
		}
		parts = append(parts, raw)
	}
	return strings.Join(parts, " ")
}

func bagStyle(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	added := styleListString(kit.Arg(ctx, 0))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	existing := ""
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if e.KeyStr == "style" {
			existing = e.Value.AsString()
			continue
		}
		zv := data.NewZVal(e.Value)
		zv.Name = e.KeyStr
		out.List = append(out.List, zv)
	}
	// 对齐 ComponentAttributeBag::merge：各段以 ';' 收尾后空格拼接；两边都空则不写 style。
	parts := make([]string, 0, 2)
	for _, s := range []string{existing, added} {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if !strings.HasSuffix(s, ";") {
			s += ";"
		}
		parts = append(parts, s)
	}
	if len(parts) > 0 {
		zv := data.NewZVal(data.NewStringValue(strings.Join(parts, " ")))
		zv.Name = "style"
		out.List = append(out.List, zv)
	}
	return newBag(ctx, out), nil
}

func bagFilterKeys(ctx data.Context, keep func(string) bool) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range kit.Entries(bagAttrs(cv)) {
		if keep(e.KeyStr) {
			zv := data.NewZVal(e.Value)
			zv.Name = e.KeyStr
			out.List = append(out.List, zv)
		}
	}
	return newBag(ctx, out), nil
}

func keySet(v data.Value) map[string]bool {
	m := map[string]bool{}
	for _, e := range kit.Entries(v) {
		if _, isInt := e.Key.(*data.IntValue); isInt || e.KeyStr == "" {
			m[e.Value.AsString()] = true
		} else {
			m[e.KeyStr] = true
		}
	}
	if len(m) == 0 && v != nil {
		m[v.AsString()] = true
	}
	return m
}

func bagOnly(ctx data.Context) (data.GetValue, data.Control) {
	keys := keySet(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return keys[k] })
}
func bagExcept(ctx data.Context) (data.GetValue, data.Control) {
	keys := keySet(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return !keys[k] })
}

func extractPropNames(props data.Value) map[string]bool {
	m := map[string]bool{}
	for _, e := range kit.Entries(props) {
		name := e.KeyStr
		if _, isInt := e.Key.(*data.IntValue); isInt || name == "" {
			name = e.Value.AsString()
		}
		m[name] = true
		m[toKebab(name)] = true
	}
	return m
}

func bagExceptProps(ctx data.Context) (data.GetValue, data.Control) {
	props := extractPropNames(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return !props[k] })
}
func bagOnlyProps(ctx data.Context) (data.GetValue, data.Control) {
	props := extractPropNames(kit.Arg(ctx, 0))
	return bagFilterKeys(ctx, func(k string) bool { return props[k] })
}

func bagWhereStartsWith(ctx data.Context) (data.GetValue, data.Control) {
	prefix := kit.Arg(ctx, 0).AsString()
	return bagFilterKeys(ctx, func(k string) bool { return strings.HasPrefix(k, prefix) })
}
func bagWhereDoesntStartWith(ctx data.Context) (data.GetValue, data.Control) {
	prefix := kit.Arg(ctx, 0).AsString()
	return bagFilterKeys(ctx, func(k string) bool { return !strings.HasPrefix(k, prefix) })
}

func bagToHtml(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewStringValue(renderBagAttributes(bagAttrs(cv))), nil
}

// renderBagAttributes 对齐 Laravel 13 ComponentAttributeBag::__toString()。
// 绑定属性在 ComponentTagCompiler::sanitizeComponentAttribute() 中已经转义；
// 此处再次 html.EscapeString 会把 &#039; 变成 &amp;#039;，导致 Alpine
// 在 DOM 中读到字面实体并将表达式解析为非法 JavaScript。
func renderBagAttributes(attrs *data.ArrayValue) string {
	var b strings.Builder
	for _, e := range kit.Entries(attrs) {
		if e.Value == nil || kit.IsNull(e.Value) {
			continue
		}
		value := e.Value.AsString()
		if bv, ok := e.Value.(*data.BoolValue); ok {
			okv, _ := bv.AsBool()
			if !okv {
				continue
			}
			if e.KeyStr == "x-data" || strings.HasPrefix(e.KeyStr, "wire:") {
				value = ""
			} else {
				value = e.KeyStr
			}
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(e.KeyStr)
		b.WriteString(`="`)
		b.WriteString(strings.ReplaceAll(strings.TrimSpace(value), `"`, `\"`))
		b.WriteByte('"')
	}
	return b.String()
}

func bagGetIterator(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	vm := ctx.GetVM()
	cls, ok := vm.GetClass("ArrayIterator")
	if !ok {
		var c2 data.Control
		cls, c2 = vm.GetOrLoadClass("ArrayIterator")
		if c2 != nil {
			return bagAttrs(cv), nil
		}
	}
	it := data.NewClassValue(cls, ctx.CreateBaseContext())
	if ctor := cls.GetConstruct(); ctor != nil {
		nctx := it.CreateContext(ctor.GetVariables())
		data.BindDeclaredArgs(nctx, ctor, []data.Value{bagAttrs(cv)})
		if _, ctl := ctor.Call(nctx); ctl != nil {
			return nil, ctl
		}
	}
	return it, nil
}

func bagOffsetGet(ctx data.Context) (data.GetValue, data.Control) { return bagGet(ctx) }

func bagOffsetSet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	val := kit.Arg(ctx, 1)
	attrs := bagAttrs(cv)
	for i, z := range attrs.List {
		if z != nil && z.Name == key {
			attrs.List[i] = data.NewZVal(val)
			attrs.List[i].Name = key
			return data.NewNullValue(), nil
		}
	}
	zv := data.NewZVal(val)
	zv.Name = key
	attrs.List = append(attrs.List, zv)
	return data.NewNullValue(), nil
}

func bagOffsetUnset(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := bagRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := kit.Arg(ctx, 0).AsString()
	attrs := bagAttrs(cv)
	out := make([]*data.ZVal, 0, len(attrs.List))
	for _, z := range attrs.List {
		if z != nil && z.Name == key {
			continue
		}
		out = append(out, z)
	}
	attrs.List = out
	return data.NewNullValue(), nil
}

func bagShouldEscape(ctx data.Context) (data.GetValue, data.Control) {
	v := kit.Arg(ctx, 1)
	if cv, ok := kit.Unwrap(v).(*data.ClassValue); ok && cv.Class != nil {
		if cv.Class.GetName() == appendableClassName {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

// AppendableClass 对齐 AppendableAttributeValue。
type AppendableClass struct {
	node.Node
	methods map[string]data.Method
}

func NewAppendableClass() data.ClassStmt {
	c := &AppendableClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = kit.InstanceMethod("__construct", []string{"value"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := kit.Receiver(ctx)
		if cv == nil {
			return data.NewNullValue(), nil
		}
		_ = cv.SetProperty("value", kit.Arg(ctx, 0))
		return cv, nil
	})
	c.methods["__tostring"] = kit.InstanceMethod("__toString", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := kit.Receiver(ctx)
		if cv == nil {
			return data.NewStringValue(""), nil
		}
		v, _ := cv.GetProperty("value")
		if v == nil {
			return data.NewStringValue(""), nil
		}
		return data.NewStringValue(v.AsString()), nil
	})
	return c
}

func (c *AppendableClass) GetName() string    { return appendableClassName }
func (c *AppendableClass) GetExtend() *string { return nil }
func (c *AppendableClass) GetImplements() []string {
	return []string{"Stringable"}
}
func (c *AppendableClass) GetProperty(name string) (data.Property, bool) {
	if name == "value" {
		return node.NewProperty(nil, "value", "public", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *AppendableClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "value", "public", false, data.NewNullValue())}
}
func (c *AppendableClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *AppendableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *AppendableClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *AppendableClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *AppendableClass) GetStaticMethod(string) (data.Method, bool) { return nil, false }

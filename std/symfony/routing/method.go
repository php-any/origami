package routing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/preg"
)

func strPtr(s string) *string { return &s }

func emptyArray() data.Value { return data.NewArrayValue(nil) }

type rtClass struct {
	node.Node
	name       string
	extend     *string
	impl       []string
	props      []data.Property
	methods    map[string]data.Method
	methodList []data.Method
	statics    map[string]data.Value
}

func newRtClass(name string, extend *string, impl []string, props []data.Property) *rtClass {
	return &rtClass{
		name:    name,
		extend:  extend,
		impl:    impl,
		props:   props,
		methods: map[string]data.Method{},
		statics: map[string]data.Value{},
	}
}

func (c *rtClass) add(m data.Method) *rtClass {
	c.methods[data.MethodLookupKey(m.GetName())] = m
	c.methodList = append(c.methodList, m)
	return c
}

func (c *rtClass) constInt(name string, v int) *rtClass {
	c.statics[name] = data.NewIntValue(v)
	return c
}

func (c *rtClass) constStr(name, v string) *rtClass {
	c.statics[name] = data.NewStringValue(v)
	return c
}

func (c *rtClass) GetName() string                  { return c.name }
func (c *rtClass) GetExtend() *string               { return c.extend }
func (c *rtClass) GetImplements() []string          { return c.impl }
func (c *rtClass) GetPropertyList() []data.Property { return c.props }
func (c *rtClass) GetMethods() []data.Method        { return c.methodList }
func (c *rtClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *rtClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.props {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *rtClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *rtClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}
func (c *rtClass) GetStaticProperty(name string) (data.Value, bool) {
	v, ok := c.statics[name]
	return v, ok
}
func (c *rtClass) GetConstruct() data.Method {
	if m, ok := c.methods["__construct"]; ok {
		return m
	}
	return nil
}

type rtMethod struct {
	name   string
	static bool
	params []data.GetValue
	vars   []data.Variable
	ret    data.Types
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *rtMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *rtMethod) GetName() string                                     { return m.name }
func (m *rtMethod) GetModifier() data.Modifier                          { return data.ModifierPublic }
func (m *rtMethod) GetIsStatic() bool                                   { return m.static }
func (m *rtMethod) GetParams() []data.GetValue                          { return m.params }
func (m *rtMethod) GetVariables() []data.Variable                       { return m.vars }
func (m *rtMethod) GetReturnType() data.Types                           { return m.ret }

func param(name string, index int, def data.GetValue, ty data.Types) data.GetValue {
	return node.NewParameter(nil, name, index, def, ty)
}

func variable(name string, index int, ty data.Types) data.Variable {
	return node.NewVariable(nil, name, index, ty)
}

func meth(name string, ps []string, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	params := make([]data.GetValue, len(ps))
	vars := make([]data.Variable, len(ps))
	for i, p := range ps {
		params[i] = param(p, i, nil, nil)
		vars[i] = variable(p, i, nil)
	}
	return &rtMethod{name: name, params: params, vars: vars, fn: fn}
}

func methDef(name string, params []data.GetValue, vars []data.Variable, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	return &rtMethod{name: name, params: params, vars: vars, fn: fn}
}

func staticMeth(name string, ps []string, fn func(data.Context) (data.GetValue, data.Control)) data.Method {
	m := meth(name, ps, fn).(*rtMethod)
	m.static = true
	return m
}

func privProp(name string, def data.GetValue) data.Property {
	return node.NewProperty(nil, name, "private", false, def)
}

func protProp(name string, def data.GetValue) data.Property {
	return node.NewProperty(nil, name, "protected", false, def)
}

func pubProp(name string, def data.GetValue) data.Property {
	return node.NewProperty(nil, name, "public", false, def)
}

func rtSelf(ctx data.Context) *data.ClassValue {
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		return c.ClassValue
	case *data.ClassValue:
		return c
	}
	return nil
}

func asClassValue(v data.GetValue) *data.ClassValue {
	switch t := v.(type) {
	case *data.ClassValue:
		return t
	case *data.ThisValue:
		if t != nil {
			return t.ClassValue
		}
	}
	return nil
}

func isNull(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func arg(ctx data.Context, i int) data.Value {
	v, ok := ctx.GetIndexValue(i)
	if !ok || v == nil {
		return nil
	}
	return v
}

func argString(ctx data.Context, i int, def string) string {
	v := arg(ctx, i)
	if v == nil || isNull(v) {
		return def
	}
	return v.AsString()
}

func argInt(ctx data.Context, i int, def int) int {
	v := arg(ctx, i)
	if v == nil || isNull(v) {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	n, err := strconv.Atoi(strings.TrimSpace(v.AsString()))
	if err != nil {
		return def
	}
	return n
}

func argBool(ctx data.Context, i int, def bool) bool {
	v := arg(ctx, i)
	if v == nil || isNull(v) {
		return def
	}
	if bv, ok := v.(data.AsBool); ok {
		b, _ := bv.AsBool()
		return b
	}
	return def
}

func argArray(ctx data.Context, i int) *data.ArrayValue {
	return valueToArray(arg(ctx, i))
}

func valueToArray(v data.Value) *data.ArrayValue {
	if v == nil || isNull(v) {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	if zv, ok := v.(*data.ZValValue); ok && zv != nil && zv.ZVal != nil {
		return valueToArray(zv.ZVal.Value)
	}
	switch t := v.(type) {
	case *data.ArrayValue:
		return unwrapAssoc(t)
	case *data.ObjectValue:
		return unwrapAssoc(objectToArray(t))
	default:
		// 关联数组在 Origami 里是 ObjectValue；未知类型不要再包一层 [0 => $v]。
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
}

func objectToArray(obj *data.ObjectValue) *data.ArrayValue {
	out := &data.ArrayValue{}
	if obj == nil {
		return out
	}
	obj.RangeProperties(func(key string, value data.Value) bool {
		out.List = append(out.List, data.NewNamedZVal(key, value))
		return true
	})
	return out
}

func unwrapAssoc(arr *data.ArrayValue) *data.ArrayValue {
	if arr == nil {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	if len(arr.List) != 1 || arr.List[0] == nil {
		return arr
	}
	z := arr.List[0]
	if z.Name != "" && z.Name != "0" {
		return arr
	}
	switch inner := z.Value.(type) {
	case *data.ObjectValue:
		return unwrapAssoc(objectToArray(inner))
	case *data.ArrayValue:
		// 仅拆「又包了一层的关联数组」，不要把单元素列表（如一条 compiled row）拆掉。
		if arrayLooksAssoc(inner) {
			return unwrapAssoc(inner)
		}
	}
	return arr
}

func arrayLooksAssoc(arr *data.ArrayValue) bool {
	if arr == nil {
		return false
	}
	for i, z := range arr.List {
		if z == nil {
			continue
		}
		if z.Name != "" && z.Name != strconv.Itoa(i) {
			return true
		}
	}
	return false
}

func toStringSlice(v data.Value) []string {
	if v == nil || isNull(v) {
		return nil
	}
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		av := valueToArray(v)
		out := make([]string, 0, len(av.List))
		for _, z := range av.List {
			if z != nil && z.Value != nil && !isNull(z.Value) {
				out = append(out, z.Value.AsString())
			}
		}
		return out
	}
	s := v.AsString()
	if s == "" {
		return nil
	}
	return []string{s}
}

func prop(cv *data.ClassValue, name string) data.Value {
	if cv == nil {
		return data.NewNullValue()
	}
	v, _ := cv.GetProperty(name)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func propString(cv *data.ClassValue, name, def string) string {
	v := prop(cv, name)
	if isNull(v) {
		return def
	}
	return v.AsString()
}

func propInt(cv *data.ClassValue, name string, def int) int {
	v := prop(cv, name)
	if isNull(v) {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	return def
}

func propArray(cv *data.ClassValue, name string) *data.ArrayValue {
	return valueToArray(prop(cv, name))
}

func valueIsTruthy(v data.Value) bool {
	if v == nil || isNull(v) {
		return false
	}
	if bv, ok := v.(data.AsBool); ok {
		b, _ := bv.AsBool()
		return b
	}
	s := strings.TrimSpace(v.AsString())
	return s != "" && s != "0"
}

func setProp(cv *data.ClassValue, name string, v data.Value) {
	if cv == nil {
		return
	}
	if v == nil {
		v = data.NewNullValue()
	}
	_ = cv.SetProperty(name, v)
}

func phpList(vals ...data.Value) *data.ArrayValue {
	return data.NewArrayValue(vals).(*data.ArrayValue)
}

func phpAssoc(pairs ...any) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		key, _ := pairs[i].(string)
		val, _ := pairs[i+1].(data.Value)
		if val == nil {
			val = data.NewNullValue()
		}
		list = append(list, data.NewNamedZVal(key, val))
	}
	return &data.ArrayValue{List: list}
}

func cloneArray(src *data.ArrayValue) *data.ArrayValue {
	if src == nil {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	out := data.CloneArrayValue(src)
	if out == nil {
		return data.NewArrayValue(nil).(*data.ArrayValue)
	}
	return out
}

func assocGet(arr *data.ArrayValue, key string) (data.Value, bool) {
	if arr == nil {
		return nil, false
	}
	for _, z := range arr.List {
		if z != nil && z.Name == key {
			return z.Value, true
		}
	}
	return nil, false
}

func assocHas(arr *data.ArrayValue, key string) bool {
	_, ok := assocGet(arr, key)
	return ok
}

func assocSet(arr *data.ArrayValue, key string, val data.Value) *data.ArrayValue {
	if arr == nil {
		arr = &data.ArrayValue{}
	}
	if val == nil {
		val = data.NewNullValue()
	}
	for _, z := range arr.List {
		if z != nil && z.Name == key {
			z.Value = val
			return arr
		}
	}
	arr.List = append(arr.List, data.NewNamedZVal(key, val))
	return arr
}

func assocUnset(arr *data.ArrayValue, key string) *data.ArrayValue {
	if arr == nil {
		return &data.ArrayValue{}
	}
	out := make([]*data.ZVal, 0, len(arr.List))
	for _, z := range arr.List {
		if z != nil && z.Name == key {
			continue
		}
		out = append(out, z)
	}
	arr.List = out
	return arr
}

func arrayKey(z *data.ZVal, idx int) string {
	if z != nil && z.Name != "" {
		return z.Name
	}
	return strconv.Itoa(idx)
}

func stringsToArray(ss []string) *data.ArrayValue {
	vals := make([]data.Value, len(ss))
	for i, s := range ss {
		vals[i] = data.NewStringValue(s)
	}
	return phpList(vals...)
}

func flipStrings(ss []string) *data.ArrayValue {
	if len(ss) == 0 {
		return nil
	}
	list := make([]*data.ZVal, 0, len(ss))
	for i, s := range ss {
		list = append(list, data.NewNamedZVal(s, data.NewIntValue(i)))
	}
	return &data.ArrayValue{List: list}
}

func bindArgs(fnCtx data.Context, vars []data.Variable, args []data.Value) data.Control {
	n := len(vars)
	if n > len(args) {
		n = len(args)
	}
	for i := 0; i < n; i++ {
		a := args[i]
		if a == nil {
			a = data.NewNullValue()
		}
		if ctl := fnCtx.SetVariableValue(vars[i], a); ctl != nil {
			return ctl
		}
	}
	if len(args) > 0 {
		fnCtx.SetFlatCallArgs(args)
	}
	return nil
}

func callNamed(cv *data.ClassValue, name string, args ...data.Value) (data.Value, data.Control) {
	if cv == nil {
		return data.NewNullValue(), nil
	}
	m, ok := cv.GetMethod(name)
	if !ok || m == nil {
		return data.NewNullValue(), nil
	}
	fnCtx := cv.CreateContext(m.GetVariables())
	if ctl := bindArgs(fnCtx, m.GetVariables(), args); ctl != nil {
		return nil, ctl
	}
	ret, ctl := m.Call(fnCtx)
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil {
		return data.NewNullValue(), nil
	}
	if v, ok := ret.(data.Value); ok {
		return v, nil
	}
	return data.NewNullValue(), nil
}

func throwNamed(ctx data.Context, className, message string) data.Control {
	return throwNamedArgs(ctx, className, data.NewStringValue(message))
}

func throwNamedArgs(ctx data.Context, className string, args ...data.Value) data.Control {
	if ctx == nil || ctx.GetVM() == nil {
		msg := className
		if len(args) > 0 && args[0] != nil {
			msg = args[0].AsString()
		}
		return data.NewErrorThrowByName(nil, fmt.Errorf("%s", msg), className)
	}
	stmt, ok := ctx.GetVM().GetClass(className)
	if !ok {
		msg := className
		if len(args) > 0 && args[0] != nil {
			msg = args[0].AsString()
		}
		return data.NewErrorThrowByName(nil, fmt.Errorf("%s", msg), className)
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return acl
	}
	cv := asClassValue(obj)
	if cv == nil {
		return data.NewErrorThrowByName(nil, fmt.Errorf("%s", className), className)
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		if m, ok := cv.GetMethod("__construct"); ok {
			method = m
		}
	}
	if method != nil {
		fnCtx := cv.CreateContext(method.GetVariables())
		if ctl := bindArgs(fnCtx, method.GetVariables(), args); ctl != nil {
			return ctl
		}
		if _, acl := method.Call(fnCtx); acl != nil {
			return acl
		}
	}
	return data.NewErrorThrowFromClassValue(nil, cv)
}

func phpPregQuote(s string) string {
	const specials = `\.+*?[^]$(){}=!<>|:-#`
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(s); i++ {
		ch := rune(s[i])
		if strings.ContainsRune(specials, ch) {
			b.WriteByte('\\')
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func phpRawURLDecode(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) && isHex(s[i+1]) && isHex(s[i+2]) {
			v, _ := strconv.ParseUint(s[i+1:i+3], 16, 8)
			b.WriteByte(byte(v))
			i += 2
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func phpRawURLEncode(s string) string {
	var b strings.Builder
	b.Grow(len(s) * 3)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func pregMatch(pattern, subject string) bool {
	m, err := preg.CompileAny(pattern)
	if err != nil || m == nil {
		return false
	}
	return m.MatchString(subject)
}

func pregMatchGroups(pattern, subject string) (map[string]string, []string, bool) {
	m, err := preg.CompileAny(pattern)
	if err != nil || m == nil {
		return nil, nil, false
	}
	loc := m.FindStringSubmatchIndex(subject)
	if loc == nil {
		return nil, nil, false
	}
	names := extractGroupNames(pattern)
	named := map[string]string{}
	indexed := make([]string, 0, len(loc)/2)
	for i := 0; i < len(loc)/2; i++ {
		start, end := loc[2*i], loc[2*i+1]
		var s string
		if start >= 0 && end >= start && end <= len(subject) {
			s = subject[start:end]
		}
		indexed = append(indexed, s)
		if i > 0 && i <= len(names) && names[i-1] != "" {
			named[names[i-1]] = s
		}
	}
	return named, indexed, true
}

func extractGroupNames(pattern string) []string {
	body := pattern
	if len(pattern) >= 2 {
		open := pattern[0]
		close := open
		switch open {
		case '{':
			close = '}'
		case '(':
			close = ')'
		case '[':
			close = ']'
		case '<':
			close = '>'
		}
		if idx := strings.LastIndexByte(pattern, close); idx > 0 {
			body = pattern[1:idx]
		}
	}
	var names []string
	for i := 0; i < len(body); i++ {
		if body[i] == '(' {
			rest := body[i:]
			if strings.HasPrefix(rest, "(?P<") {
				end := strings.IndexByte(rest, '>')
				if end > 4 {
					names = append(names, rest[4:end])
				}
			} else if strings.HasPrefix(rest, "(?<") && !strings.HasPrefix(rest, "(?<=") && !strings.HasPrefix(rest, "(?<!") {
				end := strings.IndexByte(rest, '>')
				if end > 3 {
					names = append(names, rest[3:end])
				}
			} else if strings.HasPrefix(rest, "(?:") || strings.HasPrefix(rest, "(?=") ||
				strings.HasPrefix(rest, "(?!") || strings.HasPrefix(rest, "(?<=") ||
				strings.HasPrefix(rest, "(?<!") || strings.HasPrefix(rest, "(?>") {
				continue
			} else {
				names = append(names, "")
			}
		}
	}
	return names
}

func newArrayIterator(ctx data.Context, storage data.Value) (data.GetValue, data.Control) {
	if ctx == nil || ctx.GetVM() == nil {
		return storage, nil
	}
	stmt, ok := ctx.GetVM().GetClass("ArrayIterator")
	if !ok {
		return storage, nil
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return storage, nil
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		return cv, nil
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	if storage == nil {
		storage = emptyArray()
	}
	if ctl := bindArgs(fnCtx, method.GetVariables(), []data.Value{storage}); ctl != nil {
		return nil, ctl
	}
	if _, acl = method.Call(fnCtx); acl != nil {
		return nil, acl
	}
	return cv, nil
}

func instantiate(ctx data.Context, className string, args ...data.Value) (*data.ClassValue, data.Control) {
	if ctx == nil || ctx.GetVM() == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("class %s: no vm", className))
	}
	stmt, ok := ctx.GetVM().GetClass(className)
	if !ok {
		loaded, acl := ctx.GetVM().GetOrLoadClass(className)
		if acl != nil {
			return nil, acl
		}
		stmt = loaded
	}
	if stmt == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("class %s not found", className))
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv := asClassValue(obj)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("class %s invalid instance", className))
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		if m, ok := cv.GetMethod("__construct"); ok {
			method = m
		}
	}
	if method != nil && (len(args) > 0 || len(method.GetVariables()) > 0) {
		fnCtx := cv.CreateContext(method.GetVariables())
		if ctl := bindArgs(fnCtx, method.GetVariables(), args); ctl != nil {
			return nil, ctl
		}
		if _, acl := method.Call(fnCtx); acl != nil {
			return nil, acl
		}
	}
	return cv, nil
}

func inStringSlice(ss []string, v string) bool {
	for _, s := range ss {
		if s == v {
			return true
		}
	}
	return false
}

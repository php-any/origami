package sfstring

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	unicodeStringName   = "Symfony\\Component\\String\\UnicodeString"
	byteStringName      = "Symfony\\Component\\String\\ByteString"
	codePointStringName = "Symfony\\Component\\String\\CodePointString"
	abstractUnicodeName = "Symfony\\Component\\String\\AbstractUnicodeString"
	abstractStringName  = "Symfony\\Component\\String\\AbstractString"
	invalidArgName      = "Symfony\\Component\\String\\Exception\\InvalidArgumentException"
)

const (
	unicodeTrimDefault = " \t\n\r\x00\x0B\x0C\u00A0\uFEFF"
	byteTrimDefault    = " \t\n\r\x00\x0B\x0C"
	byteAlphabet       = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type strKind int

const (
	kindGrapheme strKind = iota
	kindCodePoint
	kindByte
)

var (
	unicodeStringClass   *strClass
	byteStringClass      *strClass
	codePointStringClass *strClass
)

type strClass struct {
	node.Node
	name       string
	parent     string
	kind       strKind
	methods    map[string]data.Method
	implements []string
	abstract   bool
	props      []data.Property
	propIndex  map[string]data.Property
}

func initStringProps(c *strClass) {
	strP := node.NewProperty(nil, "string", "protected", false, data.NewStringValue(""), data.NewBaseType("string"))
	icP := node.NewProperty(nil, "ignoreCase", "protected", false, data.NewBoolValue(false), data.NewNullableType(data.NewBaseType("bool")))
	c.props = []data.Property{strP, icP}
	c.propIndex = map[string]data.Property{"string": strP, "ignoreCase": icP}
}

func (c *strClass) GetName() string { return c.name }
func (c *strClass) GetExtend() *string {
	if c.parent == "" {
		return nil
	}
	p := c.parent
	return &p
}
func (c *strClass) GetImplements() []string { return c.implements }
func (c *strClass) GetProperty(name string) (data.Property, bool) {
	p, ok := c.propIndex[name]
	return p, ok
}
func (c *strClass) GetPropertyList() []data.Property { return c.props }
func (c *strClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *strClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *strClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *strClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *strClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := c.GetMethod(name)
	if ok && m != nil && m.GetIsStatic() {
		return m, true
	}
	return nil, false
}
func (c *strClass) IsBuiltinAbstractClass() bool { return c.abstract }
func unicodeNormConsts(c *strClass) bool {
	switch c.name {
	case abstractUnicodeName, unicodeStringName, codePointStringName:
		return true
	}
	return false
}

func (c *strClass) GetStaticProperty(name string) (data.Value, bool) {
	// NFC/NFD 只在 AbstractUnicodeString 及其子类上；ByteString / AbstractString 不暴露。
	if unicodeNormConsts(c) {
		switch name {
		case "NFC", "FORM_C":
			return data.NewIntValue(4), true
		case "NFD", "FORM_D":
			return data.NewIntValue(1), true
		case "NFKC", "FORM_KC":
			return data.NewIntValue(5), true
		case "NFKD", "FORM_KD":
			return data.NewIntValue(2), true
		}
	}
	switch name {
	case "PREG_PATTERN_ORDER":
		return data.NewIntValue(1), true
	case "PREG_SET_ORDER":
		return data.NewIntValue(2), true
	case "PREG_OFFSET_CAPTURE":
		return data.NewIntValue(256), true
	case "PREG_UNMATCHED_AS_NULL":
		return data.NewIntValue(512), true
	case "PREG_SPLIT":
		return data.NewIntValue(0), true
	case "PREG_SPLIT_NO_EMPTY":
		return data.NewIntValue(1), true
	case "PREG_SPLIT_DELIM_CAPTURE":
		return data.NewIntValue(2), true
	case "PREG_SPLIT_OFFSET_CAPTURE":
		return data.NewIntValue(4), true
	}
	return nil, false
}

func (c *strClass) add(name string, static bool, params []data.GetValue, vars []data.Variable, fn func(data.Context) (data.GetValue, data.Control)) {
	c.methods[data.MethodLookupKey(name)] = &strMethod{name: name, static: static, params: params, vars: vars, fn: fn}
}

type strMethod struct {
	name   string
	static bool
	params []data.GetValue
	vars   []data.Variable
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *strMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *strMethod) GetName() string                                     { return m.name }
func (m *strMethod) GetModifier() data.Modifier                          { return data.ModifierPublic }
func (m *strMethod) GetIsStatic() bool                                   { return m.static }
func (m *strMethod) GetReturnType() data.Types                           { return nil }
func (m *strMethod) GetParams() []data.GetValue                          { return m.params }
func (m *strMethod) GetVariables() []data.Variable                       { return m.vars }

func p(name string, i int, def data.GetValue) data.GetValue {
	return node.NewParameter(nil, name, i, def, nil)
}
func v(name string, i int) data.Variable { return node.NewVariable(nil, name, i, nil) }
func variadic(name string, i int) data.GetValue {
	return node.NewParameters(nil, name, i, nil, nil)
}

func throwInvalid(msg string) data.Control {
	return data.NewErrorThrowByName(nil, fmt.Errorf("%s", msg), invalidArgName)
}

func strSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func strKindOf(cv *data.ClassValue) strKind {
	if cv == nil || cv.Class == nil {
		return kindGrapheme
	}
	if sc, ok := cv.Class.(*strClass); ok {
		return sc.kind
	}
	switch cv.Class.GetName() {
	case byteStringName:
		return kindByte
	case codePointStringName:
		return kindCodePoint
	default:
		return kindGrapheme
	}
}

func getString(cv *data.ClassValue) string {
	if cv == nil {
		return ""
	}
	val, _ := cv.GetProperty("string")
	if val == nil {
		return ""
	}
	return val.AsString()
}

func ignoreCaseOf(cv *data.ClassValue) bool {
	if cv == nil {
		return false
	}
	val, _ := cv.GetProperty("ignoreCase")
	if val == nil {
		return false
	}
	if b, ok := val.(data.AsBool); ok {
		okv, err := b.AsBool()
		return err == nil && okv
	}
	return false
}

func classOf(cv *data.ClassValue) data.ClassStmt {
	if cv != nil && cv.Class != nil {
		return cv.Class
	}
	return ensureUnicodeClass()
}

func cloneBase(cv *data.ClassValue) data.Context {
	if cv == nil {
		return nil
	}
	if cv.Context != nil {
		return cv.Context.CreateBaseContext()
	}
	return cv.CreateBaseContext()
}

// newOf 按 Symfony String 不变性 clone：新 ClassValue + 新 Context，不改原对象。
// PHP __clone 会把 ignoreCase 重置为 false。
func newOf(cv *data.ClassValue, s string) *data.ClassValue {
	out := data.NewClassValue(classOf(cv), cloneBase(cv))
	_ = out.SetProperty("string", data.NewStringValue(s))
	_ = out.SetProperty("ignoreCase", data.NewBoolValue(false))
	return out
}

func ensureUnicodeClass() *strClass {
	if unicodeStringClass == nil {
		NewUnicodeStringClass()
	}
	return unicodeStringClass
}

func ensureByteClass() *strClass {
	if byteStringClass == nil {
		NewByteStringClass()
	}
	return byteStringClass
}

func ensureCodePointClass() *strClass {
	if codePointStringClass == nil {
		NewCodePointStringClass()
	}
	return codePointStringClass
}

func newTyped(ctx data.Context, cls data.ClassStmt, s string) *data.ClassValue {
	base := ctx
	if c, ok := ctx.(*data.ClassMethodContext); ok && c.ClassValue != nil && c.ClassValue.Context != nil {
		base = c.ClassValue.Context
	}
	out := data.NewClassValue(cls, base.CreateBaseContext())
	_ = out.SetProperty("string", data.NewStringValue(s))
	_ = out.SetProperty("ignoreCase", data.NewBoolValue(false))
	return out
}

func lateClass(ctx data.Context) data.ClassStmt {
	if cm, ok := ctx.(*data.ClassMethodContext); ok {
		if cm.StaticClass != nil {
			return cm.StaticClass
		}
		if cm.ClassValue != nil && cm.ClassValue.Class != nil {
			return cm.ClassValue.Class
		}
	}
	return ensureUnicodeClass()
}

func argValue(ctx data.Context, i int) data.Value {
	v, ok := ctx.GetIndexValue(i)
	if !ok || v == nil {
		return nil
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return nil
	}
	return v
}

func argString(ctx data.Context, i int, def string) string {
	v := argValue(ctx, i)
	if v == nil {
		return def
	}
	return valueString(ctx, v)
}

func argInt(ctx data.Context, i int, def int) int {
	v := argValue(ctx, i)
	if v == nil {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err == nil {
			return n
		}
	}
	return def
}

func argBool(ctx data.Context, i int, def bool) bool {
	v := argValue(ctx, i)
	if v == nil {
		return def
	}
	if b, ok := v.(data.AsBool); ok {
		okv, err := b.AsBool()
		if err == nil {
			return okv
		}
	}
	return def
}

func argOptionalInt(ctx data.Context, i int) *int {
	v := argValue(ctx, i)
	if v == nil {
		return nil
	}
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err == nil {
			return &n
		}
	}
	return nil
}

func valueString(ctx data.Context, v data.Value) string {
	if v == nil {
		return ""
	}
	if cv, ok := v.(*data.ClassValue); ok && cv != nil {
		if isStringObj(cv) {
			return getString(cv)
		}
		if m, ok := cv.GetMethod("__toString"); ok && m != nil {
			fnCtx := cv.CreateContext(m.GetVariables())
			fnCtx.SetCallArgs([]data.GetValue{})
			ret, ctl := m.Call(fnCtx)
			if ctl == nil && ret != nil {
				if sv, ok := ret.(data.Value); ok {
					if _, isObj := sv.(*data.ClassValue); !isObj {
						return sv.AsString()
					}
				}
			}
		}
	}
	return v.AsString()
}

func isStringObj(cv *data.ClassValue) bool {
	if cv == nil || cv.Class == nil {
		return false
	}
	switch cv.Class.GetName() {
	case unicodeStringName, byteStringName, codePointStringName, abstractStringName, abstractUnicodeName:
		return true
	}
	ext := cv.Class.GetExtend()
	for ext != nil && *ext != "" {
		switch *ext {
		case abstractStringName, abstractUnicodeName, unicodeStringName, byteStringName, codePointStringName:
			return true
		}
		stmt, acl := cv.GetVM().GetOrLoadClass(*ext)
		if acl != nil || stmt == nil {
			break
		}
		ext = stmt.GetExtend()
	}
	return false
}

func needles(ctx data.Context, v data.Value) []string {
	if v == nil {
		return nil
	}
	if av, ok := v.(*data.ArrayValue); ok {
		out := make([]string, 0, len(av.List))
		for _, z := range av.List {
			if z != nil && z.Value != nil {
				out = append(out, valueString(ctx, z.Value))
			}
		}
		return out
	}
	if cv, ok := v.(*data.ClassValue); ok && isStringObj(cv) {
		return []string{getString(cv)}
	}
	return []string{valueString(ctx, v)}
}

func variadicStrings(ctx data.Context, i int) []string {
	v, ok := ctx.GetIndexValue(i)
	if ok && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			out := make([]string, 0, len(av.List))
			for _, z := range av.List {
				if z != nil && z.Value != nil {
					if _, isNull := z.Value.(*data.NullValue); isNull {
						continue
					}
					out = append(out, z.Value.AsString())
				}
			}
			return out
		}
		if _, isNull := v.(*data.NullValue); !isNull {
			return []string{v.AsString()}
		}
	}
	if g, ok := ctx.(interface{ GetFlatCallArgs() []data.Value }); ok {
		flat := g.GetFlatCallArgs()
		if len(flat) > i {
			out := make([]string, 0, len(flat)-i)
			for _, a := range flat[i:] {
				if a != nil {
					out = append(out, a.AsString())
				}
			}
			return out
		}
	}
	return nil
}

func toArray(items []*data.ClassValue) *data.ArrayValue {
	vals := make([]data.Value, len(items))
	for i, it := range items {
		vals[i] = it
	}
	return data.NewArrayValue(vals).(*data.ArrayValue)
}

func mustUTF8(s string) data.Control {
	if s != "" && !utf8.ValidString(s) {
		return throwInvalid("Invalid UTF-8 string.")
	}
	return nil
}

func units(s string, k strKind) []string {
	if s == "" {
		return nil
	}
	if k == kindByte {
		out := make([]string, len(s))
		for i := 0; i < len(s); i++ {
			out[i] = s[i : i+1]
		}
		return out
	}
	rs := []rune(s)
	out := make([]string, len(rs))
	for i, r := range rs {
		out[i] = string(r)
	}
	return out
}

func unitCount(s string, k strKind) int {
	if k == kindByte {
		return len(s)
	}
	return utf8.RuneCountInString(s)
}

func sliceUnits(s string, start int, length *int, k strKind) string {
	u := units(s, k)
	n := len(u)
	if start < 0 {
		start = n + start
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		return ""
	}
	end := n
	if length != nil {
		if *length < 0 {
			end = n + *length
			if end < start {
				return ""
			}
		} else {
			end = start + *length
			if end > n {
				end = n
			}
		}
	}
	return strings.Join(u[start:end], "")
}

func prefixUnits(s string, n int, k strKind) string {
	if n <= 0 {
		return ""
	}
	return sliceUnits(s, 0, &n, k)
}

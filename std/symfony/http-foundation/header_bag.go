package httpfoundation

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var headerTokenRe = regexp.MustCompile(`(?i)^[a-z0-9!#$%&'*.^_` + "`" + `|~-]+$`)

func normalizeHeaderKey(key string) string {
	b := []byte(key)
	for i, c := range b {
		if c == '_' {
			b[i] = '-'
		} else if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}

func headerSet(h *HeaderBagData, key string, values []*string, replace bool) {
	if h == nil {
		return
	}
	nk := normalizeHeaderKey(key)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.headers == nil {
		h.headers = make(map[string][]*string)
	}
	if replace || h.headers[nk] == nil {
		if _, ok := h.headers[nk]; !ok {
			h.keys = append(h.keys, nk)
		}
		h.headers[nk] = values
	} else {
		h.headers[nk] = append(h.headers[nk], values...)
	}
	if nk == "cache-control" {
		parts := make([]string, 0, len(h.headers[nk]))
		for _, v := range h.headers[nk] {
			if v != nil {
				parts = append(parts, *v)
			}
		}
		h.cacheControl = parseCacheControl(strings.Join(parts, ", "))
	}
}

func headerAll(h *HeaderBagData, key *string) map[string][]*string {
	if h == nil {
		return map[string][]*string{}
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if key != nil {
		nk := normalizeHeaderKey(*key)
		vs := h.headers[nk]
		if vs == nil {
			return map[string][]*string{nk: {}}
		}
		cp := make([]*string, len(vs))
		copy(cp, vs)
		return map[string][]*string{nk: cp}
	}
	out := make(map[string][]*string, len(h.headers))
	for k, vs := range h.headers {
		cp := make([]*string, len(vs))
		copy(cp, vs)
		out[k] = cp
	}
	return out
}

func headerValuesToArray(vs []*string) *data.ArrayValue {
	vals := make([]data.Value, len(vs))
	for i, v := range vs {
		if v == nil {
			vals[i] = data.NewNullValue()
		} else {
			vals[i] = data.NewStringValue(*v)
		}
	}
	return data.NewArrayValue(vals).(*data.ArrayValue)
}

func headersMapToArrayValue(headers map[string][]*string, order []string) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(headers))
	seen := map[string]bool{}
	for _, k := range order {
		if vs, ok := headers[k]; ok {
			list = append(list, &data.ZVal{Name: k, Value: headerValuesToArray(vs)})
			seen[k] = true
		}
	}
	for k, vs := range headers {
		if !seen[k] {
			list = append(list, &data.ZVal{Name: k, Value: headerValuesToArray(vs)})
		}
	}
	return &data.ArrayValue{List: list}
}

func parseCacheControl(header string) map[string]any {
	parts := headerUtilsSplit(header, ",=")
	return headerUtilsCombine(parts)
}

func getCacheControlHeader(cc map[string]any) string {
	keys := make([]string, 0, len(cc))
	for k := range cc {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	assoc := make(map[string]any, len(keys))
	for _, k := range keys {
		assoc[k] = cc[k]
	}
	return headerUtilsToString(assoc, ",")
}

func headerUtilsCombine(parts [][]string) map[string]any {
	assoc := make(map[string]any)
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		name := strings.ToLower(part[0])
		if len(part) > 1 {
			assoc[name] = part[1]
		} else {
			assoc[name] = true
		}
	}
	return assoc
}

func headerUtilsToString(assoc map[string]any, separator string) string {
	parts := make([]string, 0, len(assoc))
	keys := make([]string, 0, len(assoc))
	for k := range assoc {
		keys = append(keys, k)
	}
	// 保持调用方已排序时的顺序不稳定也没关系；getCacheControlHeader 会先排序
	sort.Strings(keys)
	for _, name := range keys {
		value := assoc[name]
		if value == true {
			parts = append(parts, name)
		} else {
			parts = append(parts, name+"="+headerUtilsQuote(fmt.Sprint(value)))
		}
	}
	return strings.Join(parts, separator+" ")
}

func headerUtilsQuote(s string) string {
	if headerTokenRe.MatchString(s) {
		return s
	}
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}

// 简化版 HeaderUtils::split，覆盖 Cache-Control / Cookie 常见形态。
func headerUtilsSplit(header, separators string) [][]string {
	header = strings.TrimSpace(header)
	if header == "" {
		return nil
	}
	sep0 := rune(separators[0])
	var groups [][]string
	var current []string
	var buf strings.Builder
	inQuote := false
	escape := false

	flushToken := func() {
		tok := strings.TrimSpace(buf.String())
		buf.Reset()
		if tok != "" || len(current) > 0 {
			current = append(current, headerUtilsUnquote(tok))
		}
	}

	for _, r := range header {
		if escape {
			buf.WriteRune(r)
			escape = false
			continue
		}
		if r == '\\' && inQuote {
			escape = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			buf.WriteRune(r)
			continue
		}
		if !inQuote && r == sep0 {
			flushToken()
			if len(current) > 0 {
				groups = append(groups, current)
				current = nil
			}
			continue
		}
		if !inQuote && strings.ContainsRune(separators[1:], r) {
			flushToken()
			continue
		}
		buf.WriteRune(r)
	}
	flushToken()
	if len(current) > 0 {
		groups = append(groups, current)
	}
	return groups
}

func headerUtilsUnquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		inner := s[1 : len(s)-1]
		inner = strings.ReplaceAll(inner, `\"`, `"`)
		inner = strings.ReplaceAll(inner, `\\`, `\`)
		return inner
	}
	return s
}

func makeDisposition(disposition, filename, filenameFallback string) (string, error) {
	if disposition != "attachment" && disposition != "inline" {
		return "", fmt.Errorf(`The disposition must be either "attachment" or "inline".`)
	}
	if filenameFallback == "" {
		filenameFallback = filename
	}
	for _, r := range filenameFallback {
		if r < 0x20 || r > 0x7e {
			return "", fmt.Errorf("The filename fallback must only contain ASCII characters.")
		}
	}
	if strings.Contains(filenameFallback, "%") {
		return "", fmt.Errorf(`The filename fallback cannot contain the "%%" character.`)
	}
	if strings.ContainsAny(filename, `/\`) || strings.ContainsAny(filenameFallback, `/\`) {
		return "", fmt.Errorf(`The filename and the fallback cannot contain the "/" and "\\" characters.`)
	}
	params := map[string]any{"filename": filenameFallback}
	if filename != filenameFallback {
		params["filename*"] = "utf-8''" + url.PathEscape(filename)
	}
	return disposition + "; " + headerUtilsToString(params, ";"), nil
}

func ucwordsDash(name string) string {
	parts := strings.Split(name, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "-")
}

// HeaderBagClass 实现 Symfony\Component\HttpFoundation\HeaderBag。
type HeaderBagClass struct {
	node.Node
	source     *HeaderBagData
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewHeaderBagClass() data.ClassStmt {
	return NewHeaderBagClassFrom(nil)
}

func NewHeaderBagClassFrom(source *HeaderBagData) data.ClassStmt {
	c := &HeaderBagClass{
		source: source,
		properties: []data.Property{
			protectedArrayProp("headers"),
			protectedArrayProp("cacheControl"),
		},
	}
	c.methods, c.methodList = headerBagMethods()
	return c
}

func (c *HeaderBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := c.source
	if src == nil {
		src = newHeaderBagData()
	} else {
		src = src.clone()
	}
	return data.NewProxyValue(NewHeaderBagClassFrom(src), ctx.CreateBaseContext()), nil
}

func (c *HeaderBagClass) GetName() string {
	return fqnHeaderBag
}
func (c *HeaderBagClass) GetExtend() *string { return nil }
func (c *HeaderBagClass) GetImplements() []string {
	return []string{"IteratorAggregate", "Countable", "Stringable"}
}
func (c *HeaderBagClass) GetSource() any                   { return c.source }
func (c *HeaderBagClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *HeaderBagClass) GetPropertyList() []data.Property { return c.properties }
func (c *HeaderBagClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *HeaderBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *HeaderBagClass) GetMethods() []data.Method { return c.methodList }

func headerBagMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{param("headers", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("headers", 0, nil)},
			nil, headerBagConstruct),
		pubMethod("__toString", nil, nil, data.NewBaseType("string"), headerBagToString),
		pubMethod("all",
			[]data.GetValue{param("key", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("key", 0, nil)},
			data.NewBaseType("array"), headerBagAll),
		pubMethod("keys", nil, nil, data.NewBaseType("array"), headerBagKeys),
		pubMethod("replace",
			[]data.GetValue{param("headers", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("headers", 0, nil)},
			nil, headerBagReplace),
		pubMethod("add",
			[]data.GetValue{param("headers", 0, nil, nil)},
			[]data.Variable{variable("headers", 0, nil)},
			nil, headerBagAdd),
		pubMethod("get",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			nil, headerBagGet),
		pubMethod("set",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("values", 1, nil, nil),
				param("replace", 2, data.NewBoolValue(true), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("values", 1, nil),
				variable("replace", 2, nil),
			},
			nil, headerBagSet),
		pubMethod("has",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			data.NewBaseType("bool"), headerBagHas),
		pubMethod("contains",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("value", 1, nil, nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("value", 1, nil),
			},
			data.NewBaseType("bool"), headerBagContains),
		pubMethod("remove",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			nil, headerBagRemove),
		pubMethod("getDate",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("default", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("default", 1, nil),
			},
			nil, headerBagGetDate),
		pubMethod("addCacheControlDirective",
			[]data.GetValue{
				param("key", 0, nil, nil),
				param("value", 1, data.NewBoolValue(true), nil),
			},
			[]data.Variable{
				variable("key", 0, nil),
				variable("value", 1, nil),
			},
			nil, headerBagAddCacheControlDirective),
		pubMethod("hasCacheControlDirective",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			data.NewBaseType("bool"), headerBagHasCacheControlDirective),
		pubMethod("getCacheControlDirective",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			nil, headerBagGetCacheControlDirective),
		pubMethod("removeCacheControlDirective",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			nil, headerBagRemoveCacheControlDirective),
		pubMethod("getIterator", nil, nil, nil, headerBagGetIterator),
		pubMethod("count", nil, nil, data.NewBaseType("int"), headerBagCount),
	}
	m := make(map[string]data.Method, len(list))
	for _, method := range list {
		m[method.GetName()] = method
	}
	return m, list
}

func headerBagConstruct(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	return nil, headerBagAddMap(h, raw)
}

func headerBagAddMap(h *HeaderBagData, raw data.Value) data.Control {
	keys, m, err := valueToOrderedAssoc(raw)
	if err != nil {
		return data.NewErrorThrow(nil, err)
	}
	for _, k := range keys {
		vals, ctl := coerceHeaderValues(m[k])
		if ctl != nil {
			return ctl
		}
		headerSet(h, k, vals, true)
	}
	return nil
}

func coerceHeaderValues(v data.Value) ([]*string, data.Control) {
	if v == nil || isNull(v) {
		return []*string{nil}, nil
	}
	if isArrayValue(v) {
		m, err := valueToAssocMap(v)
		if err != nil {
			return nil, data.NewErrorThrow(nil, err)
		}
		// 保持数值顺序：优先 ArrayValue list
		if arr, ok := v.(*data.ArrayValue); ok {
			out := make([]*string, 0, len(arr.List))
			for _, z := range arr.List {
				if z == nil || z.Value == nil || isNull(z.Value) {
					out = append(out, nil)
				} else {
					s := z.Value.AsString()
					out = append(out, &s)
				}
			}
			return out, nil
		}
		out := make([]*string, 0, len(m))
		for _, val := range m {
			if val == nil || isNull(val) {
				out = append(out, nil)
			} else {
				s := val.AsString()
				out = append(out, &s)
			}
		}
		return out, nil
	}
	s := v.AsString()
	return []*string{&s}, nil
}

func headerBagToString(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return data.NewStringValue(""), nil
	}
	all := headerAll(h, nil)
	if len(all) == 0 {
		return data.NewStringValue(""), nil
	}
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	maxLen := 0
	for _, k := range keys {
		n := len(ucwordsDash(k)) + 1
		if n > maxLen {
			maxLen = n
		}
	}
	var b strings.Builder
	for _, k := range keys {
		name := ucwordsDash(k) + ":"
		for _, v := range all[k] {
			val := ""
			if v != nil {
				val = *v
			}
			b.WriteString(fmt.Sprintf("%-*s %s\r\n", maxLen, name, val))
		}
	}
	return data.NewStringValue(b.String()), nil
}

func headerBagAll(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	key, present, isStr := optionalStringParam(ctx, 0)
	if present && isStr {
		all := headerAll(h, &key)
		return headerValuesToArray(all[normalizeHeaderKey(key)]), nil
	}
	if h == nil {
		return data.NewArrayValue(nil), nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return headersMapToArrayValue(h.headers, h.keys), nil
}

func headerBagKeys(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return data.NewArrayValue(nil), nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	vals := make([]data.Value, len(h.keys))
	for i, k := range h.keys {
		vals[i] = data.NewStringValue(k)
	}
	return data.NewArrayValue(vals), nil
}

func headerBagReplace(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return nil, nil
	}
	h.mu.Lock()
	h.keys = nil
	h.headers = make(map[string][]*string)
	h.cacheControl = make(map[string]any)
	h.mu.Unlock()
	raw, _ := ctx.GetIndexValue(0)
	return nil, headerBagAddMap(h, raw)
}

func headerBagAdd(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	return nil, headerBagAddMap(h, raw)
}

func headerBagGet(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewNullValue(), nil
	}
	def := defaultValueParam(ctx, 1, data.NewNullValue())
	all := headerAll(h, &key)
	vs := all[normalizeHeaderKey(key)]
	if len(vs) == 0 {
		return def, nil
	}
	if vs[0] == nil {
		return data.NewNullValue(), nil
	}
	return data.NewStringValue(*vs[0]), nil
}

func headerBagSet(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	values := defaultValueParam(ctx, 1, data.NewNullValue())
	replace := boolParam(ctx, 2, true)
	vals, ctl := coerceHeaderValues(values)
	if ctl != nil {
		return nil, ctl
	}
	headerSet(h, key, vals, replace)
	return nil, nil
}

func headerBagHas(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok || h == nil {
		return data.NewBoolValue(false), nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.headers[normalizeHeaderKey(key)]
	return data.NewBoolValue(exists), nil
}

func headerBagContains(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	want, _, ok := optionalStringParam(ctx, 1)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	all := headerAll(h, &key)
	for _, v := range all[normalizeHeaderKey(key)] {
		if v != nil && *v == want {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func headerBagRemove(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	nk := normalizeHeaderKey(key)
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.headers, nk)
	for i, k := range h.keys {
		if k == nk {
			h.keys = append(h.keys[:i], h.keys[i+1:]...)
			break
		}
	}
	if nk == "cache-control" {
		h.cacheControl = make(map[string]any)
	}
	return nil, nil
}

func headerBagGetDate(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewNullValue(), nil
	}
	all := headerAll(h, &key)
	vs := all[normalizeHeaderKey(key)]
	if len(vs) == 0 || vs[0] == nil {
		def := defaultValueParam(ctx, 1, data.NewNullValue())
		if def == nil || isNull(def) {
			return data.NewNullValue(), nil
		}
		return def, nil
	}
	value := *vs[0]
	t, err := time.Parse(time.RFC1123, value)
	if err != nil {
		t, err = time.Parse(time.RFC1123Z, value)
	}
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf(`The "%s" HTTP header is not parseable (%s).`, key, value))
	}
	_ = t
	// 尽量返回 DateTimeImmutable；失败则返回格式化字符串
	stmt, found := ctx.GetVM().GetClass("DateTimeImmutable")
	if !found {
		return data.NewStringValue(value), nil
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return data.NewStringValue(value), nil
	}
	if gsm, ok := stmt.(data.GetStaticMethod); ok {
		if m, ok := gsm.GetStaticMethod("createFromFormat"); ok && m != nil {
			callCtx := ctx.CreateContext(m.GetVariables())
			vars := m.GetVariables()
			if len(vars) > 0 {
				_ = callCtx.SetVariableValue(vars[0], data.NewStringValue(time.RFC1123))
			}
			if len(vars) > 1 {
				_ = callCtx.SetVariableValue(vars[1], data.NewStringValue(value))
			}
			ret, c2 := m.Call(callCtx)
			if c2 == nil && ret != nil {
				return ret, nil
			}
		}
	}
	return cv, nil
}

func headerBagAddCacheControlDirective(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	val := defaultValueParam(ctx, 1, data.NewBoolValue(true))
	h.mu.Lock()
	if h.cacheControl == nil {
		h.cacheControl = make(map[string]any)
	}
	if b, ok := val.(*data.BoolValue); ok {
		bv, _ := b.AsBool()
		h.cacheControl[key] = bv
	} else {
		h.cacheControl[key] = val.AsString()
	}
	cc := getCacheControlHeader(h.cacheControl)
	h.mu.Unlock()
	s := cc
	headerSet(h, "Cache-Control", []*string{&s}, true)
	return nil, nil
}

func headerBagHasCacheControlDirective(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok || h == nil {
		return data.NewBoolValue(false), nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.cacheControl[key]
	return data.NewBoolValue(exists), nil
}

func headerBagGetCacheControlDirective(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok || h == nil {
		return data.NewNullValue(), nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	v, exists := h.cacheControl[key]
	if !exists {
		return data.NewNullValue(), nil
	}
	switch t := v.(type) {
	case bool:
		return data.NewBoolValue(t), nil
	case string:
		return data.NewStringValue(t), nil
	default:
		return data.NewStringValue(fmt.Sprint(t)), nil
	}
}

func headerBagRemoveCacheControlDirective(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	h.mu.Lock()
	delete(h.cacheControl, key)
	cc := getCacheControlHeader(h.cacheControl)
	h.mu.Unlock()
	s := cc
	headerSet(h, "Cache-Control", []*string{&s}, true)
	return nil, nil
}

func headerBagGetIterator(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	var arr data.Value = data.NewArrayValue(nil)
	if h != nil {
		h.mu.RLock()
		arr = headersMapToArrayValue(h.headers, h.keys)
		h.mu.RUnlock()
	}
	return newArrayIterator(ctx, arr)
}

func headerBagCount(ctx data.Context) (data.GetValue, data.Control) {
	h := headerData(ctx)
	if h == nil {
		return data.NewIntValue(0), nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return data.NewIntValue(len(h.keys)), nil
}

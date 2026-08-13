package httpfoundation

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	fqnParameterBag      = "Symfony\\Component\\HttpFoundation\\ParameterBag"
	fqnHeaderBag         = "Symfony\\Component\\HttpFoundation\\HeaderBag"
	fqnInputBag          = "Symfony\\Component\\HttpFoundation\\InputBag"
	fqnServerBag         = "Symfony\\Component\\HttpFoundation\\ServerBag"
	fqnFileBag           = "Symfony\\Component\\HttpFoundation\\FileBag"
	fqnResponseHeaderBag = "Symfony\\Component\\HttpFoundation\\ResponseHeaderBag"
)

// ParamBagData 是 ParameterBag / InputBag / ServerBag / FileBag 的底层存储。
type ParamBagData struct {
	mu     sync.RWMutex
	keys   []string
	values map[string]data.Value
}

func newParamBagData() *ParamBagData {
	return &ParamBagData{values: make(map[string]data.Value)}
}

func (p *ParamBagData) clone() *ParamBagData {
	if p == nil {
		return newParamBagData()
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := &ParamBagData{
		keys:   append([]string(nil), p.keys...),
		values: make(map[string]data.Value, len(p.values)),
	}
	for k, v := range p.values {
		out.values[k] = v
	}
	return out
}

func (p *ParamBagData) all() map[string]data.Value {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make(map[string]data.Value, len(p.values))
	for _, k := range p.keys {
		out[k] = p.values[k]
	}
	return out
}

func (p *ParamBagData) orderedPairs() []struct {
	key string
	val data.Value
} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]struct {
		key string
		val data.Value
	}, 0, len(p.keys))
	for _, k := range p.keys {
		out = append(out, struct {
			key string
			val data.Value
		}{k, p.values[k]})
	}
	return out
}

func (p *ParamBagData) replace(m map[string]data.Value) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.keys = p.keys[:0]
	p.values = make(map[string]data.Value, len(m))
	for k, v := range m {
		p.keys = append(p.keys, k)
		p.values[k] = v
	}
}

func (p *ParamBagData) replaceOrdered(keys []string, values map[string]data.Value) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.keys = append([]string(nil), keys...)
	p.values = make(map[string]data.Value, len(values))
	for k, v := range values {
		p.values[k] = v
	}
}

func (p *ParamBagData) add(m map[string]data.Value) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.values == nil {
		p.values = make(map[string]data.Value)
	}
	for k, v := range m {
		if _, ok := p.values[k]; !ok {
			p.keys = append(p.keys, k)
		}
		p.values[k] = v
	}
}

func (p *ParamBagData) get(key string) (data.Value, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	v, ok := p.values[key]
	return v, ok
}

func (p *ParamBagData) set(key string, value data.Value) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.values == nil {
		p.values = make(map[string]data.Value)
	}
	if _, ok := p.values[key]; !ok {
		p.keys = append(p.keys, key)
	}
	p.values[key] = value
}

func (p *ParamBagData) has(key string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.values[key]
	return ok
}

func (p *ParamBagData) remove(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.values[key]; !ok {
		return
	}
	delete(p.values, key)
	for i, k := range p.keys {
		if k == key {
			p.keys = append(p.keys[:i], p.keys[i+1:]...)
			break
		}
	}
}

func (p *ParamBagData) count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.keys)
}

func (p *ParamBagData) keyList() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return append([]string(nil), p.keys...)
}

func (p *ParamBagData) toArrayValue() *data.ArrayValue {
	pairs := p.orderedPairs()
	list := make([]*data.ZVal, 0, len(pairs))
	for _, pair := range pairs {
		list = append(list, &data.ZVal{Name: pair.key, Value: pair.val})
	}
	return &data.ArrayValue{List: list}
}

// HeaderBagData 是 HeaderBag / ResponseHeaderBag 的底层存储。
type HeaderBagData struct {
	mu           sync.RWMutex
	keys         []string
	headers      map[string][]*string // 小写键 -> 值列表（元素可为 nil）
	cacheControl map[string]any       // string|bool
}

func newHeaderBagData() *HeaderBagData {
	return &HeaderBagData{
		headers:      make(map[string][]*string),
		cacheControl: make(map[string]any),
	}
}

func (h *HeaderBagData) clone() *HeaderBagData {
	if h == nil {
		return newHeaderBagData()
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := &HeaderBagData{
		keys:         append([]string(nil), h.keys...),
		headers:      make(map[string][]*string, len(h.headers)),
		cacheControl: make(map[string]any, len(h.cacheControl)),
	}
	for k, vs := range h.headers {
		cp := make([]*string, len(vs))
		for i, v := range vs {
			if v != nil {
				s := *v
				cp[i] = &s
			}
		}
		out.headers[k] = cp
	}
	for k, v := range h.cacheControl {
		out.cacheControl[k] = v
	}
	return out
}

// BagCookie 是 ResponseHeaderBag 内部 Cookie 表示（供同包 Request/Response 使用）。
type BagCookie struct {
	Name        string
	Value       *string
	Expire      int64
	Path        string
	Domain      *string
	Secure      bool
	HTTPOnly    bool
	Raw         bool
	SameSite    *string
	Partitioned bool
}

func (c *BagCookie) String() string {
	if c == nil {
		return ""
	}
	name := c.Name
	val := ""
	if c.Value != nil {
		val = *c.Value
	}
	parts := []string{name + "=" + val}
	if c.Expire > 0 {
		parts = append(parts, "expires="+strconv.FormatInt(c.Expire, 10))
	}
	path := c.Path
	if path == "" {
		path = "/"
	}
	parts = append(parts, "path="+path)
	if c.Domain != nil && *c.Domain != "" {
		parts = append(parts, "domain="+*c.Domain)
	}
	if c.Secure {
		parts = append(parts, "secure")
	}
	if c.HTTPOnly {
		parts = append(parts, "httponly")
	}
	if c.SameSite != nil && *c.SameSite != "" {
		parts = append(parts, "samesite="+*c.SameSite)
	}
	if c.Partitioned {
		parts = append(parts, "partitioned")
	}
	return strings.Join(parts, "; ")
}

// ResponseHeaderBagData 扩展 HeaderBagData。
type ResponseHeaderBagData struct {
	HeaderBagData
	cookies              map[string]map[string]map[string]*BagCookie // domain -> path -> name
	headerNames          map[string]string                           // 小写 -> 原始大小写
	computedCacheControl map[string]any
}

func newResponseHeaderBagData() *ResponseHeaderBagData {
	return &ResponseHeaderBagData{
		HeaderBagData:        *newHeaderBagData(),
		cookies:              make(map[string]map[string]map[string]*BagCookie),
		headerNames:          make(map[string]string),
		computedCacheControl: make(map[string]any),
	}
}

// ---------- ClassValue / source helpers（同包 Request/Response 可用） ----------

func bagClassValue(ctx data.Context) *data.ClassValue {
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		return c.ClassValue
	case *data.ClassValue:
		return c
	default:
		return nil
	}
}

func classSource(ctx data.Context) any {
	cv := bagClassValue(ctx)
	if cv == nil {
		return nil
	}
	return cv.GetSource()
}

// ParamBagFrom 从 ClassValue 取出 ParamBagData。
func ParamBagFrom(cv *data.ClassValue) *ParamBagData {
	if cv == nil {
		return nil
	}
	if p, ok := cv.GetSource().(*ParamBagData); ok {
		return p
	}
	return nil
}

// HeaderBagFrom 从 ClassValue 取出 HeaderBagData（含 ResponseHeaderBag）。
func HeaderBagFrom(cv *data.ClassValue) *HeaderBagData {
	if cv == nil {
		return nil
	}
	switch s := cv.GetSource().(type) {
	case *HeaderBagData:
		return s
	case *ResponseHeaderBagData:
		return &s.HeaderBagData
	}
	return nil
}

// ResponseHeaderBagFrom 从 ClassValue 取出 ResponseHeaderBagData。
func ResponseHeaderBagFrom(cv *data.ClassValue) *ResponseHeaderBagData {
	if cv == nil {
		return nil
	}
	if s, ok := cv.GetSource().(*ResponseHeaderBagData); ok {
		return s
	}
	return nil
}

func paramData(ctx data.Context) *ParamBagData {
	return ParamBagFrom(bagClassValue(ctx))
}

func headerData(ctx data.Context) *HeaderBagData {
	return HeaderBagFrom(bagClassValue(ctx))
}

func responseHeaderData(ctx data.Context) *ResponseHeaderBagData {
	return ResponseHeaderBagFrom(bagClassValue(ctx))
}

// NewParameterBagValue 创建已填充的 ParameterBag 实例。
func NewParameterBagValue(ctx data.Context, params map[string]data.Value) *data.ClassValue {
	src := newParamBagData()
	if params != nil {
		src.replace(params)
	}
	return newBagInstance(ctx, NewParameterBagClassFrom(src))
}

// NewInputBagValue 创建已填充的 InputBag 实例。
func NewInputBagValue(ctx data.Context, params map[string]data.Value) *data.ClassValue {
	src := newParamBagData()
	if params != nil {
		for k, v := range params {
			src.set(k, v)
		}
	}
	return newBagInstance(ctx, NewInputBagClassFrom(src))
}

// NewServerBagValue 创建已填充的 ServerBag 实例。
func NewServerBagValue(ctx data.Context, params map[string]data.Value) *data.ClassValue {
	src := newParamBagData()
	if params != nil {
		src.replace(params)
	}
	return newBagInstance(ctx, NewServerBagClassFrom(src))
}

// NewFileBagValue 创建已填充的 FileBag 实例。
func NewFileBagValue(ctx data.Context, params map[string]data.Value) *data.ClassValue {
	src := newParamBagData()
	inst := newBagInstance(ctx, NewFileBagClassFrom(src))
	if params != nil {
		_ = fileBagReplace(inst, params)
	}
	return inst
}

// NewHeaderBagValue 创建已填充的 HeaderBag 实例。
func NewHeaderBagValue(ctx data.Context, headers map[string][]string) *data.ClassValue {
	src := newHeaderBagData()
	inst := newBagInstance(ctx, NewHeaderBagClassFrom(src))
	if headers != nil {
		for k, vs := range headers {
			headerSet(src, k, stringSliceToNullable(vs), true)
		}
	}
	return inst
}

// NewResponseHeaderBagValue 创建已填充的 ResponseHeaderBag 实例。
func NewResponseHeaderBagValue(ctx data.Context, headers map[string][]string) *data.ClassValue {
	src := newResponseHeaderBagData()
	inst := newBagInstance(ctx, NewResponseHeaderBagClassFrom(src))
	if headers != nil {
		for k, vs := range headers {
			responseHeaderSet(src, k, stringSliceToNullable(vs), true)
		}
	}
	ensureResponseDefaults(src)
	return inst
}

func newBagInstance(ctx data.Context, stmt data.ClassStmt) *data.ClassValue {
	base := ctx
	if cv := bagClassValue(ctx); cv != nil && cv.Context != nil {
		base = cv.Context
	}
	return data.NewProxyValue(stmt, base.CreateBaseContext())
}

// GetParamBagMap 读取 ParameterBag 系实例的全部参数。
func GetParamBagMap(cv *data.ClassValue) map[string]data.Value {
	if p := ParamBagFrom(cv); p != nil {
		return p.all()
	}
	return map[string]data.Value{}
}

// SetParamBagMap 替换 ParameterBag 系实例的全部参数。
func SetParamBagMap(cv *data.ClassValue, params map[string]data.Value) {
	if p := ParamBagFrom(cv); p != nil {
		if params == nil {
			params = map[string]data.Value{}
		}
		p.replace(params)
	}
}

// GetHeaderBagAll 读取 HeaderBag 全部头（小写键）。
func GetHeaderBagAll(cv *data.ClassValue) map[string][]string {
	h := HeaderBagFrom(cv)
	if h == nil {
		return map[string][]string{}
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make(map[string][]string, len(h.headers))
	for k, vs := range h.headers {
		arr := make([]string, 0, len(vs))
		for _, v := range vs {
			if v == nil {
				arr = append(arr, "")
			} else {
				arr = append(arr, *v)
			}
		}
		out[k] = arr
	}
	return out
}

// SetHeaderBagHeader 设置单个头。
func SetHeaderBagHeader(cv *data.ClassValue, key string, values []string, replace bool) {
	if rh := ResponseHeaderBagFrom(cv); rh != nil {
		responseHeaderSet(rh, key, stringSliceToNullable(values), replace)
		return
	}
	if h := HeaderBagFrom(cv); h != nil {
		headerSet(h, key, stringSliceToNullable(values), replace)
	}
}

func stringSliceToNullable(vs []string) []*string {
	out := make([]*string, len(vs))
	for i, v := range vs {
		s := v
		out[i] = &s
	}
	return out
}

func nullableToStringSlice(vs []*string) []string {
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		if v == nil {
			out = append(out, "")
		} else {
			out = append(out, *v)
		}
	}
	return out
}

func valueToAssocMap(v data.Value) (map[string]data.Value, error) {
	keys, values, err := valueToOrderedAssoc(v)
	if err != nil {
		return nil, err
	}
	_ = keys
	return values, nil
}

func valueToOrderedAssoc(v data.Value) ([]string, map[string]data.Value, error) {
	if v == nil {
		return nil, map[string]data.Value{}, nil
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		keys := make([]string, 0, len(arr.List))
		out := make(map[string]data.Value, len(arr.List))
		for i, z := range arr.List {
			if z == nil {
				continue
			}
			key := z.Name
			if key == "" {
				key = strconv.Itoa(i)
			}
			if _, exists := out[key]; !exists {
				keys = append(keys, key)
			}
			out[key] = z.Value
		}
		return keys, out, nil
	case *data.ObjectValue:
		keys := make([]string, 0)
		out := make(map[string]data.Value)
		arr.RangeProperties(func(key string, value data.Value) bool {
			keys = append(keys, key)
			out[key] = value
			return true
		})
		return keys, out, nil
	case *data.NullValue:
		return nil, map[string]data.Value{}, nil
	default:
		return nil, nil, fmt.Errorf("expected array, got %T", v)
	}
}

func assocMapToArrayValue(m map[string]data.Value) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(m))
	for k, v := range m {
		list = append(list, &data.ZVal{Name: k, Value: v})
	}
	return &data.ArrayValue{List: list}
}

func orderedAssocToArrayValue(keys []string, values map[string]data.Value) *data.ArrayValue {
	list := make([]*data.ZVal, 0, len(keys))
	for _, k := range keys {
		list = append(list, &data.ZVal{Name: k, Value: values[k]})
	}
	return &data.ArrayValue{List: list}
}

func isArrayValue(v data.Value) bool {
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return true
	default:
		return false
	}
}

func isScalarOrStringable(v data.Value) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case *data.NullValue, *data.StringValue, *data.IntValue, *data.FloatValue, *data.BoolValue:
		return true
	case *data.ClassValue:
		cv := v.(*data.ClassValue)
		_, ok := cv.GetMethod("__toString")
		return ok
	default:
		return false
	}
}

func valueAsString(v data.Value) (string, error) {
	if v == nil {
		return "", nil
	}
	if !isScalarOrStringable(v) {
		return "", fmt.Errorf("cannot convert to string")
	}
	if cv, ok := v.(*data.ClassValue); ok {
		if m, ok := cv.GetMethod("__toString"); ok && m != nil {
			ret, ctl := m.Call(cv.CreateContext(m.GetVariables()))
			if ctl != nil {
				return "", fmt.Errorf("toString failed")
			}
			if ret != nil {
				if val, ok := ret.(data.Value); ok {
					return val.AsString(), nil
				}
			}
		}
	}
	return v.AsString(), nil
}

func throwNamed(name string, format string, args ...any) data.Control {
	return data.NewErrorThrowByName(nil, fmt.Errorf(format, args...), name)
}

func newArrayIterator(ctx data.Context, storage data.Value) (data.GetValue, data.Control) {
	stmt, ok := ctx.GetVM().GetClass("ArrayIterator")
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ArrayIterator class not found"))
	}
	obj, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}
	cv, ok := obj.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("ArrayIterator invalid"))
	}
	method := cv.Class.GetConstruct()
	if method == nil {
		return cv, nil
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	if len(method.GetVariables()) > 0 {
		_ = fnCtx.SetVariableValue(method.GetVariables()[0], storage)
	}
	_, acl = method.Call(fnCtx)
	if acl != nil {
		return nil, acl
	}
	return cv, nil
}

func thisClassValue(ctx data.Context) data.GetValue {
	if cv := bagClassValue(ctx); cv != nil {
		return cv
	}
	return data.NewNullValue()
}

func optionalStringParam(ctx data.Context, index int) (string, bool, bool) {
	v, ok := ctx.GetIndexValue(index)
	if !ok || v == nil {
		return "", false, false
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return "", true, false
	}
	return v.AsString(), true, true
}

func defaultValueParam(ctx data.Context, index int, fallback data.Value) data.Value {
	v, ok := ctx.GetIndexValue(index)
	if !ok || v == nil {
		return fallback
	}
	return v
}

func boolParam(ctx data.Context, index int, def bool) bool {
	v, ok := ctx.GetIndexValue(index)
	if !ok || v == nil {
		return def
	}
	if b, ok := v.(data.AsBool); ok {
		if bv, err := b.AsBool(); err == nil {
			return bv
		}
	}
	s := strings.ToLower(v.AsString())
	return s == "1" || s == "true" || s == "on" || s == "yes"
}

func intParam(ctx data.Context, index int, def int) int {
	v, ok := ctx.GetIndexValue(index)
	if !ok || v == nil {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	n, err := strconv.Atoi(v.AsString())
	if err != nil {
		return def
	}
	return n
}

// 供属性声明复用
func protectedArrayProp(name string) data.Property {
	return node.NewProperty(nil, name, "protected", false, data.NewArrayValue(nil))
}

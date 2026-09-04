package httpfoundation

import (
	"fmt"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	cookiesFlat           = "flat"
	cookiesArray          = "array"
	dispositionAttachment = "attachment"
	dispositionInline     = "inline"
)

// ResponseHeaderBagClass 实现 Symfony\Component\HttpFoundation\ResponseHeaderBag。
type ResponseHeaderBagClass struct {
	node.Node
	source     *ResponseHeaderBagData
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewResponseHeaderBagClass() data.ClassStmt {
	return NewResponseHeaderBagClassFrom(nil)
}

func NewResponseHeaderBagClassFrom(source *ResponseHeaderBagData) data.ClassStmt {
	c := &ResponseHeaderBagClass{
		source: source,
		properties: []data.Property{
			protectedArrayProp("headers"),
			protectedArrayProp("cacheControl"),
			protectedArrayProp("computedCacheControl"),
			protectedArrayProp("cookies"),
			protectedArrayProp("headerNames"),
		},
	}
	c.methods, c.methodList = responseHeaderBagMethods()
	return c
}

func (c *ResponseHeaderBagClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	src := c.source
	if src == nil {
		src = newResponseHeaderBagData()
		ensureResponseDefaults(src)
	} else {
		cloned := &ResponseHeaderBagData{
			HeaderBagData:        *src.HeaderBagData.clone(),
			cookies:              cloneCookies(src.cookies),
			headerNames:          copyStringMap(src.headerNames),
			computedCacheControl: copyAnyMap(src.computedCacheControl),
		}
		src = cloned
	}
	return data.NewProxyValue(NewResponseHeaderBagClassFrom(src), ctx.CreateBaseContext()), nil
}

func (c *ResponseHeaderBagClass) GetName() string {
	return fqnResponseHeaderBag
}
func (c *ResponseHeaderBagClass) GetExtend() *string {
	parent := fqnHeaderBag
	return &parent
}
func (c *ResponseHeaderBagClass) GetImplements() []string          { return nil }
func (c *ResponseHeaderBagClass) GetSource() any                   { return c.source }
func (c *ResponseHeaderBagClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *ResponseHeaderBagClass) GetPropertyList() []data.Property { return c.properties }
func (c *ResponseHeaderBagClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *ResponseHeaderBagClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *ResponseHeaderBagClass) GetMethods() []data.Method { return c.methodList }

func (c *ResponseHeaderBagClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "COOKIES_FLAT":
		return data.NewStringValue(cookiesFlat), true
	case "COOKIES_ARRAY":
		return data.NewStringValue(cookiesArray), true
	case "DISPOSITION_ATTACHMENT":
		return data.NewStringValue(dispositionAttachment), true
	case "DISPOSITION_INLINE":
		return data.NewStringValue(dispositionInline), true
	}
	return nil, false
}

func copyStringMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func copyAnyMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func cloneCookies(in map[string]map[string]map[string]*BagCookie) map[string]map[string]map[string]*BagCookie {
	out := make(map[string]map[string]map[string]*BagCookie, len(in))
	for d, paths := range in {
		out[d] = make(map[string]map[string]*BagCookie, len(paths))
		for p, names := range paths {
			out[d][p] = make(map[string]*BagCookie, len(names))
			for n, c := range names {
				if c == nil {
					continue
				}
				cp := *c
				if c.Value != nil {
					v := *c.Value
					cp.Value = &v
				}
				if c.Domain != nil {
					d2 := *c.Domain
					cp.Domain = &d2
				}
				if c.SameSite != nil {
					s := *c.SameSite
					cp.SameSite = &s
				}
				out[d][p][n] = &cp
			}
		}
	}
	return out
}

func ensureResponseDefaults(src *ResponseHeaderBagData) {
	if src == nil {
		return
	}
	src.mu.Lock()
	_, hasCC := src.headers["cache-control"]
	_, hasDate := src.headers["date"]
	src.mu.Unlock()
	if !hasCC {
		responseHeaderSet(src, "Cache-Control", []*string{strPtr("")}, true)
	}
	if !hasDate {
		initResponseDate(src)
	}
}

func strPtr(s string) *string { return &s }

func initResponseDate(src *ResponseHeaderBagData) {
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05") + " GMT"
	responseHeaderSet(src, "Date", []*string{&date}, true)
}

func responseHeaderBagMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{param("headers", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("headers", 0, nil)},
			nil, responseHeaderBagConstruct),
		pubMethod("allPreserveCase", nil, nil, data.NewBaseType("array"), responseHeaderBagAllPreserveCase),
		pubMethod("allPreserveCaseWithoutCookies", nil, nil, data.NewBaseType("array"), responseHeaderBagAllPreserveCaseWithoutCookies),
		pubMethod("replace",
			[]data.GetValue{param("headers", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("headers", 0, nil)},
			nil, responseHeaderBagReplace),
		pubMethod("all",
			[]data.GetValue{param("key", 0, data.NewNullValue(), nil)},
			[]data.Variable{variable("key", 0, nil)},
			data.NewBaseType("array"), responseHeaderBagAll),
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
			nil, responseHeaderBagSet),
		pubMethod("remove",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			nil, responseHeaderBagRemove),
		pubMethod("hasCacheControlDirective",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			data.NewBaseType("bool"), responseHeaderBagHasCacheControlDirective),
		pubMethod("getCacheControlDirective",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			nil, responseHeaderBagGetCacheControlDirective),
		pubMethod("setCookie",
			[]data.GetValue{param("cookie", 0, nil, nil)},
			[]data.Variable{variable("cookie", 0, nil)},
			nil, responseHeaderBagSetCookie),
		pubMethod("removeCookie",
			[]data.GetValue{
				param("name", 0, nil, nil),
				param("path", 1, data.NewStringValue("/"), nil),
				param("domain", 2, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("name", 0, nil),
				variable("path", 1, nil),
				variable("domain", 2, nil),
			},
			nil, responseHeaderBagRemoveCookie),
		pubMethod("getCookies",
			[]data.GetValue{param("format", 0, data.NewStringValue(cookiesFlat), nil)},
			[]data.Variable{variable("format", 0, nil)},
			data.NewBaseType("array"), responseHeaderBagGetCookies),
		pubMethod("clearCookie",
			[]data.GetValue{
				param("name", 0, nil, nil),
				param("path", 1, data.NewStringValue("/"), nil),
				param("domain", 2, data.NewNullValue(), nil),
				param("secure", 3, data.NewBoolValue(false), nil),
				param("httpOnly", 4, data.NewBoolValue(true), nil),
				param("sameSite", 5, data.NewNullValue(), nil),
				param("partitioned", 6, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				variable("name", 0, nil),
				variable("path", 1, nil),
				variable("domain", 2, nil),
				variable("secure", 3, nil),
				variable("httpOnly", 4, nil),
				variable("sameSite", 5, nil),
				variable("partitioned", 6, nil),
			},
			nil, responseHeaderBagClearCookie),
		pubMethod("makeDisposition",
			[]data.GetValue{
				param("disposition", 0, nil, nil),
				param("filename", 1, nil, nil),
				param("filenameFallback", 2, data.NewStringValue(""), nil),
			},
			[]data.Variable{
				variable("disposition", 0, nil),
				variable("filename", 1, nil),
				variable("filenameFallback", 2, nil),
			},
			data.NewBaseType("string"), responseHeaderBagMakeDisposition),
	}
	m := make(map[string]data.Method, len(list))
	for _, method := range list {
		m[method.GetName()] = method
	}
	return m, list
}

func responseHeaderBagConstruct(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	m, err := valueToAssocMap(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	for k, v := range m {
		vals, ctl := coerceHeaderValues(v)
		if ctl != nil {
			return nil, ctl
		}
		responseHeaderSet(src, k, vals, true)
	}
	ensureResponseDefaults(src)
	return nil, nil
}

func responseHeaderSet(src *ResponseHeaderBagData, key string, values []*string, replace bool) {
	if src == nil {
		return
	}
	uniqueKey := normalizeHeaderKey(key)
	if uniqueKey == "set-cookie" {
		if replace {
			src.cookies = make(map[string]map[string]map[string]*BagCookie)
		}
		for _, v := range values {
			if v == nil {
				continue
			}
			cookie := cookieFromString(*v)
			setBagCookie(src, cookie)
		}
		src.headerNames[uniqueKey] = key
		return
	}
	src.headerNames[uniqueKey] = key
	headerSet(&src.HeaderBagData, key, values, replace)

	if uniqueKey == "cache-control" || uniqueKey == "etag" || uniqueKey == "last-modified" || uniqueKey == "expires" {
		computed := computeCacheControlValue(src)
		if computed != "" {
			src.mu.Lock()
			src.headers["cache-control"] = []*string{&computed}
			if _, ok := keyInSlice(src.keys, "cache-control"); !ok {
				src.keys = append(src.keys, "cache-control")
			}
			src.headerNames["cache-control"] = "Cache-Control"
			src.computedCacheControl = parseCacheControl(computed)
			src.mu.Unlock()
		}
	}
}

func keyInSlice(keys []string, key string) (int, bool) {
	for i, k := range keys {
		if k == key {
			return i, true
		}
	}
	return -1, false
}

func computeCacheControlValue(src *ResponseHeaderBagData) string {
	src.mu.RLock()
	defer src.mu.RUnlock()
	if len(src.cacheControl) == 0 {
		if _, ok := src.headers["last-modified"]; ok {
			return "private, must-revalidate"
		}
		if _, ok := src.headers["expires"]; ok {
			return "private, must-revalidate"
		}
		return "no-cache, private"
	}
	header := getCacheControlHeader(src.cacheControl)
	if _, ok := src.cacheControl["public"]; ok {
		return header
	}
	if _, ok := src.cacheControl["private"]; ok {
		return header
	}
	if _, ok := src.cacheControl["s-maxage"]; !ok {
		return header + ", private"
	}
	return header
}

func setBagCookie(src *ResponseHeaderBagData, cookie *BagCookie) {
	if src == nil || cookie == nil {
		return
	}
	domain := ""
	if cookie.Domain != nil {
		domain = *cookie.Domain
	}
	path := cookie.Path
	if path == "" {
		path = "/"
	}
	if src.cookies == nil {
		src.cookies = make(map[string]map[string]map[string]*BagCookie)
	}
	if src.cookies[domain] == nil {
		src.cookies[domain] = make(map[string]map[string]*BagCookie)
	}
	if src.cookies[domain][path] == nil {
		src.cookies[domain][path] = make(map[string]*BagCookie)
	}
	src.cookies[domain][path][cookie.Name] = cookie
	src.headerNames["set-cookie"] = "Set-Cookie"
}

func cookieFromString(raw string) *BagCookie {
	parts := headerUtilsSplit(raw, ";=")
	if len(parts) == 0 || len(parts[0]) == 0 {
		return &BagCookie{Name: raw, Path: "/"}
	}
	name := parts[0][0]
	var value *string
	if len(parts[0]) > 1 {
		v := parts[0][1]
		value = &v
	}
	c := &BagCookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HTTPOnly: false,
	}
	assoc := headerUtilsCombine(parts[1:])
	if p, ok := assoc["path"].(string); ok {
		c.Path = p
	}
	if d, ok := assoc["domain"].(string); ok {
		c.Domain = &d
	}
	if _, ok := assoc["secure"]; ok {
		c.Secure = true
	}
	if _, ok := assoc["httponly"]; ok {
		c.HTTPOnly = true
	}
	if s, ok := assoc["samesite"].(string); ok {
		c.SameSite = &s
	}
	if _, ok := assoc["partitioned"]; ok {
		c.Partitioned = true
	}
	if e, ok := assoc["expires"].(string); ok {
		if t, err := time.Parse(time.RFC1123, e); err == nil {
			c.Expire = t.Unix()
		}
	}
	return c
}

func cookieFromValue(v data.Value) *BagCookie {
	if v == nil || isNull(v) {
		return nil
	}
	if cv, ok := v.(*data.ClassValue); ok {
		get := func(name string) data.Value {
			val, _ := cv.GetProperty(name)
			return val
		}
		name := ""
		if n := get("name"); n != nil {
			name = n.AsString()
		}
		// Symfony Cookie 属性多为 private，尝试方法
		if m, ok := cv.GetMethod("getName"); ok && m != nil {
			if ret, ctl := m.Call(cv.CreateContext(m.GetVariables())); ctl == nil && ret != nil {
				if val, ok := ret.(data.Value); ok {
					name = val.AsString()
				}
			}
		}
		c := &BagCookie{Name: name, Path: "/"}
		if m, ok := cv.GetMethod("getValue"); ok && m != nil {
			if ret, ctl := m.Call(cv.CreateContext(m.GetVariables())); ctl == nil && ret != nil {
				if val, ok := ret.(data.Value); ok && !isNull(val) {
					s := val.AsString()
					c.Value = &s
				}
			}
		}
		if m, ok := cv.GetMethod("getPath"); ok && m != nil {
			if ret, ctl := m.Call(cv.CreateContext(m.GetVariables())); ctl == nil && ret != nil {
				if val, ok := ret.(data.Value); ok {
					c.Path = val.AsString()
				}
			}
		}
		if m, ok := cv.GetMethod("getDomain"); ok && m != nil {
			if ret, ctl := m.Call(cv.CreateContext(m.GetVariables())); ctl == nil && ret != nil {
				if val, ok := ret.(data.Value); ok && !isNull(val) {
					s := val.AsString()
					c.Domain = &s
				}
			}
		}
		if m, ok := cv.GetMethod("__toString"); ok && m != nil {
			if ret, ctl := m.Call(cv.CreateContext(m.GetVariables())); ctl == nil && ret != nil {
				if val, ok := ret.(data.Value); ok {
					parsed := cookieFromString(val.AsString())
					if parsed.Name != "" {
						return parsed
					}
				}
			}
		}
		if c.Name != "" {
			return c
		}
	}
	return cookieFromString(v.AsString())
}

func flatCookies(src *ResponseHeaderBagData) []*BagCookie {
	var out []*BagCookie
	if src == nil {
		return out
	}
	for _, paths := range src.cookies {
		for _, names := range paths {
			for _, c := range names {
				out = append(out, c)
			}
		}
	}
	return out
}

func responseHeaderBagAllPreserveCase(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
		return data.NewArrayValue(nil), nil
	}
	src.mu.RLock()
	defer src.mu.RUnlock()
	list := make([]*data.ZVal, 0, len(src.headers))
	for _, k := range src.keys {
		name := src.headerNames[k]
		if name == "" {
			name = k
		}
		list = append(list, &data.ZVal{Name: name, Value: headerValuesToArray(src.headers[k])})
	}
	return &data.ArrayValue{List: list}, nil
}

func responseHeaderBagAllPreserveCaseWithoutCookies(ctx data.Context) (data.GetValue, data.Control) {
	ret, ctl := responseHeaderBagAllPreserveCase(ctx)
	if ctl != nil {
		return nil, ctl
	}
	arr, ok := ret.(*data.ArrayValue)
	if !ok {
		return ret, nil
	}
	src := responseHeaderData(ctx)
	cookieName := "Set-Cookie"
	if src != nil {
		if n, ok := src.headerNames["set-cookie"]; ok {
			cookieName = n
		}
	}
	filtered := make([]*data.ZVal, 0, len(arr.List))
	for _, z := range arr.List {
		if z != nil && z.Name == cookieName {
			continue
		}
		filtered = append(filtered, z)
	}
	return &data.ArrayValue{List: filtered}, nil
}

func responseHeaderBagReplace(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
		return nil, nil
	}
	src.headerNames = make(map[string]string)
	src.cookies = make(map[string]map[string]map[string]*BagCookie)
	src.computedCacheControl = make(map[string]any)
	src.mu.Lock()
	src.keys = nil
	src.headers = make(map[string][]*string)
	src.cacheControl = make(map[string]any)
	src.mu.Unlock()

	raw, _ := ctx.GetIndexValue(0)
	m, err := valueToAssocMap(raw)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	for k, v := range m {
		vals, ctl := coerceHeaderValues(v)
		if ctl != nil {
			return nil, ctl
		}
		responseHeaderSet(src, k, vals, true)
	}
	ensureResponseDefaults(src)
	return nil, nil
}

func responseHeaderBagAll(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	key, present, isStr := optionalStringParam(ctx, 0)
	if present && isStr {
		nk := normalizeHeaderKey(key)
		if nk == "set-cookie" {
			cookies := flatCookies(src)
			vals := make([]data.Value, len(cookies))
			for i, c := range cookies {
				vals[i] = data.NewStringValue(c.String())
			}
			return data.NewArrayValue(vals), nil
		}
		return headerBagAll(ctx)
	}
	h := headerData(ctx)
	var arr *data.ArrayValue
	if h != nil {
		h.mu.RLock()
		arr = headersMapToArrayValue(h.headers, h.keys)
		h.mu.RUnlock()
	} else {
		arr = data.NewArrayValue(nil).(*data.ArrayValue)
	}
	cookies := flatCookies(src)
	if len(cookies) > 0 {
		vals := make([]data.Value, len(cookies))
		for i, c := range cookies {
			vals[i] = data.NewStringValue(c.String())
		}
		// 合并/追加 set-cookie
		found := false
		for _, z := range arr.List {
			if z != nil && z.Name == "set-cookie" {
				z.Value = data.NewArrayValue(vals)
				found = true
				break
			}
		}
		if !found {
			arr.List = append(arr.List, &data.ZVal{Name: "set-cookie", Value: data.NewArrayValue(vals)})
		}
	}
	return arr, nil
}

func responseHeaderBagSet(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
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
	responseHeaderSet(src, key, vals, replace)
	return nil, nil
}

func responseHeaderBagRemove(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
		return nil, nil
	}
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	uniqueKey := normalizeHeaderKey(key)
	delete(src.headerNames, uniqueKey)
	if uniqueKey == "set-cookie" {
		src.cookies = make(map[string]map[string]map[string]*BagCookie)
		return nil, nil
	}
	_, _ = headerBagRemove(ctx)
	if uniqueKey == "cache-control" {
		src.computedCacheControl = make(map[string]any)
	}
	if uniqueKey == "date" {
		initResponseDate(src)
	}
	return nil, nil
}

func responseHeaderBagHasCacheControlDirective(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok || src == nil {
		return data.NewBoolValue(false), nil
	}
	_, exists := src.computedCacheControl[key]
	return data.NewBoolValue(exists), nil
}

func responseHeaderBagGetCacheControlDirective(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	key, _, ok := optionalStringParam(ctx, 0)
	if !ok || src == nil {
		return data.NewNullValue(), nil
	}
	v, exists := src.computedCacheControl[key]
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

func responseHeaderBagSetCookie(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
		return nil, nil
	}
	raw, _ := ctx.GetIndexValue(0)
	cookie := cookieFromValue(raw)
	if cookie == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("setCookie expects a Cookie instance or string"))
	}
	setBagCookie(src, cookie)
	return nil, nil
}

func responseHeaderBagRemoveCookie(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
		return nil, nil
	}
	name, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	path := "/"
	if p, present, isStr := optionalStringParam(ctx, 1); present && isStr {
		path = p
	}
	domain := ""
	if d, present, isStr := optionalStringParam(ctx, 2); present && isStr {
		domain = d
	}
	if paths, ok := src.cookies[domain]; ok {
		if names, ok := paths[path]; ok {
			delete(names, name)
			if len(names) == 0 {
				delete(paths, path)
			}
		}
		if len(paths) == 0 {
			delete(src.cookies, domain)
		}
	}
	if len(src.cookies) == 0 {
		delete(src.headerNames, "set-cookie")
	}
	return nil, nil
}

func responseHeaderBagGetCookies(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	format := cookiesFlat
	if f, present, isStr := optionalStringParam(ctx, 0); present && isStr {
		format = f
	}
	if format != cookiesFlat && format != cookiesArray {
		return nil, data.NewErrorThrow(nil, fmt.Errorf(`Format "%s" invalid (%s, %s).`, format, cookiesFlat, cookiesArray))
	}
	if format == cookiesArray {
		// domain -> path -> name -> cookie string（无 Cookie 类时用字符串代替）
		outer := make([]*data.ZVal, 0)
		if src != nil {
			for domain, paths := range src.cookies {
				pathList := make([]*data.ZVal, 0, len(paths))
				for path, names := range paths {
					nameList := make([]*data.ZVal, 0, len(names))
					for name, c := range names {
						nameList = append(nameList, &data.ZVal{Name: name, Value: data.NewStringValue(c.String())})
					}
					pathList = append(pathList, &data.ZVal{Name: path, Value: &data.ArrayValue{List: nameList}})
				}
				outer = append(outer, &data.ZVal{Name: domain, Value: &data.ArrayValue{List: pathList}})
			}
		}
		return &data.ArrayValue{List: outer}, nil
	}
	cookies := flatCookies(src)
	vals := make([]data.Value, 0, len(cookies))
	for _, c := range cookies {
		if obj := createCookiePHPObject(ctx, c); obj != nil {
			vals = append(vals, obj)
		} else {
			vals = append(vals, data.NewStringValue(c.String()))
		}
	}
	return data.NewArrayValue(vals), nil
}

// createCookiePHPObject 尝试通过 PHP 的 Symfony Cookie 类创建 Cookie 对象。
// 返回 nil 表示无法创建（此时调用方可退化为字符串形式）。
func createCookiePHPObject(ctx data.Context, c *BagCookie) data.Value {
	if c == nil || c.Name == "" {
		return nil
	}
	vm := ctx.GetVM()
	stmt, control := vm.GetOrLoadClass("Symfony\\Component\\HttpFoundation\\Cookie")
	if control != nil || stmt == nil {
		return nil
	}
	gsm, ok := stmt.(data.GetStaticMethod)
	if !ok {
		return nil
	}
	method, ok := gsm.GetStaticMethod("create")
	if !ok || method == nil {
		return nil
	}
	// 创建 Cookie 类的 ClassValue，使 self:: 能正确解析
	classValue := data.NewClassValue(stmt, ctx)
	fnCtx := classValue.CreateContext(method.GetVariables())
	// 绑定类和 static
	if cmc, ok := fnCtx.(*data.ClassMethodContext); ok {
		cmc.StaticClass = stmt
		cmc.SelfClass = stmt
	}

	params := method.GetParams()
	vars := method.GetVariables()

	// 按参数索引设置值
	values := make([]data.Value, len(vars))
	// 设置 name
	if len(vars) > 0 {
		values[0] = data.NewStringValue(c.Name)
	}
	// 设置 value
	if len(vars) > 1 {
		val := ""
		if c.Value != nil {
			val = *c.Value
		}
		values[1] = data.NewStringValue(val)
	}
	// 设置 expire
	if len(vars) > 2 {
		values[2] = data.NewIntValue(int(c.Expire))
	}
	// 设置 path
	if len(vars) > 3 {
		p := c.Path
		if p == "" {
			p = "/"
		}
		values[3] = data.NewStringValue(p)
	}
	// 设置 domain
	if len(vars) > 4 {
		if c.Domain != nil && *c.Domain != "" {
			values[4] = data.NewStringValue(*c.Domain)
		} else {
			values[4] = data.NewNullValue()
		}
	}
	// 设置 secure
	if len(vars) > 5 {
		values[5] = data.NewBoolValue(c.Secure)
	}
	// 设置 httpOnly
	if len(vars) > 6 {
		values[6] = data.NewBoolValue(c.HTTPOnly)
	}
	// 设置 raw
	if len(vars) > 7 {
		values[7] = data.NewBoolValue(c.Raw)
	}
	// 设置 sameSite
	if len(vars) > 8 {
		if c.SameSite != nil && *c.SameSite != "" {
			values[8] = data.NewStringValue(*c.SameSite)
		} else {
			values[8] = data.NewNullValue()
		}
	}
	// 设置 partitioned
	if len(vars) > 9 {
		values[9] = data.NewBoolValue(c.Partitioned)
	}

	// 通过 ZVal 正确设置参数
	for i := 0; i < len(vars) && i < len(values); i++ {
		if values[i] != nil {
			zv := data.NewNamedZVal(vars[i].GetName(), values[i])
			fnCtx.SetIndexZVal(i, zv)
		}
	}

	// 对缺少的参数使用默认值
	for i := len(values); i < len(params); i++ {
		if i < len(vars) {
			if _, acl := params[i].GetValue(fnCtx); acl != nil {
				return nil
			}
		}
	}

	// 记录实参
	fnCtx.SetCallArgs(nil)
	fnCtx.SetFlatCallArgs(values)

	ret, ctl := method.Call(fnCtx)
	if ctl != nil {
		return nil
	}
	if v, ok := ret.(data.Value); ok {
		return v
	}
	return nil
}


func responseHeaderBagClearCookie(ctx data.Context) (data.GetValue, data.Control) {
	src := responseHeaderData(ctx)
	if src == nil {
		return nil, nil
	}
	name, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return nil, nil
	}
	path := "/"
	if p, present, isStr := optionalStringParam(ctx, 1); present && isStr {
		path = p
	}
	var domain *string
	if d, present, isStr := optionalStringParam(ctx, 2); present && isStr {
		domain = &d
	}
	secure := boolParam(ctx, 3, false)
	httpOnly := boolParam(ctx, 4, true)
	var sameSite *string
	if s, present, isStr := optionalStringParam(ctx, 5); present && isStr {
		sameSite = &s
	}
	partitioned := boolParam(ctx, 6, false)
	setBagCookie(src, &BagCookie{
		Name:        name,
		Value:       nil,
		Expire:      1,
		Path:        path,
		Domain:      domain,
		Secure:      secure,
		HTTPOnly:    httpOnly,
		SameSite:    sameSite,
		Partitioned: partitioned,
	})
	return nil, nil
}

func responseHeaderBagMakeDisposition(ctx data.Context) (data.GetValue, data.Control) {
	disposition, _, ok := optionalStringParam(ctx, 0)
	if !ok {
		return data.NewStringValue(""), nil
	}
	filename, _, ok := optionalStringParam(ctx, 1)
	if !ok {
		return data.NewStringValue(""), nil
	}
	fallback := ""
	if f, present, isStr := optionalStringParam(ctx, 2); present && isStr {
		fallback = f
	}
	s, err := makeDisposition(disposition, filename, fallback)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewStringValue(s), nil
}

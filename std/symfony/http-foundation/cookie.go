package httpfoundation

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const fqnCookie = "Symfony\\Component\\HttpFoundation\\Cookie"

// CookieClass 实现 Symfony\Component\HttpFoundation\Cookie。
type CookieClass struct {
	node.Node
	methods map[string]data.Method
}

func NewCookieClass() data.ClassStmt {
	c := &CookieClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = pubMethod("__construct",
		[]data.GetValue{
			param("name", 0, nil, nil),
			param("value", 1, data.NewNullValue(), nil),
			param("expire", 2, data.NewIntValue(0), nil),
			param("path", 3, data.NewStringValue("/"), nil),
			param("domain", 4, data.NewNullValue(), nil),
			param("secure", 5, data.NewNullValue(), nil),
			param("httpOnly", 6, data.NewBoolValue(true), nil),
			param("raw", 7, data.NewBoolValue(false), nil),
			param("sameSite", 8, data.NewStringValue("lax"), nil),
			param("partitioned", 9, data.NewBoolValue(false), nil),
		},
		[]data.Variable{
			variable("name", 0, nil), variable("value", 1, nil), variable("expire", 2, nil),
			variable("path", 3, nil), variable("domain", 4, nil), variable("secure", 5, nil),
			variable("httpOnly", 6, nil), variable("raw", 7, nil), variable("sameSite", 8, nil),
			variable("partitioned", 9, nil),
		},
		nil, cookieConstruct)
	c.methods["getname"] = pubMethod("getName", nil, nil, data.NewBaseType("string"), cookieGetName)
	c.methods["getvalue"] = pubMethod("getValue", nil, nil, nil, cookieGetValue)
	c.methods["getexpirestime"] = pubMethod("getExpiresTime", nil, nil, data.NewBaseType("int"), cookieGetExpires)
	c.methods["getpath"] = pubMethod("getPath", nil, nil, data.NewBaseType("string"), cookieGetPath)
	c.methods["getdomain"] = pubMethod("getDomain", nil, nil, nil, cookieGetDomain)
	c.methods["issecure"] = pubMethod("isSecure", nil, nil, data.NewBaseType("bool"), cookieIsSecure)
	c.methods["ishttponly"] = pubMethod("isHttpOnly", nil, nil, data.NewBaseType("bool"), cookieIsHTTPOnly)
	c.methods["israw"] = pubMethod("isRaw", nil, nil, data.NewBaseType("bool"), cookieIsRaw)
	c.methods["getsamesite"] = pubMethod("getSameSite", nil, nil, nil, cookieGetSameSite)
	c.methods["__tostring"] = pubMethod("__toString", nil, nil, data.NewBaseType("string"), cookieToString)
	c.methods["withvalue"] = pubMethod("withValue",
		[]data.GetValue{param("value", 0, data.NewNullValue(), nil)},
		[]data.Variable{variable("value", 0, nil)},
		data.NewBaseType(fqnCookie), cookieWithValue)
	return c
}

func (c *CookieClass) GetName() string                          { return fqnCookie }
func (c *CookieClass) GetExtend() *string                       { return nil }
func (c *CookieClass) GetImplements() []string                  { return nil }
func (c *CookieClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *CookieClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "name", "private", false, data.NewStringValue("")),
		node.NewProperty(nil, "value", "private", false, data.NewNullValue()),
		node.NewProperty(nil, "expire", "private", false, data.NewIntValue(0)),
		node.NewProperty(nil, "path", "private", false, data.NewStringValue("/")),
		node.NewProperty(nil, "domain", "private", false, data.NewNullValue()),
		node.NewProperty(nil, "secure", "private", false, data.NewBoolValue(false)),
		node.NewProperty(nil, "httpOnly", "private", false, data.NewBoolValue(true)),
		node.NewProperty(nil, "raw", "private", false, data.NewBoolValue(false)),
		node.NewProperty(nil, "sameSite", "private", false, data.NewStringValue("lax")),
		node.NewProperty(nil, "partitioned", "private", false, data.NewBoolValue(false)),
	}
}
func (c *CookieClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *CookieClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *CookieClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *CookieClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *CookieClass) GetStaticMethod(name string) (data.Method, bool) {
	if strings.EqualFold(name, "create") {
		m := pubMethod("create",
			[]data.GetValue{
				param("name", 0, nil, nil),
				param("value", 1, data.NewNullValue(), nil),
				param("expire", 2, data.NewIntValue(0), nil),
				param("path", 3, data.NewStringValue("/"), nil),
				param("domain", 4, data.NewNullValue(), nil),
				param("secure", 5, data.NewNullValue(), nil),
				param("httpOnly", 6, data.NewBoolValue(true), nil),
				param("raw", 7, data.NewBoolValue(false), nil),
				param("sameSite", 8, data.NewStringValue("lax"), nil),
				param("partitioned", 9, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				variable("name", 0, nil), variable("value", 1, nil), variable("expire", 2, nil),
				variable("path", 3, nil), variable("domain", 4, nil), variable("secure", 5, nil),
				variable("httpOnly", 6, nil), variable("raw", 7, nil), variable("sameSite", 8, nil),
				variable("partitioned", 9, nil),
			},
			nil, cookieCreate).(*bagMethod)
		m.static = true
		return m, true
	}
	return nil, false
}

func cookieSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return bagClassValue(ctx)
}

func cookieConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := cookieSelf(ctx)
	set := func(i int, prop string) {
		if v, ok := ctx.GetIndexValue(i); ok && v != nil {
			_ = cv.SetProperty(prop, v)
		}
	}
	set(0, "name")
	set(1, "value")
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		_ = cv.SetProperty("expire", data.NewIntValue(cookieExpireToInt(v)))
	}
	set(3, "path")
	set(4, "domain")
	if v, ok := ctx.GetIndexValue(5); ok && v != nil {
		if _, isNull := v.(*data.NullValue); !isNull {
			_ = cv.SetProperty("secure", v)
		}
	}
	set(6, "httpOnly")
	set(7, "raw")
	set(8, "sameSite")
	set(9, "partitioned")
	return data.NewNullValue(), nil
}

func cookieCreate(ctx data.Context) (data.GetValue, data.Control) {
	cls := NewCookieClass()
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	construct := cls.GetConstruct()
	cctx := cv.CreateContext(construct.GetVariables())
	vars := construct.GetVariables()
	for i := 0; i < len(vars); i++ {
		if v, ok := ctx.GetIndexValue(i); ok && v != nil {
			_ = cctx.SetVariableValue(vars[i], v)
		}
	}
	_, ctl := construct.Call(cctx)
	if ctl != nil {
		return nil, ctl
	}
	return cv, nil
}

func cookieExpireToInt(v data.Value) int {
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	s := strings.TrimSpace(v.AsString())
	if s == "" {
		return 0
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	if t, err := time.Parse(time.RFC1123, s); err == nil {
		return int(t.Unix())
	}
	return 0
}

func cookiePropString(cv *data.ClassValue, name, def string) string {
	v, _ := cv.GetProperty(name)
	if v == nil {
		return def
	}
	if _, ok := v.(*data.NullValue); ok {
		return def
	}
	return v.AsString()
}

func cookiePropBool(cv *data.ClassValue, name string, def bool) bool {
	v, _ := cv.GetProperty(name)
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

func cookieGetName(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(cookiePropString(cookieSelf(ctx), "name", "")), nil
}
func cookieGetValue(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := cookieSelf(ctx).GetProperty("value")
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}
func cookieGetExpires(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := cookieSelf(ctx).GetProperty("expire")
	return data.NewIntValue(cookieExpireToInt(v)), nil
}
func cookieGetPath(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(cookiePropString(cookieSelf(ctx), "path", "/")), nil
}
func cookieGetDomain(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := cookieSelf(ctx).GetProperty("domain")
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}
func cookieIsSecure(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(cookiePropBool(cookieSelf(ctx), "secure", false)), nil
}
func cookieIsHTTPOnly(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(cookiePropBool(cookieSelf(ctx), "httpOnly", true)), nil
}
func cookieIsRaw(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(cookiePropBool(cookieSelf(ctx), "raw", false)), nil
}
func cookieGetSameSite(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := cookieSelf(ctx).GetProperty("sameSite")
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}

func cookieClone(cv *data.ClassValue, ctx data.Context) *data.ClassValue {
	clone := data.NewClassValue(NewCookieClass(), ctx.CreateBaseContext())
	if cv == nil {
		return clone
	}
	for _, name := range []string{"name", "value", "expire", "path", "domain", "secure", "httpOnly", "raw", "sameSite", "partitioned"} {
		if v, ctl := cv.GetProperty(name); ctl == nil && v != nil {
			_ = clone.SetProperty(name, v)
		}
	}
	return clone
}

func cookieWithValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := cookieSelf(ctx)
	clone := cookieClone(cv, ctx)
	if v, ok := ctx.GetIndexValue(0); ok {
		_ = clone.SetProperty("value", v)
	} else {
		_ = clone.SetProperty("value", data.NewNullValue())
	}
	return clone, nil
}

func cookieToString(ctx data.Context) (data.GetValue, data.Control) {
	cv := cookieSelf(ctx)
	name := cookiePropString(cv, "name", "")
	val := cookiePropString(cv, "value", "")
	raw := cookiePropBool(cv, "raw", false)
	if !raw {
		val = url.QueryEscape(val)
	}
	parts := []string{name + "=" + val}
	expire := cookieExpireToInt(func() data.Value { v, _ := cv.GetProperty("expire"); return v }())
	if expire > 0 {
		parts = append(parts, "expires="+time.Unix(int64(expire), 0).UTC().Format(time.RFC1123))
		parts = append(parts, "Max-Age="+strconv.Itoa(expire-int(time.Now().Unix())))
	}
	parts = append(parts, "path="+cookiePropString(cv, "path", "/"))
	if d := cookiePropString(cv, "domain", ""); d != "" {
		parts = append(parts, "domain="+d)
	}
	if cookiePropBool(cv, "secure", false) {
		parts = append(parts, "secure")
	}
	if cookiePropBool(cv, "httpOnly", true) {
		parts = append(parts, "httponly")
	}
	if ss := cookiePropString(cv, "sameSite", ""); ss != "" {
		parts = append(parts, "samesite="+ss)
	}
	return data.NewStringValue(strings.Join(parts, "; ")), nil
}

// NewCookieValueFromBag 从内部 BagCookie 构造 PHP Cookie 对象。
func NewCookieValueFromBag(ctx data.Context, c *BagCookie) *data.ClassValue {
	if c == nil {
		return nil
	}
	cv := data.NewClassValue(NewCookieClass(), ctx.CreateBaseContext())
	_ = cv.SetProperty("name", data.NewStringValue(c.Name))
	if c.Value != nil {
		_ = cv.SetProperty("value", data.NewStringValue(*c.Value))
	} else {
		_ = cv.SetProperty("value", data.NewNullValue())
	}
	_ = cv.SetProperty("expire", data.NewIntValue(int(c.Expire)))
	path := c.Path
	if path == "" {
		path = "/"
	}
	_ = cv.SetProperty("path", data.NewStringValue(path))
	if c.Domain != nil {
		_ = cv.SetProperty("domain", data.NewStringValue(*c.Domain))
	}
	_ = cv.SetProperty("secure", data.NewBoolValue(c.Secure))
	_ = cv.SetProperty("httpOnly", data.NewBoolValue(c.HTTPOnly))
	_ = cv.SetProperty("raw", data.NewBoolValue(c.Raw))
	if c.SameSite != nil {
		_ = cv.SetProperty("sameSite", data.NewStringValue(*c.SameSite))
	}
	return cv
}

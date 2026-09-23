package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
)

const fqnIlluminateRedirectResponse = "Illuminate\\Http\\RedirectResponse"

// IlluminateRedirectResponseClass 对齐 Illuminate\Http\RedirectResponse。
type IlluminateRedirectResponseClass struct {
	node.Node
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewIlluminateRedirectResponseClass() data.ClassStmt {
	c := &IlluminateRedirectResponseClass{
		properties: []data.Property{
			httpfoundation.PublicProp("original", data.NewNullValue()),
			httpfoundation.PublicProp("exception", data.NewNullValue()),
			node.NewProperty(nil, "request", "protected", false, data.NewNullValue()),
			node.NewProperty(nil, "session", "protected", false, data.NewNullValue()),
		},
	}
	c.methods, c.methodList = illuminateRedirectResponseMethods()
	return c
}

func (c *IlluminateRedirectResponseClass) GetName() string { return fqnIlluminateRedirectResponse }
func (c *IlluminateRedirectResponseClass) GetExtend() *string {
	parent := "Symfony\\Component\\HttpFoundation\\RedirectResponse"
	return &parent
}
func (c *IlluminateRedirectResponseClass) GetImplements() []string          { return nil }
func (c *IlluminateRedirectResponseClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *IlluminateRedirectResponseClass) GetPropertyList() []data.Property { return c.properties }
func (c *IlluminateRedirectResponseClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *IlluminateRedirectResponseClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *IlluminateRedirectResponseClass) GetMethods() []data.Method { return c.methodList }
func (c *IlluminateRedirectResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func illuminateRedirectResponseMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		httpfoundation.PubMethod("__construct",
			[]data.GetValue{
				httpfoundation.Param("url", 0, nil, nil),
				httpfoundation.Param("status", 1, data.NewIntValue(302), nil),
				httpfoundation.Param("headers", 2, data.NewArrayValue(nil), nil),
			},
			[]data.Variable{
				httpfoundation.Variable("url", 0, nil),
				httpfoundation.Variable("status", 1, nil),
				httpfoundation.Variable("headers", 2, nil),
			},
			nil, illuminateRedirectConstruct),
		httpfoundation.PubMethod("getTargetUrl", nil, nil, data.NewBaseType("string"), illuminateRedirectGetTarget),
		httpfoundation.PubMethod("setTargetUrl",
			[]data.GetValue{httpfoundation.Param("url", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("url", 0, nil)},
			nil, illuminateRedirectSetTarget),
		httpfoundation.PubMethod("status", nil, nil, data.NewBaseType("int"), illuminateResponseStatus),
		httpfoundation.PubMethod("with",
			[]data.GetValue{
				httpfoundation.Param("key", 0, nil, nil),
				httpfoundation.Param("value", 1, data.NewNullValue(), nil),
			},
			[]data.Variable{
				httpfoundation.Variable("key", 0, nil),
				httpfoundation.Variable("value", 1, nil),
			},
			nil, illuminateRedirectWith),
	}
	out := make(map[string]data.Method, len(list))
	for _, m := range list {
		out[m.GetName()] = m
	}
	return out, list
}

func illuminateRedirectConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	urlVal, _ := ctx.GetIndexValue(0)
	status := httpfoundation.IntParam(ctx, 1, 302)
	headersVal, _ := ctx.GetIndexValue(2)
	headers := httpfoundation.CreateResponseHeaders(ctx, headersVal)
	httpfoundation.ResponseSetProp(cv, "headers", headers)
	httpfoundation.ResponseSetProp(cv, "version", data.NewStringValue("1.0"))
	if ctl := httpfoundation.ApplyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	url := ""
	if urlVal != nil {
		url = urlVal.AsString()
	}
	httpfoundation.ResponseSetProp(cv, "targetUrl", data.NewStringValue(url))
	httpfoundation.RespHeaderSet(httpfoundation.ResponseHeaders(cv), "Location", []string{url}, true)
	httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(redirectHTML(url)))
	httpfoundation.ResponseSetProp(cv, "original", data.NewNullValue())
	httpfoundation.ResponseSetProp(cv, "exception", data.NewNullValue())
	return httpfoundation.ResponseSelf(ctx), nil
}

func redirectHTML(url string) string {
	return `<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8" />
        <meta http-equiv="refresh" content="0;url='` + url + `'" />
        <title>Redirecting to ` + url + `</title>
    </head>
    <body>
        Redirecting to <a href="` + url + `">` + url + `</a>.
    </body>
</html>`
}

func illuminateRedirectGetTarget(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := httpfoundation.ResponseClassValue(ctx).GetProperty("targetUrl")
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(v.AsString()), nil
}

func illuminateRedirectSetTarget(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	urlVal, _ := ctx.GetIndexValue(0)
	url := ""
	if urlVal != nil {
		url = urlVal.AsString()
	}
	httpfoundation.ResponseSetProp(cv, "targetUrl", data.NewStringValue(url))
	httpfoundation.RespHeaderSet(httpfoundation.ResponseHeaders(cv), "Location", []string{url}, true)
	httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(redirectHTML(url)))
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateRedirectWith(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	session, _ := cv.GetProperty("session")
	if sess, ok := session.(*data.ClassValue); ok && sess != nil {
		key, _ := ctx.GetIndexValue(0)
		val, _ := ctx.GetIndexValue(1)
		if arr, ok := key.(*data.ArrayValue); ok {
			for _, z := range arr.List {
				if z != nil && z.Name != "" {
					_, _ = httpfoundation.CallObjMethod(sess, "flash", data.NewStringValue(z.Name), z.Value)
				}
			}
		} else if key != nil {
			_, _ = httpfoundation.CallObjMethod(sess, "flash", key, val)
		}
	}
	return httpfoundation.ResponseSelf(ctx), nil
}

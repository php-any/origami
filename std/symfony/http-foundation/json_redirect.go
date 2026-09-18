package httpfoundation

import (
	"encoding/json"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const fqnJsonResponse = "Symfony\\Component\\HttpFoundation\\JsonResponse"
const fqnRedirectResponse = "Symfony\\Component\\HttpFoundation\\RedirectResponse"

// JsonResponseClass 实现 Symfony JsonResponse。
type JsonResponseClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
}

func NewJsonResponseClass() data.ClassStmt {
	c := &JsonResponseClass{}
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("data", 0, data.NewArrayValue(nil), nil),
				param("status", 1, data.NewIntValue(200), nil),
				param("headers", 2, data.NewArrayValue(nil), nil),
				param("json", 3, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				variable("data", 0, nil), variable("status", 1, nil),
				variable("headers", 2, nil), variable("json", 3, nil),
			},
			nil, jsonResponseConstruct),
		pubMethod("setData",
			[]data.GetValue{param("data", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{variable("data", 0, nil)},
			nil, jsonResponseSetData),
		pubMethod("getContent", nil, nil, nil, func(ctx data.Context) (data.GetValue, data.Control) {
			return data.NewStringValue(responseContent(responseClassValue(ctx))), nil
		}),
	}
	c.methods = make(map[string]data.Method, len(list))
	for _, m := range list {
		c.methods[m.GetName()] = m
	}
	c.methodList = list
	return c
}

func (c *JsonResponseClass) GetName() string { return fqnJsonResponse }
func (c *JsonResponseClass) GetExtend() *string {
	parent := fqnSymfonyResponse
	return &parent
}
func (c *JsonResponseClass) GetImplements() []string                  { return nil }
func (c *JsonResponseClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *JsonResponseClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "data", "protected", false, data.NewNullValue()),
		node.NewProperty(nil, "callback", "protected", false, data.NewNullValue()),
	}
}
func (c *JsonResponseClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *JsonResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *JsonResponseClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *JsonResponseClass) GetMethods() []data.Method { return c.methodList }

func jsonResponseConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	dataVal, _ := ctx.GetIndexValue(0)
	status := intParam(ctx, 1, 200)
	headersVal, _ := ctx.GetIndexValue(2)
	headers := createResponseHeaders(ctx, headersVal)
	_ = cv.SetProperty("headers", headers)
	_ = cv.SetProperty("version", data.NewStringValue("1.0"))
	if ctl := applyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	respHeaderSet(headers, "Content-Type", []string{"application/json"}, true)
	isJSON := false
	if v, ok := ctx.GetIndexValue(3); ok && v != nil {
		if b, ok := v.(data.AsBool); ok {
			isJSON, _ = b.AsBool()
		}
	}
	if isJSON && dataVal != nil {
		_ = cv.SetProperty("content", data.NewStringValue(dataVal.AsString()))
		_ = cv.SetProperty("data", dataVal)
		return data.NewNullValue(), nil
	}
	_, ctl := jsonResponseSetDataWith(cv, dataVal)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), nil
}

func jsonResponseSetData(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	dataVal, _ := ctx.GetIndexValue(0)
	return jsonResponseSetDataWith(cv, dataVal)
}

func jsonResponseSetDataWith(cv *data.ClassValue, dataVal data.Value) (data.GetValue, data.Control) {
	_ = cv.SetProperty("data", dataVal)
	encoded, ok := encodeJSONValue(dataVal)
	if !ok {
		encoded = "null"
	}
	_ = cv.SetProperty("content", data.NewStringValue(encoded))
	return cv, nil
}

func encodeJSONValue(v data.Value) (string, bool) {
	goVal := phpToGoJSON(v)
	b, err := json.Marshal(goVal)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func phpToGoJSON(v data.Value) any {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case *data.NullValue:
		return nil
	case *data.BoolValue:
		return t.Value
	case *data.IntValue:
		return t.Value
	case *data.FloatValue:
		return t.Value
	case *data.StringValue:
		return t.Value
	case *data.ObjectValue:
		// Origami 关联数组是 ObjectValue
		out := map[string]any{}
		t.RangeProperties(func(key string, value data.Value) bool {
			out[key] = phpToGoJSON(value)
			return true
		})
		return out
	case *data.ArrayValue:
		isList := true
		for i, z := range t.List {
			if z == nil {
				continue
			}
			if z.Name != "" {
				if n, ok := data.ParseIntArrayKeyName(z.Name); !ok || n != i {
					isList = false
					break
				}
			}
		}
		if isList {
			out := make([]any, 0, len(t.List))
			for _, z := range t.List {
				if z == nil {
					out = append(out, nil)
				} else {
					out = append(out, phpToGoJSON(z.Value))
				}
			}
			return out
		}
		out := map[string]any{}
		for i, z := range t.List {
			if z == nil {
				continue
			}
			k := z.Name
			if k == "" {
				k = data.IntArrayKeyName(i)
			}
			out[k] = phpToGoJSON(z.Value)
		}
		return out
	default:
		return v.AsString()
	}
}

// RedirectResponseClass 实现 Symfony RedirectResponse。
type RedirectResponseClass struct {
	node.Node
	methods    map[string]data.Method
	methodList []data.Method
}

func NewRedirectResponseClass() data.ClassStmt {
	c := &RedirectResponseClass{}
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("url", 0, nil, nil),
				param("status", 1, data.NewIntValue(302), nil),
				param("headers", 2, data.NewArrayValue(nil), nil),
			},
			[]data.Variable{
				variable("url", 0, nil), variable("status", 1, nil), variable("headers", 2, nil),
			},
			nil, redirectConstruct),
		pubMethod("getTargetUrl", nil, nil, data.NewBaseType("string"), redirectGetTarget),
		pubMethod("setTargetUrl",
			[]data.GetValue{param("url", 0, nil, nil)},
			[]data.Variable{variable("url", 0, nil)},
			nil, redirectSetTarget),
	}
	c.methods = make(map[string]data.Method, len(list))
	for _, m := range list {
		c.methods[m.GetName()] = m
	}
	c.methodList = list
	return c
}

func (c *RedirectResponseClass) GetName() string { return fqnRedirectResponse }
func (c *RedirectResponseClass) GetExtend() *string {
	parent := fqnSymfonyResponse
	return &parent
}
func (c *RedirectResponseClass) GetImplements() []string                  { return nil }
func (c *RedirectResponseClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *RedirectResponseClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "targetUrl", "protected", false, data.NewStringValue("")),
	}
}
func (c *RedirectResponseClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *RedirectResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *RedirectResponseClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *RedirectResponseClass) GetMethods() []data.Method { return c.methodList }

func redirectConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	urlVal, _ := ctx.GetIndexValue(0)
	status := intParam(ctx, 1, 302)
	headersVal, _ := ctx.GetIndexValue(2)
	headers := createResponseHeaders(ctx, headersVal)
	_ = cv.SetProperty("headers", headers)
	_ = cv.SetProperty("version", data.NewStringValue("1.0"))
	if ctl := applyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	url := ""
	if urlVal != nil {
		url = urlVal.AsString()
	}
	_ = cv.SetProperty("targetUrl", data.NewStringValue(url))
	respHeaderSet(headers, "Location", []string{url}, true)
	_ = cv.SetProperty("content", data.NewStringValue(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8" />
        <meta http-equiv="refresh" content="0;url='`+url+`'" />
        <title>Redirecting to `+url+`</title>
    </head>
    <body>
        Redirecting to <a href="`+url+`">`+url+`</a>.
    </body>
</html>`))
	return data.NewNullValue(), nil
}

func redirectGetTarget(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := responseClassValue(ctx).GetProperty("targetUrl")
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(v.AsString()), nil
}

func redirectSetTarget(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	urlVal, _ := ctx.GetIndexValue(0)
	url := ""
	if urlVal != nil {
		url = urlVal.AsString()
	}
	_ = cv.SetProperty("targetUrl", data.NewStringValue(url))
	respHeaderSet(responseHeaders(cv), "Location", []string{url}, true)
	return cv, nil
}

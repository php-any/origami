package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
)

const fqnIlluminateJsonResponse = "Illuminate\\Http\\JsonResponse"

// IlluminateJsonResponseClass 对齐 Illuminate\Http\JsonResponse（父类 Symfony JsonResponse）。
type IlluminateJsonResponseClass struct {
	node.Node
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

func NewIlluminateJsonResponseClass() data.ClassStmt {
	c := &IlluminateJsonResponseClass{
		properties: []data.Property{
			httpfoundation.PublicProp("original", data.NewNullValue()),
			httpfoundation.PublicProp("exception", data.NewNullValue()),
			node.NewProperty(nil, "encodingOptions", "public", false, data.NewIntValue(0)),
		},
	}
	c.methods, c.methodList = illuminateJsonResponseMethods()
	return c
}

func (c *IlluminateJsonResponseClass) GetName() string { return fqnIlluminateJsonResponse }
func (c *IlluminateJsonResponseClass) GetExtend() *string {
	parent := "Symfony\\Component\\HttpFoundation\\JsonResponse"
	return &parent
}
func (c *IlluminateJsonResponseClass) GetImplements() []string          { return nil }
func (c *IlluminateJsonResponseClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *IlluminateJsonResponseClass) GetPropertyList() []data.Property { return c.properties }
func (c *IlluminateJsonResponseClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *IlluminateJsonResponseClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *IlluminateJsonResponseClass) GetMethods() []data.Method { return c.methodList }
func (c *IlluminateJsonResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func illuminateJsonResponseMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		httpfoundation.PubMethod("__construct",
			[]data.GetValue{
				httpfoundation.Param("data", 0, data.NewNullValue(), nil),
				httpfoundation.Param("status", 1, data.NewIntValue(200), nil),
				httpfoundation.Param("headers", 2, data.NewArrayValue(nil), nil),
				httpfoundation.Param("options", 3, data.NewIntValue(0), nil),
				httpfoundation.Param("json", 4, data.NewBoolValue(false), nil),
			},
			[]data.Variable{
				httpfoundation.Variable("data", 0, nil),
				httpfoundation.Variable("status", 1, nil),
				httpfoundation.Variable("headers", 2, nil),
				httpfoundation.Variable("options", 3, nil),
				httpfoundation.Variable("json", 4, nil),
			},
			nil, illuminateJsonResponseConstruct),
		httpfoundation.PubMethod("setData",
			[]data.GetValue{httpfoundation.Param("data", 0, data.NewArrayValue(nil), nil)},
			[]data.Variable{httpfoundation.Variable("data", 0, nil)},
			nil, illuminateJsonResponseSetData),
		httpfoundation.PubMethod("getData",
			[]data.GetValue{
				httpfoundation.Param("assoc", 0, data.NewBoolValue(false), nil),
				httpfoundation.Param("depth", 1, data.NewIntValue(512), nil),
			},
			[]data.Variable{
				httpfoundation.Variable("assoc", 0, nil),
				httpfoundation.Variable("depth", 1, nil),
			},
			nil, illuminateJsonResponseGetData),
		httpfoundation.PubMethod("getContent", nil, nil, nil, func(ctx data.Context) (data.GetValue, data.Control) {
			return data.NewStringValue(httpfoundation.ResponseContent(httpfoundation.ResponseClassValue(ctx))), nil
		}),
		httpfoundation.PubMethod("status", nil, nil, data.NewBaseType("int"), illuminateResponseStatus),
		httpfoundation.PubMethod("header",
			[]data.GetValue{
				httpfoundation.Param("key", 0, nil, nil),
				httpfoundation.Param("values", 1, nil, nil),
				httpfoundation.Param("replace", 2, data.NewBoolValue(true), nil),
			},
			[]data.Variable{
				httpfoundation.Variable("key", 0, nil),
				httpfoundation.Variable("values", 1, nil),
				httpfoundation.Variable("replace", 2, nil),
			},
			nil, illuminateResponseHeader),
	}
	out := make(map[string]data.Method, len(list))
	for _, m := range list {
		out[m.GetName()] = m
	}
	return out, list
}

func illuminateJsonResponseConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	options := httpfoundation.IntParam(ctx, 3, 0)
	httpfoundation.ResponseSetProp(cv, "encodingOptions", data.NewIntValue(options))
	dataVal, _ := ctx.GetIndexValue(0)
	status := httpfoundation.IntParam(ctx, 1, 200)
	headersVal, _ := ctx.GetIndexValue(2)
	isJSON := httpfoundation.BoolParam(ctx, 4, false)

	headers := httpfoundation.CreateResponseHeaders(ctx, headersVal)
	httpfoundation.ResponseSetProp(cv, "headers", headers)
	httpfoundation.ResponseSetProp(cv, "version", data.NewStringValue("1.0"))
	if ctl := httpfoundation.ApplyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	httpfoundation.RespHeaderSet(httpfoundation.ResponseHeaders(cv), "Content-Type", []string{"application/json"}, true)

	if isJSON && dataVal != nil {
		httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(dataVal.AsString()))
		httpfoundation.ResponseSetProp(cv, "data", dataVal)
		httpfoundation.ResponseSetProp(cv, "original", dataVal)
		return httpfoundation.ResponseSelf(ctx), nil
	}
	httpfoundation.ResponseSetProp(cv, "original", dataVal)
	_, ctl := illuminateJsonResponseSetData(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateJsonResponseSetData(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	dataVal, _ := ctx.GetIndexValue(0)
	httpfoundation.ResponseSetProp(cv, "original", dataVal)
	encoded, ok, ctl := encodeJsonResponseData(ctx, cv, dataVal)
	if ctl != nil {
		return nil, ctl
	}
	if !ok {
		return nil, httpfoundation.ThrowNamed("InvalidArgumentException", "Unable to encode JSON response data.")
	}
	httpfoundation.ResponseSetProp(cv, "data", data.NewStringValue(encoded))
	httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(encoded))
	return httpfoundation.ResponseSelf(ctx), nil
}

func encodeJsonResponseData(ctx data.Context, cv *data.ClassValue, dataVal data.Value) (string, bool, data.Control) {
	if dataVal == nil {
		return "null", true, nil
	}
	if obj, ok := dataVal.(*data.ClassValue); ok {
		if implementsOrExtends(obj, "Illuminate\\Contracts\\Support\\Jsonable") {
			ret, ctl := httpfoundation.CallObjMethod(obj, "toJson")
			if ctl != nil || ret == nil {
				return "", false, ctl
			}
			return ret.(data.Value).AsString(), true, nil
		}
		if implementsOrExtends(obj, "Illuminate\\Contracts\\Support\\Arrayable") {
			ret, ctl := httpfoundation.CallObjMethod(obj, "toArray")
			if ctl != nil || ret == nil {
				return "", false, ctl
			}
			s, ok := jsonEncodeValue(ctx, ret.(data.Value))
			return s, ok, nil
		}
	}
	s, ok := jsonEncodeValue(ctx, dataVal)
	return s, ok, nil
}

func illuminateJsonResponseGetData(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	if orig, _ := cv.GetProperty("original"); orig != nil && !httpfoundation.IsNull(orig) {
		return orig, nil
	}
	raw, _ := cv.GetProperty("content")
	if raw == nil {
		return data.NewNullValue(), nil
	}
	return raw, nil
}

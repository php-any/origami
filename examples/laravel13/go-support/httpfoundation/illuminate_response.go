package httpfoundation

import (
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const fqnIlluminateResponse = "Illuminate\\Http\\Response"

// IlluminateResponseClass 实现 Illuminate\Http\Response（含 ResponseTrait）。
type IlluminateResponseClass struct {
	node.Node
	properties []data.Property
	methods    map[string]data.Method
	methodList []data.Method
}

// NewIlluminateResponseClass 导出 Illuminate Response 构造器。
func NewIlluminateResponseClass() data.ClassStmt {
	c := &IlluminateResponseClass{
		properties: []data.Property{
			publicProp("original", data.NewNullValue()),
			publicProp("exception", data.NewNullValue()),
		},
	}
	c.methods, c.methodList = illuminateResponseMethods()
	return c
}

func (c *IlluminateResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *IlluminateResponseClass) GetName() string { return fqnIlluminateResponse }
func (c *IlluminateResponseClass) GetExtend() *string {
	parent := fqnSymfonyResponse
	return &parent
}
func (c *IlluminateResponseClass) GetImplements() []string          { return nil }
func (c *IlluminateResponseClass) GetConstruct() data.Method        { return c.methods["__construct"] }
func (c *IlluminateResponseClass) GetPropertyList() []data.Property { return c.properties }
func (c *IlluminateResponseClass) GetProperty(name string) (data.Property, bool) {
	for _, p := range c.properties {
		if p.GetName() == name {
			return p, true
		}
	}
	return nil, false
}
func (c *IlluminateResponseClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[name]
	return m, ok
}
func (c *IlluminateResponseClass) GetMethods() []data.Method { return c.methodList }

func illuminateResponseMethods() (map[string]data.Method, []data.Method) {
	list := []data.Method{
		pubMethod("__construct",
			[]data.GetValue{
				param("content", 0, data.NewStringValue(""), nil),
				param("status", 1, data.NewIntValue(200), nil),
				param("headers", 2, data.NewArrayValue(nil), nil),
			},
			[]data.Variable{
				variable("content", 0, nil),
				variable("status", 1, nil),
				variable("headers", 2, nil),
			},
			nil, illuminateResponseConstruct),
		pubMethod("getContent", nil, nil, nil, illuminateResponseGetContent),
		pubMethod("setContent",
			[]data.GetValue{param("content", 0, nil, nil)},
			[]data.Variable{variable("content", 0, nil)},
			nil, illuminateResponseSetContent),
		// ResponseTrait
		pubMethod("status", nil, nil, data.NewBaseType("int"), illuminateResponseStatus),
		pubMethod("statusText", nil, nil, data.NewBaseType("string"), illuminateResponseStatusText),
		pubMethod("content", nil, nil, data.NewBaseType("string"), illuminateResponseContent),
		pubMethod("getOriginalContent", nil, nil, nil, illuminateResponseGetOriginalContent),
		pubMethod("header",
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
			nil, illuminateResponseHeader),
		pubMethod("withHeaders",
			[]data.GetValue{param("headers", 0, nil, nil)},
			[]data.Variable{variable("headers", 0, nil)},
			nil, illuminateResponseWithHeaders),
		pubMethod("withoutHeader",
			[]data.GetValue{param("key", 0, nil, nil)},
			[]data.Variable{variable("key", 0, nil)},
			nil, illuminateResponseWithoutHeader),
		pubMethod("cookie",
			[]data.GetValue{param("cookie", 0, nil, nil)},
			[]data.Variable{variable("cookie", 0, nil)},
			nil, illuminateResponseWithCookie),
		pubMethod("withCookie",
			[]data.GetValue{param("cookie", 0, nil, nil)},
			[]data.Variable{variable("cookie", 0, nil)},
			nil, illuminateResponseWithCookie),
		pubMethod("withCookies",
			[]data.GetValue{param("cookies", 0, nil, nil)},
			[]data.Variable{variable("cookies", 0, nil)},
			nil, illuminateResponseWithCookies),
		pubMethod("withoutCookie",
			[]data.GetValue{
				param("cookie", 0, nil, nil),
				param("path", 1, data.NewNullValue(), nil),
				param("domain", 2, data.NewNullValue(), nil),
			},
			[]data.Variable{
				variable("cookie", 0, nil),
				variable("path", 1, nil),
				variable("domain", 2, nil),
			},
			nil, illuminateResponseWithoutCookie),
		pubMethod("getCallback", nil, nil, nil, illuminateResponseGetCallback),
		pubMethod("withException",
			[]data.GetValue{param("e", 0, nil, nil)},
			[]data.Variable{variable("e", 0, nil)},
			nil, illuminateResponseWithException),
		pubMethod("throwResponse", nil, nil, nil, illuminateResponseThrowResponse),
		// protected helpers used by setContent path (also callable)
		pubMethod("shouldBeJson",
			[]data.GetValue{param("content", 0, nil, nil)},
			[]data.Variable{variable("content", 0, nil)},
			data.NewBaseType("bool"), illuminateResponseShouldBeJson),
		pubMethod("morphToJson",
			[]data.GetValue{param("content", 0, nil, nil)},
			[]data.Variable{variable("content", 0, nil)},
			nil, illuminateResponseMorphToJson),
	}
	// shouldBeJson / morphToJson 在 PHP 为 protected
	for _, name := range []string{"shouldBeJson", "morphToJson"} {
		if m, ok := listMethodByName(list, name); ok {
			m.modifier = data.ModifierProtected
		}
	}
	out := make(map[string]data.Method, len(list))
	for _, m := range list {
		out[m.GetName()] = m
	}
	return out, list
}

func listMethodByName(list []data.Method, name string) (*bagMethod, bool) {
	for _, m := range list {
		if m.GetName() == name {
			if bm, ok := m.(*bagMethod); ok {
				return bm, true
			}
		}
	}
	return nil, false
}

func illuminateResponseConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Illuminate\\Http\\Response::__construct: invalid context"))
	}
	content, _ := ctx.GetIndexValue(0)
	status := intParam(ctx, 1, 200)
	headersVal, _ := ctx.GetIndexValue(2)

	headers := createResponseHeaders(ctx, headersVal)
	responseSetProp(cv, "headers", headers)
	responseSetProp(cv, "version", data.NewStringValue("1.0"))
	if ctl := applyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	responseSetProp(cv, "content", data.NewStringValue(""))
	responseSetProp(cv, "original", data.NewNullValue())
	responseSetProp(cv, "exception", data.NewNullValue())

	// 走 Illuminate setContent（支持 JSON / Renderable）
	setCtx := cv.CreateContext([]data.Variable{variable("content", 0, nil)})
	if content == nil {
		content = data.NewStringValue("")
	}
	setCtx.SetVariableValue(variable("content", 0, nil), content)
	return illuminateResponseSetContent(setCtx)
}

func illuminateResponseGetContent(ctx data.Context) (data.GetValue, data.Control) {
	s := responseContent(responseClassValue(ctx))
	return data.NewStringValue(s), nil
}

func illuminateResponseSetContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	content, _ := ctx.GetIndexValue(0)
	if content == nil {
		content = data.NewNullValue()
	}
	responseSetProp(cv, "original", content)

	if shouldBeJSON(content) {
		respHeaderSet(responseHeaders(cv), "Content-Type", []string{"application/json"}, true)
		encoded, ok := morphToJSON(ctx, content)
		if !ok {
			return nil, throwNamed("InvalidArgumentException", "json_encode error")
		}
		responseSetProp(cv, "content", data.NewStringValue(encoded))
		return responseSelf(ctx), nil
	}

	if rendered, ok := tryRender(ctx, content); ok {
		responseSetProp(cv, "content", data.NewStringValue(rendered))
		return responseSelf(ctx), nil
	}

	if _, isNull := content.(*data.NullValue); isNull {
		responseSetProp(cv, "content", data.NewStringValue(""))
	} else {
		responseSetProp(cv, "content", data.NewStringValue(content.AsString()))
	}
	return responseSelf(ctx), nil
}

func shouldBeJSON(content data.Value) bool {
	if content == nil {
		return false
	}
	if _, ok := content.(*data.ArrayValue); ok {
		return true
	}
	if _, ok := content.(*data.ObjectValue); ok {
		return true
	}
	cv, ok := content.(*data.ClassValue)
	if !ok {
		return false
	}
	name := cv.Class.GetName()
	switch name {
	case "ArrayObject", "JsonSerializable":
		return true
	}
	for _, iface := range cv.Class.GetImplements() {
		switch iface {
		case "Illuminate\\Contracts\\Support\\Arrayable",
			"Illuminate\\Contracts\\Support\\Jsonable",
			"JsonSerializable":
			return true
		}
	}
	// instanceof via inheritance chain names is best-effort
	if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Arrayable") ||
		implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Jsonable") ||
		implementsOrExtends(cv, "JsonSerializable") ||
		implementsOrExtends(cv, "ArrayObject") {
		return true
	}
	return false
}

func implementsOrExtends(cv *data.ClassValue, name string) bool {
	if cv == nil || cv.Class == nil {
		return false
	}
	if cv.Class.GetName() == name {
		return true
	}
	for _, iface := range cv.Class.GetImplements() {
		if iface == name {
			return true
		}
	}
	vm := cv.GetVM()
	last := cv.Class
	for last != nil && last.GetExtend() != nil {
		ext := last.GetExtend()
		next, ok := vm.GetClass(*ext)
		if !ok || next == nil {
			break
		}
		if next.GetName() == name {
			return true
		}
		for _, iface := range next.GetImplements() {
			if iface == name {
				return true
			}
		}
		last = next
	}
	return false
}

func morphToJSON(ctx data.Context, content data.Value) (string, bool) {
	if content == nil {
		return "null", true
	}
	if cv, ok := content.(*data.ClassValue); ok {
		if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Jsonable") {
			ret, ctl := callObjMethod(cv, "toJson")
			if ctl != nil || ret == nil {
				return "", false
			}
			return ret.(data.Value).AsString(), true
		}
		if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Arrayable") {
			ret, ctl := callObjMethod(cv, "toArray")
			if ctl != nil || ret == nil {
				return "", false
			}
			return jsonEncodeValue(ret.(data.Value))
		}
	}
	return jsonEncodeValue(content)
}

func jsonEncodeValue(v data.Value) (string, bool) {
	goVal := phpValueToGo(v)
	b, err := json.Marshal(goVal)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func phpValueToGo(v data.Value) any {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case *data.NullValue:
		return nil
	case *data.BoolValue:
		b, _ := t.AsBool()
		return b
	case *data.IntValue:
		n, _ := t.AsInt()
		return n
	case *data.FloatValue:
		f, _ := t.AsFloat()
		return f
	case *data.StringValue:
		return t.Value
	case *data.ArrayValue:
		// 尝试关联数组
		isAssoc := false
		for i, z := range t.List {
			if z == nil {
				continue
			}
			if z.Name != "" && z.Name != data.IntArrayKeyName(i) {
				isAssoc = true
				break
			}
		}
		if isAssoc {
			out := map[string]any{}
			for i, z := range t.List {
				if z == nil {
					continue
				}
				key := z.Name
				if key == "" {
					key = fmt.Sprintf("%d", i)
				}
				out[key] = phpValueToGo(z.Value)
			}
			return out
		}
		out := make([]any, 0, len(t.List))
		for _, z := range t.List {
			if z == nil {
				out = append(out, nil)
			} else {
				out = append(out, phpValueToGo(z.Value))
			}
		}
		return out
	case *data.ClassValue:
		return t.AsString()
	default:
		return v.AsString()
	}
}

func tryRender(ctx data.Context, content data.Value) (string, bool) {
	cv, ok := content.(*data.ClassValue)
	if !ok {
		return "", false
	}
	if !implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Renderable") {
		// 也尝试直接有 render 方法
		if _, has := cv.GetMethod("render"); !has {
			return "", false
		}
	}
	ret, ctl := callObjMethod(cv, "render")
	if ctl != nil || ret == nil {
		return "", false
	}
	return ret.(data.Value).AsString(), true
}

func illuminateResponseShouldBeJson(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return data.NewBoolValue(shouldBeJSON(v)), nil
}

func illuminateResponseMorphToJson(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	s, ok := morphToJSON(ctx, v)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(s), nil
}

func illuminateResponseStatus(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(responseStatusCode(responseClassValue(ctx))), nil
}

func illuminateResponseStatusText(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(responseStatusText(responseClassValue(ctx))), nil
}

func illuminateResponseContent(ctx data.Context) (data.GetValue, data.Control) {
	return illuminateResponseGetContent(ctx)
}

func illuminateResponseGetOriginalContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := cv.GetProperty("original")
	if nested, ok := v.(*data.ClassValue); ok {
		if nested.Class != nil && nested.Class.GetName() == fqnIlluminateResponse {
			ret, ctl := callObjMethod(nested, "getOriginalContent")
			if ctl != nil {
				return nil, ctl
			}
			return ret, nil
		}
	}
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}

func illuminateResponseHeader(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	key, _ := ctx.GetIndexValue(0)
	values, _ := ctx.GetIndexValue(1)
	replace := boolParam(ctx, 2, true)
	if key == nil {
		return responseSelf(ctx), nil
	}
	var vals []string
	if arr, ok := values.(*data.ArrayValue); ok {
		for _, z := range arr.List {
			if z != nil && z.Value != nil {
				vals = append(vals, z.Value.AsString())
			}
		}
	} else if values != nil {
		vals = []string{values.AsString()}
	}
	respHeaderSet(responseHeaders(cv), key.AsString(), vals, replace)
	return responseSelf(ctx), nil
}

func illuminateResponseWithHeaders(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	headersVal, _ := ctx.GetIndexValue(0)
	var m map[string]data.Value
	if hv, ok := headersVal.(*data.ClassValue); ok {
		if HeaderBagFrom(hv) != nil {
			all := GetHeaderBagAll(hv)
			m = make(map[string]data.Value, len(all))
			for k, vs := range all {
				if len(vs) == 1 {
					m[k] = data.NewStringValue(vs[0])
				} else {
					arr := make([]data.Value, len(vs))
					for i, s := range vs {
						arr[i] = data.NewStringValue(s)
					}
					m[k] = data.NewArrayValue(arr)
				}
			}
		}
	}
	if m == nil {
		var err error
		m, err = valueToAssocMap(headersVal)
		if err != nil {
			m = map[string]data.Value{}
		}
	}
	for k, v := range m {
		var vals []string
		if arr, ok := v.(*data.ArrayValue); ok {
			for _, z := range arr.List {
				if z != nil && z.Value != nil {
					vals = append(vals, z.Value.AsString())
				}
			}
		} else if v != nil {
			vals = []string{v.AsString()}
		}
		respHeaderSet(responseHeaders(cv), k, vals, true)
	}
	return responseSelf(ctx), nil
}

func illuminateResponseWithoutHeader(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	key, _ := ctx.GetIndexValue(0)
	if arr, ok := key.(*data.ArrayValue); ok {
		for _, z := range arr.List {
			if z != nil && z.Value != nil {
				respHeaderRemove(responseHeaders(cv), z.Value.AsString())
			}
		}
	} else if key != nil {
		respHeaderRemove(responseHeaders(cv), key.AsString())
	}
	return responseSelf(ctx), nil
}

func illuminateResponseWithCookie(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	cookie, _ := ctx.GetIndexValue(0)
	headers := responseHeaders(cv)
	if rh := ResponseHeaderBagFrom(headers); rh != nil && cookie != nil {
		// 优先走 ResponseHeaderBag::setCookie
		_, ctl := callObjMethod(headers, "setCookie", cookie)
		if ctl != nil {
			return nil, ctl
		}
	}
	return responseSelf(ctx), nil
}

func illuminateResponseWithCookies(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	cookies, _ := ctx.GetIndexValue(0)
	arr, ok := cookies.(*data.ArrayValue)
	if !ok {
		return responseSelf(ctx), nil
	}
	headers := responseHeaders(cv)
	for _, z := range arr.List {
		if z == nil || z.Value == nil {
			continue
		}
		_, ctl := callObjMethod(headers, "setCookie", z.Value)
		if ctl != nil {
			return nil, ctl
		}
	}
	return responseSelf(ctx), nil
}

func illuminateResponseWithoutCookie(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	cookie, _ := ctx.GetIndexValue(0)
	headers := responseHeaders(cv)
	if cookie == nil {
		return responseSelf(ctx), nil
	}
	// 字符串 cookie 名：构造过期 cookie（简化）
	if _, ok := cookie.(*data.StringValue); ok {
		name := cookie.AsString()
		path := "/"
		if p, present, ok := optionalStringParam(ctx, 1); present && ok && p != "" {
			path = p
		}
		domain := ""
		if d, present, ok := optionalStringParam(ctx, 2); present && ok {
			domain = d
		}
		val := ""
		bc := &BagCookie{
			Name:     name,
			Value:    &val,
			Expire:   -2628000,
			Path:     path,
			HTTPOnly: true,
		}
		if domain != "" {
			bc.Domain = &domain
		}
		if rh := ResponseHeaderBagFrom(headers); rh != nil {
			setBagCookie(rh, bc)
			return responseSelf(ctx), nil
		}
	}
	_, ctl := callObjMethod(headers, "setCookie", cookie)
	if ctl != nil {
		return nil, ctl
	}
	return responseSelf(ctx), nil
}

func illuminateResponseGetCallback(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	v, _ := cv.GetProperty("callback")
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}

func illuminateResponseWithException(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	e, _ := ctx.GetIndexValue(0)
	if e == nil {
		e = data.NewNullValue()
	}
	responseSetProp(cv, "exception", e)
	return responseSelf(ctx), nil
}

func illuminateResponseThrowResponse(ctx data.Context) (data.GetValue, data.Control) {
	cv := responseClassValue(ctx)
	// 尽量构造 HttpResponseException($this)
	vm := ctx.GetVM()
	if vm != nil {
		cls, ctl := vm.GetOrLoadClass("Illuminate\\Http\\Exceptions\\HttpResponseException")
		if ctl == nil && cls != nil {
			obj := data.NewClassValue(cls, ctx.CreateBaseContext())
			if ctor := cls.GetConstruct(); ctor != nil {
				cctx := obj.CreateContext(ctor.GetVariables())
				vars := ctor.GetVariables()
				if len(vars) > 0 {
					cctx.SetVariableValue(vars[0], cv)
				}
				_, _ = ctor.Call(cctx)
			}
			return nil, data.NewErrorThrow(nil, fmt.Errorf("HttpResponseException"))
		}
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Illuminate\\Http\\Exceptions\\HttpResponseException"))
}

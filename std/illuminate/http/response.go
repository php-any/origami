package http

import (
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
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
			httpfoundation.PublicProp("original", data.NewNullValue()),
			httpfoundation.PublicProp("exception", data.NewNullValue()),
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
	parent := httpfoundation.FqnSymfonyResponse
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
		httpfoundation.PubMethod("__construct",
			[]data.GetValue{
				httpfoundation.Param("content", 0, data.NewStringValue(""), nil),
				httpfoundation.Param("status", 1, data.NewIntValue(200), nil),
				httpfoundation.Param("headers", 2, data.NewArrayValue(nil), nil),
			},
			[]data.Variable{
				httpfoundation.Variable("content", 0, nil),
				httpfoundation.Variable("status", 1, nil),
				httpfoundation.Variable("headers", 2, nil),
			},
			nil, illuminateResponseConstruct),
		httpfoundation.PubMethod("getContent", nil, nil, nil, illuminateResponseGetContent),
		httpfoundation.PubMethod("setContent",
			[]data.GetValue{httpfoundation.Param("content", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("content", 0, nil)},
			nil, illuminateResponseSetContent),
		// ResponseTrait
		httpfoundation.PubMethod("status", nil, nil, data.NewBaseType("int"), illuminateResponseStatus),
		httpfoundation.PubMethod("statusText", nil, nil, data.NewBaseType("string"), illuminateResponseStatusText),
		httpfoundation.PubMethod("content", nil, nil, data.NewBaseType("string"), illuminateResponseContent),
		httpfoundation.PubMethod("getOriginalContent", nil, nil, nil, illuminateResponseGetOriginalContent),
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
		httpfoundation.PubMethod("withHeaders",
			[]data.GetValue{httpfoundation.Param("headers", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("headers", 0, nil)},
			nil, illuminateResponseWithHeaders),
		httpfoundation.PubMethod("withoutHeader",
			[]data.GetValue{httpfoundation.Param("key", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("key", 0, nil)},
			nil, illuminateResponseWithoutHeader),
		httpfoundation.PubMethod("cookie",
			[]data.GetValue{httpfoundation.Param("cookie", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("cookie", 0, nil)},
			nil, illuminateResponseWithCookie),
		httpfoundation.PubMethod("withCookie",
			[]data.GetValue{httpfoundation.Param("cookie", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("cookie", 0, nil)},
			nil, illuminateResponseWithCookie),
		httpfoundation.PubMethod("withCookies",
			[]data.GetValue{httpfoundation.Param("cookies", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("cookies", 0, nil)},
			nil, illuminateResponseWithCookies),
		httpfoundation.PubMethod("withoutCookie",
			[]data.GetValue{
				httpfoundation.Param("cookie", 0, nil, nil),
				httpfoundation.Param("path", 1, data.NewNullValue(), nil),
				httpfoundation.Param("domain", 2, data.NewNullValue(), nil),
			},
			[]data.Variable{
				httpfoundation.Variable("cookie", 0, nil),
				httpfoundation.Variable("path", 1, nil),
				httpfoundation.Variable("domain", 2, nil),
			},
			nil, illuminateResponseWithoutCookie),
		httpfoundation.PubMethod("getCallback", nil, nil, nil, illuminateResponseGetCallback),
		httpfoundation.PubMethod("withException",
			[]data.GetValue{httpfoundation.Param("e", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("e", 0, nil)},
			nil, illuminateResponseWithException),
		httpfoundation.PubMethod("throwResponse", nil, nil, nil, illuminateResponseThrowResponse),
		// protected helpers used by setContent path (also callable)
		httpfoundation.PubMethod("shouldBeJson",
			[]data.GetValue{httpfoundation.Param("content", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("content", 0, nil)},
			data.NewBaseType("bool"), illuminateResponseShouldBeJson),
		httpfoundation.PubMethod("morphToJson",
			[]data.GetValue{httpfoundation.Param("content", 0, nil, nil)},
			[]data.Variable{httpfoundation.Variable("content", 0, nil)},
			nil, illuminateResponseMorphToJson),
	}
	// shouldBeJson / morphToJson 在 PHP 为 protected
	for _, name := range []string{"shouldBeJson", "morphToJson"} {
		if m, ok := listMethodByName(list, name); ok {
			httpfoundation.SetMethodModifier(m, data.ModifierProtected)
		}
	}
	out := make(map[string]data.Method, len(list))
	for _, m := range list {
		out[m.GetName()] = m
	}
	return out, list
}

func listMethodByName(list []data.Method, name string) (data.Method, bool) {
	for _, m := range list {
		if m.GetName() == name {
			return m, true
		}
	}
	return nil, false
}

func illuminateResponseConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Illuminate\\Http\\Response::__construct: invalid context"))
	}
	content, _ := ctx.GetIndexValue(0)
	status := httpfoundation.IntParam(ctx, 1, 200)
	headersVal, _ := ctx.GetIndexValue(2)

	headers := httpfoundation.CreateResponseHeaders(ctx, headersVal)
	httpfoundation.ResponseSetProp(cv, "headers", headers)
	httpfoundation.ResponseSetProp(cv, "version", data.NewStringValue("1.0"))
	if ctl := httpfoundation.ApplyStatusCode(cv, status, nil); ctl != nil {
		return nil, ctl
	}
	httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(""))
	httpfoundation.ResponseSetProp(cv, "original", data.NewNullValue())
	httpfoundation.ResponseSetProp(cv, "exception", data.NewNullValue())

	// 走 Illuminate setContent（支持 JSON / Renderable）
	setCtx := cv.CreateContext([]data.Variable{httpfoundation.Variable("content", 0, nil)})
	if content == nil {
		content = data.NewStringValue("")
	}
	setCtx.SetVariableValue(httpfoundation.Variable("content", 0, nil), content)
	return illuminateResponseSetContent(setCtx)
}

func illuminateResponseGetContent(ctx data.Context) (data.GetValue, data.Control) {
	s := httpfoundation.ResponseContent(httpfoundation.ResponseClassValue(ctx))
	return data.NewStringValue(s), nil
}

func illuminateResponseSetContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	content, _ := ctx.GetIndexValue(0)
	if content == nil {
		content = data.NewNullValue()
	}
	httpfoundation.ResponseSetProp(cv, "original", content)

	if shouldBeJSON(content) {
		httpfoundation.RespHeaderSet(httpfoundation.ResponseHeaders(cv), "Content-Type", []string{"application/json"}, true)
		encoded, ok := morphToJSON(ctx, content)
		if !ok {
			return nil, httpfoundation.ThrowNamed("InvalidArgumentException", "json_encode error")
		}
		httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(encoded))
		return httpfoundation.ResponseSelf(ctx), nil
	}

	rendered, ok, ctl := tryRender(ctx, content)
	if ctl != nil {
		return nil, ctl
	}
	if ok {
		httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(rendered))
		return httpfoundation.ResponseSelf(ctx), nil
	}

	if _, isNull := content.(*data.NullValue); isNull {
		httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(""))
	} else if _, isObj := content.(*data.ClassValue); isObj {
		return nil, httpfoundation.ThrowNamed("InvalidArgumentException",
			"The Response content must be a string or object implementing __toString(), object given.")
	} else {
		httpfoundation.ResponseSetProp(cv, "content", data.NewStringValue(content.AsString()))
	}
	return httpfoundation.ResponseSelf(ctx), nil
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
	if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Renderable") ||
		implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Htmlable") {
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
	vm := cv.GetVM()
	// 检查类直接实现的接口及其继承链
	for _, iface := range cv.Class.GetImplements() {
		if interfaceMatches(vm, iface, name) {
			return true
		}
	}
	// 向上遍历父类链
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
			if interfaceMatches(vm, iface, name) {
				return true
			}
		}
		last = next
	}
	return false
}

// interfaceMatches 检查接口名或接口继承链中是否包含目标接口
func interfaceMatches(vm data.VM, ifaceName string, target string) bool {
	if ifaceName == target {
		return true
	}
	if vm == nil {
		return false
	}
	iface, ok := vm.GetInterface(ifaceName)
	if !ok || iface == nil {
		// 尝试加载接口（可能尚未注册到 VM）
		loaded, acl := vm.GetOrLoadInterface(ifaceName)
		if acl != nil || loaded == nil {
			return false
		}
		iface = loaded
	}
	// 检查该接口继承的父接口
	for _, parent := range iface.GetExtends() {
		if parent == target {
			return true
		}
		if interfaceMatches(vm, parent, target) {
			return true
		}
	}
	return false
}

func morphToJSON(ctx data.Context, content data.Value) (string, bool) {
	if content == nil {
		return "null", true
	}
	if cv, ok := content.(*data.ClassValue); ok {
		if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Jsonable") {
			ret, ctl := httpfoundation.CallObjMethod(cv, "toJson")
			if ctl != nil || ret == nil {
				return "", false
			}
			return ret.(data.Value).AsString(), true
		}
		if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Arrayable") {
			ret, ctl := httpfoundation.CallObjMethod(cv, "toArray")
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
	case *data.ObjectValue:
		// 将 ObjectValue（关联数组）转换为 Go map
		out := map[string]any{}
		t.RangeProperties(func(key string, value data.Value) bool {
			out[key] = phpValueToGo(value)
			return true
		})
		return out
	case *data.ClassValue:
		return t.AsString()
	default:
		return v.AsString()
	}
}

func tryRender(ctx data.Context, content data.Value) (string, bool, data.Control) {
	cv, ok := content.(*data.ClassValue)
	if !ok {
		return "", false, nil
	}

	callStringMethod := func(name string) (string, bool, data.Control) {
		if _, has := cv.GetMethod(name); !has {
			return "", false, nil
		}
		ret, ctl := httpfoundation.CallObjMethod(cv, name)
		if ctl != nil {
			return "", false, ctl
		}
		s, ok := stringifyRenderedValue(ret)
		return s, ok, nil
	}

	if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Renderable") {
		if s, ok, ctl := callStringMethod("render"); ctl != nil || ok {
			return s, ok, ctl
		}
	} else if s, ok, ctl := callStringMethod("render"); ctl != nil || ok {
		return s, ok, ctl
	}

	if implementsOrExtends(cv, "Illuminate\\Contracts\\Support\\Htmlable") {
		if s, ok, ctl := callStringMethod("toHtml"); ctl != nil || ok {
			return s, ok, ctl
		}
	} else if s, ok, ctl := callStringMethod("toHtml"); ctl != nil || ok {
		return s, ok, ctl
	}

	if s, ok, ctl := callStringMethod("__toString"); ctl != nil || ok {
		return s, ok, ctl
	}

	return "", false, nil
}

// stringifyRenderedValue 把 render()/toHtml() 的返回值转成响应正文。
// ClassValue 不能走 AsString()：那是对象调试 dump，不是 HTML。
func stringifyRenderedValue(ret data.GetValue) (string, bool) {
	if ret == nil {
		return "", false
	}
	switch t := ret.(type) {
	case *data.StringValue:
		return t.Value, true
	case *data.NullValue:
		return "", true
	case *data.ClassValue:
		return "", false
	case data.Value:
		return t.AsString(), true
	default:
		return "", false
	}
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
	return data.NewIntValue(httpfoundation.ResponseStatusCode(httpfoundation.ResponseClassValue(ctx))), nil
}

func illuminateResponseStatusText(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(httpfoundation.ResponseStatusText(httpfoundation.ResponseClassValue(ctx))), nil
}

func illuminateResponseContent(ctx data.Context) (data.GetValue, data.Control) {
	return illuminateResponseGetContent(ctx)
}

func illuminateResponseGetOriginalContent(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	v, _ := cv.GetProperty("original")
	if nested, ok := v.(*data.ClassValue); ok {
		if nested.Class != nil && nested.Class.GetName() == fqnIlluminateResponse {
			ret, ctl := httpfoundation.CallObjMethod(nested, "getOriginalContent")
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
	cv := httpfoundation.ResponseClassValue(ctx)
	key, _ := ctx.GetIndexValue(0)
	values, _ := ctx.GetIndexValue(1)
	replace := httpfoundation.BoolParam(ctx, 2, true)
	if key == nil {
		return httpfoundation.ResponseSelf(ctx), nil
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
	httpfoundation.RespHeaderSet(httpfoundation.ResponseHeaders(cv), key.AsString(), vals, replace)
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateResponseWithHeaders(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	headersVal, _ := ctx.GetIndexValue(0)
	var m map[string]data.Value
	if hv, ok := headersVal.(*data.ClassValue); ok {
		if httpfoundation.HeaderBagFrom(hv) != nil {
			all := httpfoundation.GetHeaderBagAll(hv)
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
		m, err = httpfoundation.ValueToAssocMap(headersVal)
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
		httpfoundation.RespHeaderSet(httpfoundation.ResponseHeaders(cv), k, vals, true)
	}
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateResponseWithoutHeader(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	key, _ := ctx.GetIndexValue(0)
	if arr, ok := key.(*data.ArrayValue); ok {
		for _, z := range arr.List {
			if z != nil && z.Value != nil {
				httpfoundation.RespHeaderRemove(httpfoundation.ResponseHeaders(cv), z.Value.AsString())
			}
		}
	} else if key != nil {
		httpfoundation.RespHeaderRemove(httpfoundation.ResponseHeaders(cv), key.AsString())
	}
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateResponseWithCookie(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	cookie, _ := ctx.GetIndexValue(0)
	headers := httpfoundation.ResponseHeaders(cv)
	if rh := httpfoundation.ResponseHeaderBagFrom(headers); rh != nil && cookie != nil {
		// 优先走 ResponseHeaderBag::setCookie
		_, ctl := httpfoundation.CallObjMethod(headers, "setCookie", cookie)
		if ctl != nil {
			return nil, ctl
		}
	}
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateResponseWithCookies(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	cookies, _ := ctx.GetIndexValue(0)
	arr, ok := cookies.(*data.ArrayValue)
	if !ok {
		return httpfoundation.ResponseSelf(ctx), nil
	}
	headers := httpfoundation.ResponseHeaders(cv)
	for _, z := range arr.List {
		if z == nil || z.Value == nil {
			continue
		}
		_, ctl := httpfoundation.CallObjMethod(headers, "setCookie", z.Value)
		if ctl != nil {
			return nil, ctl
		}
	}
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateResponseWithoutCookie(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	cookie, _ := ctx.GetIndexValue(0)
	headers := httpfoundation.ResponseHeaders(cv)
	if cookie == nil {
		return httpfoundation.ResponseSelf(ctx), nil
	}
	// 字符串 cookie 名：构造过期 cookie（简化）
	if _, ok := cookie.(*data.StringValue); ok {
		name := cookie.AsString()
		path := "/"
		if p, present, ok := httpfoundation.OptionalStringParam(ctx, 1); present && ok && p != "" {
			path = p
		}
		domain := ""
		if d, present, ok := httpfoundation.OptionalStringParam(ctx, 2); present && ok {
			domain = d
		}
		val := ""
		bc := &httpfoundation.BagCookie{
			Name:     name,
			Value:    &val,
			Expire:   -2628000,
			Path:     path,
			HTTPOnly: true,
		}
		if domain != "" {
			bc.Domain = &domain
		}
		if rh := httpfoundation.ResponseHeaderBagFrom(headers); rh != nil {
			httpfoundation.SetBagCookie(rh, bc)
			return httpfoundation.ResponseSelf(ctx), nil
		}
	}
	_, ctl := httpfoundation.CallObjMethod(headers, "setCookie", cookie)
	if ctl != nil {
		return nil, ctl
	}
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateResponseGetCallback(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	v, _ := cv.GetProperty("callback")
	if v == nil {
		return data.NewNullValue(), nil
	}
	return v, nil
}

func illuminateResponseWithException(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
	e, _ := ctx.GetIndexValue(0)
	if e == nil {
		e = data.NewNullValue()
	}
	httpfoundation.ResponseSetProp(cv, "exception", e)
	return httpfoundation.ResponseSelf(ctx), nil
}

func illuminateResponseThrowResponse(ctx data.Context) (data.GetValue, data.Control) {
	cv := httpfoundation.ResponseClassValue(ctx)
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

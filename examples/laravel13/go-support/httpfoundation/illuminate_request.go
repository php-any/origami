package httpfoundation

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const illuminateRequestClassName = "Illuminate\\Http\\Request"

var illuminateRequestMacros = struct {
	sync.RWMutex
	values map[string]data.Value
}{values: make(map[string]data.Value)}

// illuminateRequestState only contains state introduced by Illuminate Request.
// Symfony request state remains owned by the parent ClassStmt/source.
type illuminateRequestState struct {
	source        *http.Request
	json          *data.ClassValue
	userResolver  data.Value
	routeResolver data.Value
	session       data.Value
	locale        string
	defaultLocale string
}

// IlluminateRequestClass implements Illuminate\Http\Request and its four HTTP
// concerns. It deliberately obtains the active request through Symfony's
// createFromGlobals implementation; it never reads node HTTP globals.
type IlluminateRequestClass struct {
	node.Node
	methods map[string]data.Method
}

func NewIlluminateRequestClass() data.ClassStmt {
	c := &IlluminateRequestClass{methods: make(map[string]data.Method)}
	c.addMethods()
	return c
}

// NewIlluminateRequestValue 直接从 Go HTTP 请求构造 Laravel Request。
// Go HTTP 内核使用此入口，无需再执行 public/index.php 或经过 PHP capture 桥接。
func NewIlluminateRequestValue(ctx data.Context, request *http.Request) *data.ClassValue {
	stmt := NewIlluminateRequestClass()
	value := data.NewClassValue(stmt, ctx.CreateBaseContext())
	initializeIlluminateRequest(value, request)
	return value
}

func (c *IlluminateRequestClass) GetName() string { return illuminateRequestClassName }
func (c *IlluminateRequestClass) GetExtend() *string {
	parent := "Symfony\\Component\\HttpFoundation\\Request"
	return &parent
}
func (c *IlluminateRequestClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Support\\Arrayable", "ArrayAccess"}
}
func (c *IlluminateRequestClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *IlluminateRequestClass) GetPropertyList() []data.Property         { return nil }
func (c *IlluminateRequestClass) GetConstruct() data.Method                { return &illuminateRequestConstruct{} }

func (c *IlluminateRequestClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	initializeIlluminateRequest(cv, nil)
	return cv, nil
}

func (c *IlluminateRequestClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}

func (c *IlluminateRequestClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GetName() < out[j].GetName() })
	return out
}

func (c *IlluminateRequestClass) GetStaticMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "capture":
		return newIlluminateMethod("capture", true, nil, illuminateCapture), true
	case "create":
		return newIlluminateMethod("create", true, []string{"uri", "method", "parameters", "cookies", "files", "server", "content"}, illuminateCreate), true
	case "macro":
		return newIlluminateMethod("macro", true, []string{"name", "macro"}, illuminateMacro), true
	case "mixin":
		return newIlluminateMethod("mixin", true, []string{"mixin", "replace"}, illuminateMixin), true
	case "hasmacro":
		return newIlluminateMethod("hasMacro", true, []string{"name"}, illuminateHasMacro), true
	case "flushmacros":
		return newIlluminateMethod("flushMacros", true, nil, illuminateFlushMacros), true
	case "__callstatic":
		return newIlluminateMethod("__callStatic", true, []string{"method", "parameters"}, illuminateCallStatic), true
	case "createfrom":
		return newIlluminateMethod("createFrom", true, []string{"from", "to"}, illuminateCreateFrom), true
	case "createfrombase":
		return newIlluminateMethod("createFromBase", true, []string{"request"}, illuminateCreateFromBase), true
	case "matchestype":
		return newIlluminateMethod("matchesType", true, []string{"actual", "type"}, illuminateMatchesType), true
	}
	return nil, false
}

func (c *IlluminateRequestClass) add(name string, params []string, call illuminateCall) {
	c.methods[strings.ToLower(name)] = newIlluminateMethod(name, false, params, call)
}

func (c *IlluminateRequestClass) addVariadic(name string, params []string, call illuminateCall) {
	m := newIlluminateMethod(name, false, params, call)
	m.variadic = true
	c.methods[strings.ToLower(name)] = m
}

func (c *IlluminateRequestClass) addMethods() {
	// Request.php
	for _, name := range []string{"instance", "method", "uri", "root", "url", "fullUrl", "path",
		"decodedPath", "segments", "host", "httpHost", "schemeAndHttpHost", "ajax", "pjax",
		"prefetch", "secure", "ip", "ips", "userAgent", "getAcceptableContentTypes",
		"hasSession", "getSession", "session", "getUserResolver", "getRouteResolver", "toArray"} {
		c.add(name, nil, illuminateDispatch)
	}
	c.add("fullUrlWithQuery", []string{"query"}, illuminateDispatch)
	c.add("fullUrlWithoutQuery", []string{"keys"}, illuminateDispatch)
	c.add("segment", []string{"index", "default"}, illuminateDispatch)
	c.addVariadic("is", []string{"patterns"}, illuminateDispatch)
	c.addVariadic("routeIs", []string{"patterns"}, illuminateDispatch)
	c.addVariadic("fullUrlIs", []string{"patterns"}, illuminateDispatch)
	c.add("merge", []string{"input"}, illuminateDispatch)
	c.add("mergeIfMissing", []string{"input"}, illuminateDispatch)
	c.add("replace", []string{"input"}, illuminateDispatch)
	c.add("get", []string{"key", "default"}, illuminateDispatch)
	c.add("json", []string{"key", "default"}, illuminateDispatch)
	c.add("duplicate", []string{"query", "request", "attributes", "cookies", "files", "server"}, illuminateDispatch)
	c.add("setLaravelSession", []string{"session"}, illuminateDispatch)
	c.add("setRequestLocale", []string{"locale"}, illuminateDispatch)
	c.add("setDefaultRequestLocale", []string{"locale"}, illuminateDispatch)
	c.add("user", []string{"guard"}, illuminateDispatch)
	c.add("route", []string{"param", "default"}, illuminateDispatch)
	c.add("fingerprint", nil, illuminateDispatch)
	c.add("setJson", []string{"json"}, illuminateDispatch)
	c.add("setUserResolver", []string{"callback"}, illuminateDispatch)
	c.add("setRouteResolver", []string{"callback"}, illuminateDispatch)
	c.add("offsetExists", []string{"offset"}, illuminateDispatch)
	c.add("offsetGet", []string{"offset"}, illuminateDispatch)
	c.add("offsetSet", []string{"offset", "value"}, illuminateDispatch)
	c.add("offsetUnset", []string{"offset"}, illuminateDispatch)
	c.add("__isset", []string{"key"}, illuminateDispatch)
	c.add("__get", []string{"key"}, illuminateDispatch)

	// InteractsWithInput.
	for _, name := range []string{"server", "header", "query", "post", "cookie", "file"} {
		c.add(name, []string{"key", "default"}, illuminateDispatch)
	}
	for _, name := range []string{"hasHeader", "hasCookie", "hasFile"} {
		c.add(name, []string{"key"}, illuminateDispatch)
	}
	for _, name := range []string{"bearerToken", "keys", "allFiles"} {
		c.add(name, nil, illuminateDispatch)
	}
	c.addVariadic("all", []string{"keys"}, illuminateDispatch)
	c.add("input", []string{"key", "default"}, illuminateDispatch)
	c.add("fluent", []string{"key", "default"}, illuminateDispatch)
	c.add("image", []string{"key"}, illuminateDispatch)
	c.addVariadic("dump", []string{"keys"}, illuminateDispatch)

	// InteractsWithData, pulled in by InteractsWithInput.
	for _, name := range []string{"exists", "has", "hasAny", "filled", "isNotFilled", "anyFilled", "missing"} {
		c.addVariadic(name, []string{"keys"}, illuminateDispatch)
	}
	for _, name := range []string{"whenHas", "whenFilled", "whenMissing"} {
		c.add(name, []string{"key", "callback", "default"}, illuminateDispatch)
	}
	c.add("whenEnum", []string{"key", "enumClass", "callback", "default"}, illuminateDispatch)
	for _, name := range []string{"str", "string", "boolean", "integer", "float"} {
		c.add(name, []string{"key", "default"}, illuminateDispatch)
	}
	c.add("clamp", []string{"key", "min", "max", "default"}, illuminateDispatch)
	c.add("date", []string{"key", "format", "tz"}, illuminateDispatch)
	c.add("interval", []string{"key", "unit"}, illuminateDispatch)
	c.add("enum", []string{"key", "enumClass", "default"}, illuminateDispatch)
	c.add("enums", []string{"key", "enumClass"}, illuminateDispatch)
	c.add("array", []string{"key"}, illuminateDispatch)
	c.add("collect", []string{"key"}, illuminateDispatch)
	c.addVariadic("only", []string{"keys"}, illuminateDispatch)
	c.addVariadic("except", []string{"keys"}, illuminateDispatch)

	// Flash data.
	c.add("old", []string{"key", "default"}, illuminateDispatch)
	for _, name := range []string{"flash", "flush"} {
		c.add(name, nil, illuminateDispatch)
	}
	c.addVariadic("flashOnly", []string{"keys"}, illuminateDispatch)
	c.addVariadic("flashExcept", []string{"keys"}, illuminateDispatch)

	// Content negotiation and precognition.
	for _, name := range []string{"isJson", "expectsJson", "wantsJson", "wantsMarkdown",
		"acceptsAnyContentType", "acceptsJson", "acceptsMarkdown", "acceptsHtml",
		"isAttemptingPrecognition", "isPrecognitive"} {
		c.add(name, nil, illuminateDispatch)
	}
	c.add("accepts", []string{"contentTypes"}, illuminateDispatch)
	c.add("prefers", []string{"contentTypes"}, illuminateDispatch)
	c.add("format", []string{"default"}, illuminateDispatch)
	c.add("filterPrecognitiveRules", []string{"rules"}, illuminateDispatch)
	c.add("__call", []string{"method", "parameters"}, illuminateCallMacro)
}

type illuminateCall func(*illuminateRequestMethod, data.Context) (data.GetValue, data.Control)

type illuminateRequestMethod struct {
	name     string
	static   bool
	params   []string
	variadic bool
	call     illuminateCall
}

func newIlluminateMethod(name string, static bool, params []string, call illuminateCall) *illuminateRequestMethod {
	return &illuminateRequestMethod{name: name, static: static, params: params, call: call}
}

func (m *illuminateRequestMethod) GetName() string            { return m.name }
func (m *illuminateRequestMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *illuminateRequestMethod) GetIsStatic() bool          { return m.static }
func (m *illuminateRequestMethod) GetReturnType() data.Types  { return nil }
func (m *illuminateRequestMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, name := range m.params {
		if m.variadic && i == len(m.params)-1 {
			out[i] = data.NewParameters(name, i)
		} else {
			out[i] = data.NewParameterDefault(name, i, data.NewNullValue(), nil)
		}
	}
	return out
}
func (m *illuminateRequestMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, name := range m.params {
		out[i] = data.NewVariable(name, i, nil)
	}
	return out
}
func (m *illuminateRequestMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return m.call(m, ctx)
}

type illuminateRequestConstruct struct{}

func (m *illuminateRequestConstruct) GetName() string            { return "__construct" }
func (m *illuminateRequestConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *illuminateRequestConstruct) GetIsStatic() bool          { return false }
func (m *illuminateRequestConstruct) GetReturnType() data.Types  { return nil }
func (m *illuminateRequestConstruct) GetParams() []data.GetValue {
	names := []string{"query", "request", "attributes", "cookies", "files", "server", "content"}
	out := make([]data.GetValue, len(names))
	for i, name := range names {
		out[i] = data.NewParameterDefault(name, i, data.NewNullValue(), nil)
	}
	return out
}
func (m *illuminateRequestConstruct) GetVariables() []data.Variable {
	names := []string{"query", "request", "attributes", "cookies", "files", "server", "content"}
	out := make([]data.Variable, len(names))
	for i, name := range names {
		out[i] = data.NewVariable(name, i, nil)
	}
	return out
}
func (m *illuminateRequestConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := requestClassValue(ctx)
	if cv == nil {
		return nil, requestError("Request constructor missing object context")
	}
	initializeIlluminateRequest(cv, nil)
	bags := []string{"query", "request", "attributes", "cookies", "files", "server"}
	for i, property := range bags {
		if v, ok := ctx.GetIndexValue(i); ok {
			if values, err := valueToAssocMap(v); err == nil {
				setRequestBag(cv, property, values)
			}
		}
	}
	if content, ok := ctx.GetIndexValue(6); ok {
		_ = cv.SetProperty("content", content)
	}
	return data.NewNullValue(), nil
}

func initializeIlluminateRequest(cv *data.ClassValue, source *http.Request) {
	if cv == nil {
		return
	}
	state := &illuminateRequestState{source: source, defaultLocale: "en", locale: "en"}
	_ = cv.SetProperty("__illuminate_request_state", data.NewAnyValue(state))
	if _, ctl := cv.GetProperty("query"); ctl != nil {
		_ = cv.SetProperty("query", NewInputBagValue(cv, map[string]data.Value{}))
	}
	if source != nil {
		populateRequestFromHTTP(cv, source)
	}
	for name, maker := range map[string]func() data.Value{
		"query":      func() data.Value { return NewInputBagValue(cv, nil) },
		"request":    func() data.Value { return NewInputBagValue(cv, nil) },
		"attributes": func() data.Value { return NewParameterBagValue(cv, nil) },
		"cookies":    func() data.Value { return NewInputBagValue(cv, nil) },
		"files":      func() data.Value { return NewFileBagValue(cv, nil) },
		"server":     func() data.Value { return NewServerBagValue(cv, nil) },
		"headers":    func() data.Value { return NewHeaderBagValue(cv, nil) },
	} {
		if v, _ := cv.GetProperty(name); isNull(v) {
			_ = cv.SetProperty(name, maker())
		}
	}
}

func requestState(cv *data.ClassValue) *illuminateRequestState {
	if cv == nil {
		return nil
	}
	if value, _ := cv.GetProperty("__illuminate_request_state"); value != nil {
		if anyValue, ok := value.(*data.AnyValue); ok {
			if state, ok := anyValue.Value.(*illuminateRequestState); ok {
				return state
			}
		}
	}
	state := &illuminateRequestState{defaultLocale: "en", locale: "en"}
	if source, ok := requestHTTPSource(cv); ok {
		state.source = source
	}
	_ = cv.SetProperty("__illuminate_request_state", data.NewAnyValue(state))
	return state
}

func requestHTTPSource(cv *data.ClassValue) (*http.Request, bool) {
	if cv == nil {
		return nil, false
	}
	if state := requestStateNoCreate(cv); state != nil && state.source != nil {
		return state.source, true
	}
	if sourceClass, ok := cv.Class.(data.GetSource); ok {
		if source, ok := sourceClass.GetSource().(*http.Request); ok {
			return source, source != nil
		}
	}
	if class, ok := cv.Class.(*SymfonyRequestClass); ok && class.source != nil && class.source.request != nil {
		return class.source.request, true
	}
	return nil, false
}

func requestStateNoCreate(cv *data.ClassValue) *illuminateRequestState {
	value, _ := cv.GetProperty("__illuminate_request_state")
	if anyValue, ok := value.(*data.AnyValue); ok {
		state, _ := anyValue.Value.(*illuminateRequestState)
		return state
	}
	return nil
}

func populateRequestFromHTTP(cv *data.ClassValue, r *http.Request) {
	query := map[string]data.Value{}
	for key, values := range r.URL.Query() {
		query[key] = stringsToValue(values)
	}
	_ = r.ParseMultipartForm(32 << 20)
	form := map[string]data.Value{}
	for key, values := range r.PostForm {
		form[key] = stringsToValue(values)
	}
	cookies := map[string]data.Value{}
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = data.NewStringValue(cookie.Value)
	}
	server := map[string]data.Value{
		"REQUEST_METHOD":  data.NewStringValue(r.Method),
		"REQUEST_URI":     data.NewStringValue(r.URL.RequestURI()),
		"REMOTE_ADDR":     data.NewStringValue(r.RemoteAddr),
		"SERVER_PROTOCOL": data.NewStringValue(r.Proto),
	}
	_ = cv.SetProperty("query", NewInputBagValue(cv, query))
	_ = cv.SetProperty("request", NewInputBagValue(cv, form))
	_ = cv.SetProperty("attributes", NewParameterBagValue(cv, nil))
	_ = cv.SetProperty("cookies", NewInputBagValue(cv, cookies))
	_ = cv.SetProperty("files", NewFileBagValue(cv, nil))
	_ = cv.SetProperty("server", NewServerBagValue(cv, server))
	_ = cv.SetProperty("headers", NewHeaderBagValue(cv, r.Header))
}

func stringsToValue(values []string) data.Value {
	if len(values) == 1 {
		return data.NewStringValue(values[0])
	}
	list := make([]data.Value, len(values))
	for i, value := range values {
		list[i] = data.NewStringValue(value)
	}
	return data.NewArrayValue(list)
}

func illuminateCapture(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	parent, control := ctx.GetVM().GetOrLoadClass("Symfony\\Component\\HttpFoundation\\Request")
	if control != nil {
		return nil, control
	}
	getter, ok := parent.(data.GetStaticMethod)
	if !ok {
		return nil, requestError("Symfony Request does not expose createFromGlobals")
	}
	method, ok := getter.GetStaticMethod("createFromGlobals")
	if !ok {
		return nil, requestError("Symfony Request::createFromGlobals is unavailable")
	}
	base, control := method.Call(ctx.CreateContext(method.GetVariables()))
	if control != nil {
		return nil, control
	}
	baseValue, ok := base.(data.Value)
	if !ok {
		return nil, requestError("Symfony Request::createFromGlobals returned an invalid request")
	}
	return createIlluminateFromValue(ctx, baseValue, nil)
}

func illuminateMacro(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	name := strings.ToLower(argString(ctx, 0, ""))
	macro := argValue(ctx, 1, data.NewNullValue())
	illuminateRequestMacros.Lock()
	illuminateRequestMacros.values[name] = macro
	illuminateRequestMacros.Unlock()
	return data.NewNullValue(), nil
}

func illuminateMixin(_ *illuminateRequestMethod, _ data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func illuminateHasMacro(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	name := strings.ToLower(argString(ctx, 0, ""))
	illuminateRequestMacros.RLock()
	_, ok := illuminateRequestMacros.values[name]
	illuminateRequestMacros.RUnlock()
	return data.NewBoolValue(ok), nil
}

func illuminateFlushMacros(_ *illuminateRequestMethod, _ data.Context) (data.GetValue, data.Control) {
	illuminateRequestMacros.Lock()
	illuminateRequestMacros.values = make(map[string]data.Value)
	illuminateRequestMacros.Unlock()
	return data.NewNullValue(), nil
}

func illuminateCallStatic(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	return callIlluminateMacro(ctx, nil)
}

func illuminateCallMacro(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	return callIlluminateMacro(ctx, requestClassValue(ctx))
}

func callIlluminateMacro(ctx data.Context, receiver *data.ClassValue) (data.GetValue, data.Control) {
	name := strings.ToLower(argString(ctx, 0, ""))
	illuminateRequestMacros.RLock()
	macro := illuminateRequestMacros.values[name]
	illuminateRequestMacros.RUnlock()
	if macro == nil {
		return nil, requestError("Method %s::%s does not exist.", illuminateRequestClassName, name)
	}
	args := make([]data.Value, 0)
	if array, ok := argValue(ctx, 1, data.NewArrayValue(nil)).(*data.ArrayValue); ok {
		for _, item := range array.List {
			if item != nil && item.Value != nil {
				args = append(args, item.Value)
			}
		}
	}
	if receiver != nil {
		if fn, ok := macro.(*data.FuncValue); ok {
			macro = data.NewBoundFuncValue(fn.Value, illuminateRequestClassName, receiver)
		}
	}
	return invokeCallable(macro, args, ctx)
}

func illuminateCreate(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	parent, control := ctx.GetVM().GetOrLoadClass("Symfony\\Component\\HttpFoundation\\Request")
	if control != nil {
		return nil, control
	}
	getter, ok := parent.(data.GetStaticMethod)
	if !ok {
		return nil, requestError("Symfony Request does not expose create")
	}
	method, ok := getter.GetStaticMethod("create")
	if !ok {
		return nil, requestError("Symfony Request::create is unavailable")
	}
	base, control := method.Call(ctx)
	if control != nil {
		return nil, control
	}
	baseValue, ok := base.(data.Value)
	if !ok {
		return nil, requestError("Symfony Request::create returned an invalid request")
	}
	return createIlluminateFromValue(ctx, baseValue, nil)
}

func illuminateCreateFrom(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	from, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, requestError("Request::createFrom expects a source request")
	}
	var target *data.ClassValue
	if value, ok := ctx.GetIndexValue(1); ok {
		target, _ = value.(*data.ClassValue)
	}
	return createIlluminateFromValue(ctx, from, target)
}

func illuminateCreateFromBase(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	from, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, requestError("Request::createFromBase expects a Symfony request")
	}
	return createIlluminateFromValue(ctx, from, nil)
}

func createIlluminateFromValue(ctx data.Context, from data.Value, target *data.ClassValue) (data.GetValue, data.Control) {
	sourceCV, ok := from.(*data.ClassValue)
	if !ok {
		return nil, requestError("source request is not an object")
	}
	if target == nil {
		stmt := NewIlluminateRequestClass()
		target = data.NewClassValue(stmt, ctx.CreateBaseContext())
	}
	var source *http.Request
	if request, ok := requestHTTPSource(sourceCV); ok {
		source = request
	}
	initializeIlluminateRequest(target, source)
	for name, value := range sourceCV.GetProperties() {
		if name == "__illuminate_request_state" {
			continue
		}
		_ = target.SetProperty(name, value)
	}
	state := requestState(target)
	if old := requestStateNoCreate(sourceCV); old != nil {
		state.userResolver = old.userResolver
		state.routeResolver = old.routeResolver
		state.session = old.session
		state.locale = old.locale
		state.defaultLocale = old.defaultLocale
		state.json = old.json
	}
	if isJSONRequest(target) {
		jsonValue, _ := requestJSON(target, nil, data.NewNullValue())
		if bag, ok := jsonValue.(*data.ClassValue); ok {
			_ = target.SetProperty("request", bag)
			state.json = bag
		}
	}
	return target, nil
}

func illuminateMatchesType(_ *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	actual := argString(ctx, 0, "")
	typ := argString(ctx, 1, "")
	return data.NewBoolValue(matchesContentType(actual, typ)), nil
}

func illuminateDispatch(method *illuminateRequestMethod, ctx data.Context) (data.GetValue, data.Control) {
	cv := requestClassValue(ctx)
	if cv == nil {
		return nil, requestError("Request::%s missing object context", method.name)
	}
	switch strings.ToLower(method.name) {
	case "instance":
		return cv, nil
	case "method":
		return data.NewStringValue(requestMethod(cv)), nil
	case "uri":
		return data.NewStringValue(requestFullURL(cv)), nil
	case "root":
		return data.NewStringValue(strings.TrimRight(requestSchemeHost(cv)+requestBaseURL(cv), "/")), nil
	case "url":
		return data.NewStringValue(strings.TrimRight(requestURL(cv), "/")), nil
	case "fullurl":
		return data.NewStringValue(requestFullURL(cv)), nil
	case "fullurlwithquery":
		return data.NewStringValue(fullURLWithQuery(cv, argMap(ctx, 0))), nil
	case "fullurlwithoutquery":
		return data.NewStringValue(fullURLWithoutQuery(cv, argumentKeys(ctx, 0))), nil
	case "path":
		p := strings.Trim(requestPath(cv), "/")
		if p == "" {
			p = "/"
		}
		return data.NewStringValue(p), nil
	case "decodedpath":
		p, err := url.PathUnescape(strings.Trim(requestPath(cv), "/"))
		if err != nil || p == "" {
			p = "/"
		}
		return data.NewStringValue(p), nil
	case "segments":
		return stringArray(requestSegments(cv)), nil
	case "segment":
		index := argInt(ctx, 0, 0) - 1
		segments := requestSegments(cv)
		if index >= 0 && index < len(segments) {
			return data.NewStringValue(segments[index]), nil
		}
		return argValue(ctx, 1, data.NewNullValue()), nil
	case "is":
		return data.NewBoolValue(anyPattern(requestDecodedPath(cv), variadicStrings(ctx, 0))), nil
	case "fullurlis":
		return data.NewBoolValue(anyPattern(requestFullURL(cv), variadicStrings(ctx, 0))), nil
	case "routeis":
		route, control := requestRoute(cv, nil, data.NewNullValue(), ctx)
		if control != nil || isNullValue(route) {
			return data.NewBoolValue(false), control
		}
		return invokeObject(route, "named", valuesFromStrings(variadicStrings(ctx, 0)), ctx)
	case "host":
		return data.NewStringValue(requestHost(cv)), nil
	case "httphost":
		return data.NewStringValue(requestHTTPHost(cv)), nil
	case "schemeandhttphost":
		return data.NewStringValue(requestSchemeHost(cv)), nil
	case "ajax":
		return data.NewBoolValue(strings.EqualFold(requestHeader(cv, "X-Requested-With"), "XMLHttpRequest")), nil
	case "pjax":
		return data.NewBoolValue(requestHeader(cv, "X-PJAX") != ""), nil
	case "prefetch":
		return data.NewBoolValue(strings.EqualFold(requestHeader(cv, "X-Moz"), "prefetch") ||
			strings.EqualFold(requestHeader(cv, "Purpose"), "prefetch") ||
			strings.EqualFold(requestHeader(cv, "Sec-Purpose"), "prefetch")), nil
	case "secure":
		return data.NewBoolValue(requestSecure(cv)), nil
	case "ip":
		return nullableString(requestIP(cv)), nil
	case "ips":
		return stringArray(requestIPs(cv)), nil
	case "useragent":
		return nullableString(requestHeader(cv, "User-Agent")), nil
	case "getacceptablecontenttypes":
		return stringArray(acceptableTypes(cv)), nil
	case "merge", "mergeifmissing":
		values := argMap(ctx, 0)
		current := requestInputSource(cv)
		if strings.EqualFold(method.name, "mergeIfMissing") {
			for key := range values {
				if _, ok := current[key]; ok {
					delete(values, key)
				}
			}
		}
		for key, value := range values {
			setNested(current, key, value)
		}
		setInputSource(cv, current)
		return cv, nil
	case "replace":
		setInputSource(cv, argMap(ctx, 0))
		return cv, nil
	case "get":
		key := argString(ctx, 0, "")
		for _, bag := range []string{"attributes", "query", "request"} {
			if value, ok := requestBagMap(cv, bag)[key]; ok {
				return value, nil
			}
		}
		return argValue(ctx, 1, data.NewNullValue()), nil
	case "json":
		return requestJSON(cv, argOptionalString(ctx, 0), argValue(ctx, 1, data.NewNullValue()))
	case "duplicate":
		return duplicateIlluminateRequest(cv, ctx)
	case "hassession":
		return data.NewBoolValue(requestState(cv).session != nil && !isNull(requestState(cv).session)), nil
	case "getsession":
		if requestState(cv).session == nil || isNull(requestState(cv).session) {
			return nil, requestError("Session not found")
		}
		return requestState(cv).session, nil
	case "session":
		if requestState(cv).session == nil || isNull(requestState(cv).session) {
			return nil, requestError("Session store not set on request.")
		}
		return requestState(cv).session, nil
	case "setlaravelsession":
		requestState(cv).session = argValue(ctx, 0, data.NewNullValue())
		return data.NewNullValue(), nil
	case "setrequestlocale":
		requestState(cv).locale = argString(ctx, 0, "en")
		_ = cv.SetProperty("locale", data.NewStringValue(requestState(cv).locale))
		return data.NewNullValue(), nil
	case "setdefaultrequestlocale":
		requestState(cv).defaultLocale = argString(ctx, 0, "en")
		_ = cv.SetProperty("defaultLocale", data.NewStringValue(requestState(cv).defaultLocale))
		return data.NewNullValue(), nil
	case "user":
		resolver := requestState(cv).userResolver
		if resolver == nil {
			return data.NewNullValue(), nil
		}
		return invokeCallable(resolver, []data.Value{argValue(ctx, 0, data.NewNullValue())}, ctx)
	case "route":
		return requestRoute(cv, argOptionalString(ctx, 0), argValue(ctx, 1, data.NewNullValue()), ctx)
	case "fingerprint":
		return requestFingerprint(cv, ctx)
	case "setjson":
		if bag, ok := argValue(ctx, 0, data.NewNullValue()).(*data.ClassValue); ok {
			requestState(cv).json = bag
		}
		return cv, nil
	case "getuserresolver":
		if resolver := requestState(cv).userResolver; resolver != nil {
			return resolver, nil
		}
		return data.NewNullValue(), nil
	case "setuserresolver":
		requestState(cv).userResolver = argValue(ctx, 0, data.NewNullValue())
		return cv, nil
	case "getrouteresolver":
		if resolver := requestState(cv).routeResolver; resolver != nil {
			return resolver, nil
		}
		return data.NewNullValue(), nil
	case "setrouteresolver":
		requestState(cv).routeResolver = argValue(ctx, 0, data.NewNullValue())
		return cv, nil
	case "toarray":
		return requestAll(cv), nil
	case "offsetexists", "__isset":
		_, ok := getNested(requestAllMap(cv), argString(ctx, 0, ""))
		return data.NewBoolValue(ok), nil
	case "offsetget", "__get":
		key := argString(ctx, 0, "")
		if value, ok := getNested(requestAllMap(cv), key); ok {
			return value, nil
		}
		return requestRoute(cv, &key, data.NewNullValue(), ctx)
	case "offsetset":
		input := requestInputSource(cv)
		setNested(input, argString(ctx, 0, ""), argValue(ctx, 1, data.NewNullValue()))
		setInputSource(cv, input)
		return data.NewNullValue(), nil
	case "offsetunset":
		input := requestInputSource(cv)
		deleteNested(input, argString(ctx, 0, ""))
		setInputSource(cv, input)
		return data.NewNullValue(), nil
	}
	return illuminateConcernDispatch(method, ctx, cv)
}

func illuminateConcernDispatch(method *illuminateRequestMethod, ctx data.Context, cv *data.ClassValue) (data.GetValue, data.Control) {
	name := strings.ToLower(method.name)
	switch name {
	case "server", "header", "query", "post", "cookie":
		source := map[string]string{"post": "request", "cookie": "cookies"}[name]
		if source == "" {
			source = name
			if name == "header" {
				source = "headers"
			}
		}
		return retrieveRequestItem(cv, source, argOptionalString(ctx, 0), argValue(ctx, 1, data.NewNullValue())), nil
	case "hasheader":
		return data.NewBoolValue(requestHeader(cv, argString(ctx, 0, "")) != ""), nil
	case "hascookie":
		_, ok := requestBagMap(cv, "cookies")[argString(ctx, 0, "")]
		return data.NewBoolValue(ok), nil
	case "bearertoken":
		header := requestHeader(cv, "Authorization")
		index := strings.LastIndex(strings.ToLower(header), "bearer ")
		if index < 0 {
			return data.NewNullValue(), nil
		}
		token := header[index+7:]
		if comma := strings.Index(token, ","); comma >= 0 {
			token = token[:comma]
		}
		return data.NewStringValue(token), nil
	case "keys":
		keys := make([]string, 0)
		for key := range requestInputMap(cv) {
			keys = append(keys, key)
		}
		for key := range requestBagMap(cv, "files") {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return stringArray(keys), nil
	case "all":
		all := requestAllMap(cv)
		keys := variadicStrings(ctx, 0)
		if len(keys) == 0 {
			return assocMapToArrayValue(all), nil
		}
		return assocMapToArrayValue(selectKeys(all, keys)), nil
	case "input":
		key := argOptionalString(ctx, 0)
		input := requestInputMap(cv)
		if key == nil {
			return assocMapToArrayValue(input), nil
		}
		if value, ok := getNested(input, *key); ok {
			return value, nil
		}
		return argValue(ctx, 1, data.NewNullValue()), nil
	case "fluent":
		if key := argOptionalString(ctx, 0); key != nil {
			if value, ok := getNested(requestInputMap(cv), *key); ok {
				return value, nil
			}
		}
		return assocMapToArrayValue(requestInputMap(cv)), nil
	case "allfiles":
		return assocMapToArrayValue(requestBagMap(cv, "files")), nil
	case "file":
		return retrieveRequestItem(cv, "files", argOptionalString(ctx, 0), argValue(ctx, 1, data.NewNullValue())), nil
	case "hasfile":
		value, ok := getNested(requestBagMap(cv, "files"), argString(ctx, 0, ""))
		return data.NewBoolValue(ok && !isNull(value)), nil
	case "image":
		return data.NewNullValue(), nil
	case "dump":
		return cv, nil
	case "exists", "has":
		return data.NewBoolValue(allKeysExist(requestAllMap(cv), variadicStrings(ctx, 0))), nil
	case "hasany":
		return data.NewBoolValue(anyKeyExists(requestAllMap(cv), variadicStrings(ctx, 0))), nil
	case "filled", "isnotfilled", "anyfilled":
		keys := variadicStrings(ctx, 0)
		filled := 0
		for _, key := range keys {
			if value, ok := getNested(requestAllMap(cv), key); ok && !emptyString(value) {
				filled++
			}
		}
		result := filled == len(keys)
		if name == "isnotfilled" {
			result = filled == 0
		} else if name == "anyfilled" {
			result = filled > 0
		}
		return data.NewBoolValue(result), nil
	case "missing":
		return data.NewBoolValue(!allKeysExist(requestAllMap(cv), variadicStrings(ctx, 0))), nil
	case "whenhas", "whenfilled", "whenmissing":
		key := argString(ctx, 0, "")
		value, exists := getNested(requestAllMap(cv), key)
		condition := exists
		if name == "whenfilled" {
			condition = exists && !emptyString(value)
		} else if name == "whenmissing" {
			condition = !exists
		}
		if condition {
			return callbackOrThis(argValue(ctx, 1, data.NewNullValue()), []data.Value{valueOrNull(value)}, cv, ctx)
		}
		if fallback := argValue(ctx, 2, data.NewNullValue()); !isNull(fallback) {
			return invokeCallable(fallback, nil, ctx)
		}
		return cv, nil
	case "whenenum":
		return callbackOrThis(argValue(ctx, 2, data.NewNullValue()), []data.Value{requestData(cv, argString(ctx, 0, ""), data.NewNullValue())}, cv, ctx)
	case "str", "string":
		return data.NewStringValue(requestData(cv, argString(ctx, 0, ""), argValue(ctx, 1, data.NewNullValue())).AsString()), nil
	case "boolean":
		return data.NewBoolValue(valueBool(requestData(cv, argString(ctx, 0, ""), argValue(ctx, 1, data.NewBoolValue(false))))), nil
	case "integer":
		return data.NewIntValue(valueInt(requestData(cv, argString(ctx, 0, ""), argValue(ctx, 1, data.NewIntValue(0))))), nil
	case "float":
		return data.NewFloatValue(valueFloat(requestData(cv, argString(ctx, 0, ""), argValue(ctx, 1, data.NewFloatValue(0))))), nil
	case "clamp":
		value := valueFloat(requestData(cv, argString(ctx, 0, ""), argValue(ctx, 3, data.NewIntValue(0))))
		minimum, maximum := valueFloat(argValue(ctx, 1, data.NewIntValue(0))), valueFloat(argValue(ctx, 2, data.NewIntValue(0)))
		if value < minimum {
			value = minimum
		}
		if value > maximum {
			value = maximum
		}
		return data.NewFloatValue(value), nil
	case "date", "interval":
		value := requestData(cv, argString(ctx, 0, ""), data.NewNullValue())
		if emptyString(value) {
			return data.NewNullValue(), nil
		}
		return value, nil
	case "enum":
		value := requestData(cv, argString(ctx, 0, ""), data.NewNullValue())
		if emptyString(value) {
			return argValue(ctx, 2, data.NewNullValue()), nil
		}
		return value, nil
	case "enums":
		value := requestData(cv, argString(ctx, 0, ""), data.NewNullValue())
		if arr, ok := value.(*data.ArrayValue); ok {
			return arr, nil
		}
		return data.NewArrayValue(nil), nil
	case "array", "collect":
		key := argOptionalString(ctx, 0)
		if key == nil {
			return requestAll(cv), nil
		}
		value := requestData(cv, *key, data.NewArrayValue(nil))
		if _, ok := value.(*data.ArrayValue); ok {
			return value, nil
		}
		if obj, ok := value.(*data.ObjectValue); ok {
			return obj, nil
		}
		return data.NewArrayValue([]data.Value{value}), nil
	case "only":
		return assocMapToArrayValue(selectKeys(requestAllMap(cv), variadicStrings(ctx, 0))), nil
	case "except":
		result := requestAllMap(cv)
		for _, key := range variadicStrings(ctx, 0) {
			deleteNested(result, key)
		}
		return assocMapToArrayValue(result), nil
	case "old":
		if requestState(cv).session == nil {
			return argValue(ctx, 1, data.NewNullValue()), nil
		}
		return invokeObject(requestState(cv).session, "getOldInput", []data.Value{
			argValue(ctx, 0, data.NewNullValue()), argValue(ctx, 1, data.NewNullValue()),
		}, ctx)
	case "flash", "flashonly", "flashexcept", "flush":
		return flashInput(name, cv, ctx)
	case "isjson":
		return data.NewBoolValue(isJSONRequest(cv)), nil
	case "expectsjson":
		ajax := strings.EqualFold(requestHeader(cv, "X-Requested-With"), "XMLHttpRequest")
		return data.NewBoolValue((ajax && requestHeader(cv, "X-PJAX") == "" && acceptsAny(cv)) || wantsJSON(cv)), nil
	case "wantsjson":
		return data.NewBoolValue(wantsJSON(cv)), nil
	case "wantsmarkdown":
		types := acceptableTypes(cv)
		return data.NewBoolValue(len(types) > 0 && strings.HasPrefix(strings.ToLower(types[0]), "text/markdown")), nil
	case "accepts":
		return data.NewBoolValue(requestAccepts(cv, valueStrings(argValue(ctx, 0, data.NewNullValue())))), nil
	case "prefers":
		for _, accepted := range acceptableTypes(cv) {
			for _, offered := range valueStrings(argValue(ctx, 0, data.NewNullValue())) {
				if matchesContentType(accepted, offered) || accepted == "*/*" || accepted == "*" {
					return data.NewStringValue(offered), nil
				}
			}
		}
		return data.NewNullValue(), nil
	case "acceptsanycontenttype":
		return data.NewBoolValue(acceptsAny(cv)), nil
	case "acceptsjson":
		return data.NewBoolValue(requestAccepts(cv, []string{"application/json"})), nil
	case "acceptsmarkdown":
		return data.NewBoolValue(requestAccepts(cv, []string{"text/markdown"})), nil
	case "acceptshtml":
		return data.NewBoolValue(requestAccepts(cv, []string{"text/html"})), nil
	case "format":
		for _, accepted := range acceptableTypes(cv) {
			if extensions, _ := mime.ExtensionsByType(accepted); len(extensions) > 0 {
				return data.NewStringValue(strings.TrimPrefix(extensions[0], ".")), nil
			}
			if accepted == "application/json" {
				return data.NewStringValue("json"), nil
			}
		}
		return data.NewStringValue(argString(ctx, 0, "html")), nil
	case "filterprecognitiverules":
		rules := argMap(ctx, 0)
		header := requestHeader(cv, "Precognition-Validate-Only")
		if header == "" {
			return assocMapToArrayValue(rules), nil
		}
		patterns := strings.Split(header, ",")
		out := map[string]data.Value{}
		for attribute, rule := range rules {
			if anyPattern(attribute, patterns) {
				out[attribute] = rule
			}
		}
		return assocMapToArrayValue(out), nil
	case "isattemptingprecognition":
		return data.NewBoolValue(requestHeader(cv, "Precognition") == "true"), nil
	case "isprecognitive":
		value, ok := requestBagMap(cv, "attributes")["precognitive"]
		return data.NewBoolValue(ok && valueBool(value)), nil
	}
	return nil, requestError("Request::%s is not implemented", method.name)
}

func duplicateIlluminateRequest(cv *data.ClassValue, ctx data.Context) (data.GetValue, data.Control) {
	target := data.NewClassValue(NewIlluminateRequestClass(), cv.Context.CreateBaseContext())
	initializeIlluminateRequest(target, requestState(cv).source)
	for name, value := range cv.GetProperties() {
		_ = target.SetProperty(name, value)
	}
	for index, property := range []string{"query", "request", "attributes", "cookies", "files", "server"} {
		if value, ok := ctx.GetIndexValue(index); ok && !isNull(value) {
			if values, err := valueToAssocMap(value); err == nil {
				setRequestBag(target, property, values)
			}
		}
	}
	old, state := requestState(cv), requestState(target)
	state.userResolver, state.routeResolver, state.session = old.userResolver, old.routeResolver, old.session
	state.locale, state.defaultLocale, state.json = old.locale, old.defaultLocale, old.json
	return target, nil
}

func requestJSON(cv *data.ClassValue, key *string, fallback data.Value) (data.GetValue, data.Control) {
	state := requestState(cv)
	if state.json == nil {
		values := map[string]data.Value{}
		content := requestContent(cv)
		var decoded any
		if strings.TrimSpace(content) != "" && json.Unmarshal([]byte(content), &decoded) == nil {
			values = goMapToValues(decoded)
		}
		state.json = NewInputBagValue(cv, values)
	}
	if key == nil {
		return state.json, nil
	}
	if value, ok := getNested(GetParamBagMap(state.json), *key); ok {
		return value, nil
	}
	return fallback, nil
}

func requestFingerprint(cv *data.ClassValue, ctx data.Context) (data.GetValue, data.Control) {
	route, control := requestRoute(cv, nil, data.NewNullValue(), ctx)
	if control != nil {
		return nil, control
	}
	if isNullValue(route) {
		return nil, requestError("Unable to generate fingerprint. Route unavailable.")
	}
	parts := []string{requestIP(cv)}
	for _, method := range []string{"getDomain", "uri"} {
		value, control := invokeObject(route, method, nil, ctx)
		if control == nil && value != nil {
			if scalar, ok := value.(data.Value); ok {
				parts = append(parts, scalar.AsString())
			}
		}
	}
	sum := sha1.Sum([]byte(strings.Join(parts, "|")))
	return data.NewStringValue(hex.EncodeToString(sum[:])), nil
}

func requestRoute(cv *data.ClassValue, parameter *string, fallback data.Value, ctx data.Context) (data.GetValue, data.Control) {
	resolver := requestState(cv).routeResolver
	if resolver == nil || isNull(resolver) {
		return data.NewNullValue(), nil
	}
	route, control := invokeCallable(resolver, nil, ctx)
	if control != nil || parameter == nil || isNullValue(route) {
		return route, control
	}
	return invokeObject(route, "parameter", []data.Value{data.NewStringValue(*parameter), fallback}, ctx)
}

func flashInput(name string, cv *data.ClassValue, ctx data.Context) (data.GetValue, data.Control) {
	session := requestState(cv).session
	if session == nil || isNull(session) {
		return nil, requestError("Session store not set on request.")
	}
	input := requestInputMap(cv)
	switch name {
	case "flashonly":
		input = selectKeys(input, variadicStrings(ctx, 0))
	case "flashexcept":
		for _, key := range variadicStrings(ctx, 0) {
			deleteNested(input, key)
		}
	case "flush":
		input = map[string]data.Value{}
	}
	_, control := invokeObject(session, "flashInput", []data.Value{assocMapToArrayValue(input)}, ctx)
	return data.NewNullValue(), control
}

func invokeCallable(callable data.Value, args []data.Value, ctx data.Context) (data.GetValue, data.Control) {
	switch fn := callable.(type) {
	case *data.FuncValue:
		callCtx := ctx.CreateContext(fn.Value.GetVariables())
		for i, value := range args {
			_ = callCtx.SetVariableValue(data.NewVariable("", i, nil), value)
		}
		return fn.Call(callCtx)
	case *data.BoundFuncValue:
		callCtx := ctx.CreateContext(fn.Value.GetVariables())
		for i, value := range args {
			_ = callCtx.SetVariableValue(data.NewVariable("", i, nil), value)
		}
		return fn.Call(callCtx)
	default:
		return data.NewNullValue(), nil
	}
}

func invokeObject(object data.GetValue, methodName string, args []data.Value, ctx data.Context) (data.GetValue, data.Control) {
	cv, ok := object.(*data.ClassValue)
	if !ok {
		return data.NewNullValue(), nil
	}
	method, ok := cv.GetMethod(methodName)
	if !ok {
		return data.NewNullValue(), nil
	}
	callCtx := cv.CreateContext(method.GetVariables())
	variables := method.GetVariables()
	for i, value := range args {
		if i < len(variables) {
			_ = callCtx.SetVariableValue(variables[i], value)
		}
	}
	return method.Call(callCtx)
}

func callbackOrThis(callback data.Value, args []data.Value, cv *data.ClassValue, ctx data.Context) (data.GetValue, data.Control) {
	result, control := invokeCallable(callback, args, ctx)
	if control != nil {
		return nil, control
	}
	if result == nil || isNullValue(result) {
		return cv, nil
	}
	return result, nil
}

func requestBag(cv *data.ClassValue, name string) *data.ClassValue {
	value, _ := cv.GetProperty(name)
	bag, _ := value.(*data.ClassValue)
	return bag
}

func requestBagMap(cv *data.ClassValue, name string) map[string]data.Value {
	return GetParamBagMap(requestBag(cv, name))
}

func setRequestBag(cv *data.ClassValue, name string, values map[string]data.Value) {
	bag := requestBag(cv, name)
	if bag != nil {
		SetParamBagMap(bag, values)
		return
	}
	switch name {
	case "query", "request", "cookies":
		_ = cv.SetProperty(name, NewInputBagValue(cv, values))
	case "attributes":
		_ = cv.SetProperty(name, NewParameterBagValue(cv, values))
	case "files":
		_ = cv.SetProperty(name, NewFileBagValue(cv, values))
	case "server":
		_ = cv.SetProperty(name, NewServerBagValue(cv, values))
	}
}

func retrieveRequestItem(cv *data.ClassValue, source string, key *string, fallback data.Value) data.GetValue {
	if source == "headers" {
		headers := GetHeaderBagAll(requestBag(cv, "headers"))
		if key == nil {
			out := map[string]data.Value{}
			for name, values := range headers {
				out[name] = stringsToValue(values)
			}
			return assocMapToArrayValue(out)
		}
		for name, values := range headers {
			if strings.EqualFold(name, *key) {
				if len(values) == 0 {
					return fallback
				}
				return data.NewStringValue(values[0])
			}
		}
		return fallback
	}
	values := requestBagMap(cv, source)
	if key == nil {
		return assocMapToArrayValue(values)
	}
	if value, ok := getNested(values, *key); ok {
		return value
	}
	return fallback
}

func requestInputSource(cv *data.ClassValue) map[string]data.Value {
	if isJSONRequest(cv) {
		if jsonBag, _ := requestJSON(cv, nil, data.NewNullValue()); jsonBag != nil {
			return GetParamBagMap(jsonBag.(*data.ClassValue))
		}
	}
	method := strings.ToUpper(requestMethod(cv))
	if method == "GET" || method == "HEAD" {
		return requestBagMap(cv, "query")
	}
	return requestBagMap(cv, "request")
}

func setInputSource(cv *data.ClassValue, values map[string]data.Value) {
	if isJSONRequest(cv) {
		bag := NewInputBagValue(cv, values)
		requestState(cv).json = bag
		return
	}
	method := strings.ToUpper(requestMethod(cv))
	if method == "GET" || method == "HEAD" {
		setRequestBag(cv, "query", values)
	} else {
		setRequestBag(cv, "request", values)
	}
}

func requestInputMap(cv *data.ClassValue) map[string]data.Value {
	out := cloneMap(requestBagMap(cv, "query"))
	for key, value := range requestInputSource(cv) {
		if _, exists := out[key]; !exists {
			out[key] = value
		}
	}
	return out
}

func requestAllMap(cv *data.ClassValue) map[string]data.Value {
	out := requestInputMap(cv)
	for key, value := range requestBagMap(cv, "files") {
		out[key] = value
	}
	return out
}

func requestAll(cv *data.ClassValue) data.GetValue { return assocMapToArrayValue(requestAllMap(cv)) }
func requestData(cv *data.ClassValue, key string, fallback data.Value) data.Value {
	if value, ok := getNested(requestInputMap(cv), key); ok {
		return value
	}
	return fallback
}

func requestMethod(cv *data.ClassValue) string {
	if r, ok := requestHTTPSource(cv); ok && r.Method != "" {
		return r.Method
	}
	if value, ok := requestBagMap(cv, "server")["REQUEST_METHOD"]; ok {
		return value.AsString()
	}
	return "GET"
}

func requestDecodedPath(cv *data.ClassValue) string {
	decoded, err := url.PathUnescape(strings.Trim(requestPath(cv), "/"))
	if err != nil || decoded == "" {
		return "/"
	}
	return decoded
}

func requestSegments(cv *data.ClassValue) []string {
	raw := strings.Trim(requestDecodedPath(cv), "/")
	if raw == "" {
		return nil
	}
	return strings.Split(raw, "/")
}

func requestSchemeHost(cv *data.ClassValue) string {
	scheme := "http"
	if requestSecure(cv) {
		scheme = "https"
	}
	return scheme + "://" + requestHTTPHost(cv)
}

func requestURL(cv *data.ClassValue) string {
	return requestSchemeHost(cv) + requestBaseURL(cv) + requestPath(cv)
}

func requestFullURL(cv *data.ClassValue) string {
	raw := requestURL(cv)
	if r, ok := requestHTTPSource(cv); ok && r.URL != nil && r.URL.RawQuery != "" {
		return raw + "?" + r.URL.RawQuery
	}
	query := requestBagMap(cv, "query")
	if len(query) > 0 {
		return raw + "?" + encodeValues(query)
	}
	return raw
}

func fullURLWithQuery(cv *data.ClassValue, additions map[string]data.Value) string {
	query := requestBagMap(cv, "query")
	for key, value := range additions {
		query[key] = value
	}
	if len(query) == 0 {
		return requestURL(cv)
	}
	return requestURL(cv) + "?" + encodeValues(query)
}

func fullURLWithoutQuery(cv *data.ClassValue, keys []string) string {
	query := requestBagMap(cv, "query")
	for _, key := range keys {
		delete(query, key)
	}
	if len(query) == 0 {
		return requestURL(cv)
	}
	return requestURL(cv) + "?" + encodeValues(query)
}

func requestContent(cv *data.ClassValue) string {
	if value, _ := cv.GetProperty("content"); value != nil && !isNull(value) {
		return value.AsString()
	}
	if r, ok := requestHTTPSource(cv); ok && r.Body != nil {
		// Symfony's source helper is expected to cache content. Avoid consuming an
		// uncached Go body here because doing so would alter subsequent readers.
		if getter, ok := any(r.Body).(interface{ String() string }); ok {
			return getter.String()
		}
	}
	return ""
}

func requestIP(cv *data.ClassValue) string {
	ips := requestIPs(cv)
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

func requestIPs(cv *data.ClassValue) []string {
	if forwarded := requestHeader(cv, "X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	if r, ok := requestHTTPSource(cv); ok {
		host := r.RemoteAddr
		if parsed, _, err := strings.Cut(host, ":"); err {
			host = parsed
		}
		if host != "" {
			return []string{host}
		}
	}
	return nil
}

func acceptableTypes(cv *data.ClassValue) []string {
	header := requestHeader(cv, "Accept")
	if strings.TrimSpace(header) == "" {
		return nil
	}
	type accept struct {
		value string
		q     float64
		order int
	}
	items := make([]accept, 0)
	for order, part := range strings.Split(header, ",") {
		segments := strings.Split(strings.TrimSpace(part), ";")
		item := accept{value: strings.TrimSpace(segments[0]), q: 1, order: order}
		for _, segment := range segments[1:] {
			if strings.HasPrefix(strings.TrimSpace(segment), "q=") {
				item.q, _ = strconv.ParseFloat(strings.TrimPrefix(strings.TrimSpace(segment), "q="), 64)
			}
		}
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].q > items[j].q })
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.value
	}
	return out
}

func requestAccepts(cv *data.ClassValue, offered []string) bool {
	accepted := acceptableTypes(cv)
	if len(accepted) == 0 {
		return true
	}
	for _, actual := range accepted {
		actual = strings.TrimSpace(strings.Split(actual, ";")[0])
		for _, typ := range offered {
			if actual == "*" || actual == "*/*" || matchesContentType(actual, typ) {
				return true
			}
		}
	}
	return false
}

func acceptsAny(cv *data.ClassValue) bool {
	types := acceptableTypes(cv)
	return len(types) == 0 || types[0] == "*" || types[0] == "*/*"
}

func wantsJSON(cv *data.ClassValue) bool {
	types := acceptableTypes(cv)
	return len(types) > 0 && (strings.Contains(strings.ToLower(types[0]), "/json") ||
		strings.Contains(strings.ToLower(types[0]), "+json"))
}

func isJSONRequest(cv *data.ClassValue) bool {
	contentType := strings.ToLower(requestHeader(cv, "CONTENT_TYPE"))
	return strings.Contains(contentType, "/json") || strings.Contains(contentType, "+json")
}

func matchesContentType(actual, typ string) bool {
	actual = strings.ToLower(strings.TrimSpace(strings.Split(actual, ";")[0]))
	typ = strings.ToLower(strings.TrimSpace(strings.Split(typ, ";")[0]))
	if actual == typ || actual == "*" || actual == "*/*" {
		return true
	}
	if strings.HasSuffix(actual, "/*") {
		return strings.HasPrefix(typ, strings.TrimSuffix(actual, "*"))
	}
	actualParts, typeParts := strings.Split(actual, "/"), strings.Split(typ, "/")
	if len(actualParts) != 2 || len(typeParts) != 2 {
		return false
	}
	return actualParts[0] == typeParts[0] && strings.HasSuffix(typeParts[1], "+"+actualParts[1])
}

func anyPattern(value string, patterns []string) bool {
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		quoted := regexp.QuoteMeta(pattern)
		quoted = strings.ReplaceAll(quoted, `\*`, ".*")
		if matched, _ := regexp.MatchString("^"+quoted+"$", value); matched {
			return true
		}
	}
	return false
}

func getNested(values map[string]data.Value, dotted string) (data.Value, bool) {
	if value, ok := values[dotted]; ok {
		return value, true
	}
	parts := strings.Split(dotted, ".")
	var current data.Value = assocMapToArrayValue(values)
	for _, part := range parts {
		switch container := current.(type) {
		case *data.ArrayValue:
			found := false
			for index, item := range container.List {
				if item == nil {
					continue
				}
				key := item.Name
				if key == "" {
					key = strconv.Itoa(index)
				}
				if key == part {
					current, found = item.Value, true
					break
				}
			}
			if !found {
				return nil, false
			}
		case *data.ObjectValue:
			value, control := container.GetProperty(part)
			if control != nil {
				return nil, false
			}
			current = value
		default:
			return nil, false
		}
	}
	return current, true
}

func setNested(values map[string]data.Value, dotted string, value data.Value) {
	// Preserve Laravel's dotted-key contract for top-level request consumers.
	// Nested arrays are resolved by getNested; retaining the dotted alias avoids
	// lossy conversion between ObjectValue and PHP arrays.
	values[dotted] = value
}

func deleteNested(values map[string]data.Value, dotted string) { delete(values, dotted) }

func selectKeys(values map[string]data.Value, keys []string) map[string]data.Value {
	out := map[string]data.Value{}
	for _, key := range keys {
		if value, ok := getNested(values, key); ok {
			out[key] = value
		}
	}
	return out
}

func cloneMap(values map[string]data.Value) map[string]data.Value {
	out := make(map[string]data.Value, len(values))
	for key, value := range values {
		out[key] = value
	}
	return out
}

func allKeysExist(values map[string]data.Value, keys []string) bool {
	if len(keys) == 0 {
		return false
	}
	for _, key := range keys {
		if _, ok := getNested(values, key); !ok {
			return false
		}
	}
	return true
}

func anyKeyExists(values map[string]data.Value, keys []string) bool {
	for _, key := range keys {
		if _, ok := getNested(values, key); ok {
			return true
		}
	}
	return false
}

func encodeValues(values map[string]data.Value) string {
	query := url.Values{}
	for key, value := range values {
		for _, stringValue := range valueStrings(value) {
			query.Add(key, stringValue)
		}
	}
	return query.Encode()
}

func goMapToValues(value any) map[string]data.Value {
	out := map[string]data.Value{}
	values, ok := value.(map[string]any)
	if !ok {
		return out
	}
	for key, item := range values {
		out[key] = goValue(item)
	}
	return out
}

func goValue(value any) data.Value {
	switch item := value.(type) {
	case nil:
		return data.NewNullValue()
	case string:
		return data.NewStringValue(item)
	case bool:
		return data.NewBoolValue(item)
	case float64:
		return data.NewFloatValue(item)
	case []any:
		out := make([]data.Value, len(item))
		for i, value := range item {
			out[i] = goValue(value)
		}
		return data.NewArrayValue(out)
	case map[string]any:
		return assocMapToArrayValue(goMapToValues(item))
	default:
		return data.NewAnyValue(item)
	}
}

func argValue(ctx data.Context, index int, fallback data.Value) data.Value {
	if value, ok := ctx.GetIndexValue(index); ok && value != nil && !isNull(value) {
		return value
	}
	return fallback
}

func argString(ctx data.Context, index int, fallback string) string {
	if value, ok := ctx.GetIndexValue(index); ok && value != nil && !isNull(value) {
		return value.AsString()
	}
	return fallback
}

func argOptionalString(ctx data.Context, index int) *string {
	if value, ok := ctx.GetIndexValue(index); ok && value != nil && !isNull(value) {
		result := value.AsString()
		return &result
	}
	return nil
}

func argInt(ctx data.Context, index, fallback int) int {
	return valueInt(argValue(ctx, index, data.NewIntValue(fallback)))
}

func argMap(ctx data.Context, index int) map[string]data.Value {
	value, ok := ctx.GetIndexValue(index)
	if !ok || isNull(value) {
		return map[string]data.Value{}
	}
	result, err := valueToAssocMap(value)
	if err != nil {
		return map[string]data.Value{}
	}
	return result
}

func argumentKeys(ctx data.Context, index int) []string {
	value, ok := ctx.GetIndexValue(index)
	if !ok {
		return nil
	}
	return valueStrings(value)
}

func variadicStrings(ctx data.Context, index int) []string {
	value, ok := ctx.GetIndexValue(index)
	if !ok || isNull(value) {
		return nil
	}
	return valueStrings(value)
}

func valuesFromStrings(values []string) []data.Value {
	out := make([]data.Value, len(values))
	for i, value := range values {
		out[i] = data.NewStringValue(value)
	}
	return out
}

func stringArray(values []string) data.GetValue {
	out := make([]data.Value, len(values))
	for i, value := range values {
		out[i] = data.NewStringValue(value)
	}
	return data.NewArrayValue(out)
}

func nullableString(value string) data.GetValue {
	if value == "" {
		return data.NewNullValue()
	}
	return data.NewStringValue(value)
}

func isNullValue(value data.GetValue) bool {
	if value == nil {
		return true
	}
	_, ok := value.(*data.NullValue)
	return ok
}

func valueOrNull(value data.Value) data.Value {
	if value == nil {
		return data.NewNullValue()
	}
	return value
}

func emptyString(value data.Value) bool {
	if value == nil || isNull(value) {
		return true
	}
	switch value.(type) {
	case *data.BoolValue, *data.ArrayValue, *data.ObjectValue:
		return false
	}
	return strings.TrimSpace(value.AsString()) == ""
}

func valueBool(value data.Value) bool {
	if value == nil || isNull(value) {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(value.AsString())) {
	case "1", "true", "on", "yes":
		return true
	}
	if boolean, ok := value.(interface{ AsBool() (bool, error) }); ok {
		result, _ := boolean.AsBool()
		return result
	}
	return false
}

func valueInt(value data.Value) int {
	if integer, ok := value.(interface{ AsInt() (int, error) }); ok {
		result, _ := integer.AsInt()
		return result
	}
	result, _ := strconv.Atoi(value.AsString())
	return result
}

func valueFloat(value data.Value) float64 {
	if number, ok := value.(interface{ AsFloat() (float64, error) }); ok {
		result, _ := number.AsFloat()
		return result
	}
	result, _ := strconv.ParseFloat(value.AsString(), 64)
	return result
}

func requestError(format string, args ...any) data.Control {
	return data.NewErrorThrow(nil, fmt.Errorf(format, args...))
}

var _ data.ClassStmt = (*IlluminateRequestClass)(nil)
var _ data.GetStaticMethod = (*IlluminateRequestClass)(nil)

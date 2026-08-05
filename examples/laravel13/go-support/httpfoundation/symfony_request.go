package httpfoundation

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const symfonyRequestFQN = "Symfony\\Component\\HttpFoundation\\Request"

var requestStaticState = struct {
	sync.RWMutex
	trustedProxies        []string
	trustedHosts          []string
	trustedHeaderSet      int
	methodOverride        bool
	allowedMethodOverride []string
	formats               map[string][]string
}{
	trustedHeaderSet: -1,
	formats: map[string][]string{
		"html":    {"text/html", "application/xhtml+xml"},
		"txt":     {"text/plain"},
		"js":      {"application/javascript", "application/x-javascript", "text/javascript"},
		"css":     {"text/css"},
		"json":    {"application/json", "application/x-json"},
		"jsonld":  {"application/ld+json"},
		"xml":     {"text/xml", "application/xml", "application/x-xml"},
		"rdf":     {"application/rdf+xml"},
		"atom":    {"application/atom+xml"},
		"rss":     {"application/rss+xml"},
		"form":    {"application/x-www-form-urlencoded", "multipart/form-data"},
		"soap":    {"application/soap+xml"},
		"problem": {"application/problem+json"},
		"hal":     {"application/hal+json", "application/hal+xml"},
		"jsonapi": {"application/vnd.api+json"},
		"yaml":    {"text/yaml", "application/x-yaml"},
		"wbxml":   {"application/vnd.wap.wbxml"},
		"pdf":     {"application/pdf"},
		"csv":     {"text/csv"},
	},
}

type symfonyRequestSource struct {
	request       *http.Request
	content       string
	requestFormat *string
	defaultLocale string
	locale        string
	session       data.Value
}

// SymfonyRequestClass 是 Symfony HttpFoundation Request 的 Go 实现。
type SymfonyRequestClass struct {
	node.Node
	source *symfonyRequestSource
}

// NewSymfonyRequestClass 返回可注册到 VM 的 Request 类。
func NewSymfonyRequestClass() data.ClassStmt {
	return &SymfonyRequestClass{}
}

// NewSymfonyRequestClassFrom 返回封装指定 net/http 请求的类结构。
// 同包 Illuminate Request 可复用该构造路径。
func NewSymfonyRequestClassFrom(request *http.Request) *SymfonyRequestClass {
	return &SymfonyRequestClass{source: &symfonyRequestSource{
		request:       cloneHTTPRequest(request),
		content:       bodyFromRequest(request),
		defaultLocale: "en",
		locale:        "en",
	}}
}

// NewSymfonyRequestValue 创建一个已从 net/http.Request 初始化的 PHP Request 对象。
func NewSymfonyRequestValue(ctx data.Context, request *http.Request) (*data.ClassValue, data.Control) {
	return newRequestValue(ctx, request, nil, nil, nil, nil, nil, nil)
}

// RequestSource 返回 Request 对象封装的 Go 请求，供同包 Illuminate 适配层复用。
func RequestSource(value *data.ClassValue) *http.Request {
	if value == nil {
		return nil
	}
	if class, ok := value.Class.(*SymfonyRequestClass); ok && class.source != nil {
		return class.source.request
	}
	if source, ok := value.GetSource().(*http.Request); ok {
		return source
	}
	return nil
}

func (class *SymfonyRequestClass) GetName() string           { return symfonyRequestFQN }
func (class *SymfonyRequestClass) GetExtend() *string        { return nil }
func (class *SymfonyRequestClass) GetImplements() []string   { return nil }
func (class *SymfonyRequestClass) GetConstruct() data.Method { return newRequestMethod("__construct") }
func (class *SymfonyRequestClass) GetSource() any {
	if class.source == nil {
		return nil
	}
	return class.source.request
}

func (class *SymfonyRequestClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	instanceClass := NewSymfonyRequestClassFrom(nil)
	return data.NewProxyValue(instanceClass, ctx.CreateBaseContext()), nil
}

func (class *SymfonyRequestClass) GetProperty(name string) (data.Property, bool) {
	for _, property := range class.GetPropertyList() {
		if property.GetName() == name {
			return property, true
		}
	}
	return nil, false
}

func (class *SymfonyRequestClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "attributes", "public", false, data.NewNullValue()),
		node.NewProperty(nil, "request", "public", false, data.NewNullValue()),
		node.NewProperty(nil, "query", "public", false, data.NewNullValue()),
		node.NewProperty(nil, "server", "public", false, data.NewNullValue()),
		node.NewProperty(nil, "files", "public", false, data.NewNullValue()),
		node.NewProperty(nil, "cookies", "public", false, data.NewNullValue()),
		node.NewProperty(nil, "headers", "public", false, data.NewNullValue()),
	}
}

var requestMethodNames = []string{
	"__construct", "initialize", "createFromGlobals", "create", "setFactory", "duplicate", "__clone",
	"__toString", "overrideGlobals", "setTrustedProxies", "getTrustedProxies", "getTrustedHeaderSet",
	"setTrustedHosts", "getTrustedHosts", "normalizeQueryString", "enableHttpMethodParameterOverride",
	"getHttpMethodParameterOverride", "setAllowedHttpMethodOverride", "getAllowedHttpMethodOverride",
	"getSession", "hasPreviousSession", "hasSession", "setSession", "setSessionFactory", "getClientIps",
	"getClientIp", "getScriptName", "getPathInfo", "getBasePath", "getBaseUrl", "getScheme", "getPort",
	"getUser", "getPassword", "getUserInfo", "getHttpHost", "getRequestUri", "getSchemeAndHttpHost",
	"getUri", "getUriForPath", "getRelativeUriForPath", "getQueryString", "isSecure", "getHost",
	"setMethod", "getMethod", "getRealMethod", "getMimeType", "getMimeTypes", "getFormat", "setFormat",
	"getRequestFormat", "setRequestFormat", "getContentTypeFormat", "setDefaultLocale", "getDefaultLocale",
	"setLocale", "getLocale", "isMethod", "isMethodSafe", "isMethodIdempotent", "isMethodCacheable",
	"getProtocolVersion", "getContent", "getPayload", "toArray", "getETags", "isNoCache",
	"getPreferredFormat", "getPreferredLanguage", "getLanguages", "getCharsets", "getEncodings",
	"getAcceptableContentTypes", "isXmlHttpRequest", "preferSafeContent", "isFromTrustedProxy",
}

func (class *SymfonyRequestClass) GetMethod(name string) (data.Method, bool) {
	for _, methodName := range requestMethodNames {
		if methodName == name {
			return newRequestMethod(name), true
		}
	}
	return nil, false
}

func (class *SymfonyRequestClass) GetMethods() []data.Method {
	methods := make([]data.Method, len(requestMethodNames))
	for index, name := range requestMethodNames {
		methods[index] = newRequestMethod(name)
	}
	return methods
}

func (class *SymfonyRequestClass) GetStaticMethod(name string) (data.Method, bool) {
	method, ok := class.GetMethod(name)
	if !ok || !method.GetIsStatic() {
		return nil, false
	}
	return method, true
}

func (class *SymfonyRequestClass) GetStaticProperty(name string) (data.Value, bool) {
	constants := map[string]data.Value{
		"HEADER_FORWARDED":           data.NewIntValue(1),
		"HEADER_X_FORWARDED_FOR":     data.NewIntValue(2),
		"HEADER_X_FORWARDED_HOST":    data.NewIntValue(4),
		"HEADER_X_FORWARDED_PROTO":   data.NewIntValue(8),
		"HEADER_X_FORWARDED_PORT":    data.NewIntValue(16),
		"HEADER_X_FORWARDED_PREFIX":  data.NewIntValue(32),
		"HEADER_X_FORWARDED_AWS_ELB": data.NewIntValue(26),
		"HEADER_X_FORWARDED_TRAEFIK": data.NewIntValue(62),
		"METHOD_HEAD":                data.NewStringValue("HEAD"),
		"METHOD_GET":                 data.NewStringValue("GET"),
		"METHOD_POST":                data.NewStringValue("POST"),
		"METHOD_PUT":                 data.NewStringValue("PUT"),
		"METHOD_PATCH":               data.NewStringValue("PATCH"),
		"METHOD_DELETE":              data.NewStringValue("DELETE"),
		"METHOD_PURGE":               data.NewStringValue("PURGE"),
		"METHOD_OPTIONS":             data.NewStringValue("OPTIONS"),
		"METHOD_TRACE":               data.NewStringValue("TRACE"),
		"METHOD_CONNECT":             data.NewStringValue("CONNECT"),
		"METHOD_QUERY":               data.NewStringValue("QUERY"),
	}
	value, ok := constants[name]
	return value, ok
}

type symfonyRequestMethod struct {
	name      string
	static    bool
	variables []data.Variable
	params    []data.GetValue
}

func newRequestMethod(name string) *symfonyRequestMethod {
	names := requestParameterNames(name)
	method := &symfonyRequestMethod{name: name, static: requestStaticMethod(name)}
	for index, variableName := range names {
		variable := node.NewVariable(nil, variableName, index, nil)
		method.variables = append(method.variables, variable)
		method.params = append(method.params, node.NewParameter(nil, variableName, index, requestParameterDefault(name, index), nil))
	}
	return method
}

func requestParameterDefault(name string, index int) data.GetValue {
	switch name {
	case "__construct", "initialize":
		if index < 6 {
			return data.NewArrayValue(nil)
		}
		return data.NewNullValue()
	case "create":
		switch index {
		case 0:
			return nil
		case 1:
			return data.NewStringValue("GET")
		case 2, 3, 4, 5:
			return data.NewArrayValue(nil)
		default:
			return data.NewNullValue()
		}
	case "duplicate":
		return data.NewNullValue()
	case "hasSession":
		return data.NewBoolValue(false)
	case "getFormat":
		if index == 1 {
			return data.NewBoolValue(false)
		}
	case "getRequestFormat", "getPreferredFormat":
		return data.NewStringValue("html")
	case "getContent":
		return data.NewBoolValue(false)
	case "getPreferredLanguage":
		return data.NewNullValue()
	}
	return nil
}

func requestStaticMethod(name string) bool {
	switch name {
	case "createFromGlobals", "create", "setFactory", "setTrustedProxies", "getTrustedProxies",
		"getTrustedHeaderSet", "setTrustedHosts", "getTrustedHosts", "normalizeQueryString",
		"enableHttpMethodParameterOverride", "getHttpMethodParameterOverride",
		"setAllowedHttpMethodOverride", "getAllowedHttpMethodOverride", "getMimeTypes":
		return true
	}
	return false
}

func requestParameterNames(name string) []string {
	switch name {
	case "__construct", "initialize":
		return []string{"query", "request", "attributes", "cookies", "files", "server", "content"}
	case "create":
		return []string{"uri", "method", "parameters", "cookies", "files", "server", "content"}
	case "duplicate":
		return []string{"query", "request", "attributes", "cookies", "files", "server"}
	case "setTrustedProxies":
		return []string{"proxies", "trustedHeaderSet"}
	case "setTrustedHosts":
		return []string{"hostPatterns"}
	case "normalizeQueryString":
		return []string{"qs"}
	case "setAllowedHttpMethodOverride":
		return []string{"methods"}
	case "hasSession":
		return []string{"skipIfUninitialized"}
	case "setSession":
		return []string{"session"}
	case "setSessionFactory", "setFactory":
		return []string{"factory"}
	case "getUriForPath", "getRelativeUriForPath":
		return []string{"path"}
	case "setMethod", "isMethod":
		return []string{"method"}
	case "getMimeType", "getMimeTypes":
		return []string{"format"}
	case "getFormat":
		return []string{"mimeType", "subtypeFallback"}
	case "setFormat":
		return []string{"format", "mimeTypes"}
	case "getRequestFormat", "getPreferredFormat":
		return []string{"default"}
	case "setRequestFormat":
		return []string{"format"}
	case "setDefaultLocale", "setLocale":
		return []string{"locale"}
	case "getContent":
		return []string{"asResource"}
	case "getPreferredLanguage":
		return []string{"locales"}
	}
	return nil
}

func (method *symfonyRequestMethod) GetName() string               { return method.name }
func (method *symfonyRequestMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (method *symfonyRequestMethod) GetIsStatic() bool             { return method.static }
func (method *symfonyRequestMethod) GetParams() []data.GetValue    { return method.params }
func (method *symfonyRequestMethod) GetVariables() []data.Variable { return method.variables }
func (method *symfonyRequestMethod) GetReturnType() data.Types     { return nil }

func requestClassValue(ctx data.Context) *data.ClassValue {
	if methodCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return methodCtx.ClassValue
	}
	if value, ok := ctx.(*data.ClassValue); ok {
		return value
	}
	return nil
}

func requestSource(ctx data.Context) (*symfonyRequestSource, data.Control) {
	value := requestClassValue(ctx)
	if value == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Request method requires an object"))
	}
	if class, ok := value.Class.(*SymfonyRequestClass); ok && class.source != nil {
		return class.source, nil
	}
	if _, ok := value.Class.(*IlluminateRequestClass); ok {
		if state := requestStateNoCreate(value); state != nil && state.source != nil {
			return &symfonyRequestSource{
				request:       state.source,
				content:       bodyFromRequest(state.source),
				defaultLocale: state.defaultLocale,
				locale:        state.locale,
				session:       state.session,
			}, nil
		}
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Request object is not initialized"))
}

func (method *symfonyRequestMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	switch method.name {
	case "__construct", "initialize":
		return requestInitialize(ctx)
	case "createFromGlobals":
		return newRequestValue(ctx, ActiveRequestFromContext(ctx), nil, nil, nil, nil, nil, nil)
	case "create":
		return requestCreate(ctx)
	case "setFactory":
		return data.NewNullValue(), nil
	case "duplicate":
		return requestDuplicate(ctx)
	case "__clone":
		return data.NewNullValue(), nil
	case "__toString":
		return requestToString(ctx)
	case "overrideGlobals":
		return data.NewNullValue(), nil
	case "setTrustedProxies":
		requestStaticState.Lock()
		requestStaticState.trustedProxies = arrayStrings(indexValue(ctx, 0))
		requestStaticState.trustedHeaderSet = intParam(ctx, 1, -1)
		requestStaticState.Unlock()
		return data.NewNullValue(), nil
	case "getTrustedProxies":
		requestStaticState.RLock()
		result := valuesArray(requestStaticState.trustedProxies)
		requestStaticState.RUnlock()
		return result, nil
	case "getTrustedHeaderSet":
		requestStaticState.RLock()
		result := data.NewIntValue(requestStaticState.trustedHeaderSet)
		requestStaticState.RUnlock()
		return result, nil
	case "setTrustedHosts":
		requestStaticState.Lock()
		requestStaticState.trustedHosts = arrayStrings(indexValue(ctx, 0))
		requestStaticState.Unlock()
		return data.NewNullValue(), nil
	case "getTrustedHosts":
		requestStaticState.RLock()
		result := valuesArray(requestStaticState.trustedHosts)
		requestStaticState.RUnlock()
		return result, nil
	case "normalizeQueryString":
		return data.NewStringValue(normalizeQuery(indexString(ctx, 0, ""))), nil
	case "enableHttpMethodParameterOverride":
		requestStaticState.Lock()
		requestStaticState.methodOverride = true
		requestStaticState.Unlock()
		return data.NewNullValue(), nil
	case "getHttpMethodParameterOverride":
		requestStaticState.RLock()
		result := data.NewBoolValue(requestStaticState.methodOverride)
		requestStaticState.RUnlock()
		return result, nil
	case "setAllowedHttpMethodOverride":
		requestStaticState.Lock()
		requestStaticState.allowedMethodOverride = arrayStrings(indexValue(ctx, 0))
		requestStaticState.Unlock()
		return data.NewNullValue(), nil
	case "getAllowedHttpMethodOverride":
		requestStaticState.RLock()
		result := valuesArray(requestStaticState.allowedMethodOverride)
		requestStaticState.RUnlock()
		return result, nil
	case "getSession":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		if source.session == nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("Session has not been set"))
		}
		return source.session, nil
	case "hasPreviousSession":
		return data.NewBoolValue(requestCookie(ctx, "MOCKSESSID") != ""), nil
	case "hasSession":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		return data.NewBoolValue(source.session != nil), nil
	case "setSession", "setSessionFactory":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		source.session = indexValue(ctx, 0)
		return data.NewNullValue(), nil
	case "getClientIps":
		return valuesArray(requestClientIPs(ctx)), nil
	case "getClientIp":
		ips := requestClientIPs(ctx)
		if len(ips) == 0 {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(ips[0]), nil
	case "getScriptName":
		return data.NewStringValue(requestServer(ctx, "SCRIPT_NAME")), nil
	case "getPathInfo":
		return data.NewStringValue(requestPath(ctx)), nil
	case "getBasePath":
		return data.NewStringValue(requestBasePath(ctx)), nil
	case "getBaseUrl":
		return data.NewStringValue(requestBaseURL(ctx)), nil
	case "getScheme":
		if requestSecure(ctx) {
			return data.NewStringValue("https"), nil
		}
		return data.NewStringValue("http"), nil
	case "getPort":
		return data.NewIntValue(requestPort(ctx)), nil
	case "getUser":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		if source.request.URL != nil && source.request.URL.User != nil {
			return data.NewStringValue(source.request.URL.User.Username()), nil
		}
		return nullableServer(ctx, "PHP_AUTH_USER"), nil
	case "getPassword":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		if source.request.URL != nil && source.request.URL.User != nil {
			if password, ok := source.request.URL.User.Password(); ok {
				return data.NewStringValue(password), nil
			}
		}
		return nullableServer(ctx, "PHP_AUTH_PW"), nil
	case "getUserInfo":
		user := requestUserInfo(ctx)
		return data.NewStringValue(user), nil
	case "getHttpHost":
		return data.NewStringValue(requestHTTPHost(ctx)), nil
	case "getRequestUri":
		return data.NewStringValue(requestURI(ctx)), nil
	case "getSchemeAndHttpHost":
		scheme := "http"
		if requestSecure(ctx) {
			scheme = "https"
		}
		return data.NewStringValue(scheme + "://" + requestHTTPHost(ctx)), nil
	case "getUri":
		return data.NewStringValue(requestFullURI(ctx)), nil
	case "getUriForPath":
		return data.NewStringValue(requestOrigin(ctx) + indexString(ctx, 0, "")), nil
	case "getRelativeUriForPath":
		return data.NewStringValue(relativeURI(requestPath(ctx), indexString(ctx, 0, ""))), nil
	case "getQueryString":
		query := requestQueryString(ctx)
		if query == "" {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(query), nil
	case "isSecure":
		return data.NewBoolValue(requestSecure(ctx)), nil
	case "getHost":
		return data.NewStringValue(requestHost(ctx)), nil
	case "setMethod":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		source.request.Method = strings.ToUpper(indexString(ctx, 0, "GET"))
		return data.NewNullValue(), nil
	case "getMethod":
		return data.NewStringValue(requestMethodName(ctx, true)), nil
	case "getRealMethod":
		return data.NewStringValue(requestMethodName(ctx, false)), nil
	case "getMimeType":
		mimes := requestFormats(indexString(ctx, 0, ""))
		if len(mimes) == 0 {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(mimes[0]), nil
	case "getMimeTypes":
		return valuesArray(requestFormats(indexString(ctx, 0, ""))), nil
	case "getFormat":
		format := formatForMime(indexString(ctx, 0, ""))
		if format == "" {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(format), nil
	case "setFormat":
		requestStaticState.Lock()
		requestStaticState.formats[indexString(ctx, 0, "")] = valueStrings(indexValue(ctx, 1))
		requestStaticState.Unlock()
		return data.NewNullValue(), nil
	case "getRequestFormat":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		if source.requestFormat != nil {
			return data.NewStringValue(*source.requestFormat), nil
		}
		if format := requestAttribute(ctx, "_format"); format != "" {
			return data.NewStringValue(format), nil
		}
		value := indexValue(ctx, 0)
		if value == nil {
			return data.NewStringValue("html"), nil
		}
		return value, nil
	case "setRequestFormat":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		value := indexValue(ctx, 0)
		if _, ok := value.(*data.NullValue); ok || value == nil {
			source.requestFormat = nil
		} else {
			format := value.AsString()
			source.requestFormat = &format
		}
		return data.NewNullValue(), nil
	case "getContentTypeFormat":
		format := formatForMime(contentType(requestHeader(ctx, "Content-Type")))
		if format == "" {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(format), nil
	case "setDefaultLocale":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		source.defaultLocale = indexString(ctx, 0, "en")
		return data.NewNullValue(), nil
	case "getDefaultLocale":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		return data.NewStringValue(source.defaultLocale), nil
	case "setLocale":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		source.locale = indexString(ctx, 0, source.defaultLocale)
		return data.NewNullValue(), nil
	case "getLocale":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		if source.locale == "" {
			source.locale = source.defaultLocale
		}
		return data.NewStringValue(source.locale), nil
	case "isMethod":
		return data.NewBoolValue(strings.EqualFold(requestMethodName(ctx, true), indexString(ctx, 0, ""))), nil
	case "isMethodSafe":
		value := requestMethodName(ctx, true)
		return data.NewBoolValue(value == "GET" || value == "HEAD" || value == "OPTIONS" || value == "TRACE"), nil
	case "isMethodIdempotent":
		value := requestMethodName(ctx, true)
		return data.NewBoolValue(value == "HEAD" || value == "GET" || value == "PUT" || value == "DELETE" || value == "OPTIONS" || value == "TRACE" || value == "PURGE"), nil
	case "isMethodCacheable":
		value := requestMethodName(ctx, true)
		return data.NewBoolValue(value == "GET" || value == "HEAD"), nil
	case "getProtocolVersion":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		return data.NewStringValue(strings.TrimPrefix(source.request.Proto, "HTTP/")), nil
	case "getContent":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		return data.NewStringValue(source.content), nil
	case "getPayload":
		if requestBag := requestBagProperty(ctx, "request"); requestBag != nil && len(GetParamBagMap(requestBag)) > 0 {
			return NewInputBagValue(ctx, GetParamBagMap(requestBag)), nil
		}
		return requestJSONBag(ctx)
	case "toArray":
		source, control := requestSource(ctx)
		if control != nil {
			return nil, control
		}
		value, err := decodeJSON(source.content)
		if err != nil {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("Could not decode request body: %w", err))
		}
		if _, ok := value.(*data.ArrayValue); !ok {
			return nil, data.NewErrorThrow(nil, fmt.Errorf("JSON content was expected to decode to an array"))
		}
		return value, nil
	case "getETags":
		var tags []string
		for _, tag := range strings.Split(requestHeader(ctx, "If-None-Match"), ",") {
			if tag = strings.TrimSpace(tag); tag != "" {
				tags = append(tags, tag)
			}
		}
		return valuesArray(tags), nil
	case "isNoCache":
		cache := strings.ToLower(requestHeader(ctx, "Cache-Control"))
		pragma := strings.ToLower(requestHeader(ctx, "Pragma"))
		return data.NewBoolValue(strings.Contains(cache, "no-cache") || pragma == "no-cache"), nil
	case "getPreferredFormat":
		if format := requestFormatValue(ctx); format != "" {
			return data.NewStringValue(format), nil
		}
		for _, mimeType := range parseAccept(requestHeader(ctx, "Accept")) {
			if format := formatForMime(mimeType); format != "" {
				return data.NewStringValue(format), nil
			}
		}
		value := indexValue(ctx, 0)
		if value == nil {
			return data.NewStringValue("html"), nil
		}
		return value, nil
	case "getPreferredLanguage":
		return preferredLanguage(ctx), nil
	case "getLanguages":
		return valuesArray(parseAccept(requestHeader(ctx, "Accept-Language"))), nil
	case "getCharsets":
		return valuesArray(parseAccept(requestHeader(ctx, "Accept-Charset"))), nil
	case "getEncodings":
		return valuesArray(parseAccept(requestHeader(ctx, "Accept-Encoding"))), nil
	case "getAcceptableContentTypes":
		return valuesArray(parseAccept(requestHeader(ctx, "Accept"))), nil
	case "isXmlHttpRequest":
		return data.NewBoolValue(strings.EqualFold(requestHeader(ctx, "X-Requested-With"), "XMLHttpRequest")), nil
	case "preferSafeContent":
		return data.NewBoolValue(strings.Contains(strings.ToLower(requestHeader(ctx, "Prefer")), "safe")), nil
	case "isFromTrustedProxy":
		return data.NewBoolValue(requestFromTrustedProxy(ctx)), nil
	}
	return data.NewNullValue(), nil
}

func requestInitialize(ctx data.Context) (data.GetValue, data.Control) {
	value := requestClassValue(ctx)
	if value == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Request initialization requires an object"))
	}
	source, _ := requestSource(ctx)
	if source == nil {
		source = &symfonyRequestSource{request: cloneHTTPRequest(nil), defaultLocale: "en", locale: "en"}
		if class, ok := value.Class.(*SymfonyRequestClass); ok {
			class.source = source
		}
	}
	query := valueMap(indexValue(ctx, 0))
	requestData := valueMap(indexValue(ctx, 1))
	attributes := valueMap(indexValue(ctx, 2))
	cookies := valueMap(indexValue(ctx, 3))
	files := valueMap(indexValue(ctx, 4))
	server := valueMap(indexValue(ctx, 5))
	if content := indexValue(ctx, 6); content != nil {
		source.content = content.AsString()
	}
	installRequestBags(value, ctx, query, requestData, attributes, cookies, files, server, nil)
	applyServerToRequest(source.request, server)
	return data.NewNullValue(), nil
}

func newRequestValue(ctx data.Context, request *http.Request, query, requestData, attributes, cookies, files, server map[string]data.Value) (*data.ClassValue, data.Control) {
	class := NewSymfonyRequestClassFrom(request)
	value := data.NewProxyValue(class, ctx.CreateBaseContext())
	if query == nil && class.source.request.URL != nil {
		query = urlValuesToMap(class.source.request.URL.Query())
	}
	if requestData == nil {
		requestData = map[string]data.Value{}
		contentType := contentType(class.source.request.Header.Get("Content-Type"))
		if contentType == "application/x-www-form-urlencoded" {
			if values, err := url.ParseQuery(class.source.content); err == nil {
				requestData = urlValuesToMap(values)
			}
		}
	}
	if attributes == nil {
		attributes = map[string]data.Value{}
	}
	if cookies == nil {
		cookies = cookiesFromRequest(class.source.request)
	}
	if files == nil {
		files = map[string]data.Value{}
	}
	if server == nil {
		server = serverFromRequest(class.source.request)
	}
	installRequestBags(value, ctx, query, requestData, attributes, cookies, files, server, headersFromRequest(class.source.request))
	return value, nil
}

func installRequestBags(value *data.ClassValue, ctx data.Context, query, requestData, attributes, cookies, files, server map[string]data.Value, headers map[string][]string) {
	if headers == nil {
		headers = headersFromServer(server)
	}
	_ = value.SetProperty("query", NewInputBagValue(ctx, query))
	_ = value.SetProperty("request", NewInputBagValue(ctx, requestData))
	_ = value.SetProperty("attributes", NewParameterBagValue(ctx, attributes))
	_ = value.SetProperty("cookies", NewInputBagValue(ctx, cookies))
	_ = value.SetProperty("files", NewFileBagValue(ctx, files))
	_ = value.SetProperty("server", NewServerBagValue(ctx, server))
	_ = value.SetProperty("headers", NewHeaderBagValue(ctx, headers))
}

func requestCreate(ctx data.Context) (data.GetValue, data.Control) {
	rawURI := indexString(ctx, 0, "/")
	method := strings.ToUpper(indexString(ctx, 1, "GET"))
	parameters := valueMap(indexValue(ctx, 2))
	cookies := valueMap(indexValue(ctx, 3))
	files := valueMap(indexValue(ctx, 4))
	server := valueMap(indexValue(ctx, 5))
	content := indexString(ctx, 6, "")
	parsed, err := url.Parse(rawURI)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Invalid URI: %w", err))
	}
	if parsed.Host == "" {
		parsed.Scheme = "http"
		parsed.Host = "localhost"
	}
	query, body := urlValuesToMap(parsed.Query()), map[string]data.Value{}
	switch method {
	case "POST", "PUT", "DELETE", "PATCH", "QUERY":
		body = parameters
		if content == "" && method != "PATCH" {
			content = mapToURLValues(parameters).Encode()
		}
	default:
		for key, value := range parameters {
			query[key] = value
		}
	}
	parsed.RawQuery = mapToURLValues(query).Encode()
	request, err := http.NewRequest(method, parsed.String(), strings.NewReader(content))
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	request.Host = parsed.Host
	if method == "POST" || method == "PUT" || method == "DELETE" || method == "QUERY" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for key, value := range server {
		if strings.HasPrefix(key, "HTTP_") {
			request.Header.Set(strings.ReplaceAll(strings.TrimPrefix(key, "HTTP_"), "_", "-"), value.AsString())
		}
	}
	defaultServer := serverFromRequest(request)
	for key, value := range server {
		defaultServer[key] = value
	}
	value, control := newRequestValue(ctx, request, query, body, nil, cookies, files, defaultServer)
	if control == nil {
		if class, ok := value.Class.(*SymfonyRequestClass); ok {
			class.source.content = content
		}
	}
	return value, control
}

func requestDuplicate(ctx data.Context) (data.GetValue, data.Control) {
	source, control := requestSource(ctx)
	if control != nil {
		return nil, control
	}
	current := requestClassValue(ctx)
	names := []string{"query", "request", "attributes", "cookies", "files", "server"}
	maps := make([]map[string]data.Value, len(names))
	for index, name := range names {
		if value := indexValue(ctx, index); value != nil {
			if _, null := value.(*data.NullValue); !null {
				maps[index] = valueMap(value)
				continue
			}
		}
		maps[index] = GetParamBagMap(requestBagPropertyFrom(current, name))
	}
	return newRequestValue(ctx, source.request, maps[0], maps[1], maps[2], maps[3], maps[4], maps[5])
}

func requestToString(ctx data.Context) (data.GetValue, data.Control) {
	source, control := requestSource(ctx)
	if control != nil {
		return nil, control
	}
	var builder strings.Builder
	fmt.Fprintf(&builder, "%s %s %s\r\n", requestMethodName(ctx, true), requestURI(ctx), source.request.Proto)
	headers := GetHeaderBagAll(requestBagProperty(ctx, "headers"))
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		for _, value := range headers[key] {
			fmt.Fprintf(&builder, "%s: %s\r\n", http.CanonicalHeaderKey(key), value)
		}
	}
	builder.WriteString("\r\n")
	builder.WriteString(source.content)
	return data.NewStringValue(builder.String()), nil
}

func indexValue(ctx data.Context, index int) data.Value {
	value, _ := ctx.GetIndexValue(index)
	return value
}

func indexString(ctx data.Context, index int, fallback string) string {
	value := indexValue(ctx, index)
	if value == nil {
		return fallback
	}
	if _, null := value.(*data.NullValue); null {
		return fallback
	}
	return value.AsString()
}

func arrayStrings(value data.Value) []string {
	if _, null := value.(*data.NullValue); value == nil || null {
		return nil
	}
	return valueStrings(value)
}

func valueStrings(value data.Value) []string {
	if array, ok := value.(*data.ArrayValue); ok {
		out := make([]string, 0, len(array.List))
		for _, item := range array.List {
			if item != nil && item.Value != nil {
				out = append(out, item.Value.AsString())
			}
		}
		return out
	}
	return []string{value.AsString()}
}

func normalizeQuery(query string) string {
	values, err := url.ParseQuery(query)
	if err != nil {
		return query
	}
	return values.Encode()
}

func headersFromServer(server map[string]data.Value) map[string][]string {
	headers := make(map[string][]string)
	for key, value := range server {
		header := ""
		switch {
		case strings.HasPrefix(key, "HTTP_"):
			header = strings.ReplaceAll(strings.TrimPrefix(key, "HTTP_"), "_", "-")
		case key == "CONTENT_TYPE", key == "CONTENT_LENGTH":
			header = strings.ReplaceAll(key, "_", "-")
		}
		if header != "" {
			headers[header] = []string{value.AsString()}
		}
	}
	return headers
}

func applyServerToRequest(request *http.Request, server map[string]data.Value) {
	if request == nil {
		return
	}
	if value, ok := server["REQUEST_METHOD"]; ok {
		request.Method = strings.ToUpper(value.AsString())
	}
	if value, ok := server["HTTP_HOST"]; ok {
		request.Host = value.AsString()
	}
	if value, ok := server["SERVER_PROTOCOL"]; ok {
		request.Proto = value.AsString()
	}
	if value, ok := server["REMOTE_ADDR"]; ok {
		request.RemoteAddr = value.AsString()
	}
	if request.URL == nil {
		request.URL = &url.URL{}
	}
	if value, ok := server["REQUEST_URI"]; ok {
		if parsed, err := url.ParseRequestURI(value.AsString()); err == nil {
			request.URL.Path, request.URL.RawPath, request.URL.RawQuery = parsed.Path, parsed.RawPath, parsed.RawQuery
		}
	}
}

func requestBagProperty(ctx data.Context, name string) *data.ClassValue {
	return requestBagPropertyFrom(requestClassValue(ctx), name)
}

func requestBagPropertyFrom(value *data.ClassValue, name string) *data.ClassValue {
	if value == nil {
		return nil
	}
	property, _ := value.ObjectValue.GetProperty(name)
	bag, _ := property.(*data.ClassValue)
	return bag
}

func requestServer(ctx data.Context, key string) string {
	bag := requestBagProperty(ctx, "server")
	if value, ok := GetParamBagMap(bag)[key]; ok && value != nil {
		return value.AsString()
	}
	return ""
}

func nullableServer(ctx data.Context, key string) data.GetValue {
	if value := requestServer(ctx, key); value != "" {
		return data.NewStringValue(value)
	}
	return data.NewNullValue()
}

func requestAttribute(ctx data.Context, key string) string {
	bag := requestBagProperty(ctx, "attributes")
	if value, ok := GetParamBagMap(bag)[key]; ok && value != nil {
		return value.AsString()
	}
	return ""
}

func requestCookie(ctx data.Context, key string) string {
	bag := requestBagProperty(ctx, "cookies")
	if value, ok := GetParamBagMap(bag)[key]; ok && value != nil {
		return value.AsString()
	}
	return ""
}

func requestHeader(ctx data.Context, key string) string {
	values := GetHeaderBagAll(requestBagProperty(ctx, "headers"))
	items := values[strings.ToLower(key)]
	if len(items) == 0 {
		items = values[key]
	}
	return strings.Join(items, ", ")
}

func requestPath(ctx data.Context) string {
	source, _ := requestSource(ctx)
	if source != nil && source.request.URL != nil {
		if value := source.request.URL.EscapedPath(); value != "" {
			return value
		}
	}
	if value := requestServer(ctx, "PATH_INFO"); value != "" {
		return value
	}
	return "/"
}

func requestBaseURL(ctx data.Context) string {
	script := requestServer(ctx, "SCRIPT_NAME")
	if script == "" || script == "/" {
		return ""
	}
	uri := requestURI(ctx)
	if strings.HasPrefix(uri, script) {
		return script
	}
	return ""
}

func requestBasePath(ctx data.Context) string {
	base := requestBaseURL(ctx)
	if base == "" {
		return ""
	}
	return strings.TrimSuffix(path.Dir(base), "/")
}

func requestSecure(ctx data.Context) bool {
	source, _ := requestSource(ctx)
	if source != nil && (source.request.TLS != nil || source.request.URL != nil && source.request.URL.Scheme == "https") {
		return true
	}
	value := strings.ToLower(requestServer(ctx, "HTTPS"))
	return value != "" && value != "off" && value != "0"
}

func requestPort(ctx data.Context) int {
	if port := requestServer(ctx, "SERVER_PORT"); port != "" {
		if parsed, err := strconv.Atoi(port); err == nil {
			return parsed
		}
	}
	_, port, err := net.SplitHostPort(requestHTTPHost(ctx))
	if err == nil {
		parsed, _ := strconv.Atoi(port)
		return parsed
	}
	if requestSecure(ctx) {
		return 443
	}
	return 80
}

func requestHost(ctx data.Context) string {
	host := requestHTTPHost(ctx)
	if parsed, _, err := net.SplitHostPort(host); err == nil {
		host = parsed
	}
	return strings.Trim(host, "[]")
}

func requestHTTPHost(ctx data.Context) string {
	if host := requestHeader(ctx, "Host"); host != "" {
		return host
	}
	if host := requestServer(ctx, "HTTP_HOST"); host != "" {
		return host
	}
	host := requestServer(ctx, "SERVER_NAME")
	port := requestPort(ctx)
	if port != 80 && port != 443 {
		return net.JoinHostPort(host, strconv.Itoa(port))
	}
	return host
}

func requestURI(ctx data.Context) string {
	source, _ := requestSource(ctx)
	if source != nil && source.request.URL != nil {
		return source.request.URL.RequestURI()
	}
	if value := requestServer(ctx, "REQUEST_URI"); value != "" {
		return value
	}
	return "/"
}

func requestQueryString(ctx data.Context) string {
	source, _ := requestSource(ctx)
	if source != nil && source.request.URL != nil {
		return normalizeQuery(source.request.URL.RawQuery)
	}
	return normalizeQuery(requestServer(ctx, "QUERY_STRING"))
}

func requestOrigin(ctx data.Context) string {
	scheme := "http"
	if requestSecure(ctx) {
		scheme = "https"
	}
	return scheme + "://" + requestHTTPHost(ctx) + requestBaseURL(ctx)
}

func requestFullURI(ctx data.Context) string {
	uri := requestOrigin(ctx) + requestPath(ctx)
	if query := requestQueryString(ctx); query != "" {
		uri += "?" + query
	}
	return uri
}

func requestUserInfo(ctx data.Context) string {
	userValue, _ := newRequestMethod("getUser").Call(ctx)
	if userValue == nil {
		return ""
	}
	user := userValue.(data.Value).AsString()
	if user == "" {
		return ""
	}
	passwordValue, _ := newRequestMethod("getPassword").Call(ctx)
	if passwordValue != nil {
		if password := passwordValue.(data.Value).AsString(); password != "" {
			return user + ":" + password
		}
	}
	return user
}

func requestMethodName(ctx data.Context, override bool) string {
	source, _ := requestSource(ctx)
	if source == nil || source.request == nil {
		return "GET"
	}
	method := strings.ToUpper(source.request.Method)
	if !override || method != "POST" {
		return method
	}
	overrideMethod := requestHeader(ctx, "X-HTTP-METHOD-OVERRIDE")
	requestStaticState.RLock()
	enabled := requestStaticState.methodOverride
	allowed := append([]string(nil), requestStaticState.allowedMethodOverride...)
	requestStaticState.RUnlock()
	if overrideMethod == "" && enabled {
		if value, ok := GetParamBagMap(requestBagProperty(ctx, "request"))["_method"]; ok {
			overrideMethod = value.AsString()
		}
	}
	overrideMethod = strings.ToUpper(overrideMethod)
	if overrideMethod == "" {
		return method
	}
	if len(allowed) > 0 {
		for _, candidate := range allowed {
			if strings.EqualFold(candidate, overrideMethod) {
				return overrideMethod
			}
		}
		return method
	}
	return overrideMethod
}

func requestClientIPs(ctx data.Context) []string {
	source, _ := requestSource(ctx)
	remote := ""
	if source != nil {
		remote = source.request.RemoteAddr
		if host, _, err := net.SplitHostPort(remote); err == nil {
			remote = host
		}
	}
	if requestFromTrustedProxy(ctx) {
		if forwarded := requestHeader(ctx, "X-Forwarded-For"); forwarded != "" {
			var ips []string
			for _, value := range strings.Split(forwarded, ",") {
				ips = append(ips, strings.TrimSpace(value))
			}
			if remote != "" {
				ips = append(ips, remote)
			}
			return ips
		}
	}
	if remote == "" {
		return nil
	}
	return []string{remote}
}

func requestFromTrustedProxy(ctx data.Context) bool {
	source, _ := requestSource(ctx)
	if source == nil {
		return false
	}
	remote := source.request.RemoteAddr
	if host, _, err := net.SplitHostPort(remote); err == nil {
		remote = host
	}
	requestStaticState.RLock()
	defer requestStaticState.RUnlock()
	for _, trusted := range requestStaticState.trustedProxies {
		if trusted == remote || trusted == "REMOTE_ADDR" && remote == requestServer(ctx, "REMOTE_ADDR") {
			return true
		}
		if _, network, err := net.ParseCIDR(trusted); err == nil && network.Contains(net.ParseIP(remote)) {
			return true
		}
	}
	return false
}

func requestFormats(format string) []string {
	requestStaticState.RLock()
	defer requestStaticState.RUnlock()
	return append([]string(nil), requestStaticState.formats[format]...)
}

func formatForMime(mimeType string) string {
	mimeType = contentType(strings.ToLower(mimeType))
	requestStaticState.RLock()
	defer requestStaticState.RUnlock()
	for format, types := range requestStaticState.formats {
		for _, candidate := range types {
			if candidate == mimeType {
				return format
			}
		}
	}
	if plus := strings.LastIndex(mimeType, "+"); plus >= 0 {
		suffix := mimeType[plus+1:]
		if suffix == "json" {
			return "json"
		}
		if suffix == "xml" {
			return "xml"
		}
	}
	return ""
}

func requestFormatValue(ctx data.Context) string {
	source, _ := requestSource(ctx)
	if source != nil && source.requestFormat != nil {
		return *source.requestFormat
	}
	return requestAttribute(ctx, "_format")
}

func requestJSONBag(ctx data.Context) (data.GetValue, data.Control) {
	source, control := requestSource(ctx)
	if control != nil {
		return nil, control
	}
	if source.content == "" {
		return NewInputBagValue(ctx, nil), nil
	}
	value, err := decodeJSON(source.content)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Could not decode request body: %w", err))
	}
	return NewInputBagValue(ctx, valueMap(value)), nil
}

func preferredLanguage(ctx data.Context) data.GetValue {
	accepted := parseAccept(requestHeader(ctx, "Accept-Language"))
	allowedValue := indexValue(ctx, 0)
	if allowedValue == nil {
		if len(accepted) == 0 {
			return data.NewNullValue()
		}
		return data.NewStringValue(accepted[0])
	}
	allowed := arrayStrings(allowedValue)
	for _, language := range accepted {
		for _, candidate := range allowed {
			if strings.EqualFold(language, candidate) || strings.HasPrefix(strings.ToLower(language), strings.ToLower(candidate)+"-") {
				return data.NewStringValue(candidate)
			}
		}
	}
	if len(allowed) > 0 {
		return data.NewStringValue(allowed[0])
	}
	return data.NewNullValue()
}

func relativeURI(from, target string) string {
	fromURL, _ := url.Parse(from)
	targetURL, _ := url.Parse(target)
	fromParts := strings.Split(strings.Trim(path.Dir(fromURL.Path), "/"), "/")
	targetParts := strings.Split(strings.Trim(targetURL.Path, "/"), "/")
	for len(fromParts) > 0 && len(targetParts) > 0 && fromParts[0] == targetParts[0] {
		fromParts, targetParts = fromParts[1:], targetParts[1:]
	}
	relative := strings.Repeat("../", len(fromParts)) + strings.Join(targetParts, "/")
	if relative == "" {
		relative = "./"
	}
	if targetURL.RawQuery != "" {
		relative += "?" + targetURL.RawQuery
	}
	return relative
}

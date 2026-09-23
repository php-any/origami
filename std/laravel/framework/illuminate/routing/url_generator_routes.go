package routing

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
)

var (
	urlValidPattern   = regexp.MustCompile(`^(#|//|https?://|(mailto|tel|sms):)`)
	routeParamPattern = regexp.MustCompile(`\{(.*?)(\?)?\}`)
)

func (c *UrlGeneratorClass) registerRouteMethods() {
	c.methods["route"] = kit.InstanceMethodOpt("route", []string{"name", "parameters", "absolute"}, 1, urlRoute)
	c.methods["toroute"] = kit.InstanceMethodOpt("toRoute", []string{"route", "parameters", "absolute"}, 1, urlToRoute)
	c.methods["action"] = kit.InstanceMethodOpt("action", []string{"action", "parameters", "absolute"}, 1, urlAction)
	c.methods["previous"] = kit.InstanceMethodOpt("previous", []string{"fallback"}, 0, urlPrevious)
	c.methods["previouspath"] = kit.InstanceMethodOpt("previousPath", []string{"fallback"}, 0, urlPreviousPath)
	c.methods["query"] = kit.InstanceMethodOpt("query", []string{"path", "query", "extra", "secure"}, 1, urlQuery)
	c.methods["isvalidurl"] = kit.InstanceMethod("isValidUrl", []string{"path"}, urlIsValidURL)
	c.methods["setroutes"] = kit.InstanceMethod("setRoutes", []string{"routes"}, urlSetRoutes)
	c.methods["getrequest"] = kit.InstanceMethod("getRequest", nil, urlGetRequest)
	c.methods["forcehttps"] = kit.InstanceMethodOpt("forceHttps", []string{"force"}, 0, urlForceHTTPS)
	c.methods["defaults"] = kit.InstanceMethod("defaults", []string{"defaults"}, urlDefaults)
	c.methods["formatparameters"] = kit.InstanceMethod("formatParameters", []string{"parameters"}, urlFormatParameters)
	c.methods["formatroot"] = kit.InstanceMethodOpt("formatRoot", []string{"scheme", "root"}, 1, urlFormatRoot)
	c.methods["setsessionresolver"] = kit.InstanceMethod("setSessionResolver", []string{"sessionResolver"}, urlSetSessionResolver)
	c.methods["setkeyresolver"] = kit.InstanceMethod("setKeyResolver", []string{"keyResolver"}, urlSetKeyResolver)
	c.methods["withkeyresolver"] = kit.InstanceMethod("withKeyResolver", []string{"keyResolver"}, urlWithKeyResolver)
	c.methods["setrootcontrollernamespace"] = kit.InstanceMethod("setRootControllerNamespace", []string{"rootNamespace"}, urlSetRootControllerNamespace)
}

func urlSetSessionResolver(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("sessionResolver", kit.Arg(ctx, 0))
	return cv, nil
}

func urlSetKeyResolver(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("keyResolver", kit.Arg(ctx, 0))
	return cv, nil
}

func urlWithKeyResolver(ctx data.Context) (data.GetValue, data.Control) {
	return urlSetKeyResolver(ctx)
}

func urlSetRootControllerNamespace(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("rootNamespace", kit.Arg(ctx, 0))
	return cv, nil
}

func urlRoutesCollection(cv *data.ClassValue) *data.ClassValue {
	routes, _ := cv.GetProperty("routes")
	if rcv, ok := kit.Unwrap(routes).(*data.ClassValue); ok {
		return rcv
	}
	return nil
}

func urlRoute(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	nameVal := kit.Arg(ctx, 0)
	name := nameVal.AsString()
	if cvEnum, ok := kit.Unwrap(nameVal).(*data.ClassValue); ok && cvEnum.Class != nil {
		if m, ok := cvEnum.GetMethod("value"); ok && m != nil {
			nctx := cvEnum.CreateContext(m.GetVariables())
			if ret, ctl := m.Call(nctx); ctl == nil && ret != nil {
				if v, ok := ret.(data.Value); ok {
					name = v.AsString()
				}
			}
		}
	}
	params := kit.Arg(ctx, 1)
	absolute := true
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		if b, ok := v.(data.AsBool); ok {
			absolute, _ = b.AsBool()
		}
	}
	routes := urlRoutesCollection(cv)
	if routes != nil {
		ret, ctl := kit.CallInstanceMethod(ctx, routes, "getByName", data.NewStringValue(name))
		if ctl != nil {
			return nil, ctl
		}
		if ret != nil && !kit.IsNull(ret.(data.Value)) {
			return urlToRouteWith(ctx, cv, ret.(data.Value), params, absolute)
		}
	}
	if resolver, _ := cv.GetProperty("missingNamedRouteResolver"); resolver != nil && !kit.IsNull(resolver) {
		ret, ctl := kit.Call(ctx, resolver, data.NewStringValue(name), params, data.NewBoolValue(absolute))
		if ctl != nil {
			return nil, ctl
		}
		if ret != nil && !kit.IsNull(ret.(data.Value)) {
			return ret, nil
		}
	}
	return nil, httpfoundation.ThrowNamed("Symfony\\Component\\Routing\\Exception\\RouteNotFoundException", "Route [%s] not defined.", name)
}

func urlToRoute(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	route := kit.Arg(ctx, 0)
	params := kit.Arg(ctx, 1)
	absolute := true
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		if b, ok := v.(data.AsBool); ok {
			absolute, _ = b.AsBool()
		}
	}
	return urlToRouteWith(ctx, cv, route, params, absolute)
}

func urlToRouteWith(ctx data.Context, cv *data.ClassValue, route, params data.Value, absolute bool) (data.GetValue, data.Control) {
	if route == nil || kit.IsNull(route) {
		return data.NewStringValue("/"), nil
	}
	routeCV, ok := kit.Unwrap(route).(*data.ClassValue)
	if !ok {
		return data.NewStringValue("/"), nil
	}
	uriRet, ctl := kit.CallInstanceMethod(ctx, routeCV, "uri")
	if ctl != nil {
		return nil, ctl
	}
	uri := ""
	if uriRet != nil {
		uri = uriRet.(data.Value).AsString()
	}
	formatted := urlFormatParametersOnCtx(ctx, cv, params)
	path, q := replaceRouteURI(uri, formatted, cv)
	fullPath := strings.TrimLeft(path, "/")
	if q != "" {
		fullPath += "?" + q
	}
	if !absolute {
		base := "/"
		if req, _ := cv.GetProperty("request"); req != nil {
			if rcv, ok := kit.Unwrap(req).(*data.ClassValue); ok {
				if ret, ctl := kit.CallInstanceMethod(ctx, rcv, "getBaseUrl"); ctl == nil && ret != nil {
					base = ret.(data.Value).AsString()
				}
			}
		}
		out := "/" + strings.TrimLeft(strings.TrimPrefix(fullPath, "?"), "/")
		if base != "" && base != "/" {
			out = strings.TrimPrefix(out, base)
			if out == "" {
				out = "/"
			}
		}
		return data.NewStringValue(out), nil
	}
	return urlToPath(cv, fullPath, nil, nil)
}

func urlAction(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	actionVal := kit.Arg(ctx, 0)
	action := urlFormatAction(cv, actionVal)
	params := kit.Arg(ctx, 1)
	absolute := true
	if v := kit.Arg(ctx, 2); v != nil && !kit.IsNull(v) {
		if b, ok := v.(data.AsBool); ok {
			absolute, _ = b.AsBool()
		}
	}
	routes := urlRoutesCollection(cv)
	if routes == nil {
		return nil, httpfoundation.ThrowNamed("InvalidArgumentException", "Action %s not defined.", action)
	}
	ret, ctl := kit.CallInstanceMethod(ctx, routes, "getByAction", data.NewStringValue(action))
	if ctl != nil {
		return nil, ctl
	}
	if ret == nil || kit.IsNull(ret.(data.Value)) {
		return nil, httpfoundation.ThrowNamed("InvalidArgumentException", "Action %s not defined.", action)
	}
	return urlToRouteWith(ctx, cv, ret.(data.Value), params, absolute)
}

func urlFormatAction(cv *data.ClassValue, action data.Value) string {
	if arr, ok := kit.Unwrap(action).(*data.ArrayValue); ok {
		parts := make([]string, 0, len(arr.List))
		for _, z := range arr.List {
			if z != nil && z.Value != nil {
				parts = append(parts, z.Value.AsString())
			}
		}
		if len(parts) >= 2 {
			action = data.NewStringValue("\\" + parts[0] + "@" + parts[1])
		}
	}
	s := action.AsString()
	ns, _ := cv.GetProperty("rootNamespace")
	if nsVal, ok := kit.Unwrap(ns).(*data.StringValue); ok && nsVal.Value != "" {
		if !strings.HasPrefix(s, "\\") {
			s = nsVal.Value + "\\" + s
		}
	}
	return strings.Trim(s, "\\")
}

func urlPrevious(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	referrer := requestHeaderValue(cv, "referer")
	var urlStr string
	if referrer != "" {
		ret, ctl := urlToPath(cv, referrer, nil, nil)
		if ctl != nil {
			return nil, ctl
		}
		urlStr = ret.(data.Value).AsString()
	}
	if urlStr == "" {
		urlStr = urlPreviousFromSession(ctx, cv)
	}
	if urlStr != "" {
		return data.NewStringValue(urlStr), nil
	}
	fallback := kit.Arg(ctx, 0)
	if fallback != nil && !kit.IsNull(fallback) && kit.Truthy(fallback) {
		return urlToPath(cv, fallback.AsString(), nil, nil)
	}
	return urlToPath(cv, "/", nil, nil)
}

func urlPreviousFromSession(ctx data.Context, cv *data.ClassValue) string {
	resolver, _ := cv.GetProperty("sessionResolver")
	if resolver == nil || kit.IsNull(resolver) {
		return ""
	}
	ret, ctl := kit.Call(ctx, resolver)
	if ctl != nil || ret == nil {
		return ""
	}
	sess, ok := kit.Unwrap(ret.(data.Value)).(*data.ClassValue)
	if !ok {
		return ""
	}
	prev, ctl := kit.CallInstanceMethod(ctx, sess, "previousUrl")
	if ctl != nil || prev == nil {
		return ""
	}
	return prev.(data.Value).AsString()
}

func urlPreviousPath(ctx data.Context) (data.GetValue, data.Control) {
	prevRet, ctl := urlPrevious(ctx)
	if ctl != nil {
		return nil, ctl
	}
	prev := prevRet.(data.Value).AsString()
	u, err := url.Parse(prev)
	if err != nil || u.Path == "" {
		return data.NewStringValue("/"), nil
	}
	cv, _ := urlRecv(ctx)
	basePath := requestBasePath(ctx, cv)
	previousPath := u.Path
	if basePath != "" && basePath != "/" {
		previousPath = strings.TrimPrefix(previousPath, basePath)
	}
	previousPath = strings.TrimRight(previousPath, "/")
	if previousPath == "" {
		previousPath = "/"
	}
	return data.NewStringValue(previousPath), nil
}

func urlQuery(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	path := kit.Arg(ctx, 0).AsString()
	queryVal := kit.Arg(ctx, 1)
	extra := kit.Arg(ctx, 2)
	secure := kit.Arg(ctx, 3)
	pathPart, existingQS := extractQueryString(path)
	existing := parseQueryString(existingQS)
	merged := mergeQueryStringMaps(existing, paramsAssocMap(queryVal))
	qs := buildQueryString(merged)
	newPath := pathPart
	if qs != "" {
		newPath += "?" + qs
	}
	ret, ctl := urlToPath(cv, newPath, extra, secure)
	if ctl != nil {
		return nil, ctl
	}
	s := strings.TrimRight(ret.(data.Value).AsString(), "?")
	return data.NewStringValue(s), nil
}

func urlIsValidURL(ctx data.Context) (data.GetValue, data.Control) {
	path := kit.Arg(ctx, 0).AsString()
	return data.NewBoolValue(urlPathIsValid(path)), nil
}

func urlPathIsValid(path string) bool {
	if urlValidPattern.MatchString(path) {
		return true
	}
	u, err := url.Parse(path)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func urlSetRoutes(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("routes", kit.Arg(ctx, 0))
	return cv, nil
}

func urlGetRequest(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	req, _ := cv.GetProperty("request")
	if req == nil {
		return data.NewNullValue(), nil
	}
	return req, nil
}

func urlForceHTTPS(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	force := true
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		if b, ok := v.(data.AsBool); ok {
			force, _ = b.AsBool()
		}
	}
	if force {
		_ = cv.SetProperty("forceScheme", data.NewStringValue("https://"))
		_ = cv.SetProperty("cachedScheme", data.NewNullValue())
	}
	return data.NewNullValue(), nil
}

func urlDefaults(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	incoming := paramsAssocMap(kit.Arg(ctx, 0))
	existing, _ := cv.GetProperty("defaultParameters")
	base := paramsAssocMap(existing)
	for k, v := range incoming {
		base[k] = v
	}
	_ = cv.SetProperty("defaultParameters", assocMapToArray(base))
	return data.NewNullValue(), nil
}

func urlFormatParameters(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := urlFormatParametersOnCtx(ctx, cv, kit.Arg(ctx, 0))
	return assocMapToArray(out), nil
}

func urlFormatParametersOnCtx(ctx data.Context, cv *data.ClassValue, parameters data.Value) map[string]data.Value {
	wrapped := paramsAssocMap(parameters)
	for k, v := range wrapped {
		if cvObj, ok := kit.Unwrap(v).(*data.ClassValue); ok {
			if m, ok := cvObj.GetMethod("getRouteKey"); ok && m != nil {
				nctx := cvObj.CreateContext(m.GetVariables())
				if ret, ctl := m.Call(nctx); ctl == nil && ret != nil {
					wrapped[k] = ret.(data.Value)
				}
			}
		}
	}
	return wrapped
}

func urlFormatRoot(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	scheme := kit.Arg(ctx, 0).AsString()
	rootVal := kit.Arg(ctx, 1)
	root := ""
	if rootVal != nil && !kit.IsNull(rootVal) {
		root = strings.TrimRight(rootVal.AsString(), "/")
	} else {
		root = strings.TrimRight(requestRoot(cv), "/")
	}
	if root == "" {
		return data.NewStringValue(scheme), nil
	}
	start := "http://"
	if strings.HasPrefix(root, "https://") {
		start = "https://"
	} else if strings.HasPrefix(root, "http://") {
		start = "http://"
	}
	if scheme == "" {
		return data.NewStringValue(root), nil
	}
	out := strings.Replace(root, start, scheme, 1)
	return data.NewStringValue(out), nil
}

func urlToPath(cv *data.ClassValue, path string, extra, secure data.Value) (data.GetValue, data.Control) {
	if urlPathIsValid(path) {
		return data.NewStringValue(path), nil
	}
	tailParts := formatParameterSegments(cv, extra)
	tail := strings.Join(tailParts, "/")
	pathPart, qs := extractQueryString(path)
	if tail != "" {
		pathPart = strings.Trim(pathPart+"/"+tail, "/")
	}
	root := requestRoot(cv)
	if force, _ := cv.GetProperty("forceScheme"); force != nil && !kit.IsNull(force) {
		scheme := force.AsString()
		if root != "" && scheme != "" {
			if strings.HasPrefix(root, "http://") {
				root = scheme + strings.TrimPrefix(root, "http://")
			} else if strings.HasPrefix(root, "https://") {
				root = scheme + strings.TrimPrefix(root, "https://")
			}
		}
	}
	if secure != nil && !kit.IsNull(secure) {
		if b, ok := secure.(data.AsBool); ok {
			okb, _ := b.AsBool()
			if okb && root != "" {
				root = strings.Replace(root, "http://", "https://", 1)
			}
		}
	}
	rel := "/" + strings.TrimLeft(pathPart, "/")
	if root == "" {
		if qs != "" {
			return data.NewStringValue(rel + qs), nil
		}
		return data.NewStringValue(rel), nil
	}
	full := strings.TrimRight(root, "/") + rel
	if qs != "" {
		full += qs
	}
	return data.NewStringValue(full), nil
}

func replaceRouteURI(uri string, params map[string]data.Value, cv *data.ClassValue) (string, string) {
	remaining := make(map[string]data.Value, len(params))
	for k, v := range params {
		remaining[k] = v
	}
	defaults := paramsAssocMap(nil)
	if d, _ := cv.GetProperty("defaultParameters"); d != nil {
		defaults = paramsAssocMap(d)
	}
	path := routeParamPattern.ReplaceAllStringFunc(uri, func(m string) string {
		sub := routeParamPattern.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		name := sub[1]
		optional := len(sub) > 2 && sub[2] == "?"
		if v, ok := remaining[name]; ok {
			delete(remaining, name)
			if kit.IsNull(v) || v.AsString() == "" {
				if optional {
					return ""
				}
				return m
			}
			return url.PathEscape(v.AsString())
		}
		if dv, ok := defaults[name]; ok && dv != nil && !kit.IsNull(dv) && dv.AsString() != "" {
			return url.PathEscape(dv.AsString())
		}
		if optional {
			return ""
		}
		return m
	})
	path = strings.Trim(strings.ReplaceAll(path, "//", "/"), "/")
	q := buildQueryString(stringParamsFromRemaining(remaining))
	return path, q
}

func stringParamsFromRemaining(m map[string]data.Value) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		if v == nil || kit.IsNull(v) {
			continue
		}
		out[k] = v.AsString()
	}
	return out
}

func paramsAssocMap(v data.Value) map[string]data.Value {
	if v == nil || kit.IsNull(v) {
		return map[string]data.Value{}
	}
	entries := kit.Entries(v)
	if len(entries) == 0 {
		return map[string]data.Value{}
	}
	out := make(map[string]data.Value, len(entries))
	for _, e := range entries {
		key := e.KeyStr
		if key == "" {
			key = kit.KeyString(e.Key)
		}
		if key == "" {
			continue
		}
		out[key] = e.Value
	}
	return out
}

func assocMapToArray(m map[string]data.Value) data.Value {
	arr := data.NewArrayValue(nil)
	av, ok := arr.(*data.ArrayValue)
	if !ok {
		return arr
	}
	for k, v := range m {
		av.SetStringKey(k, v)
	}
	return av
}

func formatParameterSegments(cv *data.ClassValue, extra data.Value) []string {
	if extra == nil || kit.IsNull(extra) {
		return nil
	}
	formatted := urlFormatParametersOnCtx(nil, cv, extra)
	if len(formatted) == 0 {
		if arr, ok := kit.Unwrap(extra).(*data.ArrayValue); ok {
			out := make([]string, 0, len(arr.List))
			for _, z := range arr.List {
				if z != nil && z.Value != nil {
					out = append(out, url.QueryEscape(z.Value.AsString()))
				}
			}
			return out
		}
		return nil
	}
	out := make([]string, 0, len(formatted))
	for _, v := range formatted {
		if v != nil {
			out = append(out, url.QueryEscape(v.AsString()))
		}
	}
	return out
}

func extractQueryString(path string) (string, string) {
	if i := strings.Index(path, "?"); i >= 0 {
		return path[:i], path[i:]
	}
	return path, ""
}

func parseQueryString(qs string) map[string]string {
	qs = strings.TrimPrefix(qs, "?")
	out := map[string]string{}
	if qs == "" {
		return out
	}
	for _, part := range strings.Split(qs, "&") {
		if part == "" {
			continue
		}
		k, v, _ := strings.Cut(part, "=")
		k, _ = url.QueryUnescape(k)
		v, _ = url.QueryUnescape(v)
		out[k] = v
	}
	return out
}

func mergeQueryStringMaps(a map[string]string, b map[string]data.Value) map[string]string {
	out := make(map[string]string, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v != nil {
			out[k] = v.AsString()
		}
	}
	return out
}

func buildQueryString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	return strings.Join(parts, "&")
}

func requestBasePath(ctx data.Context, cv *data.ClassValue) string {
	req, _ := cv.GetProperty("request")
	rcv, ok := kit.Unwrap(req).(*data.ClassValue)
	if !ok {
		return ""
	}
	ret, ctl := kit.CallInstanceMethod(ctx, rcv, "getBaseUrl")
	if ctl != nil || ret == nil {
		return ""
	}
	return ret.(data.Value).AsString()
}

func requestHeaderValue(cv *data.ClassValue, key string) string {
	req, _ := cv.GetProperty("request")
	rcv, ok := kit.Unwrap(req).(*data.ClassValue)
	if !ok {
		return ""
	}
	headers, _ := rcv.GetProperty("headers")
	hcv, ok := kit.Unwrap(headers).(*data.ClassValue)
	if !ok {
		return ""
	}
	ret, ctl := kit.CallInstanceMethod(nil, hcv, "get", data.NewStringValue(key))
	if ctl != nil || ret == nil {
		return ""
	}
	return ret.(data.Value).AsString()
}

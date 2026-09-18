package routing

import (
	"net/url"
	"strings"

	"github.com/php-any/origami/data"
)

const requestContextName = "Symfony\\Component\\Routing\\RequestContext"

func newRequestContextClass() *rtClass {
	c := newRtClass(requestContextName, nil, nil, []data.Property{
		privProp("baseUrl", data.NewStringValue("")),
		privProp("pathInfo", data.NewStringValue("/")),
		privProp("method", data.NewStringValue("GET")),
		privProp("host", data.NewStringValue("localhost")),
		privProp("scheme", data.NewStringValue("http")),
		privProp("httpPort", data.NewIntValue(80)),
		privProp("httpsPort", data.NewIntValue(443)),
		privProp("queryString", data.NewStringValue("")),
		privProp("parameters", emptyArray()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{
			param("baseUrl", 0, data.NewStringValue(""), nil),
			param("method", 1, data.NewStringValue("GET"), nil),
			param("host", 2, data.NewStringValue("localhost"), nil),
			param("scheme", 3, data.NewStringValue("http"), nil),
			param("httpPort", 4, data.NewIntValue(80), nil),
			param("httpsPort", 5, data.NewIntValue(443), nil),
			param("path", 6, data.NewStringValue("/"), nil),
			param("queryString", 7, data.NewStringValue(""), nil),
			param("parameters", 8, data.NewNullValue(), nil),
		},
		[]data.Variable{
			variable("baseUrl", 0, nil),
			variable("method", 1, nil),
			variable("host", 2, nil),
			variable("scheme", 3, nil),
			variable("httpPort", 4, nil),
			variable("httpsPort", 5, nil),
			variable("path", 6, nil),
			variable("queryString", 7, nil),
			variable("parameters", 8, nil),
		},
		contextConstruct))
	c.add(staticMeth("fromUri", []string{"uri", "host", "scheme", "httpPort", "httpsPort"}, contextFromUri))
	c.add(meth("fromRequest", []string{"request"}, contextFromRequest))
	c.add(meth("getBaseUrl", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "baseUrl", "")), nil
	}))
	c.add(meth("setBaseUrl", []string{"baseUrl"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "baseUrl", data.NewStringValue(strings.TrimRight(argString(ctx, 0, ""), "/")))
		return cv, nil
	}))
	c.add(meth("getPathInfo", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "pathInfo", "/")), nil
	}))
	c.add(meth("setPathInfo", []string{"pathInfo"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "pathInfo", data.NewStringValue(argString(ctx, 0, "/")))
		return cv, nil
	}))
	c.add(meth("getMethod", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "method", "GET")), nil
	}))
	c.add(meth("setMethod", []string{"method"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "method", data.NewStringValue(strings.ToUpper(argString(ctx, 0, "GET"))))
		return cv, nil
	}))
	c.add(meth("getHost", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "host", "localhost")), nil
	}))
	c.add(meth("setHost", []string{"host"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "host", data.NewStringValue(strings.ToLower(argString(ctx, 0, "localhost"))))
		return cv, nil
	}))
	c.add(meth("getScheme", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "scheme", "http")), nil
	}))
	c.add(meth("setScheme", []string{"scheme"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "scheme", data.NewStringValue(strings.ToLower(argString(ctx, 0, "http"))))
		return cv, nil
	}))
	c.add(meth("getHttpPort", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewIntValue(propInt(rtSelf(ctx), "httpPort", 80)), nil
	}))
	c.add(meth("setHttpPort", []string{"httpPort"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "httpPort", data.NewIntValue(argInt(ctx, 0, 80)))
		return cv, nil
	}))
	c.add(meth("getHttpsPort", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewIntValue(propInt(rtSelf(ctx), "httpsPort", 443)), nil
	}))
	c.add(meth("setHttpsPort", []string{"httpsPort"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "httpsPort", data.NewIntValue(argInt(ctx, 0, 443)))
		return cv, nil
	}))
	c.add(meth("getQueryString", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "queryString", "")), nil
	}))
	c.add(meth("setQueryString", []string{"queryString"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "queryString", data.NewStringValue(argString(ctx, 0, "")))
		return cv, nil
	}))
	c.add(meth("getParameters", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "parameters"), nil
	}))
	c.add(meth("setParameters", []string{"parameters"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		setProp(cv, "parameters", argArray(ctx, 0))
		return cv, nil
	}))
	c.add(meth("getParameter", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		v, ok := assocGet(propArray(rtSelf(ctx), "parameters"), argString(ctx, 0, ""))
		if !ok {
			return data.NewNullValue(), nil
		}
		return v, nil
	}))
	c.add(meth("hasParameter", []string{"name"}, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewBoolValue(assocHas(propArray(rtSelf(ctx), "parameters"), argString(ctx, 0, ""))), nil
	}))
	c.add(meth("setParameter", []string{"name", "parameter"}, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		params := cloneArray(propArray(cv, "parameters"))
		assocSet(params, argString(ctx, 0, ""), arg(ctx, 1))
		setProp(cv, "parameters", params)
		return cv, nil
	}))
	c.add(meth("isSecure", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewBoolValue(propString(rtSelf(ctx), "scheme", "http") == "https"), nil
	}))
	return c
}

func contextConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	setProp(cv, "baseUrl", data.NewStringValue(strings.TrimRight(argString(ctx, 0, ""), "/")))
	setProp(cv, "method", data.NewStringValue(strings.ToUpper(argString(ctx, 1, "GET"))))
	setProp(cv, "host", data.NewStringValue(strings.ToLower(argString(ctx, 2, "localhost"))))
	setProp(cv, "scheme", data.NewStringValue(strings.ToLower(argString(ctx, 3, "http"))))
	setProp(cv, "httpPort", data.NewIntValue(argInt(ctx, 4, 80)))
	setProp(cv, "httpsPort", data.NewIntValue(argInt(ctx, 5, 443)))
	path := argString(ctx, 6, "/")
	if path == "" {
		path = "/"
	}
	setProp(cv, "pathInfo", data.NewStringValue(path))
	setProp(cv, "queryString", data.NewStringValue(argString(ctx, 7, "")))
	params := arg(ctx, 8)
	if params == nil || isNull(params) {
		setProp(cv, "parameters", phpList())
	} else {
		setProp(cv, "parameters", argArray(ctx, 8))
	}
	return data.NewNullValue(), nil
}

func contextFromUri(ctx data.Context) (data.GetValue, data.Control) {
	uri := argString(ctx, 0, "")
	host := argString(ctx, 1, "localhost")
	scheme := argString(ctx, 2, "http")
	httpPort := argInt(ctx, 3, 80)
	httpsPort := argInt(ctx, 4, 443)
	if i := strings.IndexByte(uri, '\\'); i >= 0 && i < strings.IndexAny(uri+"?#", "?#") {
		uri = ""
	}
	if uri != "" {
		if uri[0] <= 32 || uri[len(uri)-1] <= 32 || strings.ContainsAny(uri, "\r\n\t") {
			uri = ""
		}
	}
	u, err := url.Parse(uri)
	path := ""
	if err == nil && u != nil {
		if u.Scheme != "" {
			scheme = u.Scheme
		}
		if u.Hostname() != "" {
			host = u.Hostname()
		}
		if u.Port() != "" {
			if scheme == "http" {
				httpPort = atoiDef(u.Port(), httpPort)
			} else if scheme == "https" {
				httpsPort = atoiDef(u.Port(), httpsPort)
			}
		}
		path = u.EscapedPath()
		if path == "" {
			path = u.Path
		}
	}
	return instantiate(ctx, requestContextName,
		data.NewStringValue(path),
		data.NewStringValue("GET"),
		data.NewStringValue(host),
		data.NewStringValue(scheme),
		data.NewIntValue(httpPort),
		data.NewIntValue(httpsPort),
	)
}

func contextFromRequest(ctx data.Context) (data.GetValue, data.Control) {
	cv := rtSelf(ctx)
	req := asClassValue(arg(ctx, 0))
	if req == nil {
		return cv, nil
	}
	if v, ctl := callNamed(req, "getBaseUrl"); ctl != nil {
		return nil, ctl
	} else if v != nil {
		setProp(cv, "baseUrl", data.NewStringValue(strings.TrimRight(v.AsString(), "/")))
	}
	if v, ctl := callNamed(req, "getPathInfo"); ctl != nil {
		return nil, ctl
	} else if v != nil {
		setProp(cv, "pathInfo", v)
	}
	if v, ctl := callNamed(req, "getMethod"); ctl != nil {
		return nil, ctl
	} else if v != nil {
		setProp(cv, "method", data.NewStringValue(strings.ToUpper(v.AsString())))
	}
	if v, ctl := callNamed(req, "getHost"); ctl != nil {
		return nil, ctl
	} else if v != nil {
		setProp(cv, "host", data.NewStringValue(strings.ToLower(v.AsString())))
	}
	if v, ctl := callNamed(req, "getScheme"); ctl != nil {
		return nil, ctl
	} else if v != nil {
		setProp(cv, "scheme", data.NewStringValue(strings.ToLower(v.AsString())))
	}
	secure := false
	if v, _ := callNamed(req, "isSecure"); v != nil {
		if bv, ok := v.(data.AsBool); ok {
			secure, _ = bv.AsBool()
		}
	}
	portVal, _ := callNamed(req, "getPort")
	if !secure && portVal != nil && !isNull(portVal) {
		if iv, ok := portVal.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				setProp(cv, "httpPort", data.NewIntValue(n))
			}
		}
	}
	if secure && portVal != nil && !isNull(portVal) {
		if iv, ok := portVal.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				setProp(cv, "httpsPort", data.NewIntValue(n))
			}
		}
	}
	if server := prop(req, "server"); server != nil {
		if scv := asClassValue(server); scv != nil {
			if v, _ := callNamed(scv, "get", data.NewStringValue("QUERY_STRING"), data.NewStringValue("")); v != nil {
				setProp(cv, "queryString", data.NewStringValue(v.AsString()))
			}
		}
	}
	return cv, nil
}

func cloneRequestContext(ctx data.Context, src *data.ClassValue) (*data.ClassValue, data.Control) {
	return instantiate(ctx, requestContextName,
		data.NewStringValue(propString(src, "baseUrl", "")),
		data.NewStringValue(propString(src, "method", "GET")),
		data.NewStringValue(propString(src, "host", "localhost")),
		data.NewStringValue(propString(src, "scheme", "http")),
		data.NewIntValue(propInt(src, "httpPort", 80)),
		data.NewIntValue(propInt(src, "httpsPort", 443)),
		data.NewStringValue(propString(src, "pathInfo", "/")),
		data.NewStringValue(propString(src, "queryString", "")),
		cloneArray(propArray(src, "parameters")),
	)
}

func atoiDef(s string, def int) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return def
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}

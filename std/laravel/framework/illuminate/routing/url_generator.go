package routing

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const urlGeneratorClassName = "Illuminate\\Routing\\UrlGenerator"

// UrlGeneratorClass 对齐 Illuminate\Routing\UrlGenerator 热路径 API。
type UrlGeneratorClass struct {
	node.Node
	methods map[string]data.Method
}

func NewUrlGeneratorClass() data.ClassStmt {
	c := &UrlGeneratorClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *UrlGeneratorClass) GetName() string    { return urlGeneratorClassName }
func (c *UrlGeneratorClass) GetExtend() *string { return nil }
func (c *UrlGeneratorClass) GetImplements() []string {
	return []string{"Illuminate\\Contracts\\Routing\\UrlGenerator"}
}
func (c *UrlGeneratorClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "request", "routes", "assetRoot", "rootNamespace", "sessionResolver", "formatHostUsing", "formatPathUsing", "forcedRoot", "forceScheme", "cachedRoot", "cachedScheme", "defaultParameters", "routeGenerator", "missingNamedRouteResolver", "keyResolver":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *UrlGeneratorClass) GetPropertyList() []data.Property { return nil }
func (c *UrlGeneratorClass) GetConstruct() data.Method {
	return c.methods["__construct"]
}
func (c *UrlGeneratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *UrlGeneratorClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *UrlGeneratorClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *UrlGeneratorClass) GetStaticMethod(name string) (data.Method, bool) {
	return c.GetMethod(name)
}

func (c *UrlGeneratorClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"routes", "request", "assetRoot"}, 2, urlConstruct)
	c.methods["to"] = kit.InstanceMethodOpt("to", []string{"path", "extra", "secure"}, 1, urlTo)
	c.methods["asset"] = kit.InstanceMethodOpt("asset", []string{"path", "secure"}, 1, urlAsset)
	c.methods["secure"] = kit.InstanceMethodOpt("secure", []string{"path", "parameters"}, 1, urlSecure)
	c.methods["full"] = kit.InstanceMethod("full", nil, urlFull)
	c.methods["current"] = kit.InstanceMethod("current", nil, urlCurrent)
	c.methods["format"] = kit.InstanceMethod("format", []string{"root", "path"}, urlFormat)
	c.methods["forcescheme"] = kit.InstanceMethod("forceScheme", []string{"scheme"}, urlForceScheme)
	c.methods["forceroot"] = kit.InstanceMethod("forceRootUrl", []string{"root"}, urlForceRoot)
	c.methods["setrequest"] = kit.InstanceMethod("setRequest", []string{"request"}, urlSetRequest)
	c.registerRouteMethods()
}

func urlRecv(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("UrlGenerator missing $this"))
}

func urlConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("routes", kit.Arg(ctx, 0))
	_ = cv.SetProperty("request", kit.Arg(ctx, 1))
	if ar := kit.Arg(ctx, 2); ar != nil {
		_ = cv.SetProperty("assetRoot", ar)
	}
	return data.NewNullValue(), nil
}

func requestRoot(cv *data.ClassValue) string {
	if v, _ := cv.GetProperty("forcedRoot"); v != nil && !kit.IsNull(v) {
		return strings.TrimRight(v.AsString(), "/")
	}
	req, _ := cv.GetProperty("request")
	if rcv, ok := kit.Unwrap(req).(*data.ClassValue); ok {
		if m, ok := rcv.GetMethod("root"); ok && m != nil {
			nctx := rcv.CreateContext(m.GetVariables())
			ret, ctl := m.Call(nctx)
			if ctl == nil && ret != nil {
				if val, ok := ret.(data.Value); ok {
					return strings.TrimRight(val.AsString(), "/")
				}
			}
		}
		if m, ok := rcv.GetMethod("getSchemeAndHttpHost"); ok && m != nil {
			nctx := rcv.CreateContext(m.GetVariables())
			ret, ctl := m.Call(nctx)
			if ctl == nil && ret != nil {
				if val, ok := ret.(data.Value); ok {
					return strings.TrimRight(val.AsString(), "/")
				}
			}
		}
	}
	return ""
}

func urlTo(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return urlToPath(cv, kit.Arg(ctx, 0).AsString(), kit.Arg(ctx, 1), kit.Arg(ctx, 2))
}

func urlAsset(ctx data.Context) (data.GetValue, data.Control) {
	return urlTo(ctx)
}

func urlSecure(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	path := kit.Arg(ctx, 0).AsString()
	root := requestRoot(cv)
	if u, err := url.Parse(root); err == nil && u.Scheme != "" {
		u.Scheme = "https"
		root = strings.TrimRight(u.String(), "/")
	} else if root != "" {
		root = "https://" + strings.TrimPrefix(strings.TrimPrefix(root, "http://"), "https://")
	}
	path = strings.TrimLeft(path, "/")
	return data.NewStringValue(root + "/" + path), nil
}

func urlFull(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	req, _ := cv.GetProperty("request")
	if rcv, ok := kit.Unwrap(req).(*data.ClassValue); ok {
		if m, ok := rcv.GetMethod("fullUrl"); ok && m != nil {
			nctx := rcv.CreateContext(m.GetVariables())
			return m.Call(nctx)
		}
	}
	return data.NewStringValue(""), nil
}

func urlCurrent(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	req, _ := cv.GetProperty("request")
	if rcv, ok := kit.Unwrap(req).(*data.ClassValue); ok {
		if m, ok := rcv.GetMethod("url"); ok && m != nil {
			nctx := rcv.CreateContext(m.GetVariables())
			return m.Call(nctx)
		}
	}
	return data.NewStringValue(""), nil
}

func urlFormat(ctx data.Context) (data.GetValue, data.Control) {
	root := strings.TrimRight(kit.Arg(ctx, 0).AsString(), "/")
	path := strings.TrimLeft(kit.Arg(ctx, 1).AsString(), "/")
	if path == "" {
		return data.NewStringValue(root), nil
	}
	return data.NewStringValue(root + "/" + path), nil
}

func urlForceScheme(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	scheme := kit.Arg(ctx, 0).AsString()
	if scheme != "" {
		_ = cv.SetProperty("forceScheme", data.NewStringValue(scheme+"://"))
	} else {
		_ = cv.SetProperty("forceScheme", data.NewNullValue())
	}
	_ = cv.SetProperty("cachedScheme", data.NewNullValue())
	return data.NewNullValue(), nil
}

func urlForceRoot(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("forcedRoot", kit.Arg(ctx, 0))
	return data.NewNullValue(), nil
}

func urlSetRequest(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := urlRecv(ctx)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("request", kit.Arg(ctx, 0))
	return data.NewNullValue(), nil
}

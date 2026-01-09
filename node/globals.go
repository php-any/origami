package node

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/php-any/origami/data"
)

var groups sync.Map

// GlobalsNode 表示超级全局变量节点
// 支持：$_SERVER, $_GET, $_POST, $_COOKIE, $_SESSION, $_REQUEST, $_ENV, $GLOBALS
type GlobalsNode struct {
	*Node `pp:"-"`
	Name  string
}

func NewGlobalsNode(from data.From, name string) *GlobalsNode {
	return &GlobalsNode{Node: NewNode(from), Name: name}
}

func (g *GlobalsNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 根据变量名返回对应的全局变量
	switch g.Name {
	case "$GLOBALS":
		return g.getGlobals(ctx)
	case "$_ENV":
		return g.getEnv()
	case "$_SERVER":
		return g.getServer(ctx)
	case "$_GET":
		return g.getGet(ctx)
	case "$_POST":
		return g.getPost(ctx)
	case "$_COOKIE":
		return g.getCookie(ctx)
	case "$_SESSION":
		return g.getSession(ctx)
	case "$_REQUEST":
		return g.getRequest(ctx)
	default:
		// 默认返回空对象
		v, ok := groups.Load(g.Name)
		if !ok {
			v = data.NewObjectValue()
			groups.Store(g.Name, v)
		}
		return v.(data.GetValue), nil
	}
}

// getGlobals 获取 $GLOBALS
func (g *GlobalsNode) getGlobals(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load("$GLOBALS")
	if !ok {
		v = data.NewObjectValue()
		groups.Store("$GLOBALS", v)
	}
	return v.(data.GetValue), nil
}

// getEnv 获取 $_ENV（环境变量）
func (g *GlobalsNode) getEnv() (data.GetValue, data.Control) {
	v, ok := groups.Load("$_ENV")
	var envObj *data.ObjectValue
	if !ok {
		// 如果 $_ENV 不存在，创建新的对象
		envObj = data.NewObjectValue()
		groups.Store("$_ENV", envObj)
	} else {
		// 如果 $_ENV 已存在，保留现有值
		if obj, ok := v.(*data.ObjectValue); ok {
			envObj = obj
		} else {
			// 如果不是 ObjectValue 类型，创建新的
			envObj = data.NewObjectValue()
			groups.Store("$_ENV", envObj)
		}
	}

	// 合并系统环境变量（保留已有的 $_ENV 值，但用系统环境变量更新或新增）
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envObj.SetProperty(parts[0], data.NewStringValue(parts[1]))
		}
	}

	return envObj, nil
}

// getServer 获取 $_SERVER（服务器和执行环境信息）
func (g *GlobalsNode) getServer(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load("$_SERVER")
	if !ok {
		serverObj := data.NewObjectValue()
		// 尝试从上下文获取 HTTP 请求
		if httpReq := g.getHTTPRequest(ctx); httpReq != nil {
			// 填充 HTTP 请求相关的服务器信息
			serverObj.SetProperty("REQUEST_METHOD", data.NewStringValue(httpReq.Method))
			serverObj.SetProperty("REQUEST_URI", data.NewStringValue(httpReq.RequestURI))
			serverObj.SetProperty("QUERY_STRING", data.NewStringValue(httpReq.URL.RawQuery))
			serverObj.SetProperty("HTTP_HOST", data.NewStringValue(httpReq.Host))
			serverObj.SetProperty("SERVER_NAME", data.NewStringValue(httpReq.Host))
			serverObj.SetProperty("SERVER_PORT", data.NewStringValue(httpReq.URL.Port()))
			serverObj.SetProperty("REMOTE_ADDR", data.NewStringValue(httpReq.RemoteAddr))
			serverObj.SetProperty("SCRIPT_NAME", data.NewStringValue(httpReq.URL.Path))
			serverObj.SetProperty("PATH_INFO", data.NewStringValue(httpReq.URL.Path))

			// 添加所有请求头，格式为 HTTP_HEADER_NAME
			for key, values := range httpReq.Header {
				if len(values) > 0 {
					headerKey := "HTTP_" + strings.ReplaceAll(strings.ToUpper(key), "-", "_")
					serverObj.SetProperty(headerKey, data.NewStringValue(values[0]))
				}
			}
		} else {
			// 没有 HTTP 请求时，填充一些基本系统信息
			serverObj.SetProperty("SERVER_SOFTWARE", data.NewStringValue("Origami"))
		}
		groups.Store("$_SERVER", serverObj)
		return serverObj, nil
	}
	return v.(data.GetValue), nil
}

// getGet 获取 $_GET（URL 查询参数）
func (g *GlobalsNode) getGet(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load("$_GET")
	if !ok {
		getObj := data.NewObjectValue()
		if httpReq := g.getHTTPRequest(ctx); httpReq != nil {
			for key, values := range httpReq.URL.Query() {
				if len(values) == 1 {
					getObj.SetProperty(key, data.NewStringValue(values[0]))
				} else {
					// 多个值，使用第一个
					getObj.SetProperty(key, data.NewStringValue(values[0]))
				}
			}
		}
		groups.Store("$_GET", getObj)
		return getObj, nil
	}
	return v.(data.GetValue), nil
}

// getPost 获取 $_POST（POST 表单数据）
func (g *GlobalsNode) getPost(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load("$_POST")
	if !ok {
		postObj := data.NewObjectValue()
		if httpReq := g.getHTTPRequest(ctx); httpReq != nil {
			if httpReq.Form != nil {
				for key, values := range httpReq.Form {
					if len(values) > 0 {
						postObj.SetProperty(key, data.NewStringValue(values[0]))
					}
				}
			}
		}
		groups.Store("$_POST", postObj)
		return postObj, nil
	}
	return v.(data.GetValue), nil
}

// getCookie 获取 $_COOKIE
func (g *GlobalsNode) getCookie(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load("$_COOKIE")
	if !ok {
		cookieObj := data.NewObjectValue()
		if httpReq := g.getHTTPRequest(ctx); httpReq != nil {
			for _, cookie := range httpReq.Cookies() {
				cookieObj.SetProperty(cookie.Name, data.NewStringValue(cookie.Value))
			}
		}
		groups.Store("$_COOKIE", cookieObj)
		return cookieObj, nil
	}
	return v.(data.GetValue), nil
}

// getSession 获取 $_SESSION
func (g *GlobalsNode) getSession(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load("$_SESSION")
	if !ok {
		sessionObj := data.NewObjectValue()
		groups.Store("$_SESSION", sessionObj)
		return sessionObj, nil
	}
	return v.(data.GetValue), nil
}

// getRequest 获取 $_REQUEST（合并 GET、POST、COOKIE）
func (g *GlobalsNode) getRequest(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load("$_REQUEST")
	if !ok {
		requestObj := data.NewObjectValue()

		// 先添加 GET 参数
		if getVal, _ := g.getGet(ctx); getVal != nil {
			if getObj, ok := getVal.(*data.ObjectValue); ok {
				getObj.RangeProperties(func(key string, value data.Value) bool {
					requestObj.SetProperty(key, value)
					return true
				})
			}
		}

		// 再添加 POST 参数（会覆盖 GET 中的同名参数）
		if postVal, _ := g.getPost(ctx); postVal != nil {
			if postObj, ok := postVal.(*data.ObjectValue); ok {
				postObj.RangeProperties(func(key string, value data.Value) bool {
					requestObj.SetProperty(key, value)
					return true
				})
			}
		}

		// 最后添加 COOKIE 参数（会覆盖前面的同名参数）
		if cookieVal, _ := g.getCookie(ctx); cookieVal != nil {
			if cookieObj, ok := cookieVal.(*data.ObjectValue); ok {
				cookieObj.RangeProperties(func(key string, value data.Value) bool {
					requestObj.SetProperty(key, value)
					return true
				})
			}
		}

		groups.Store("$_REQUEST", requestObj)
		return requestObj, nil
	}
	return v.(data.GetValue), nil
}

// getHTTPRequest 尝试从上下文中获取 HTTP 请求
func (g *GlobalsNode) getHTTPRequest(ctx data.Context) *http.Request {
	// 尝试从 Go context 中获取 HTTP 请求
	if goCtx := ctx.GoContext(); goCtx != nil {
		if req, ok := goCtx.Value("http_request").(*http.Request); ok {
			return req
		}
	}

	// 尝试从上下文中查找名为 "r" 的变量（HTTP handler 中通常使用这个名称）
	// 由于 Context 接口没有提供按名称查找的方法，我们尝试通过索引查找
	// HTTP handler 通常将 request 放在索引 0
	if reqVal, ok := ctx.GetIndexValue(0); ok {
		// 尝试从 ProxyValue 中提取 RequestClass
		if proxyVal, ok := reqVal.(*data.ProxyValue); ok {
			if classStmt, ok := proxyVal.Class.(interface{ GetSource() any }); ok {
				if source := classStmt.GetSource(); source != nil {
					if httpReq, ok := source.(*http.Request); ok {
						return httpReq
					}
				}
			}
		}
	}

	return nil
}

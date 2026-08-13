package node

import (
	"net/http"

	"github.com/php-any/origami/data"
)

// HTTPResponseWriter 保留旧 API；请求级响应必须由 Context 绑定的 VM 提供。
func HTTPResponseWriter() http.ResponseWriter {
	return nil
}

// getHTTPRequest 尝试从上下文中获取 HTTP 请求（供多个超全局节点复用）
func getHTTPRequest(ctx data.Context) *http.Request {
	if goCtx := ctx.GoContext(); goCtx != nil {
		if req, ok := goCtx.Value("http_request").(*http.Request); ok {
			return req
		}
	}

	if reqVal, ok := ctx.GetIndexValue(0); ok {
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

	if ctx != nil {
		if host, ok := ctx.GetVM().(interface{ HTTPRequest() *http.Request }); ok {
			if req := host.HTTPRequest(); req != nil {
				return req
			}
		}
	}

	return nil
}

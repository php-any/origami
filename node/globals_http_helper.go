package node

import (
	"net/http"

	"github.com/php-any/origami/data"
)

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

	return nil
}

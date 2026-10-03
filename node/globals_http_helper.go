package node

import (
	"bytes"
	"fmt"
	"io"
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

// superglobalArray initializes once per VM/request; mutations are retained in
// that request's global slot instead of rebuilding host data on every read.
func superglobalArray(ctx data.Context, name string, init func() *data.ArrayValue) data.Value {
	if ctx != nil && ctx.GetVM() != nil {
		slot := ctx.GetVM().EnsureGlobalZVal(name)
		if slot.Defined {
			return slot.ReadValue()
		}
		array := init()
		data.CowAssign(slot, array)
		return array
	}
	return init()
}

func setSuperglobalArray(ctx data.Context, name string, value data.Value) data.Control {
	if ctx != nil && ctx.GetVM() != nil {
		data.CowAssign(ctx.GetVM().EnsureGlobalZVal(name), value)
		return nil
	}
	return data.NewErrorThrow(nil, fmt.Errorf("superglobal assignment requires a VM"))
}

// HTTPInputBody preserves the SAPI input stream for repeated reads and POST
// decoding. Host adapters expose the request through the execution VM.
func HTTPInputBody(ctx data.Context) ([]byte, bool, error) {
	request := getHTTPRequest(ctx)
	if request == nil {
		return nil, false, nil
	}
	if request.Body == nil {
		return nil, true, nil
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, true, err
	}
	request.Body = io.NopCloser(bytes.NewReader(body))
	return body, true, nil
}

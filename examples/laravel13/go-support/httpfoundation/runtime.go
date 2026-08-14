package httpfoundation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/php-any/origami/data"
)

// SendResponse 调用 Symfony Response::send，将 Header 与 Body 写入当前 Go writer。
// 适用于通过 VM 输出通道（OutputSink）发送响应的场景。
func SendResponse(response data.GetValue) (*data.ClassValue, data.Control) {
	value, ok := response.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpfoundation: Kernel 返回值不是 Response"))
	}
	method, ok := value.GetMethod("send")
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpfoundation: Response 缺少 send 方法"))
	}
	callCtx := value.CreateContext(method.GetVariables())
	result, control := method.Call(callCtx)
	if control != nil {
		return nil, control
	}
	if sent, ok := result.(*data.ClassValue); ok {
		return sent, nil
	}
	return value, nil
}

// SendResponseTo 直接将 Response 的 Header 与 Body 写入指定 http.ResponseWriter。
// 适用于普通模式 web 服务（不依赖 VM 输出通道 / RequestVM）的场景。
func SendResponseTo(w http.ResponseWriter, response data.GetValue) (*data.ClassValue, data.Control) {
	value, ok := response.(*data.ClassValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpfoundation: Kernel 返回值不是 Response"))
	}

	// 写入 headers
	if headers := responseHeaders(value); headers != nil {
		all := GetHeaderBagAll(headers)
		for name, values := range all {
			if strings.EqualFold(name, "set-cookie") {
				continue
			}
			replace := strings.EqualFold(name, "Content-Type")
			for _, v := range values {
				if replace {
					w.Header().Set(http.CanonicalHeaderKey(name), v)
					replace = false
				} else {
					w.Header().Add(http.CanonicalHeaderKey(name), v)
				}
			}
		}
		// cookies
		if rh := ResponseHeaderBagFrom(headers); rh != nil {
			rh.mu.RLock()
			for _, byPath := range rh.cookies {
				for _, byName := range byPath {
					for _, cookie := range byName {
						if cookie != nil {
							w.Header().Add("Set-Cookie", cookie.String())
						}
					}
				}
			}
			rh.mu.RUnlock()
		}
	}

	// 设置 status code 并写入 content
	statusCode := responseStatusCode(value)
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	content := responseContent(value)
	w.WriteHeader(statusCode)
	if _, err := io.WriteString(w, content); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}

	return value, nil
}

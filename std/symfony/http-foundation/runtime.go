package httpfoundation

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

func classValueFromGetValue(v data.GetValue) *data.ClassValue {
	switch t := v.(type) {
	case *data.ClassValue:
		return t
	case *data.ThisValue:
		if t != nil {
			return t.ClassValue
		}
	}
	return nil
}

func callResponseMethod(cv *data.ClassValue, name string) (data.GetValue, data.Control) {
	if cv == nil {
		return nil, nil
	}
	method, ok := cv.GetMethod(name)
	if !ok {
		return nil, nil
	}
	return method.Call(cv.CreateContext(method.GetVariables()))
}

func binaryFilePath(cv *data.ClassValue) string {
	if cv == nil || cv.Class == nil {
		return ""
	}
	if !strings.Contains(cv.Class.GetName(), "BinaryFileResponse") {
		return ""
	}
	fileVal, ctl := callResponseMethod(cv, "getFile")
	if ctl != nil || fileVal == nil {
		return ""
	}
	fileCV := classValueFromGetValue(fileVal)
	if fileCV == nil {
		if s, ok := fileVal.(data.Value); ok {
			p := s.AsString()
			if p != "" {
				return p
			}
		}
		return ""
	}
	pathVal, ctl := callResponseMethod(fileCV, "getPathname")
	if ctl != nil || pathVal == nil {
		if s, ok := fileVal.(data.Value); ok {
			return s.AsString()
		}
		return ""
	}
	if s, ok := pathVal.(data.Value); ok {
		return s.AsString()
	}
	return ""
}

// SendResponse 调用 Symfony Response::send，将 Header 与 Body 写入当前 Go writer。
// 适用于通过 VM 输出通道（OutputSink）发送响应的场景。
func SendResponse(response data.GetValue) (*data.ClassValue, data.Control) {
	value := classValueFromGetValue(response)
	if value == nil {
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
	if sent := classValueFromGetValue(result); sent != nil {
		return sent, nil
	}
	return value, nil
}

// SendResponseTo 直接将 Response 的 Header 与 Body 写入指定 http.ResponseWriter。
// leftover 对齐 Octane/Swoole：请求级 ob_start 残留的 echo，写在 getContent() 之前。
func SendResponseTo(w http.ResponseWriter, response data.GetValue, leftover string) (*data.ClassValue, data.Control) {
	value := classValueFromGetValue(response)
	if value == nil {
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

	if filePath := binaryFilePath(value); filePath != "" {
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			http.ServeFile(w, &http.Request{Method: http.MethodGet}, filePath)
			return value, nil
		}
	}

	// 设置 status code 并写入 content
	statusCode := responseStatusCode(value)
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	content := responseContent(value)
	if content == "" {
		if got, ctl := callResponseMethod(value, "getContent"); ctl == nil && got != nil {
			if s, ok := got.(data.Value); ok {
				if _, isNull := s.(*data.NullValue); !isNull {
					if b, ok := s.(*data.BoolValue); !ok || b.Value {
						content = s.AsString()
					}
				}
			}
		}
	}
	if leftover != "" {
		content = leftover + content
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(statusCode)
	if _, err := io.WriteString(w, content); err != nil {
		if isClientAbortError(err) {
			return value, nil
		}
		return nil, data.NewErrorThrow(nil, err)
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	return value, nil
}

func isClientAbortError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "wsasend") ||
		strings.Contains(msg, "forcibly closed") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "aborted by the software")
}

// IsClientAbortControl 判断控制流是否来自客户端断开（写响应时连接已关）。
func IsClientAbortControl(ctl data.Control) bool {
	if ctl == nil {
		return false
	}
	return isClientAbortError(fmt.Errorf("%s", ctl.AsString()))
}

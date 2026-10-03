package httpfoundation

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/utils"
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
	if _, ok := cv.GetMethod(name); !ok {
		return nil, nil
	}
	return node.NewObjectMethod(nil, cv, name, nil).GetValue(cv.Context)
}

func binaryFilePath(cv *data.ClassValue) string {
	if cv == nil || cv.Class == nil {
		return ""
	}
	if !responseIsA(cv, fqnBinaryFileResponse) {
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
// leftover 对齐 Octane/Swoole：请求级 ob_start 残留的 echo，写在响应内容之前。
func SendResponseTo(w http.ResponseWriter, response data.GetValue, leftover string) (*data.ClassValue, data.Control) {
	return SendResponseForRequest(w, nil, response, leftover)
}

// SendResponseForRequest retains the real HTTP method for HEAD and file responses.
func SendResponseForRequest(w http.ResponseWriter, request *http.Request, response data.GetValue, leftover string) (*data.ClassValue, data.Control) {
	value := classValueFromGetValue(response)
	if value == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("httpfoundation: Kernel 返回值不是 Response"))
	}

	// 写入 headers
	if headers := responseHeaders(value); headers != nil {
		allValue, ctl := callResponseMethod(headers, "all")
		if ctl != nil {
			return value, ctl
		}
		all := headersMapFromValue(asResponseValue(allValue))
		for name, values := range all {
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
	}

	// 设置 status code 并写入 content
	statusCode := responseStatusCode(value)
	if status, ctl := callResponseMethod(value, "getStatusCode"); ctl != nil {
		return value, ctl
	} else if integer, ok := status.(data.AsInt); ok {
		statusCode, _ = integer.AsInt()
	}
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	if statusCode < 200 || statusCode == http.StatusNoContent || statusCode == http.StatusNotModified {
		w.Header().Del("Content-Length")
		w.Header().Del("Transfer-Encoding")
		w.WriteHeader(statusCode)
		return value, nil
	}
	if request != nil && request.Method == http.MethodHead {
		w.WriteHeader(statusCode)
		return value, nil
	}
	_, officialPHP := value.Class.(*node.ClassStatement)
	phpFile := officialPHP && responseIsA(value, fqnBinaryFileResponse)
	if filePath := binaryFilePath(value); filePath != "" && !phpFile {
		fileContext := runtime.RequestContext()
		if request != nil {
			fileContext = request.Context()
		}
		file, closeFile, err := utils.OpenRequestFile(fileContext, filePath, os.O_RDONLY, 0)
		if err != nil {
			return nil, data.NewErrorThrow(nil, err)
		}
		defer closeFile()
		defer data.CheckRequest(fileContext)
		if cookiePropBool(value, "deleteFileAfterSend", false) {
			defer os.Remove(filePath)
		}
		info, err := file.Stat()
		if err != nil {
			return nil, data.NewErrorThrow(nil, err)
		}
		if request == nil {
			request = &http.Request{Method: http.MethodGet, Header: make(http.Header)}
		}
		// ServeContent consumes Range/If-Range and conditional request headers.
		// A non-200 status belongs to Symfony and must not be replaced with 200.
		if statusCode == http.StatusOK {
			http.ServeContent(w, request, info.Name(), info.ModTime(), file)
		} else {
			w.WriteHeader(statusCode)
			if statusCode >= 200 && statusCode < 300 {
				if _, err := io.Copy(w, file); err != nil {
					return value, data.NewErrorThrow(nil, err)
				}
			}
		}
		return value, nil
	}
	// PHP responses own their body semantics, including subclass sendContent overrides.
	if officialPHP || responseIsA(value, fqnStreamedResponse) {
		if !phpFile {
			w.Header().Del("Content-Length")
		}
		w.WriteHeader(statusCode)
		if _, err := io.WriteString(w, leftover); err != nil {
			return value, data.NewErrorThrow(nil, err)
		}
		host, ok := value.GetVM().(data.OutputTargetHost)
		if !ok {
			return value, data.NewErrorThrow(nil, fmt.Errorf("httpfoundation: VM lacks a response output target"))
		}
		restore := host.BindOutputTarget(func(s string) data.Control {
			_, err := io.WriteString(w, s)
			if err != nil {
				return data.NewErrorThrow(nil, err)
			}
			return nil
		}, func() {
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		})
		defer restore()
		_, ctl := callResponseMethod(value, "sendContent")
		return value, ctl
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

func asResponseValue(value data.GetValue) data.Value {
	if actual, ok := value.(data.Value); ok {
		return actual
	}
	return nil
}

func responseIsA(value *data.ClassValue, name string) bool {
	if value == nil {
		return false
	}
	for class := value.Class; class != nil; {
		if strings.EqualFold(strings.TrimPrefix(class.GetName(), "\\"), name) {
			return true
		}
		parent := class.GetExtend()
		if parent == nil || value.GetVM() == nil {
			return false
		}
		class, _ = value.GetVM().GetClass(*parent)
	}
	return false
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

package httpfoundation

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SendHeaderFunc 写出单条 HTTP 头（对齐 PHP header($header, $replace, $responseCode)）。
// 默认走 node.HTTPResponseWriter；测试或嵌入方可覆盖。
var SendHeaderFunc = defaultSendHeader

// SendHeaderContext 优先从当前 VM（如 LaravelRequestVM）取 ResponseWriter。
func SendHeaderContext(ctx data.Context, header string, replace bool, responseCode int) {
	if w := responseWriterFromContext(ctx); w != nil {
		writeHeaderLine(w, header, replace, responseCode)
		return
	}
	SendHeaderFunc(header, replace, responseCode)
}

func defaultSendHeader(header string, replace bool, responseCode int) {
	writer := node.HTTPResponseWriter()
	if writer == nil {
		return
	}
	writeHeaderLine(writer, header, replace, responseCode)
}

type httpResponseHost interface {
	HTTPResponseWriter() http.ResponseWriter
}

func responseWriterFromContext(ctx data.Context) http.ResponseWriter {
	if ctx != nil {
		if host, ok := ctx.GetVM().(httpResponseHost); ok {
			if w := host.HTTPResponseWriter(); w != nil {
				return w
			}
		}
	}
	return node.HTTPResponseWriter()
}

func writeHeaderLine(writer http.ResponseWriter, header string, replace bool, responseCode int) {
	line := header
	if strings.HasPrefix(strings.ToUpper(line), "HTTP/") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			if code, err := strconv.Atoi(parts[1]); err == nil {
				writer.WriteHeader(code)
			}
		}
	} else if name, value, ok := strings.Cut(line, ":"); ok {
		name = http.CanonicalHeaderKey(strings.TrimSpace(name))
		value = strings.TrimSpace(value)
		if replace {
			writer.Header().Set(name, value)
		} else {
			writer.Header().Add(name, value)
		}
	}
	if responseCode > 0 {
		writer.WriteHeader(responseCode)
	}
}

func responseClassValue(ctx data.Context) *data.ClassValue {
	return bagClassValue(ctx)
}

func responseSelf(ctx data.Context) data.GetValue {
	return thisClassValue(ctx)
}

func responseHeaders(cv *data.ClassValue) *data.ClassValue {
	if cv == nil {
		return nil
	}
	v, _ := cv.GetProperty("headers")
	if h, ok := v.(*data.ClassValue); ok {
		return h
	}
	return nil
}

func responseGetStringProp(cv *data.ClassValue, name, def string) string {
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	if v == nil {
		return def
	}
	if _, ok := v.(*data.NullValue); ok {
		return def
	}
	return v.AsString()
}

func responseGetIntProp(cv *data.ClassValue, name string, def int) int {
	if cv == nil {
		return def
	}
	v, _ := cv.GetProperty(name)
	if v == nil {
		return def
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	n, err := strconv.Atoi(strings.TrimSpace(v.AsString()))
	if err != nil {
		return def
	}
	return n
}

func responseSetProp(cv *data.ClassValue, name string, value data.Value) {
	if cv != nil {
		cv.SetProperty(name, value)
	}
}

func responseStatusCode(cv *data.ClassValue) int {
	return responseGetIntProp(cv, "statusCode", 200)
}

func responseStatusText(cv *data.ClassValue) string {
	return responseGetStringProp(cv, "statusText", "")
}

func responseVersion(cv *data.ClassValue) string {
	return responseGetStringProp(cv, "version", "1.0")
}

func responseContent(cv *data.ClassValue) string {
	return responseGetStringProp(cv, "content", "")
}

func responseCharset(cv *data.ClassValue) (string, bool) {
	if cv == nil {
		return "", false
	}
	v, _ := cv.GetProperty("charset")
	if v == nil {
		return "", false
	}
	if _, ok := v.(*data.NullValue); ok {
		return "", false
	}
	s := v.AsString()
	return s, true
}

func headersMapFromValue(v data.Value) map[string][]string {
	out := map[string][]string{}
	if v == nil {
		return out
	}
	if cv, ok := v.(*data.ClassValue); ok {
		if ResponseHeaderBagFrom(cv) != nil || HeaderBagFrom(cv) != nil {
			return GetHeaderBagAll(cv)
		}
	}
	m, err := valueToAssocMap(v)
	if err != nil {
		return out
	}
	for k, val := range m {
		if arr, ok := val.(*data.ArrayValue); ok {
			vals := make([]string, 0, len(arr.List))
			for _, z := range arr.List {
				if z != nil && z.Value != nil {
					vals = append(vals, z.Value.AsString())
				}
			}
			out[k] = vals
		} else if val != nil {
			out[k] = []string{val.AsString()}
		}
	}
	return out
}

func createResponseHeaders(ctx data.Context, v data.Value) *data.ClassValue {
	if cv, ok := v.(*data.ClassValue); ok {
		if ResponseHeaderBagFrom(cv) != nil {
			return cv
		}
	}
	return NewResponseHeaderBagValue(ctx, headersMapFromValue(v))
}

func respHeaderHas(headers *data.ClassValue, key string) bool {
	h := HeaderBagFrom(headers)
	if h == nil {
		return false
	}
	lk := strings.ToLower(key)
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.headers[lk]
	return ok
}

func respHeaderGet(headers *data.ClassValue, key string, def string) string {
	h := HeaderBagFrom(headers)
	if h == nil {
		return def
	}
	lk := strings.ToLower(key)
	h.mu.RLock()
	defer h.mu.RUnlock()
	vs, ok := h.headers[lk]
	if !ok || len(vs) == 0 {
		return def
	}
	if vs[0] == nil {
		return def
	}
	return *vs[0]
}

func respHeaderGetAll(headers *data.ClassValue, key string) []string {
	h := HeaderBagFrom(headers)
	if h == nil {
		return nil
	}
	lk := strings.ToLower(key)
	h.mu.RLock()
	defer h.mu.RUnlock()
	vs, ok := h.headers[lk]
	if !ok {
		return nil
	}
	return nullableToStringSlice(vs)
}

func respHeaderSet(headers *data.ClassValue, key string, values []string, replace bool) {
	SetHeaderBagHeader(headers, key, values, replace)
}

func respHeaderRemove(headers *data.ClassValue, key string) {
	lk := strings.ToLower(key)
	if rh := ResponseHeaderBagFrom(headers); rh != nil {
		rh.mu.Lock()
		delete(rh.headers, lk)
		delete(rh.headerNames, lk)
		for i, k := range rh.keys {
			if strings.ToLower(k) == lk || k == lk {
				rh.keys = append(rh.keys[:i], rh.keys[i+1:]...)
				break
			}
		}
		rh.mu.Unlock()
		return
	}
	if h := HeaderBagFrom(headers); h != nil {
		h.mu.Lock()
		delete(h.headers, lk)
		for i, k := range h.keys {
			if strings.ToLower(k) == lk || k == lk {
				h.keys = append(h.keys[:i], h.keys[i+1:]...)
				break
			}
		}
		h.mu.Unlock()
	}
}

func respHasCC(headers *data.ClassValue, directive string) bool {
	h := HeaderBagFrom(headers)
	if h == nil {
		return false
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.cacheControl[directive]
	return ok
}

func respGetCC(headers *data.ClassValue, directive string) (any, bool) {
	h := HeaderBagFrom(headers)
	if h == nil {
		return nil, false
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	v, ok := h.cacheControl[directive]
	return v, ok
}

func respAddCC(headers *data.ClassValue, directive string, value any) {
	h := HeaderBagFrom(headers)
	if h == nil {
		return
	}
	h.mu.Lock()
	h.cacheControl[directive] = value
	h.mu.Unlock()
	respRefreshCacheControl(headers)
}

func respRemoveCC(headers *data.ClassValue, directive string) {
	h := HeaderBagFrom(headers)
	if h == nil {
		return
	}
	h.mu.Lock()
	delete(h.cacheControl, directive)
	h.mu.Unlock()
	respRefreshCacheControl(headers)
}

func respRefreshCacheControl(headers *data.ClassValue) {
	h := HeaderBagFrom(headers)
	if h == nil {
		return
	}
	h.mu.RLock()
	parts := make([]string, 0, len(h.cacheControl))
	for k, v := range h.cacheControl {
		switch t := v.(type) {
		case bool:
			if t {
				parts = append(parts, k)
			}
		case string:
			if t == "" {
				parts = append(parts, k)
			} else {
				parts = append(parts, k+"="+t)
			}
		case int:
			parts = append(parts, k+"="+strconv.Itoa(t))
		default:
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
	}
	h.mu.RUnlock()
	if len(parts) == 0 {
		respHeaderRemove(headers, "Cache-Control")
		return
	}
	respHeaderSet(headers, "Cache-Control", []string{strings.Join(parts, ", ")}, true)
}

func respFormatHTTPDate(t time.Time) string {
	return t.UTC().Format(http.TimeFormat)
}

func respParseHTTPDate(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	if t, err := http.ParseTime(s); err == nil {
		return t, true
	}
	layouts := []string{
		time.RFC1123,
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Monday, 02-Jan-06 15:04:05 MST",
		"Mon Jan 2 15:04:05 2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}

func statusTextFor(code int) string {
	if t, ok := responseStatusTexts[code]; ok {
		return t
	}
	return "unknown status"
}

func optionalIntParam(ctx data.Context, index int) (int, bool) {
	v, ok := ctx.GetIndexValue(index)
	if !ok || v == nil {
		return 0, false
	}
	if _, isNull := v.(*data.NullValue); isNull {
		return 0, false
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n, true
		}
	}
	n, err := strconv.Atoi(strings.TrimSpace(v.AsString()))
	if err != nil {
		return 0, false
	}
	return n, true
}

func optionalBoolParam(ctx data.Context, index int, def bool) bool {
	return boolParam(ctx, index, def)
}

func callObjMethod(obj *data.ClassValue, name string, args ...data.Value) (data.GetValue, data.Control) {
	if obj == nil {
		return data.NewNullValue(), nil
	}
	m, ok := obj.GetMethod(name)
	if !ok || m == nil {
		return data.NewNullValue(), nil
	}
	vars := m.GetVariables()
	fnCtx := obj.CreateContext(vars)
	for i, arg := range args {
		if i < len(vars) {
			fnCtx.SetVariableValue(vars[i], arg)
		}
	}
	return m.Call(fnCtx)
}

func valueIsTruthy(v data.GetValue) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(data.AsBool); ok {
		okv, err := b.AsBool()
		return err == nil && okv
	}
	if sv, ok := v.(data.AsString); ok {
		s := sv.AsString()
		return s != "" && s != "0" && s != "false"
	}
	return true
}

func protectedProp(name string, def data.GetValue) data.Property {
	return node.NewProperty(nil, name, "protected", false, def)
}

func publicProp(name string, def data.GetValue) data.Property {
	return node.NewProperty(nil, name, "public", false, def)
}

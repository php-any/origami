package httpfoundation

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

// ActiveRequestFromContext 从当前请求 VM 取出 Go 请求。
func ActiveRequestFromContext(ctx data.Context) *http.Request {
	if ctx != nil {
		if host, ok := ctx.GetVM().(interface{ HTTPRequest() *http.Request }); ok {
			if req := host.HTTPRequest(); req != nil {
				return req
			}
		}
	}
	return nil
}

func cloneHTTPRequest(request *http.Request) *http.Request {
	if request == nil {
		request = &http.Request{
			Method:     http.MethodGet,
			URL:        &url.URL{Path: "/"},
			Header:     make(http.Header),
			Host:       "localhost",
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
	}
	out := request.Clone(request.Context())
	out.Header = request.Header.Clone()
	if request.URL != nil {
		u := *request.URL
		out.URL = &u
	}
	if request.Body != nil {
		body, _ := io.ReadAll(request.Body)
		request.Body = io.NopCloser(bytes.NewReader(body))
		out.Body = io.NopCloser(bytes.NewReader(body))
	}
	return out
}

func valueMap(value data.Value) map[string]data.Value {
	result, err := valueToAssocMap(value)
	if err != nil {
		return map[string]data.Value{}
	}
	return result
}

func urlValuesToMap(values url.Values) map[string]data.Value {
	out := make(map[string]data.Value, len(values))
	for key, items := range values {
		if len(items) == 1 {
			out[key] = data.NewStringValue(items[0])
			continue
		}
		vals := make([]data.Value, len(items))
		for i, item := range items {
			vals[i] = data.NewStringValue(item)
		}
		out[key] = data.NewArrayValue(vals)
	}
	return out
}

func stringMapToValues(values map[string]string) map[string]data.Value {
	out := make(map[string]data.Value, len(values))
	for key, value := range values {
		out[key] = data.NewStringValue(value)
	}
	return out
}

func mapToURLValues(values map[string]data.Value) url.Values {
	out := make(url.Values, len(values))
	for key, value := range values {
		if array, ok := value.(*data.ArrayValue); ok {
			for _, item := range array.List {
				if item != nil && item.Value != nil {
					out.Add(key, item.Value.AsString())
				}
			}
		} else if value != nil {
			out.Set(key, value.AsString())
		}
	}
	return out
}

func serverFromRequest(request *http.Request) map[string]data.Value {
	if request == nil {
		return map[string]data.Value{}
	}
	host, port := request.Host, ""
	if parsedHost, parsedPort, err := strings.Cut(request.Host, ":"); err {
		host, port = parsedHost, parsedPort
	}
	if host == "" && request.URL != nil {
		host = request.URL.Hostname()
		port = request.URL.Port()
	}
	if port == "" {
		if request.URL != nil && request.URL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	uri, query, path := "/", "", "/"
	if request.URL != nil {
		uri = request.URL.RequestURI()
		query = request.URL.RawQuery
		path = request.URL.EscapedPath()
		if path == "" {
			path = "/"
		}
	}
	server := map[string]data.Value{
		"REQUEST_METHOD":  data.NewStringValue(strings.ToUpper(request.Method)),
		"REQUEST_URI":     data.NewStringValue(uri),
		"QUERY_STRING":    data.NewStringValue(query),
		"PATH_INFO":       data.NewStringValue(path),
		"SERVER_NAME":     data.NewStringValue(host),
		"HTTP_HOST":       data.NewStringValue(request.Host),
		"SERVER_PORT":     data.NewStringValue(port),
		"SERVER_PROTOCOL": data.NewStringValue(request.Proto),
		"REMOTE_ADDR":     data.NewStringValue(request.RemoteAddr),
	}
	if request.TLS != nil || (request.URL != nil && request.URL.Scheme == "https") {
		server["HTTPS"] = data.NewStringValue("on")
	}
	for name, values := range request.Header {
		key := "HTTP_" + strings.ReplaceAll(strings.ToUpper(name), "-", "_")
		if strings.EqualFold(name, "Content-Type") || strings.EqualFold(name, "Content-Length") {
			key = strings.ReplaceAll(strings.ToUpper(name), "-", "_")
		}
		server[key] = data.NewStringValue(strings.Join(values, ", "))
	}
	return server
}

func headersFromRequest(request *http.Request) map[string][]string {
	out := make(map[string][]string)
	if request == nil {
		return out
	}
	for key, values := range request.Header {
		out[key] = append([]string(nil), values...)
	}
	if request.Host != "" {
		out["Host"] = []string{request.Host}
	}
	return out
}

func cookiesFromRequest(request *http.Request) map[string]data.Value {
	out := make(map[string]data.Value)
	if request != nil {
		for _, cookie := range request.Cookies() {
			out[cookie.Name] = data.NewStringValue(cookie.Value)
		}
	}
	return out
}

func bodyFromRequest(request *http.Request) string {
	if request == nil || request.Body == nil {
		return ""
	}
	body, _ := io.ReadAll(request.Body)
	request.Body = io.NopCloser(bytes.NewReader(body))
	return string(body)
}

func phpValue(value any) data.Value {
	switch typed := value.(type) {
	case nil:
		return data.NewNullValue()
	case string:
		return data.NewStringValue(typed)
	case bool:
		return data.NewBoolValue(typed)
	case float64:
		return data.NewFloatValue(typed)
	case json.Number:
		if integer, err := strconv.Atoi(typed.String()); err == nil {
			return data.NewIntValue(integer)
		}
		return data.NewStringValue(typed.String())
	case []any:
		items := make([]data.Value, len(typed))
		for i, item := range typed {
			items[i] = phpValue(item)
		}
		return data.NewArrayValue(items)
	case map[string]any:
		out := &data.ArrayValue{}
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			out.List = append(out.List, data.NewNamedZVal(key, phpValue(typed[key])))
		}
		return out
	default:
		return data.NewAnyValue(typed)
	}
}

func decodeJSON(content string) (data.Value, error) {
	var decoded any
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return nil, err
	}
	return phpValue(decoded), nil
}

func parseAccept(header string) []string {
	type choice struct {
		value string
		q     float64
		order int
	}
	var choices []choice
	for order, part := range strings.Split(header, ",") {
		segments := strings.Split(strings.TrimSpace(part), ";")
		if segments[0] == "" {
			continue
		}
		q := 1.0
		for _, parameter := range segments[1:] {
			if key, value, ok := strings.Cut(strings.TrimSpace(parameter), "="); ok && key == "q" {
				if parsed, err := strconv.ParseFloat(value, 64); err == nil {
					q = parsed
				}
			}
		}
		choices = append(choices, choice{segments[0], q, order})
	}
	sort.SliceStable(choices, func(i, j int) bool {
		return choices[i].q > choices[j].q
	})
	out := make([]string, len(choices))
	for i, choice := range choices {
		out[i] = choice.value
	}
	return out
}

func valuesArray(values []string) data.Value {
	out := make([]data.Value, len(values))
	for i, value := range values {
		out[i] = data.NewStringValue(value)
	}
	return data.NewArrayValue(out)
}

func contentType(value string) string {
	contentType, _, err := mime.ParseMediaType(value)
	if err != nil {
		return strings.TrimSpace(strings.Split(value, ";")[0])
	}
	return contentType
}

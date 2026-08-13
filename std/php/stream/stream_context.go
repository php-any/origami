package stream

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/core"
)

var nextStreamContextID int64 = 100000

func allocStreamContextID() int {
	return int(atomic.AddInt64(&nextStreamContextID, 1))
}

// StreamContext 保存 stream_context_create 的 options / params。
type StreamContext struct {
	Options map[string]map[string]string
	Params  map[string]string
}

// NewStreamContext 创建流上下文。
func NewStreamContext(options map[string]map[string]string, params map[string]string) *StreamContext {
	if options == nil {
		options = make(map[string]map[string]string)
	}
	if params == nil {
		params = make(map[string]string)
	}
	return &StreamContext{
		Options: options,
		Params:  params,
	}
}

// WrapperOptions 返回指定 wrapper 的选项副本。
func (sc *StreamContext) WrapperOptions(wrapper string) map[string]string {
	if sc == nil {
		return nil
	}
	opts, ok := sc.Options[wrapper]
	if !ok {
		return nil
	}
	copyOpts := make(map[string]string, len(opts))
	for k, v := range opts {
		copyOpts[k] = v
	}
	return copyOpts
}

// ContextFromResource 从 PHP 资源值提取 StreamContext。
func ContextFromResource(val data.Value) *StreamContext {
	if val == nil {
		return nil
	}
	if rv, ok := val.(*core.ResourceValue); ok {
		if sc, ok := rv.GetResource().(*StreamContext); ok {
			return sc
		}
	}
	return nil
}

// ParseStreamContextOptions 将 PHP 数组解析为 stream context options。
func ParseStreamContextOptions(v data.Value) map[string]map[string]string {
	result := make(map[string]map[string]string)
	if v == nil {
		return result
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		for i, zval := range arr.List {
			if zval == nil {
				continue
			}
			wrapperName := entryKeyName(zval, i)
			if opts := parseFlatOptions(zval.Value); opts != nil {
				result[wrapperName] = opts
			}
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(wrapperName string, val data.Value) bool {
			if opts := parseFlatOptions(val); opts != nil {
				result[wrapperName] = opts
			}
			return true
		})
	}
	return result
}

// ParseStreamContextParams 将 PHP 数组解析为 stream context params。
func ParseStreamContextParams(v data.Value) map[string]string {
	result := make(map[string]string)
	if v == nil {
		return result
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		for i, zval := range arr.List {
			if zval == nil {
				continue
			}
			result[entryKeyName(zval, i)] = optionValueToString(zval.Value)
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(key string, val data.Value) bool {
			result[key] = optionValueToString(val)
			return true
		})
	}
	return result
}

func parseFlatOptions(v data.Value) map[string]string {
	result := make(map[string]string)
	switch arr := v.(type) {
	case *data.ArrayValue:
		for i, zval := range arr.List {
			if zval == nil {
				continue
			}
			result[entryKeyName(zval, i)] = optionValueToString(zval.Value)
		}
	case *data.ObjectValue:
		arr.RangeProperties(func(key string, val data.Value) bool {
			result[key] = optionValueToString(val)
			return true
		})
	default:
		return nil
	}
	return result
}

func entryKeyName(zval *data.ZVal, index int) string {
	if zval.Name != "" {
		return zval.Name
	}
	return strconv.Itoa(index)
}

func optionValueToString(v data.Value) string {
	if v == nil {
		return ""
	}
	if _, ok := v.(*data.NullValue); ok {
		return ""
	}
	if b, ok := v.(*data.BoolValue); ok {
		if b.Value {
			return "1"
		}
		return "0"
	}
	return v.AsString()
}

func parseBoolOption(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "on", "yes":
		return true
	default:
		return false
	}
}

func parseHeaderBlock(block string) http.Header {
	headers := make(http.Header)
	block = strings.ReplaceAll(block, "\r\n", "\n")
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		name := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		headers.Add(name, value)
	}
	return headers
}

// HTTPGetContents 通过 http wrapper 读取远程 URL 内容。
func HTTPGetContents(url string, sc *StreamContext) (string, bool) {
	opts := sc.WrapperOptions("http")

	method := "GET"
	if m := opts["method"]; m != "" {
		method = strings.ToUpper(m)
	}

	var body io.Reader
	if content := opts["content"]; content != "" {
		body = strings.NewReader(content)
		if method == "GET" {
			method = "POST"
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", false
	}
	if headerBlock := opts["header"]; headerBlock != "" {
		req.Header = parseHeaderBlock(headerBlock)
	}

	timeout := 60 * time.Second
	if raw := opts["timeout"]; raw != "" {
		if seconds, err := strconv.ParseFloat(raw, 64); err == nil && seconds > 0 {
			timeout = time.Duration(seconds * float64(time.Second))
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}

	ignoreErrors := parseBoolOption(opts["ignore_errors"])
	if !ignoreErrors && resp.StatusCode >= 400 {
		return "", false
	}
	return string(respBody), true
}

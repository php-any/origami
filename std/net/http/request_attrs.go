package http

import (
	"regexp"
	"sync"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

var requestAttrBags sync.Map

func attachRequestAttrs(r *httpsrc.Request) {
	if r == nil {
		return
	}
	requestAttrBags.LoadOrStore(r, make(map[string]data.Value))
}

var requestFormatterSlots sync.Map

func attachRequestFormatter(r *httpsrc.Request, slot *formatHandlerSlot) {
	if r == nil || slot == nil {
		return
	}
	requestFormatterSlots.Store(r, slot)
}

func requestFormatterFor(r *httpsrc.Request) *formatHandlerSlot {
	if r == nil {
		return nil
	}
	if v, ok := requestFormatterSlots.Load(r); ok {
		return v.(*formatHandlerSlot)
	}
	return nil
}

func detachRequestAttrs(r *httpsrc.Request) {
	if r != nil {
		requestAttrBags.Delete(r)
		requestFormatterSlots.Delete(r)
		pathValueKeysStore.Delete(r)
	}
}

func requestAttrs(r *httpsrc.Request) map[string]data.Value {
	if r == nil {
		return nil
	}
	if v, ok := requestAttrBags.Load(r); ok {
		return v.(map[string]data.Value)
	}
	bag := make(map[string]data.Value)
	requestAttrBags.Store(r, bag)
	return bag
}

// pathValueKeysStore 存储每个请求的路由参数键名列表
var pathValueKeysStore sync.Map

// pathParamPattern 匹配路由路径中的 {paramName} 模式
var pathParamPattern = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

// extractPathParamKeys 从路由路径模式中提取参数名列表
// 例如 "/users/{id}/posts/{postId}" → ["id", "postId"]
func extractPathParamKeys(pathPattern string) []string {
	matches := pathParamPattern.FindAllStringSubmatch(pathPattern, -1)
	keys := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			keys = append(keys, m[1])
		}
	}
	return keys
}

// setPathValueKeys 为请求设置路由参数键名列表
func setPathValueKeys(r *httpsrc.Request, keys []string) {
	if r == nil || len(keys) == 0 {
		return
	}
	pathValueKeysStore.Store(r, keys)
}

// getPathValueKeys 获取请求的路由参数键名列表
func getPathValueKeys(r *httpsrc.Request) []string {
	if r == nil {
		return nil
	}
	if v, ok := pathValueKeysStore.Load(r); ok {
		return v.([]string)
	}
	return nil
}

// collectPathValues 收集请求的所有路由参数值
// 返回 map[key]value，如果请求为 nil 则返回 nil
func collectPathValues(r *httpsrc.Request) map[string]string {
	if r == nil {
		return nil
	}
	keys := getPathValueKeys(r)
	if len(keys) == 0 {
		return nil
	}
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		val := r.PathValue(key)
		if val != "" {
			result[key] = val
		}
	}
	return result
}

func beginRequest(r *httpsrc.Request) (*httpsrc.Request, data.ClassStmt) {
	attachRequestAttrs(r)
	return r, NewRequestClassFrom(r)
}

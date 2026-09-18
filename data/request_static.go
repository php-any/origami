package data

import (
	"sync"
	"sync/atomic"
)

// RequestGoid 由 runtime 注入：请求级静态属性 overlay 按 goroutine 隔离。
var RequestGoid func() uint64

var (
	requestStaticActive atomic.Int32
	requestStatics      sync.Map // uint64 goid -> *requestStaticMap
)

type requestStaticMap struct {
	m atomic.Value // map[string]Value
}

func newRequestStaticMap() *requestStaticMap {
	s := &requestStaticMap{}
	s.m.Store(map[string]Value{})
	return s
}

func BeginRequestStaticOverlay() {
	if RequestGoid == nil {
		requestStaticActive.Add(1)
		return
	}
	requestStatics.Store(RequestGoid(), newRequestStaticMap())
	requestStaticActive.Add(1)
}

func EndRequestStaticOverlay() {
	if RequestGoid != nil {
		requestStatics.Delete(RequestGoid())
	}
	requestStaticActive.Add(-1)
}

func requestStaticKey(class, name string) string {
	return class + "\x00" + name
}

// LoadRequestStatic 读当前 goroutine 对某类静态属性的请求覆盖。
func LoadRequestStatic(class, name string) (Value, bool) {
	if requestStaticActive.Load() == 0 || RequestGoid == nil {
		return nil, false
	}
	raw, ok := requestStatics.Load(RequestGoid())
	if !ok {
		return nil, false
	}
	mp, _ := raw.(*requestStaticMap).m.Load().(map[string]Value)
	v, ok := mp[requestStaticKey(class, name)]
	return v, ok
}

// StoreRequestStatic 把静态属性写到当前请求 overlay，避免改进程级 StaticProperty。
// overlay 未启用时返回 false，调用方应写入类上的全局槽。
func StoreRequestStatic(class, name string, v Value) bool {
	if requestStaticActive.Load() == 0 || RequestGoid == nil {
		return false
	}
	raw, ok := requestStatics.Load(RequestGoid())
	if !ok {
		return false
	}
	slot := raw.(*requestStaticMap)
	key := requestStaticKey(class, name)
	old, _ := slot.m.Load().(map[string]Value)
	next := make(map[string]Value, len(old)+1)
	for k, val := range old {
		next[k] = val
	}
	next[key] = v
	slot.m.Store(next)
	return true
}

// CowRequestStatic 请求 overlay 下第一次读到可变静态数组时拷贝进 overlay。
// 否则 static::$arr[] = 会改进程级数组，而 static::$arr = [] 只写 overlay，
// 下一请求又看到泄漏的全局栈（Livewire ExtendBlade::$livewireComponents）。
func CowRequestStatic(class, name string, v Value) Value {
	if v == nil || requestStaticActive.Load() == 0 || RequestGoid == nil {
		return v
	}
	if _, ok := LoadRequestStatic(class, name); ok {
		return v
	}
	var cloned Value
	switch val := v.(type) {
	case *ArrayValue:
		cloned = CloneArrayValue(val)
	case *ObjectValue:
		cloned = CloneObjectValue(val)
	default:
		return v
	}
	if StoreRequestStatic(class, name, cloned) {
		return cloned
	}
	return v
}

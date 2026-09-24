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
	// m 用 sync.Map 而不是「整表拷贝 + atomic.Value 换指针」：
	// 每次静态属性写入都 copy 整张表是 O(现有条目数)，Telescope/容器这类
	// 每请求写几十次静态属性的路径上会放大成大量临时 map，也是 GC 压力的主要来源之一。
	// sync.Map 的读路径无锁、写路径不拷贝，条目随请求结束整体丢弃，不会无限增长。
	m sync.Map // string -> Value
}

func newRequestStaticMap() *requestStaticMap {
	return &requestStaticMap{}
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
	mp := &raw.(*requestStaticMap).m
	if v, ok := mp.Load(requestStaticKey(class, name)); ok {
		val, _ := v.(Value)
		return val, true
	}
	return nil, false
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
	slot.m.Store(requestStaticKey(class, name), v)
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

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
	m     sync.Map // requestStaticProperty -> Value
	scope *RequestObjectScope
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

type requestStaticProperty struct{ class, name string }

func requestStaticKey(class, name string) requestStaticProperty {
	return requestStaticProperty{class, name}
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
		if ref, ok := val.(*ZValValue); ok {
			return ref.ZVal.ReadValue(), true
		}
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
	if previous, ok := slot.m.Load(requestStaticKey(class, name)); ok {
		if ref, ok := previous.(*ZValValue); ok {
			CowAssign(ref.ZVal, v)
			return true
		}
	}
	slot.m.Store(requestStaticKey(class, name), v)
	return true
}

// StaticReferenceSlot promotes the property once, keeping aliases private to
// the request. Ordinary static values remain inline in the existing map.
func StaticReferenceSlot(class, name string, initial Value) (*ZVal, bool) {
	if requestStaticActive.Load() == 0 || RequestGoid == nil {
		return nil, false
	}
	raw, ok := requestStatics.Load(RequestGoid())
	if !ok {
		return nil, false
	}
	state := raw.(*requestStaticMap)
	key := requestStaticKey(class, name)
	if previous, ok := state.m.Load(key); ok {
		if ref, ok := previous.(*ZValValue); ok {
			return ref.ZVal, true
		}
		initial = previous.(Value)
	}
	var slot *ZVal
	if ref, ok := initial.(*ZValValue); ok && state.scope != nil {
		slot = state.scope.BindSlot(ref.ZVal)
	} else {
		slot = NewZVal(initial)
	}
	state.m.Store(key, NewZValValue(slot))
	return slot, true
}

// Reference operations replace the slot itself; ordinary assignment writes its value.
func BindRequestStaticReference(class, name string, slot *ZVal) bool {
	if requestStaticActive.Load() == 0 || RequestGoid == nil {
		return false
	}
	raw, ok := requestStatics.Load(RequestGoid())
	if !ok {
		return false
	}
	raw.(*requestStaticMap).m.Store(requestStaticKey(class, name), NewZValValue(slot))
	return true
}

func LoadRequestStaticReference(class, name string) *ZVal {
	if requestStaticActive.Load() == 0 || RequestGoid == nil {
		return nil
	}
	raw, ok := requestStatics.Load(RequestGoid())
	if !ok {
		return nil
	}
	value, _ := raw.(*requestStaticMap).m.Load(requestStaticKey(class, name))
	if ref, ok := value.(*ZValValue); ok {
		return ref.ZVal
	}
	return nil
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
	if raw, ok := requestStatics.Load(RequestGoid()); ok {
		if scope := raw.(*requestStaticMap).scope; scope != nil {
			cloned := scope.Bind(v)
			if cloned != v && StoreRequestStatic(class, name, cloned) {
				return cloned
			}
			return v
		}
	}
	var cloned Value
	switch val := v.(type) {
	case *ArrayValue:
		cloned = CloneArrayValue(val)
	default:
		return v
	}
	if StoreRequestStatic(class, name, cloned) {
		return cloned
	}
	return v
}

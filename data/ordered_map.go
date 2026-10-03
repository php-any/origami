package data

import (
	"sync"
)

type PropertyStore interface {
	Set(key string, value Value)
	Get(key string) (Value, bool)
	GetZVal(key string) (*ZVal, bool)
	Delete(key string)
	Range(fn func(key string, value Value) bool)
	Len() int
	GetByIndex(index int) (string, Value, bool)
}

// OrderedMap 是一个有序的键值对存储结构，保持插入顺序。
//
// 键与 data 平行存在 keys []string 里，不再用 map[int]string 存「下标 -> 键」：
// 那张 map 和 data 是一一对应的，热路径上每次 Set 都要多写一张哈希表、
// 每次 Range/GetByIndex 都要多查一次哈希，而按下标取键本该是 O(1) 的切片下标。
// 改成平行切片后，Set 少一次哈希写、Range/GetByIndex 各少一次哈希查，
// Delete 重建的是切片而不是又一张 map。
//
// indexMap（键 -> 下标）保留：它是按键查找的唯一路径，不能省。
// 不变量：len(keys) == len(data)，两者下标一一对应。
type OrderedMap struct {
	mu       sync.RWMutex
	data     []*ZVal
	keys     []string // 与 data 同下标：第 i 个槽的键（PHP 空字符串键就是 ""）
	indexMap map[string]int
}

// NewOrderedMap 创建新的有序映射
func NewOrderedMap() PropertyStore {
	return &OrderedMap{}
}

func newOrderedMapWithCapacity(capacity int) PropertyStore {
	if capacity == 0 {
		return NewOrderedMap()
	}
	return &OrderedMap{
		data:     make([]*ZVal, 0, capacity),
		keys:     make([]string, 0, capacity),
		indexMap: make(map[string]int, capacity),
	}
}

func (om *OrderedMap) GetZVal(key string) (*ZVal, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if idx, exists := om.indexMap[key]; exists {
		if idx >= 0 && idx < len(om.data) {
			return om.data[idx], true
		}
	}
	return nil, false
}

// Set 设置键值对，保持插入顺序
func (om *OrderedMap) Set(key string, value Value) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if idx, exists := om.indexMap[key]; exists {
		// 更新已存在的键值对
		if idx >= 0 && idx < len(om.data) && om.data[idx] != nil {
			CowAssign(om.data[idx], value)
		}
		return
	}
	// 添加新的键值对
	if om.indexMap == nil {
		om.indexMap = make(map[string]int)
	}
	index := len(om.data)
	om.data = append(om.data, NewZVal(CowAddRef(value)))
	om.keys = append(om.keys, key)
	om.indexMap[key] = index
}

func (om *OrderedMap) BindZVal(key string, slot *ZVal) {
	om.mu.Lock()
	defer om.mu.Unlock()
	if om.indexMap == nil {
		om.indexMap = make(map[string]int)
	}
	if index, found := om.indexMap[key]; found {
		om.data[index] = slot
		return
	}
	om.indexMap[key] = len(om.data)
	om.data = append(om.data, slot)
	om.keys = append(om.keys, key)
}

// Delete 删除键（不存在则无操作）
func (om *OrderedMap) Delete(key string) {
	om.mu.Lock()
	defer om.mu.Unlock()
	if _, exists := om.indexMap[key]; !exists {
		return
	}
	newData := make([]*ZVal, 0, len(om.data))
	newKeys := make([]string, 0, len(om.data))
	newIndex := make(map[string]int, len(om.indexMap))
	for i, z := range om.data {
		if i >= len(om.keys) || om.keys[i] == key {
			continue
		}
		k := om.keys[i]
		newIndex[k] = len(newData)
		newData = append(newData, z)
		newKeys = append(newKeys, k)
	}
	om.data = newData
	om.keys = newKeys
	om.indexMap = newIndex
}

// Get 获取值
func (om *OrderedMap) Get(key string) (Value, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if idx, exists := om.indexMap[key]; exists {
		if idx >= 0 && idx < len(om.data) {
			return om.data[idx].ReadValue(), true
		}
	}
	return nil, false
}

// Range 遍历所有键值对，按插入顺序。
// 必须先在读锁下快照再回调：Go 的 RWMutex 不可重入，若回调内对同一 map
// 再 Get/Set（例如 foreach 写回当前关联数组），持 RLock 调用会直接死锁。
func (om *OrderedMap) Range(fn func(key string, value Value) bool) {
	type kv struct {
		key string
		val Value
	}

	om.mu.RLock()
	items := make([]kv, 0, len(om.data))
	for i, zval := range om.data {
		if i >= len(om.keys) {
			break
		}
		var val Value
		if zval != nil {
			val = zval.ReadValue()
		}
		items = append(items, kv{key: om.keys[i], val: val})
	}
	om.mu.RUnlock()

	for _, item := range items {
		if !fn(item.key, item.val) {
			break
		}
	}
}

// Len 获取元素数量
func (om *OrderedMap) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return len(om.data)
}

// GetByIndex 按索引获取键值对
func (om *OrderedMap) GetByIndex(index int) (string, Value, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()
	if index >= 0 && index < len(om.data) && index < len(om.keys) {
		return om.keys[index], om.data[index].ReadValue(), true
	}
	return "", nil, false
}

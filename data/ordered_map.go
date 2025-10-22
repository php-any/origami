package data

import (
	"sync"
)

// OrderedMap 是一个有序的键值对存储结构，保持插入顺序
type OrderedMap struct {
	mu       sync.RWMutex
	items    []OrderedMapItem // 按插入顺序存储
	indexMap map[string]int   // 快速查找索引
}

// OrderedMapItem 表示一个有序的键值对项
type OrderedMapItem struct {
	Key   string
	Value Value
	Index int // 插入顺序索引
}

// NewOrderedMap 创建新的有序映射
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		items:    make([]OrderedMapItem, 0),
		indexMap: make(map[string]int),
	}
}

// Set 设置键值对，保持插入顺序
func (om *OrderedMap) Set(key string, value Value) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if idx, exists := om.indexMap[key]; exists {
		// 更新已存在的键值对
		om.items[idx].Value = value
	} else {
		// 添加新的键值对
		index := len(om.items)
		om.items = append(om.items, OrderedMapItem{
			Key:   key,
			Value: value,
			Index: index,
		})
		om.indexMap[key] = index
	}
}

// Get 获取值
func (om *OrderedMap) Get(key string) (Value, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if idx, exists := om.indexMap[key]; exists {
		return om.items[idx].Value, true
	}
	return nil, false
}

// Delete 删除键值对
func (om *OrderedMap) Delete(key string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if idx, exists := om.indexMap[key]; exists {
		// 从切片中删除
		om.items = append(om.items[:idx], om.items[idx+1:]...)
		delete(om.indexMap, key)

		// 更新后续元素的索引
		for i := idx; i < len(om.items); i++ {
			om.items[i].Index = i
			om.indexMap[om.items[i].Key] = i
		}
	}
}

// Has 检查键是否存在
func (om *OrderedMap) Has(key string) bool {
	om.mu.RLock()
	defer om.mu.RUnlock()

	_, exists := om.indexMap[key]
	return exists
}

// Range 遍历所有键值对，按插入顺序
func (om *OrderedMap) Range(fn func(key string, value Value) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for _, item := range om.items {
		if !fn(item.Key, item.Value) {
			break
		}
	}
}

// Load 兼容 sync.Map 的 Load 方法
func (om *OrderedMap) Load(key string) (Value, bool) {
	return om.Get(key)
}

// Store 兼容 sync.Map 的 Store 方法
func (om *OrderedMap) Store(key string, value Value) {
	om.Set(key, value)
}

// GetAll 获取所有键值对，按插入顺序返回
func (om *OrderedMap) GetAll() []OrderedMapItem {
	om.mu.RLock()
	defer om.mu.RUnlock()

	// 返回副本以避免并发修改
	result := make([]OrderedMapItem, len(om.items))
	copy(result, om.items)
	return result
}

// Len 返回键值对数量
func (om *OrderedMap) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return len(om.items)
}

// GetByIndex 根据索引获取键值对
func (om *OrderedMap) GetByIndex(index int) (string, Value, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if index < 0 || index >= len(om.items) {
		return "", nil, false
	}
	item := om.items[index]
	return item.Key, item.Value, true
}

// Clear 清空所有键值对
func (om *OrderedMap) Clear() {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.items = make([]OrderedMapItem, 0)
	om.indexMap = make(map[string]int)
}

// Keys 获取所有键，按插入顺序
func (om *OrderedMap) Keys() []string {
	om.mu.RLock()
	defer om.mu.RUnlock()

	keys := make([]string, len(om.items))
	for i, item := range om.items {
		keys[i] = item.Key
	}
	return keys
}

// Values 获取所有值，按插入顺序
func (om *OrderedMap) Values() []Value {
	om.mu.RLock()
	defer om.mu.RUnlock()

	values := make([]Value, len(om.items))
	for i, item := range om.items {
		values[i] = item.Value
	}
	return values
}

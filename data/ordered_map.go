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

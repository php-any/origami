package data

import (
	"sync"
)

type PropertyStore interface {
	Set(key string, value Value)
	Get(key string) (Value, bool)
	GetZVal(key string) (*ZVal, bool)
	Range(fn func(key string, value Value) bool)
	Len() int
	GetByIndex(index int) (string, Value, bool)
}

// OrderedMap 是一个有序的键值对存储结构，保持插入顺序
type OrderedMap struct {
	mu       sync.RWMutex
	data     []*ZVal
	indexMap map[string]int // 快速查找索引
	nameMap  map[int]string // 快速查找索引
}

// NewOrderedMap 创建新的有序映射
func NewOrderedMap() PropertyStore {
	return &OrderedMap{
		data:     make([]*ZVal, 0),
		indexMap: make(map[string]int),
		nameMap:  make(map[int]string),
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
		if idx >= 0 && idx < len(om.data) {
			om.data[idx].Value = value
		}
	} else {
		// 添加新的键值对
		index := len(om.data)
		om.data = append(om.data, NewZVal(value))
		om.indexMap[key] = index
		om.nameMap[index] = key
	}
}

// Get 获取值
func (om *OrderedMap) Get(key string) (Value, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if idx, exists := om.indexMap[key]; exists {
		if idx >= 0 && idx < len(om.data) {
			return om.data[idx].Value, true
		}
	}
	return nil, false
}

// Range 遍历所有键值对，按插入顺序
func (om *OrderedMap) Range(fn func(key string, value Value) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for i, zval := range om.data {
		if key, exists := om.nameMap[i]; exists {
			if !fn(key, zval.Value) {
				break
			}
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
	if index >= 0 && index < len(om.data) {
		if key, exists := om.nameMap[index]; exists {
			return key, om.data[index].Value, true
		}
	}
	return "", nil, false
}

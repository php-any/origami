package loop

import (
	"github.com/php-any/origami/data"
)

// HashMap 泛型哈希表结构
type HashMap struct {
	items     map[string]data.Value
	iterator  int
	keys      []string
	keyType   data.Types
	valueType data.Types
}

// NewHashMap 创建一个新的 HashMap 实例
func NewHashMap(keyType, valueType data.Types) *HashMap {
	return &HashMap{
		items:     make(map[string]data.Value),
		iterator:  0,
		keys:      make([]string, 0),
		keyType:   keyType,
		valueType: valueType,
	}
}

// Put 添加或更新键值对
func (h *HashMap) Put(key, value data.Value) {
	keyStr := key.AsString()
	if _, exists := h.items[keyStr]; !exists {
		h.keys = append(h.keys, keyStr)
	}
	h.items[keyStr] = value
}

// Get 获取值
func (h *HashMap) Get(key data.Value) (data.Value, bool) {
	keyStr := key.AsString()
	value, exists := h.items[keyStr]
	return value, exists
}

// Remove 移除键值对
func (h *HashMap) Remove(key data.Value) bool {
	keyStr := key.AsString()
	if _, exists := h.items[keyStr]; exists {
		delete(h.items, keyStr)
		// 从 keys 数组中移除
		for i, k := range h.keys {
			if k == keyStr {
				h.keys = append(h.keys[:i], h.keys[i+1:]...)
				break
			}
		}
		return true
	}
	return false
}

// ContainsKey 检查是否包含键
func (h *HashMap) ContainsKey(key data.Value) bool {
	keyStr := key.AsString()
	_, exists := h.items[keyStr]
	return exists
}

// ContainsValue 检查是否包含值
func (h *HashMap) ContainsValue(value data.Value) bool {
	for _, v := range h.items {
		if v.AsString() == value.AsString() {
			return true
		}
	}
	return false
}

// Size 获取大小
func (h *HashMap) Size() int {
	return len(h.items)
}

// IsEmpty 检查是否为空
func (h *HashMap) IsEmpty() bool {
	return len(h.items) == 0
}

// Clear 清空哈希表
func (h *HashMap) Clear() {
	h.items = make(map[string]data.Value)
	h.keys = make([]string, 0)
	h.iterator = 0
}

// Keys 获取所有键
func (h *HashMap) Keys() []string {
	result := make([]string, len(h.keys))
	copy(result, h.keys)
	return result
}

// Values 获取所有值
func (h *HashMap) Values() []data.Value {
	result := make([]data.Value, 0, len(h.items))
	for _, key := range h.keys {
		if value, exists := h.items[key]; exists {
			result = append(result, value)
		}
	}
	return result
}

// Iterator 接口实现

// Current 当前值
func (h *HashMap) Current() data.Value {
	if h.iterator >= 0 && h.iterator < len(h.keys) {
		key := h.keys[h.iterator]
		return h.items[key]
	}
	return data.NewNullValue()
}

// Key 当前键
func (h *HashMap) Key() data.Value {
	if h.iterator >= 0 && h.iterator < len(h.keys) {
		// 返回原始键值，而不是字符串
		key := h.keys[h.iterator]
		// 这里需要根据 keyType 创建对应的值
		return data.NewStringValue(key)
	}
	return data.NewNullValue()
}

// Next 下一个元素
func (h *HashMap) Next() {
	h.iterator++
}

// Rewind 重置迭代器
func (h *HashMap) Rewind() {
	h.iterator = 0
}

// Valid 检查迭代器是否有效
func (h *HashMap) Valid() bool {
	return h.iterator >= 0 && h.iterator < len(h.keys)
}

// HashMapClass 表示 HashMap 类
type HashMapClass struct {
	hashmap             *HashMap
	construct           data.Method
	putMethod           data.Method
	getMethod           data.Method
	removeMethod        data.Method
	containsKeyMethod   data.Method
	containsValueMethod data.Method
	sizeMethod          data.Method
	isEmptyMethod       data.Method
	clearMethod         data.Method
	keysMethod          data.Method
	valuesMethod        data.Method
	currentMethod       data.Method
	keyMethod           data.Method
	nextMethod          data.Method
	rewindMethod        data.Method
	validMethod         data.Method
}

// NewHashMapClass 创建一个新的 HashMap 类实例
func NewHashMapClass() data.ClassStmt {
	return (&HashMapClass{}).Clone(nil).(data.ClassStmt)
}

// Clone 实现 ClassGeneric 接口
func (h *HashMapClass) Clone(m map[string]data.Types) data.ClassGeneric {
	var source *HashMap
	if len(m) >= 2 {
		// 获取键类型和值类型
		var keyType, valueType data.Types
		index := 0
		for _, t := range m {
			if index == 0 {
				keyType = t
			} else if index == 1 {
				valueType = t
			}
			index++
		}
		source = NewHashMap(keyType, valueType)
	} else {
		// 默认类型
		source = NewHashMap(data.NewBaseType("string"), data.NewBaseType("mixed"))
	}

	// 创建共享的 HashMapClass 实例
	sharedClass := &HashMapClass{hashmap: source}

	// 设置所有方法，都使用同一个 sharedClass
	sharedClass.construct = &HashMapConstructMethod{source: sharedClass}
	sharedClass.putMethod = &HashMapPutMethod{source: sharedClass}
	sharedClass.getMethod = &HashMapGetMethod{source: sharedClass}
	sharedClass.removeMethod = &HashMapRemoveMethod{source: sharedClass}
	sharedClass.containsKeyMethod = &HashMapContainsKeyMethod{source: sharedClass}
	sharedClass.containsValueMethod = &HashMapContainsValueMethod{source: sharedClass}
	sharedClass.sizeMethod = &HashMapSizeMethod{source: sharedClass}
	sharedClass.isEmptyMethod = &HashMapIsEmptyMethod{source: sharedClass}
	sharedClass.clearMethod = &HashMapClearMethod{source: sharedClass}
	sharedClass.keysMethod = &HashMapKeysMethod{source: sharedClass}
	sharedClass.valuesMethod = &HashMapValuesMethod{source: sharedClass}
	sharedClass.currentMethod = &HashMapCurrentMethod{source: sharedClass}
	sharedClass.keyMethod = &HashMapKeyMethod{source: sharedClass}
	sharedClass.nextMethod = &HashMapNextMethod{source: sharedClass}
	sharedClass.rewindMethod = &HashMapRewindMethod{source: sharedClass}
	sharedClass.validMethod = &HashMapValidMethod{source: sharedClass}

	return sharedClass
}

// GenericList 返回泛型参数列表
func (h *HashMapClass) GenericList() []data.Types {
	return []data.Types{
		data.Generic{Name: "K"},
		data.Generic{Name: "V"},
	}
}

// GetValue 泛型会提前有一次 Clone
func (h *HashMapClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(h, ctx), nil
}

// GetFrom 返回来源信息
func (h *HashMapClass) GetFrom() data.From {
	return nil
}

// GetName 返回类名
func (h *HashMapClass) GetName() string {
	return "HashMap"
}

// GetExtend 返回父类
func (h *HashMapClass) GetExtend() *string {
	return nil
}

// GetImplements 返回实现的接口
func (h *HashMapClass) GetImplements() []string {
	return []string{"Iterator"}
}

// GetProperty 获取属性
func (h *HashMapClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取所有属性列表
func (h *HashMapClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

// GetMethod 获取方法
func (h *HashMapClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "put":
		return h.putMethod, true
	case "get":
		return h.getMethod, true
	case "remove":
		return h.removeMethod, true
	case "containsKey":
		return h.containsKeyMethod, true
	case "containsValue":
		return h.containsValueMethod, true
	case "size":
		return h.sizeMethod, true
	case "isEmpty":
		return h.isEmptyMethod, true
	case "clear":
		return h.clearMethod, true
	case "keys":
		return h.keysMethod, true
	case "values":
		return h.valuesMethod, true
	case "current":
		return h.currentMethod, true
	case "key":
		return h.keyMethod, true
	case "next":
		return h.nextMethod, true
	case "rewind":
		return h.rewindMethod, true
	case "valid":
		return h.validMethod, true
	}
	return nil, false
}

// GetMethods 获取所有方法
func (h *HashMapClass) GetMethods() []data.Method {
	return []data.Method{
		h.construct,
		h.putMethod,
		h.getMethod,
		h.removeMethod,
		h.containsKeyMethod,
		h.containsValueMethod,
		h.sizeMethod,
		h.isEmptyMethod,
		h.clearMethod,
		h.keysMethod,
		h.valuesMethod,
		h.currentMethod,
		h.keyMethod,
		h.nextMethod,
		h.rewindMethod,
		h.validMethod,
	}
}

// GetConstruct 获取构造函数
func (h *HashMapClass) GetConstruct() data.Method {
	return h.construct
}

// AddAnnotations 添加注解
func (h *HashMapClass) AddAnnotations(a *data.ClassValue) {
	// 暂时不处理注解
}

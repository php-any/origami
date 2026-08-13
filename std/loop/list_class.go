package loop

import (
	"github.com/php-any/origami/data"
)

// List 泛型列表结构
type List struct {
	items    []data.Value
	iterator int
	itemType data.Types
}

// NewList 创建一个新的 List 实例
func NewList(itemType data.Types) *List {
	return &List{
		items:    make([]data.Value, 0),
		iterator: 0,
		itemType: itemType,
	}
}

// Add 添加元素
func (l *List) Add(item data.Value) {
	l.items = append(l.items, item)
}

// Get 获取元素
func (l *List) Get(index int) (data.Value, bool) {
	if index < 0 || index >= len(l.items) {
		return nil, false
	}
	return l.items[index], true
}

// Set 设置元素
func (l *List) Set(index int, item data.Value) bool {
	if index < 0 || index >= len(l.items) {
		return false
	}
	l.items[index] = item
	return true
}

// Size 获取大小
func (l *List) Size() int {
	return len(l.items)
}

// Remove 移除元素
func (l *List) Remove(item data.Value) bool {
	for i, v := range l.items {
		if v.AsString() == item.AsString() {
			l.items = append(l.items[:i], l.items[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveAt 按索引移除元素
func (l *List) RemoveAt(index int) (data.Value, bool) {
	if index < 0 || index >= len(l.items) {
		return nil, false
	}
	item := l.items[index]
	l.items = append(l.items[:index], l.items[index+1:]...)
	return item, true
}

// Clear 清空列表
func (l *List) Clear() {
	l.items = make([]data.Value, 0)
	l.iterator = 0
}

// Contains 检查是否包含元素
func (l *List) Contains(item data.Value) bool {
	for _, v := range l.items {
		if v.AsString() == item.AsString() {
			return true
		}
	}
	return false
}

// IndexOf 获取元素索引
func (l *List) IndexOf(item data.Value) int {
	for i, v := range l.items {
		if v.AsString() == item.AsString() {
			return i
		}
	}
	return -1
}

// IsEmpty 检查是否为空
func (l *List) IsEmpty() bool {
	return len(l.items) == 0
}

// ToArray 转换为数组
func (l *List) ToArray() []data.Value {
	result := make([]data.Value, len(l.items))
	copy(result, l.items)
	return result
}

// Iterator 接口实现

// Current 当前元素
func (l *List) Current() data.Value {
	if l.iterator >= 0 && l.iterator < len(l.items) {
		return l.items[l.iterator]
	}
	return data.NewNullValue()
}

// Key 当前键
func (l *List) Key() data.Value {
	if l.iterator >= 0 && l.iterator < len(l.items) {
		return data.NewIntValue(l.iterator)
	}
	return data.NewNullValue()
}

// Next 下一个元素
func (l *List) Next() {
	l.iterator++
}

// Rewind 重置迭代器
func (l *List) Rewind() {
	l.iterator = 0
}

// Valid 检查迭代器是否有效
func (l *List) Valid() bool {
	return l.iterator >= 0 && l.iterator < len(l.items)
}

// ListClass 表示 List 类
type ListClass struct {
	list           *List
	construct      data.Method
	addMethod      data.Method
	getMethod      data.Method
	setMethod      data.Method
	sizeMethod     data.Method
	removeMethod   data.Method
	removeAtMethod data.Method
	clearMethod    data.Method
	containsMethod data.Method
	indexOfMethod  data.Method
	isEmptyMethod  data.Method
	toArrayMethod  data.Method
	currentMethod  data.Method
	keyMethod      data.Method
	nextMethod     data.Method
	rewindMethod   data.Method
	validMethod    data.Method
}

// NewListClass 创建一个新的 List 类实例
func NewListClass() data.ClassStmt {
	return (&ListClass{}).Clone(nil).(data.ClassStmt)
}

// Clone 实现 ClassGeneric 接口
func (l *ListClass) Clone(m map[string]data.Types) data.ClassGeneric {
	var source *List
	if m != nil {
		// 从泛型参数中获取类型
		if len(m) > 0 {
			// 获取第一个泛型参数作为元素类型
			for _, itemType := range m {
				source = NewList(itemType)
				break
			}
		} else {
			source = NewList(nil) // 混合类型
		}
	} else {
		source = NewList(nil) // 混合类型
	}

	// 创建共享的 ListClass 实例
	sharedClass := &ListClass{list: source}

	// 设置所有方法，都使用同一个 sharedClass
	sharedClass.construct = &ListConstructMethod{source: sharedClass}
	sharedClass.addMethod = &ListAddMethod{source: sharedClass}
	sharedClass.getMethod = &ListGetMethod{source: sharedClass}
	sharedClass.setMethod = &ListSetMethod{source: sharedClass}
	sharedClass.sizeMethod = &ListSizeMethod{source: sharedClass}
	sharedClass.removeMethod = &ListRemoveMethod{source: sharedClass}
	sharedClass.removeAtMethod = &ListRemoveAtMethod{source: sharedClass}
	sharedClass.clearMethod = &ListClearMethod{source: sharedClass}
	sharedClass.containsMethod = &ListContainsMethod{source: sharedClass}
	sharedClass.indexOfMethod = &ListIndexOfMethod{source: sharedClass}
	sharedClass.isEmptyMethod = &ListIsEmptyMethod{source: sharedClass}
	sharedClass.toArrayMethod = &ListToArrayMethod{source: sharedClass}
	sharedClass.currentMethod = &ListCurrentMethod{source: sharedClass}
	sharedClass.keyMethod = &ListKeyMethod{source: sharedClass}
	sharedClass.nextMethod = &ListNextMethod{source: sharedClass}
	sharedClass.rewindMethod = &ListRewindMethod{source: sharedClass}
	sharedClass.validMethod = &ListValidMethod{source: sharedClass}

	return sharedClass
}

// GenericList 返回泛型参数列表
func (l *ListClass) GenericList() []data.Types {
	return []data.Types{
		data.Generic{Name: "T"},
	}
}

// GetValue 泛型会提前有一次 Clone
func (l *ListClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(l, ctx), nil
}

// GetFrom 返回来源信息
func (l *ListClass) GetFrom() data.From {
	return nil
}

// GetName 返回类名
func (l *ListClass) GetName() string {
	return "List"
}

// GetExtend 返回父类
func (l *ListClass) GetExtend() *string {
	return nil
}

// GetImplements 返回实现的接口
func (l *ListClass) GetImplements() []string {
	return []string{"Iterator"}
}

// GetProperty 获取属性
func (l *ListClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取所有属性列表
func (l *ListClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

// GetMethod 获取方法
func (l *ListClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "add":
		return l.addMethod, true
	case "get":
		return l.getMethod, true
	case "set":
		return l.setMethod, true
	case "size":
		return l.sizeMethod, true
	case "remove":
		return l.removeMethod, true
	case "removeAt":
		return l.removeAtMethod, true
	case "clear":
		return l.clearMethod, true
	case "contains":
		return l.containsMethod, true
	case "indexOf":
		return l.indexOfMethod, true
	case "isEmpty":
		return l.isEmptyMethod, true
	case "toArray":
		return l.toArrayMethod, true
	case "current":
		return l.currentMethod, true
	case "key":
		return l.keyMethod, true
	case "next":
		return l.nextMethod, true
	case "rewind":
		return l.rewindMethod, true
	case "valid":
		return l.validMethod, true
	}
	return nil, false
}

// GetMethods 获取所有方法
func (l *ListClass) GetMethods() []data.Method {
	return []data.Method{
		l.addMethod,
		l.getMethod,
		l.setMethod,
		l.sizeMethod,
		l.removeMethod,
		l.removeAtMethod,
		l.clearMethod,
		l.containsMethod,
		l.indexOfMethod,
		l.isEmptyMethod,
		l.toArrayMethod,
		l.currentMethod,
		l.keyMethod,
		l.nextMethod,
		l.rewindMethod,
		l.validMethod,
	}
}

// GetConstruct 获取构造函数
func (l *ListClass) GetConstruct() data.Method {
	return l.construct
}

// AddAnnotations 添加注解
func (l *ListClass) AddAnnotations(a *data.ClassValue) {
	// 暂时不处理注解
}

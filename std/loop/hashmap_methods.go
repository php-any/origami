package loop

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// HashMapConstructMethod 构造函数
type HashMapConstructMethod struct {
	source *HashMapClass
}

func (m *HashMapConstructMethod) GetName() string {
	return "__construct"
}

func (m *HashMapConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapConstructMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *HashMapConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 构造函数不需要特殊处理，HashMap 已经初始化
	return data.NewNullValue(), nil
}

// HashMapPutMethod 添加键值对方法
type HashMapPutMethod struct {
	source *HashMapClass
}

func (m *HashMapPutMethod) GetName() string {
	return "put"
}

func (m *HashMapPutMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapPutMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapPutMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}

func (m *HashMapPutMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
		node.NewVariable(nil, "value", 1, nil),
	}
}

func (m *HashMapPutMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *HashMapPutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	key, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("HashMap::put() 需要键参数"))
	}

	value, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("HashMap::put() 需要值参数"))
	}

	// 获取 HashMap 实例
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	// 泛型类型检查
	if m.source.hashmap.keyType != nil && !m.source.hashmap.keyType.Is(key) {
		return nil, utils.NewThrow(errors.New("类型不匹配: 键类型期望 " + m.source.hashmap.keyType.String() + " 但得到 " + key.AsString()))
	}

	if m.source.hashmap.valueType != nil && !m.source.hashmap.valueType.Is(value) {
		return nil, utils.NewThrow(errors.New("类型不匹配: 值类型期望 " + m.source.hashmap.valueType.String() + " 但得到 " + value.AsString()))
	}

	// 添加键值对
	m.source.hashmap.Put(key, value)
	return data.NewNullValue(), nil
}

// HashMapGetMethod 获取值方法
type HashMapGetMethod struct {
	source *HashMapClass
}

func (m *HashMapGetMethod) GetName() string {
	return "get"
}

func (m *HashMapGetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapGetMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}

func (m *HashMapGetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}

func (m *HashMapGetMethod) GetReturnType() data.Types {
	return nil
}

func (m *HashMapGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	key, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("HashMap::get() 需要键参数"))
	}

	// 获取 HashMap 实例
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	// 泛型类型检查
	if m.source.hashmap.keyType != nil && !m.source.hashmap.keyType.Is(key) {
		return nil, utils.NewThrow(errors.New("类型不匹配: 键类型期望 " + m.source.hashmap.keyType.String() + " 但得到 " + key.AsString()))
	}

	// 获取值
	value, found := m.source.hashmap.Get(key)
	if !found {
		return data.NewNullValue(), nil
	}

	return value, nil
}

// HashMapRemoveMethod 移除键值对方法
type HashMapRemoveMethod struct {
	source *HashMapClass
}

func (m *HashMapRemoveMethod) GetName() string {
	return "remove"
}

func (m *HashMapRemoveMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapRemoveMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapRemoveMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}

func (m *HashMapRemoveMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}

func (m *HashMapRemoveMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *HashMapRemoveMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	key, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("HashMap::remove() 需要键参数"))
	}

	// 获取 HashMap 实例
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	// 泛型类型检查
	if m.source.hashmap.keyType != nil && !m.source.hashmap.keyType.Is(key) {
		return nil, utils.NewThrow(errors.New("类型不匹配: 键类型期望 " + m.source.hashmap.keyType.String() + " 但得到 " + key.AsString()))
	}

	// 移除键值对
	success := m.source.hashmap.Remove(key)
	return data.NewBoolValue(success), nil
}

// HashMapContainsKeyMethod 检查是否包含键方法
type HashMapContainsKeyMethod struct {
	source *HashMapClass
}

func (m *HashMapContainsKeyMethod) GetName() string {
	return "containsKey"
}

func (m *HashMapContainsKeyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapContainsKeyMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapContainsKeyMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}

func (m *HashMapContainsKeyMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}

func (m *HashMapContainsKeyMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *HashMapContainsKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	key, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("HashMap::containsKey() 需要键参数"))
	}

	// 获取 HashMap 实例
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	// 泛型类型检查
	if m.source.hashmap.keyType != nil && !m.source.hashmap.keyType.Is(key) {
		return nil, utils.NewThrow(errors.New("类型不匹配: 键类型期望 " + m.source.hashmap.keyType.String() + " 但得到 " + key.AsString()))
	}

	// 检查是否包含键
	contains := m.source.hashmap.ContainsKey(key)
	return data.NewBoolValue(contains), nil
}

// HashMapContainsValueMethod 检查是否包含值方法
type HashMapContainsValueMethod struct {
	source *HashMapClass
}

func (m *HashMapContainsValueMethod) GetName() string {
	return "containsValue"
}

func (m *HashMapContainsValueMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapContainsValueMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapContainsValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (m *HashMapContainsValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

func (m *HashMapContainsValueMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *HashMapContainsValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("HashMap::containsValue() 需要值参数"))
	}

	// 获取 HashMap 实例
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	// 泛型类型检查
	if m.source.hashmap.valueType != nil && !m.source.hashmap.valueType.Is(value) {
		return nil, utils.NewThrow(errors.New("类型不匹配: 值类型期望 " + m.source.hashmap.valueType.String() + " 但得到 " + value.AsString()))
	}

	// 检查是否包含值
	contains := m.source.hashmap.ContainsValue(value)
	return data.NewBoolValue(contains), nil
}

// HashMapSizeMethod 获取大小方法
type HashMapSizeMethod struct {
	source *HashMapClass
}

func (m *HashMapSizeMethod) GetName() string {
	return "size"
}

func (m *HashMapSizeMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapSizeMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapSizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapSizeMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapSizeMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

func (m *HashMapSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	size := m.source.hashmap.Size()
	return data.NewIntValue(size), nil
}

// HashMapIsEmptyMethod 检查是否为空方法
type HashMapIsEmptyMethod struct {
	source *HashMapClass
}

func (m *HashMapIsEmptyMethod) GetName() string {
	return "isEmpty"
}

func (m *HashMapIsEmptyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapIsEmptyMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapIsEmptyMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapIsEmptyMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapIsEmptyMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *HashMapIsEmptyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	isEmpty := m.source.hashmap.IsEmpty()
	return data.NewBoolValue(isEmpty), nil
}

// HashMapClearMethod 清空方法
type HashMapClearMethod struct {
	source *HashMapClass
}

func (m *HashMapClearMethod) GetName() string {
	return "clear"
}

func (m *HashMapClearMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapClearMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapClearMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapClearMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapClearMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *HashMapClearMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	m.source.hashmap.Clear()
	return data.NewNullValue(), nil
}

// HashMapKeysMethod 获取所有键方法
type HashMapKeysMethod struct {
	source *HashMapClass
}

func (m *HashMapKeysMethod) GetName() string {
	return "keys"
}

func (m *HashMapKeysMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapKeysMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapKeysMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapKeysMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapKeysMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func (m *HashMapKeysMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	keys := m.source.hashmap.Keys()
	// 转换为 data.Value 数组
	result := make([]data.Value, len(keys))
	for i, key := range keys {
		result[i] = data.NewStringValue(key)
	}
	return data.NewArrayValue(result), nil
}

// HashMapValuesMethod 获取所有值方法
type HashMapValuesMethod struct {
	source *HashMapClass
}

func (m *HashMapValuesMethod) GetName() string {
	return "values"
}

func (m *HashMapValuesMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapValuesMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapValuesMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapValuesMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapValuesMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func (m *HashMapValuesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	values := m.source.hashmap.Values()
	// 转换为 data.Value 数组
	result := make([]data.Value, len(values))
	copy(result, values)
	return data.NewArrayValue(result), nil
}

// HashMapCurrentMethod 当前值方法
type HashMapCurrentMethod struct {
	source *HashMapClass
}

func (m *HashMapCurrentMethod) GetName() string {
	return "current"
}

func (m *HashMapCurrentMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapCurrentMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapCurrentMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapCurrentMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapCurrentMethod) GetReturnType() data.Types {
	return nil
}

func (m *HashMapCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	return m.source.hashmap.Current(), nil
}

// HashMapKeyMethod 当前键方法
type HashMapKeyMethod struct {
	source *HashMapClass
}

func (m *HashMapKeyMethod) GetName() string {
	return "key"
}

func (m *HashMapKeyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapKeyMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapKeyMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapKeyMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapKeyMethod) GetReturnType() data.Types {
	return nil
}

func (m *HashMapKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	return m.source.hashmap.Key(), nil
}

// HashMapNextMethod 下一个元素方法
type HashMapNextMethod struct {
	source *HashMapClass
}

func (m *HashMapNextMethod) GetName() string {
	return "next"
}

func (m *HashMapNextMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapNextMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapNextMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapNextMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapNextMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *HashMapNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	m.source.hashmap.Next()
	return data.NewNullValue(), nil
}

// HashMapRewindMethod 重置迭代器方法
type HashMapRewindMethod struct {
	source *HashMapClass
}

func (m *HashMapRewindMethod) GetName() string {
	return "rewind"
}

func (m *HashMapRewindMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapRewindMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapRewindMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapRewindMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapRewindMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *HashMapRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	m.source.hashmap.Rewind()
	return data.NewNullValue(), nil
}

// HashMapValidMethod 检查迭代器是否有效方法
type HashMapValidMethod struct {
	source *HashMapClass
}

func (m *HashMapValidMethod) GetName() string {
	return "valid"
}

func (m *HashMapValidMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *HashMapValidMethod) GetIsStatic() bool {
	return false
}

func (m *HashMapValidMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *HashMapValidMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *HashMapValidMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *HashMapValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.hashmap == nil {
		return nil, utils.NewThrow(errors.New("HashMap 实例未初始化"))
	}

	valid := m.source.hashmap.Valid()
	return data.NewBoolValue(valid), nil
}

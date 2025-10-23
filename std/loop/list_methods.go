package loop

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ListConstructMethod 构造函数
type ListConstructMethod struct {
	source *ListClass
}

func (m *ListConstructMethod) GetName() string {
	return "__construct"
}

func (m *ListConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListConstructMethod) GetIsStatic() bool {
	return false
}

func (m *ListConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListConstructMethod) GetReturnType() data.Types {
	return nil
}

func (m *ListConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 构造函数不需要特殊处理，List 已经初始化
	return data.NewNullValue(), nil
}

// ListAddMethod 添加元素方法
type ListAddMethod struct {
	source *ListClass
}

func (m *ListAddMethod) GetName() string {
	return "add"
}

func (m *ListAddMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListAddMethod) GetIsStatic() bool {
	return false
}

func (m *ListAddMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "item", 0, nil, nil),
	}
}

func (m *ListAddMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "item", 0, nil),
	}
}

func (m *ListAddMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *ListAddMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::add() 需要至少一个参数"))
	}

	// 获取 List 实例
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	// 泛型类型检查
	if m.source.list.itemType != nil && !m.source.list.itemType.Is(value) {
		return nil, data.NewErrorThrow(nil, errors.New("类型不匹配: 期望 "+m.source.list.itemType.String()+" 但得到 "+value.AsString()))
	}

	// 添加元素
	m.source.list.Add(value)
	return data.NewNullValue(), nil
}

// ListGetMethod 获取元素方法
type ListGetMethod struct {
	source *ListClass
}

func (m *ListGetMethod) GetName() string {
	return "get"
}

func (m *ListGetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListGetMethod) GetIsStatic() bool {
	return false
}

func (m *ListGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.NewBaseType("int")),
	}
}

func (m *ListGetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
	}
}

func (m *ListGetMethod) GetReturnType() data.Types {
	return nil
}

func (m *ListGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	indexValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::get() 需要一个索引参数"))
	}

	// 获取索引
	index, ok := indexValue.(*data.IntValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::get() 的索引必须是整数"))
	}

	// 获取元素
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	item, found := m.source.list.Get(index.Value)
	if !found {
		return data.NewNullValue(), nil
	}

	return item, nil
}

// ListSetMethod 设置元素方法
type ListSetMethod struct {
	source *ListClass
}

func (m *ListSetMethod) GetName() string {
	return "set"
}

func (m *ListSetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListSetMethod) GetIsStatic() bool {
	return false
}

func (m *ListSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}

func (m *ListSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "value", 1, nil),
	}
}

func (m *ListSetMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *ListSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	indexValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::set() 需要索引参数"))
	}

	value, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::set() 需要值参数"))
	}

	// 获取索引
	index, ok := indexValue.(*data.IntValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::set() 的索引必须是整数"))
	}

	// 设置元素
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	// 泛型类型检查
	if m.source.list.itemType != nil && !m.source.list.itemType.Is(value) {
		return nil, data.NewErrorThrow(nil, errors.New("类型不匹配: 期望 "+m.source.list.itemType.String()+" 但得到 "+value.AsString()))
	}

	success := m.source.list.Set(index.Value, value)
	if !success {
		return nil, data.NewErrorThrow(nil, errors.New("索引超出范围"))
	}

	return data.NewNullValue(), nil
}

// ListSizeMethod 获取大小方法
type ListSizeMethod struct {
	source *ListClass
}

func (m *ListSizeMethod) GetName() string {
	return "size"
}

func (m *ListSizeMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListSizeMethod) GetIsStatic() bool {
	return false
}

func (m *ListSizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListSizeMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListSizeMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

func (m *ListSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	size := m.source.list.Size()
	return data.NewIntValue(size), nil
}

// ListRemoveMethod 移除元素方法
type ListRemoveMethod struct {
	source *ListClass
}

func (m *ListRemoveMethod) GetName() string {
	return "remove"
}

func (m *ListRemoveMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListRemoveMethod) GetIsStatic() bool {
	return false
}

func (m *ListRemoveMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "item", 0, nil, nil),
	}
}

func (m *ListRemoveMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "item", 0, nil),
	}
}

func (m *ListRemoveMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *ListRemoveMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::remove() 需要一个参数"))
	}

	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	success := m.source.list.Remove(value)
	return data.NewBoolValue(success), nil
}

// ListRemoveAtMethod 按索引移除元素方法
type ListRemoveAtMethod struct {
	source *ListClass
}

func (m *ListRemoveAtMethod) GetName() string {
	return "removeAt"
}

func (m *ListRemoveAtMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListRemoveAtMethod) GetIsStatic() bool {
	return false
}

func (m *ListRemoveAtMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.NewBaseType("int")),
	}
}

func (m *ListRemoveAtMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
	}
}

func (m *ListRemoveAtMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *ListRemoveAtMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	indexValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::removeAt() 需要一个索引参数"))
	}

	// 获取索引
	index, ok := indexValue.(*data.IntValue)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::removeAt() 的索引必须是整数"))
	}

	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	_, success := m.source.list.RemoveAt(index.Value)
	return data.NewBoolValue(success), nil
}

// ListClearMethod 清空列表方法
type ListClearMethod struct {
	source *ListClass
}

func (m *ListClearMethod) GetName() string {
	return "clear"
}

func (m *ListClearMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListClearMethod) GetIsStatic() bool {
	return false
}

func (m *ListClearMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListClearMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListClearMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *ListClearMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	m.source.list.Clear()
	return data.NewNullValue(), nil
}

// ListContainsMethod 检查是否包含元素方法
type ListContainsMethod struct {
	source *ListClass
}

func (m *ListContainsMethod) GetName() string {
	return "contains"
}

func (m *ListContainsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListContainsMethod) GetIsStatic() bool {
	return false
}

func (m *ListContainsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "item", 0, nil, nil),
	}
}

func (m *ListContainsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "item", 0, nil),
	}
}

func (m *ListContainsMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *ListContainsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::contains() 需要一个参数"))
	}

	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	contains := m.source.list.Contains(value)
	return data.NewBoolValue(contains), nil
}

// ListIndexOfMethod 获取元素索引方法
type ListIndexOfMethod struct {
	source *ListClass
}

func (m *ListIndexOfMethod) GetName() string {
	return "indexOf"
}

func (m *ListIndexOfMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListIndexOfMethod) GetIsStatic() bool {
	return false
}

func (m *ListIndexOfMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "item", 0, nil, nil),
	}
}

func (m *ListIndexOfMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "item", 0, nil),
	}
}

func (m *ListIndexOfMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

func (m *ListIndexOfMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("List::indexOf() 需要一个参数"))
	}

	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	index := m.source.list.IndexOf(value)
	return data.NewIntValue(index), nil
}

// ListIsEmptyMethod 检查是否为空方法
type ListIsEmptyMethod struct {
	source *ListClass
}

func (m *ListIsEmptyMethod) GetName() string {
	return "isEmpty"
}

func (m *ListIsEmptyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListIsEmptyMethod) GetIsStatic() bool {
	return false
}

func (m *ListIsEmptyMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListIsEmptyMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListIsEmptyMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *ListIsEmptyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	isEmpty := m.source.list.IsEmpty()
	return data.NewBoolValue(isEmpty), nil
}

// ListToArrayMethod 转换为数组方法
type ListToArrayMethod struct {
	source *ListClass
}

func (m *ListToArrayMethod) GetName() string {
	return "toArray"
}

func (m *ListToArrayMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListToArrayMethod) GetIsStatic() bool {
	return false
}

func (m *ListToArrayMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListToArrayMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListToArrayMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func (m *ListToArrayMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	array := m.source.list.ToArray()
	return data.NewArrayValue(array), nil
}

// Iterator 接口方法实现

// ListCurrentMethod 当前元素方法
type ListCurrentMethod struct {
	source *ListClass
}

func (m *ListCurrentMethod) GetName() string {
	return "current"
}

func (m *ListCurrentMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListCurrentMethod) GetIsStatic() bool {
	return false
}

func (m *ListCurrentMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListCurrentMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListCurrentMethod) GetReturnType() data.Types {
	return nil
}

func (m *ListCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	return m.source.list.Current(), nil
}

// ListKeyMethod 当前键方法
type ListKeyMethod struct {
	source *ListClass
}

func (m *ListKeyMethod) GetName() string {
	return "key"
}

func (m *ListKeyMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListKeyMethod) GetIsStatic() bool {
	return false
}

func (m *ListKeyMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListKeyMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListKeyMethod) GetReturnType() data.Types {
	return nil
}

func (m *ListKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	return m.source.list.Key(), nil
}

// ListNextMethod 下一个元素方法
type ListNextMethod struct {
	source *ListClass
}

func (m *ListNextMethod) GetName() string {
	return "next"
}

func (m *ListNextMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListNextMethod) GetIsStatic() bool {
	return false
}

func (m *ListNextMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListNextMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListNextMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *ListNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	m.source.list.Next()
	return data.NewNullValue(), nil
}

// ListRewindMethod 重置迭代器方法
type ListRewindMethod struct {
	source *ListClass
}

func (m *ListRewindMethod) GetName() string {
	return "rewind"
}

func (m *ListRewindMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListRewindMethod) GetIsStatic() bool {
	return false
}

func (m *ListRewindMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListRewindMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListRewindMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

func (m *ListRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	m.source.list.Rewind()
	return data.NewNullValue(), nil
}

// ListValidMethod 检查迭代器是否有效方法
type ListValidMethod struct {
	source *ListClass
}

func (m *ListValidMethod) GetName() string {
	return "valid"
}

func (m *ListValidMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ListValidMethod) GetIsStatic() bool {
	return false
}

func (m *ListValidMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *ListValidMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *ListValidMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

func (m *ListValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.list == nil {
		return nil, data.NewErrorThrow(nil, errors.New("List 实例未初始化"))
	}

	valid := m.source.list.Valid()
	return data.NewBoolValue(valid), nil
}

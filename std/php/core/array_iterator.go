package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayIteratorClass 实现 PHP 的 ArrayIterator 类
// 将数组包装为 Iterator 对象
type ArrayIteratorClass struct {
	node.Node
	array   data.Value // 存储的数组（可以是 ArrayValue 或 ObjectValue）
	pos     int        // 当前位置（仅用于 ArrayValue）
	currKey string     // 当前键（仅用于 ObjectValue）
	keys    []string   // 键列表（仅用于 ObjectValue）
}

// NewArrayIteratorClass 创建新的 ArrayIterator 实例
func NewArrayIteratorClass() *ArrayIteratorClass {
	return &ArrayIteratorClass{
		array:   &data.ArrayValue{List: []*data.ZVal{}},
		pos:     0,
		currKey: "",
		keys:    []string{},
	}
}

func (c *ArrayIteratorClass) GetName() string    { return "ArrayIterator" }
func (c *ArrayIteratorClass) GetExtend() *string { return nil }
func (c *ArrayIteratorClass) GetImplements() []string {
	return []string{"Iterator", "ArrayAccess", "Countable"}
}
func (c *ArrayIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *ArrayIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *ArrayIteratorClass) GetConstruct() data.Method {
	return &ArrayIteratorConstructMethod{instance: c}
}
func (c *ArrayIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 返回当前实例，不创建克隆
	// 因为构造函数已经在这个实例上设置了状态
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ArrayIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &ArrayIteratorConstructMethod{instance: c}, true
	case "rewind":
		return &ArrayIteratorRewindMethod{instance: c}, true
	case "valid":
		return &ArrayIteratorValidMethod{instance: c}, true
	case "current":
		return &ArrayIteratorCurrentMethod{instance: c}, true
	case "key":
		return &ArrayIteratorKeyMethod{instance: c}, true
	case "next":
		return &ArrayIteratorNextMethod{instance: c}, true
	case "offsetExists":
		return &ArrayIteratorOffsetExistsMethod{instance: c}, true
	case "offsetGet":
		return &ArrayIteratorOffsetGetMethod{instance: c}, true
	case "offsetSet":
		return &ArrayIteratorOffsetSetMethod{instance: c}, true
	case "offsetUnset":
		return &ArrayIteratorOffsetUnsetMethod{instance: c}, true
	case "count":
		return &ArrayIteratorCountMethod{instance: c}, true
	case "append":
		return &ArrayIteratorAppendMethod{instance: c}, true
	case "getArrayCopy":
		return &ArrayIteratorGetArrayCopyMethod{instance: c}, true
	}
	return nil, false
}

func (c *ArrayIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&ArrayIteratorConstructMethod{instance: c},
		&ArrayIteratorRewindMethod{instance: c},
		&ArrayIteratorValidMethod{instance: c},
		&ArrayIteratorCurrentMethod{instance: c},
		&ArrayIteratorKeyMethod{instance: c},
		&ArrayIteratorNextMethod{instance: c},
		&ArrayIteratorOffsetExistsMethod{instance: c},
		&ArrayIteratorOffsetGetMethod{instance: c},
		&ArrayIteratorOffsetSetMethod{instance: c},
		&ArrayIteratorOffsetUnsetMethod{instance: c},
		&ArrayIteratorCountMethod{instance: c},
		&ArrayIteratorAppendMethod{instance: c},
		&ArrayIteratorGetArrayCopyMethod{instance: c},
	}
}

// __construct(array $array = [], int $flags = 0)
type ArrayIteratorConstructMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *ArrayIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorConstructMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorConstructMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, data.NewArrayValue(nil), data.Mixed{}),
	}
}
func (m *ArrayIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}
func (m *ArrayIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}

	// 支持 ArrayValue 和 ObjectValue（关联数组）
	switch v := arrVal.(type) {
	case *data.ArrayValue:
		m.instance.array = v
		m.instance.pos = 0
		m.instance.keys = []string{}
	case *data.ObjectValue:
		m.instance.array = v
		// 获取所有键并保存顺序
		keys := []string{}
		v.RangeProperties(func(key string, value data.Value) bool {
			keys = append(keys, key)
			return true
		})
		m.instance.keys = keys
		m.instance.currKey = ""
		if len(keys) > 0 {
			m.instance.currKey = keys[0]
		}
	default:
		m.instance.array = &data.ArrayValue{List: []*data.ZVal{}}
		m.instance.pos = 0
		m.instance.keys = []string{}
	}

	return nil, nil
}

// rewind()
type ArrayIteratorRewindMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorRewindMethod) GetName() string               { return "rewind" }
func (m *ArrayIteratorRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorRewindMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorRewindMethod) GetReturnType() data.Types     { return nil }
func (m *ArrayIteratorRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorRewindMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 根据数组类型重置不同的状态
	if _, ok := m.instance.array.(*data.ArrayValue); ok {
		m.instance.pos = 0
	} else if _, ok := m.instance.array.(*data.ObjectValue); ok {
		if len(m.instance.keys) > 0 {
			m.instance.currKey = m.instance.keys[0]
		}
	}
	return nil, nil
}

// valid()
type ArrayIteratorValidMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorValidMethod) GetName() string               { return "valid" }
func (m *ArrayIteratorValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorValidMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *ArrayIteratorValidMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorValidMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 根据数组类型检查有效性
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		valid := m.instance.pos >= 0 && m.instance.pos < len(arr.List)
		return data.NewBoolValue(valid), nil
	} else if _, ok := m.instance.array.(*data.ObjectValue); ok {
		valid := len(m.instance.keys) > 0 && m.instance.currKey != ""
		return data.NewBoolValue(valid), nil
	}
	return data.NewBoolValue(false), nil
}

// current()
type ArrayIteratorCurrentMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorCurrentMethod) GetName() string               { return "current" }
func (m *ArrayIteratorCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorCurrentMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 根据数组类型获取当前值
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		if m.instance.pos < 0 || m.instance.pos >= len(arr.List) {
			return data.NewNullValue(), nil
		}
		return arr.List[m.instance.pos].Value, nil
	} else if arr, ok := m.instance.array.(*data.ObjectValue); ok {
		if m.instance.currKey == "" {
			return data.NewNullValue(), nil
		}
		val, _ := arr.GetProperty(m.instance.currKey)
		return val, nil
	}
	return data.NewNullValue(), nil
}

// key()
type ArrayIteratorKeyMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorKeyMethod) GetName() string               { return "key" }
func (m *ArrayIteratorKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorKeyMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorKeyMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 根据数组类型返回键
	if _, ok := m.instance.array.(*data.ArrayValue); ok {
		return data.NewIntValue(m.instance.pos), nil
	} else if _, ok := m.instance.array.(*data.ObjectValue); ok {
		if m.instance.currKey == "" {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(m.instance.currKey), nil
	}
	return data.NewNullValue(), nil
}

// next()
type ArrayIteratorNextMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorNextMethod) GetName() string               { return "next" }
func (m *ArrayIteratorNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorNextMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorNextMethod) GetReturnType() data.Types     { return nil }
func (m *ArrayIteratorNextMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorNextMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 根据数组类型推进迭代器
	if _, ok := m.instance.array.(*data.ArrayValue); ok {
		m.instance.pos++
	} else if _, ok := m.instance.array.(*data.ObjectValue); ok {
		// 找到当前键在 keys 切片中的位置
		for i, key := range m.instance.keys {
			if key == m.instance.currKey {
				if i+1 < len(m.instance.keys) {
					m.instance.currKey = m.instance.keys[i+1]
				} else {
					m.instance.currKey = ""
				}
				break
			}
		}
	}
	return nil, nil
}

// offsetExists($offset)
type ArrayIteratorOffsetExistsMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorOffsetExistsMethod) GetName() string            { return "offsetExists" }
func (m *ArrayIteratorOffsetExistsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetExistsMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetExistsMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *ArrayIteratorOffsetExistsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorOffsetExistsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayIteratorOffsetExistsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	offsetVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}

	// 根据数组类型处理
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		// 检查整数索引
		if intVal, ok := offsetVal.(*data.IntValue); ok {
			idx := int(intVal.Value)
			return data.NewBoolValue(idx >= 0 && idx < len(arr.List)), nil
		}
	} else if arr, ok := m.instance.array.(*data.ObjectValue); ok {
		// 检查字符串键
		if strVal, ok := offsetVal.(*data.StringValue); ok {
			val, _ := arr.GetProperty(strVal.Value)
			if _, isNull := val.(*data.NullValue); !isNull {
				return data.NewBoolValue(true), nil
			}
		}
		return data.NewBoolValue(false), nil
	}

	return data.NewBoolValue(false), nil
}

// offsetGet($offset)
type ArrayIteratorOffsetGetMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorOffsetGetMethod) GetName() string            { return "offsetGet" }
func (m *ArrayIteratorOffsetGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetGetMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetGetMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *ArrayIteratorOffsetGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorOffsetGetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayIteratorOffsetGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	offsetVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}

	// 根据数组类型处理
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		// 检查整数索引
		if intVal, ok := offsetVal.(*data.IntValue); ok {
			idx := int(intVal.Value)
			if idx >= 0 && idx < len(arr.List) {
				return arr.List[idx].Value, nil
			}
		}
	} else if arr, ok := m.instance.array.(*data.ObjectValue); ok {
		// 检查字符串键
		if strVal, ok := offsetVal.(*data.StringValue); ok {
			val, _ := arr.GetProperty(strVal.Value)
			return val, nil
		}
	}

	return data.NewNullValue(), nil
}

// offsetSet($offset, $value)
type ArrayIteratorOffsetSetMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorOffsetSetMethod) GetName() string            { return "offsetSet" }
func (m *ArrayIteratorOffsetSetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetSetMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetSetMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorOffsetSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "offset", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "value", 1, nil, data.Mixed{}),
	}
}
func (m *ArrayIteratorOffsetSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "offset", 0, data.Mixed{}),
		node.NewVariable(nil, "value", 1, data.Mixed{}),
	}
}
func (m *ArrayIteratorOffsetSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	offsetVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}
	valueVal, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, nil
	}

	// 根据数组类型处理
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		// 如果 offset 为 null，则追加到数组末尾
		if _, isNull := offsetVal.(*data.NullValue); isNull {
			arr.List = append(arr.List, data.NewZVal(valueVal))
			return nil, nil
		}

		// 整数索引
		if intVal, ok := offsetVal.(*data.IntValue); ok {
			idx := int(intVal.Value)
			if idx >= 0 && idx < len(arr.List) {
				arr.List[idx] = data.NewZVal(valueVal)
			} else if idx == len(arr.List) {
				arr.List = append(arr.List, data.NewZVal(valueVal))
			}
			return nil, nil
		}
	} else if arr, ok := m.instance.array.(*data.ObjectValue); ok {
		// ObjectValue - 字符串键
		if strVal, ok := offsetVal.(*data.StringValue); ok {
			arr.SetProperty(strVal.Value, valueVal)
			// 如果这个键不在 keys 列表中，添加它
			found := false
			for _, key := range m.instance.keys {
				if key == strVal.Value {
					found = true
					break
				}
			}
			if !found {
				m.instance.keys = append(m.instance.keys, strVal.Value)
				if m.instance.currKey == "" {
					m.instance.currKey = strVal.Value
				}
			}
			return nil, nil
		}
	}

	return nil, nil
}

// offsetUnset($offset)
type ArrayIteratorOffsetUnsetMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorOffsetUnsetMethod) GetName() string            { return "offsetUnset" }
func (m *ArrayIteratorOffsetUnsetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorOffsetUnsetMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorOffsetUnsetMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorOffsetUnsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "offset", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorOffsetUnsetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "offset", 0, data.Mixed{})}
}
func (m *ArrayIteratorOffsetUnsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	offsetVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}

	// 根据数组类型处理
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		// 整数索引 - 设置为 null
		if intVal, ok := offsetVal.(*data.IntValue); ok {
			idx := int(intVal.Value)
			if idx >= 0 && idx < len(arr.List) {
				arr.List[idx] = data.NewZVal(data.NewNullValue())
			}
			return nil, nil
		}
	} else if arr, ok := m.instance.array.(*data.ObjectValue); ok {
		// ObjectValue - 字符串键
		if strVal, ok := offsetVal.(*data.StringValue); ok {
			arr.SetProperty(strVal.Value, data.NewNullValue())
			// 从 keys 列表中移除
			newKeys := []string{}
			for _, key := range m.instance.keys {
				if key != strVal.Value {
					newKeys = append(newKeys, key)
				}
			}
			m.instance.keys = newKeys
			// 如果移除的是当前键，更新 currKey
			if strVal.Value == m.instance.currKey {
				if len(m.instance.keys) > 0 {
					m.instance.currKey = m.instance.keys[0]
				} else {
					m.instance.currKey = ""
				}
			}
			return nil, nil
		}
	}

	return nil, nil
}

// count()
type ArrayIteratorCountMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorCountMethod) GetName() string               { return "count" }
func (m *ArrayIteratorCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorCountMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *ArrayIteratorCountMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorCountMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 根据数组类型返回长度
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		return data.NewIntValue(len(arr.List)), nil
	} else if arr, ok := m.instance.array.(*data.ObjectValue); ok {
		props := arr.GetProperties()
		return data.NewIntValue(len(props)), nil
	}
	return data.NewIntValue(0), nil
}

// append($value)
type ArrayIteratorAppendMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorAppendMethod) GetName() string            { return "append" }
func (m *ArrayIteratorAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ArrayIteratorAppendMethod) GetIsStatic() bool          { return false }
func (m *ArrayIteratorAppendMethod) GetReturnType() data.Types  { return nil }
func (m *ArrayIteratorAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *ArrayIteratorAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *ArrayIteratorAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, nil
	}

	if val != nil {
		// 只支持 ArrayValue 的 append
		if arr, ok := m.instance.array.(*data.ArrayValue); ok {
			arr.List = append(arr.List, data.NewZVal(val))
		}
	}
	return nil, nil
}

// getArrayCopy()
type ArrayIteratorGetArrayCopyMethod struct {
	instance *ArrayIteratorClass
}

func (m *ArrayIteratorGetArrayCopyMethod) GetName() string               { return "getArrayCopy" }
func (m *ArrayIteratorGetArrayCopyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *ArrayIteratorGetArrayCopyMethod) GetIsStatic() bool             { return false }
func (m *ArrayIteratorGetArrayCopyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *ArrayIteratorGetArrayCopyMethod) GetParams() []data.GetValue    { return nil }
func (m *ArrayIteratorGetArrayCopyMethod) GetVariables() []data.Variable { return nil }
func (m *ArrayIteratorGetArrayCopyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 根据数组类型返回克隆
	if arr, ok := m.instance.array.(*data.ArrayValue); ok {
		return data.CloneArrayValue(arr), nil
	} else if arr, ok := m.instance.array.(*data.ObjectValue); ok {
		// 克隆 ObjectValue
		clone := data.NewObjectValue()
		arr.RangeProperties(func(key string, value data.Value) bool {
			clone.SetProperty(key, value)
			return true
		})
		return clone, nil
	}
	return data.NewNullValue(), nil
}

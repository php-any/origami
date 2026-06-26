package data

import (
	"fmt"
	"strconv"
)

func NewArrayValue(v []Value) Value {
	list := make([]*ZVal, len(v))
	for i, val := range v {
		list[i] = NewZVal(val)
	}
	return &ArrayValue{
		List: list,
	}
}

// CloneArrayValue 创建一个新的 ArrayValue。
// 为了性能，这里仅复制 []*ZVal 切片本身（浅拷贝），不重新分配每个 ZVal：
// - 结构不共享：两个数组的 List 是不同的 slice，结构性修改（如 array_shift/append）互不影响
// - 元素仍按 ZVal 语义工作：写入单个元素时会替换对应 ZVal，不会影响其他数组
func CloneArrayValue(src *ArrayValue) *ArrayValue {
	if src == nil {
		return nil
	}
	list := make([]*ZVal, len(src.List))
	copy(list, src.List)
	return &ArrayValue{
		List:                  list,
		IndirectOverloadClass: src.IndirectOverloadClass,
	}
}

// CloneArrayValueForCallArgs 为 __call 的 $arguments 克隆数组实参：非引用元素复制 ZVal，引用槽位共享 ZVal
func CloneArrayValueForCallArgs(src *ArrayValue) *ArrayValue {
	if src == nil {
		return nil
	}
	list := make([]*ZVal, len(src.List))
	for i, z := range src.List {
		if z == nil {
			continue
		}
		if z.RefSlotCount > 0 {
			list[i] = z
		} else {
			list[i] = NewZVal(z.Value)
		}
	}
	return &ArrayValue{List: list}
}

type ArrayValue struct {
	List     []*ZVal
	iterator int // 迭代器当前位置索引
	// IndirectOverloadClass 非空表示该数组来自 ArrayAccess::offsetGet 的副本，对其元素的间接修改无效
	IndirectOverloadClass string
}

func (a *ArrayValue) Current(ctx Context) (Value, Control) {
	if a.iterator >= len(a.List) {
		return NewNullValue(), nil
	}
	return a.List[a.iterator].Value, nil
}

func (a *ArrayValue) Key(ctx Context) (Value, Control) {
	if a.iterator >= 0 && a.iterator < len(a.List) {
		z := a.List[a.iterator]
		if z != nil && z.Name != "" {
			return NewStringValue(z.Name), nil
		}
	}
	return NewIntValue(a.iterator), nil
}

func (a *ArrayValue) Next(ctx Context) Control {
	a.iterator++
	return nil
}

func (a *ArrayValue) Rewind(ctx Context) (Value, Control) {
	a.iterator = 0
	return nil, nil
}

func (a *ArrayValue) Valid(ctx Context) (Value, Control) {
	valid := a.iterator >= 0 && a.iterator < len(a.List)
	return NewBoolValue(valid), nil
}

func (a *ArrayValue) GetValue(ctx Context) (GetValue, Control) {
	return a, nil
}

func (a *ArrayValue) AsString() string {
	str := "["
	for _, zval := range a.List {
		str = str + zval.Value.AsString() + ", "
	}
	if len(str) > 2 {
		str = str[:len(str)-2]
	}

	str = str + "]"
	return fmt.Sprintf("%s", str)
}

func (a *ArrayValue) AsBool() (bool, error) {
	return len(a.List) > 0, nil
}

func (a *ArrayValue) GetMethod(name string) (Method, bool) {
	switch name {
	case "push":
		return &ArrayValuePush{&a.List}, true
	case "pop":
		return &ArrayValuePop{&a.List}, true
	case "shift":
		return &ArrayValueShift{&a.List}, true
	case "unshift":
		return &ArrayValueUnshift{&a.List}, true
	case "slice":
		return &ArrayValueSlice{a.List}, true
	case "splice":
		return &ArrayValueSplice{&a.List}, true
	case "join":
		return &ArrayValueJoin{a.List}, true
	case "reverse":
		return &ArrayValueReverse{a.List}, true
	case "sort":
		return &ArrayValueSort{&a.List}, true
	case "indexOf":
		return &ArrayValueIndexOf{a.List}, true
	case "includes":
		return &ArrayValueIncludes{a.List}, true
	case "forEach":
		return &ArrayValueForEach{a.List}, true
	case "map":
		return &ArrayValueMap{a.List}, true
	case "filter":
		return &ArrayValueFilter{a.List}, true
	case "reduce":
		return &ArrayValueReduce{a.List}, true
	case "concat":
		return &ArrayValueConcat{a.List}, true
	case "every":
		return &ArrayValueEvery{a.List}, true
	case "some":
		return &ArrayValueSome{a.List}, true
	case "find":
		return &ArrayValueFind{a.List}, true
	case "findIndex":
		return &ArrayValueFindIndex{a.List}, true
	case "flat":
		return &ArrayValueFlat{a.List}, true
	case "flatMap":
		return &ArrayValueFlatMap{a.List}, true
	}

	return nil, false
}

func (a *ArrayValue) GetProperty(name string) (Value, Control) {
	switch name {
	case "length":
		return NewIntValue(len(a.List)), nil
	}
	return nil, NewErrorThrow(nil, fmt.Errorf("ArrayValue.GetProperty called with name %s", name))
}

func (a *ArrayValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalArray(a)
}

func (a *ArrayValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalArray(data, a)
}

func (a *ArrayValue) ToGoValue(serializer Serializer) (any, error) {
	return serializer.MarshalArray(a)
}

func (a *ArrayValue) ToValueList() []Value {
	args := make([]Value, len(a.List))
	for i, zval := range a.List {
		args[i] = zval.Value
	}
	return args
}

// IntArrayKeyName 将整数键编码为 ZVal.Name（稀疏整数键，插入顺序在 List 末尾）
func IntArrayKeyName(i int) string {
	return strconv.Itoa(i)
}

// ParseIntArrayKeyName 若 name 为纯整数字符串则返回该整数键，否则 ok=false
func ParseIntArrayKeyName(name string) (int, bool) {
	if name == "" {
		return 0, false
	}
	n, err := strconv.Atoi(name)
	if err != nil {
		return 0, false
	}
	if strconv.Itoa(n) != name {
		return 0, false
	}
	return n, true
}

// FindSlotByIntKey 按 PHP 整数键查找槽位（含稀疏键 Name=="6" 等）
func (a *ArrayValue) FindSlotByIntKey(i int) (*ZVal, int) {
	keyStr := IntArrayKeyName(i)
	for j, z := range a.List {
		if z != nil && z.Name == keyStr {
			return z, j
		}
	}
	if i >= 0 && i < len(a.List) {
		if z := a.List[i]; z != nil && z.Name == "" {
			return z, i
		}
	}
	return nil, -1
}

// SetIntKey 设置整数键（不将稀疏数组转为 ObjectValue）
func (a *ArrayValue) SetIntKey(i int, value Value) {
	if z, _ := a.FindSlotByIntKey(i); z != nil {
		z.Value = value
		return
	}
	if i < 0 {
		return
	}
	if i == len(a.List) {
		a.List = append(a.List, NewZVal(value))
		return
	}
	if i > len(a.List) {
		a.List = append(a.List, NewNamedZVal(IntArrayKeyName(i), value))
		return
	}
	if slot := a.List[i]; slot != nil && slot.RefSlotCount > 0 {
		slot.Value = value
	} else {
		a.List[i] = NewZVal(value)
	}
}

// normalizeDenseIntKeys 将 Name=="" 的连续槽位转为显式整数字符串键，避免 unset 中间元素时误压缩后续键
func (a *ArrayValue) normalizeDenseIntKeys() {
	for j, z := range a.List {
		if z != nil && z.Name == "" {
			z.Name = IntArrayKeyName(j)
		}
	}
}

// UnsetKey 删除整数或字符串键（不存在则无操作）
func (a *ArrayValue) UnsetKey(index Value) {
	if iv, ok := index.(AsInt); ok {
		i, err := iv.AsInt()
		if err != nil {
			return
		}
		a.normalizeDenseIntKeys()
		keyStr := IntArrayKeyName(i)
		for j, z := range a.List {
			if z != nil && z.Name == keyStr {
				a.List = append(a.List[:j], a.List[j+1:]...)
				return
			}
		}
		return
	}
	if sv, ok := index.(AsString); ok {
		key := sv.AsString()
		for j, z := range a.List {
			if z != nil && z.Name == key {
				a.List = append(a.List[:j], a.List[j+1:]...)
				return
			}
		}
	}
}

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
// 复制 []*ZVal 切片本身，并对每个非引用元素重新分配 ZVal（对齐 PHP 数组 copy-on-write 语义）：
// - 结构不共享：两个数组的 List 是不同的 slice，结构性修改（如 array_shift/append）互不影响
// - 元素不共享：写入/替换单个元素时不会影响其他数组（因为各自持有独立的 ZVal）
// - 通过 &$arr[i] 绑定的引用槽位（RefSlotCount > 0）仍共享同一 ZVal，保持引用语义
func CloneArrayValue(src *ArrayValue) *ArrayValue {
	if src == nil {
		return nil
	}
	list := make([]*ZVal, len(src.List))
	for i, z := range src.List {
		if z == nil {
			continue
		}
		if z.RefSlotCount > 0 {
			// 引用槽位共享，保证引用赋值（&$arr[i]）语义
			list[i] = z
		} else {
			// 复制 ZVal 时保留 Name（关联数组键）
			list[i] = CopyZValKeepName(z, z.Value)
		}
	}
	return &ArrayValue{
		List:                  list,
		IndirectOverloadClass: src.IndirectOverloadClass,
	}
}

// DeepCloneArrayValue 深度克隆一个 ArrayValue，用于 PHP clone 对象时对数组类型属性做拷贝。
// 与 PHP 语义对齐：
//   - 数组（含嵌套数组）按值拷贝，结构上与原数组互不影响
//   - 数组内的对象（*ClassValue / *ObjectValue）仍按引用共享
//
// 参数 depth 用于防御极端深度的嵌套数组（避免栈溢出）。
func DeepCloneArrayValue(src *ArrayValue) *ArrayValue {
	if src == nil {
		return nil
	}
	return deepCloneArrayValue(src, 0)
}

func deepCloneArrayValue(src *ArrayValue, depth int) *ArrayValue {
	if src == nil {
		return nil
	}
	const maxDepth = 64
	list := make([]*ZVal, len(src.List))
	for i, z := range src.List {
		if z == nil {
			continue
		}
		if depth < maxDepth {
			// 嵌套数组/关联数组按值拷贝；对象与标量保持引用共享（与 PHP clone 语义一致）
			list[i] = CopyZValKeepName(z, deepCloneValue(z.Value, depth+1))
		} else {
			list[i] = CopyZValKeepName(z, z.Value)
		}
	}
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
			// 复制 ZVal 时保留 Name（关联数组键），避免在 __call/__callStatic 参数传递中丢失键
			list[i] = CopyZValKeepName(z, z.Value)
		}
	}
	return &ArrayValue{List: list, rc: 1}
}

type ArrayValue struct {
	List     []*ZVal
	iterator int // 迭代器当前位置索引
	// IndirectOverloadClass 非空表示该数组来自 ArrayAccess::offsetGet 的副本，对其元素的间接修改无效
	IndirectOverloadClass string
	rc                    int // 指向该容器的 zval 数（copy-on-write）
	keyIndex              map[string]int
	idxLen                int
	packed                bool // 全部槽位 Name==""，整数键即 List 下标
	maxIntKey             int  // 当前最大整数键；无整数键时为 -1
}

func (a *ArrayValue) Current(ctx Context) (Value, Control) {
	if a.iterator >= len(a.List) {
		return NewNullValue(), nil
	}
	return a.List[a.iterator].Value, nil
}

func (a *ArrayValue) Key(ctx Context) (Value, Control) {
	if a.iterator >= 0 && a.iterator < len(a.List) {
		return a.List[a.iterator].PHPArrayKey(a.iterator), nil
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

func (a *ArrayValue) invalidateIndex() {
	a.keyIndex = nil
	a.idxLen = -1
	a.packed = false
	a.maxIntKey = -1
}

func (a *ArrayValue) ensureIndex() {
	if a.idxLen == len(a.List) && (a.packed || a.keyIndex != nil) {
		return
	}
	a.rebuildIndex()
}

func (a *ArrayValue) rebuildIndex() {
	n := len(a.List)
	a.idxLen = n
	packed := true
	idx := make(map[string]int, n)
	hasNamed := false
	max := -1
	for i, z := range a.List {
		if z == nil {
			continue
		}
		if z.EmptyStrKey {
			packed = false
			hasNamed = true
			idx[""] = i
			continue
		}
		if z.Name != "" {
			packed = false
			hasNamed = true
			idx[z.Name] = i
			if k, ok := ParseIntArrayKeyName(z.Name); ok && k > max {
				max = k
			}
			continue
		}
		if i > max {
			max = i
		}
	}
	a.packed = packed
	a.maxIntKey = max
	if hasNamed {
		a.keyIndex = idx
	} else {
		a.keyIndex = nil
	}
}

// NextAppendIntKey 对齐 PHP $a[]：最大整数键 + 1；没有整数键时为 0。
func (a *ArrayValue) NextAppendIntKey() int {
	a.ensureIndex()
	if a.maxIntKey < 0 {
		return 0
	}
	return a.maxIntKey + 1
}

func (a *ArrayValue) AppendValue(value Value) {
	a.SetIntKey(a.NextAppendIntKey(), value)
}

func (a *ArrayValue) AppendSlot(value Value) *ZVal {
	i := a.NextAppendIntKey()
	a.SetIntKey(i, value)
	z, _ := a.FindSlotByIntKey(i)
	return z
}

// LookupZValByStringKey 按字符串键查找槽位；纯数字字符串键会回退整数键查找（如 "0" → 列表下标 0）
func (a *ArrayValue) LookupZValByStringKey(key string) (*ZVal, bool) {
	a.ensureIndex()
	if a.keyIndex != nil {
		if i, ok := a.keyIndex[key]; ok && i >= 0 && i < len(a.List) {
			z := a.List[i]
			if z != nil && (z.Name == key || (key == "" && z.EmptyStrKey)) {
				return z, true
			}
		}
	}
	if i, ok := ParseIntArrayKeyName(key); ok {
		if z, _ := a.FindSlotByIntKey(i); z != nil {
			return z, true
		}
	}
	return nil, false
}

// SetStringKey 写入 PHP 字符串键（含 ”）；纯数字字符串走整数键。
func (a *ArrayValue) SetStringKey(key string, value Value) {
	if n, ok := ParseIntArrayKeyName(key); ok {
		a.SetIntKey(n, value)
		return
	}
	if z, ok := a.LookupZValByStringKey(key); ok {
		z.Value = value
		return
	}
	a.invalidateIndex()
	if key == "" {
		a.List = append(a.List, NewEmptyStringKeyZVal(value))
		return
	}
	a.List = append(a.List, NewNamedZVal(key, value))
}

// FindSlotByIntKey 按 PHP 整数键查找槽位（含稀疏键 Name=="6" 等）
func (a *ArrayValue) FindSlotByIntKey(i int) (*ZVal, int) {
	a.ensureIndex()
	if a.packed {
		if i >= 0 && i < len(a.List) {
			if z := a.List[i]; z.IsPackedIntSlot() {
				return z, i
			}
		}
		return nil, -1
	}
	keyStr := IntArrayKeyName(i)
	if a.keyIndex != nil {
		if j, ok := a.keyIndex[keyStr]; ok && j >= 0 && j < len(a.List) {
			z := a.List[j]
			if z != nil && z.Name == keyStr {
				return z, j
			}
		}
	}
	if i >= 0 && i < len(a.List) {
		if z := a.List[i]; z.IsPackedIntSlot() {
			return z, i
		}
	}
	return nil, -1
}

// SetIntKey 设置整数键（不将稀疏数组转为 ObjectValue）
func (a *ArrayValue) SetIntKey(i int, value Value) {
	a.ensureIndex()
	if z, _ := a.FindSlotByIntKey(i); z != nil {
		z.Value = value
		return
	}
	if i < 0 {
		return
	}
	packed, n := a.packed, len(a.List)
	a.invalidateIndex()
	// 未找到整数键 i 时只能追加。禁止用 List[i] 覆盖：该槽可能是 PHP 空字符串键 ''。
	if packed && i == n {
		a.List = append(a.List, NewZVal(value))
		return
	}
	a.List = append(a.List, NewNamedZVal(IntArrayKeyName(i), value))
}

// normalizeDenseIntKeys 将 Name=="" 的连续槽位转为显式整数字符串键，避免 unset 中间元素时误压缩后续键
func (a *ArrayValue) normalizeDenseIntKeys() {
	for j, z := range a.List {
		if z != nil && z.IsPackedIntSlot() {
			z.Name = IntArrayKeyName(j)
		}
	}
	a.invalidateIndex()
}

// UnsetKey 删除整数或字符串键（不存在则无操作）
func (a *ArrayValue) UnsetKey(index Value) {
	a.invalidateIndex()
	// 先按字符串键处理：StringValue 同时实现 AsInt，非数字字符串不能在 AsInt 失败后直接 return
	if sv, ok := index.(AsString); ok {
		key := sv.AsString()
		if _, isIntKey := ParseIntArrayKeyName(key); !isIntKey {
			for j, z := range a.List {
				if z == nil {
					continue
				}
				if key == "" {
					if z.EmptyStrKey {
						a.List = append(a.List[:j], a.List[j+1:]...)
						return
					}
					continue
				}
				if z.Name == key && !z.EmptyStrKey {
					a.List = append(a.List[:j], a.List[j+1:]...)
					return
				}
			}
			return
		}
		// 纯数字字符串键：按整数键删除
		if n, err := strconv.Atoi(key); err == nil {
			a.normalizeDenseIntKeys()
			keyStr := IntArrayKeyName(n)
			for j, z := range a.List {
				if z != nil && z.Name == keyStr {
					a.List = append(a.List[:j], a.List[j+1:]...)
					return
				}
			}
			return
		}
	}
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
	}
}

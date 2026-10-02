package data

import (
	"fmt"
	"math"
	"strconv"
)

func NewArrayValue(v []Value) Value {
	list := make([]*ZVal, len(v))
	for i, val := range v {
		list[i] = NewZVal(val)
	}
	return &ArrayValue{
		flatArrayStore: FlatArrayStore{entries: list},
	}
}

// CloneArrayValue 创建一个新的 ArrayValue。
// 复制 []*ZVal 切片本身，并对每个非引用元素重新分配 ZVal（对齐 PHP 数组 copy-on-write 语义）：
// - 结构不共享：两个数组的槽位切片不同，结构性修改（如 array_shift/append）互不影响
// - 元素不共享：写入/替换单个元素时不会影响其他数组（因为各自持有独立的 ZVal）
// - 通过 &$arr[i] 绑定的引用槽位（RefSlotCount > 0）仍共享同一 ZVal，保持引用语义
func CloneArrayValue(src *ArrayValue) *ArrayValue {
	if src == nil {
		return nil
	}
	// Borrow one structural view. The source must be request-owned or immutable;
	// a slice-header snapshot does not synchronize concurrent mutations.
	srcList := src.slots()
	list := make([]*ZVal, len(srcList))
	for i, z := range srcList {
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
		flatArrayStore:        FlatArrayStore{entries: list, appendKeyKnown: src.appendKeyKnown, intKeySeen: src.intKeySeen, nextIntKey: src.nextIntKey, iterator: src.iterator},
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
	srcList := src.slots() // 同上：只借用一次只读槽位视图
	list := make([]*ZVal, len(srcList))
	for i, z := range srcList {
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
		flatArrayStore:        FlatArrayStore{entries: list, appendKeyKnown: src.appendKeyKnown, intKeySeen: src.intKeySeen, nextIntKey: src.nextIntKey, iterator: src.iterator},
		IndirectOverloadClass: src.IndirectOverloadClass,
	}
}

// CloneArrayValueForCallArgs 为 __call 的 $arguments 克隆数组实参：非引用元素复制 ZVal，引用槽位共享 ZVal
func CloneArrayValueForCallArgs(src *ArrayValue) *ArrayValue {
	if src == nil {
		return nil
	}
	srcList := src.slots() // 同上：只借用一次只读槽位视图
	list := make([]*ZVal, len(srcList))
	for i, z := range srcList {
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
	return &ArrayValue{flatArrayStore: FlatArrayStore{entries: list, appendKeyKnown: src.appendKeyKnown, intKeySeen: src.intKeySeen, nextIntKey: src.nextIntKey, iterator: src.iterator}, rc: 1}
}

type ArrayValue struct {
	flatArrayStore
	// IndirectOverloadClass identifies an ArrayAccess::offsetGet value copy.
	IndirectOverloadClass string
	rc                    int
}

func (a *ArrayValue) Current(ctx Context) (Value, Control) {
	if a.iterator < 0 || a.iterator >= len(a.entries) {
		return NewNullValue(), nil
	}
	return a.entries[a.iterator].Value, nil
}

func (a *ArrayValue) Key(ctx Context) (Value, Control) {
	if a.iterator >= 0 && a.iterator < len(a.entries) {
		return a.entries[a.iterator].PHPArrayKey(a.iterator), nil
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
	valid := a.iterator >= 0 && a.iterator < len(a.entries)
	return NewBoolValue(valid), nil
}

func (a *ArrayValue) GetValue(ctx Context) (GetValue, Control) {
	return a, nil
}

func (a *ArrayValue) AsString() string {
	str := "["
	for _, zval := range a.entries {
		str = str + zval.Value.AsString() + ", "
	}
	if len(str) > 2 {
		str = str[:len(str)-2]
	}

	str = str + "]"
	return fmt.Sprintf("%s", str)
}

func (a *ArrayValue) AsBool() (bool, error) {
	return len(a.entries) > 0, nil
}

func (a *ArrayValue) GetMethod(name string) (Method, bool) {
	switch name {
	case "push":
		return &ArrayValuePush{a}, true
	case "pop":
		return &ArrayValuePop{a}, true
	case "shift":
		return &ArrayValueShift{a}, true
	case "unshift":
		return &ArrayValueUnshift{a}, true
	case "slice":
		return &ArrayValueSlice{a}, true
	case "splice":
		return &ArrayValueSplice{a}, true
	case "join":
		return &ArrayValueJoin{a}, true
	case "reverse":
		return &ArrayValueReverse{a}, true
	case "sort":
		return &ArrayValueSort{a}, true
	case "indexOf":
		return &ArrayValueIndexOf{a}, true
	case "includes":
		return &ArrayValueIncludes{a}, true
	case "forEach":
		return &ArrayValueForEach{a}, true
	case "map":
		return &ArrayValueMap{a}, true
	case "filter":
		return &ArrayValueFilter{a}, true
	case "reduce":
		return &ArrayValueReduce{a}, true
	case "concat":
		return &ArrayValueConcat{a}, true
	case "every":
		return &ArrayValueEvery{a}, true
	case "some":
		return &ArrayValueSome{a}, true
	case "find":
		return &ArrayValueFind{a}, true
	case "findIndex":
		return &ArrayValueFindIndex{a}, true
	case "flat":
		return &ArrayValueFlat{a}, true
	case "flatMap":
		return &ArrayValueFlatMap{a}, true
	}

	return nil, false
}

func (a *ArrayValue) GetProperty(name string) (Value, Control) {
	switch name {
	case "length":
		return NewIntValue(len(a.entries)), nil
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
	args := make([]Value, len(a.entries))
	for i, zval := range a.entries {
		args[i] = zval.Value
	}
	return args
}

// IntArrayKeyName 将整数键编码为 ZVal.Name（稀疏整数键，插入顺序在 List 末尾）
func IntArrayKeyName(i int) string {
	return strconv.Itoa(i)
}

// ParseIntArrayKeyName 若 name 为纯整数字符串则返回该整数键，否则 ok=false
//
// 语义等价于「strconv.Atoi 成功 且 strconv.Itoa(n) == name」（即 name 必须是该整数的
// 规范十进制写法：无前导 0、无 + 号、"-0" 不算），但不用 strconv：
// 热路径上绝大多数键是**非数字**字符串（"id"、"name"、类名…），Atoi 会为每次失败
// 分配一个 *NumError 并格式化错误串，成功路径上的 Itoa 还会再分配一个字符串。
// 堆剖析里这条路径占全部分配的 12.8%，所以改成零分配的字节扫描，
// 顺便也吃下 "id"、"name" 这类键第一个字节就返回的短路。
func ParseIntArrayKeyName(name string) (int, bool) {
	if name == "" {
		return 0, false
	}
	i := 0
	neg := false
	if name[0] == '-' {
		neg = true
		i = 1
		if len(name) == 1 {
			return 0, false
		}
	}
	// 规范写法下前导位出现 0 只允许 name 恰好是 "0"："007" 的 Itoa 结果是 "7"
	if name[i] == '0' && len(name)-i > 1 {
		return 0, false
	}
	// 负数下界比正数上界多一（minInt 的绝对值），所以按 uint64 累积
	limit := uint64(math.MaxInt)
	if neg {
		limit++
	}
	var u uint64
	for ; i < len(name); i++ {
		c := name[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		d := uint64(c - '0')
		if u > (limit-d)/10 {
			return 0, false // 溢出：对齐 strconv 的 ErrRange → ok=false
		}
		u = u*10 + d
	}
	if neg {
		if u == 0 {
			return 0, false // "-0"：Itoa(0) 得到 "0"，不是规范写法
		}
		// u 可能等于 MaxInt+1，int(u) 按补码回绕成 minInt，再取负仍是 minInt：正确值
		return -int(u), true
	}
	return int(u), true
}

func (a *FlatArrayStore) invalidateIndex() {
	a.keyIndex = nil
	a.idxLen = -1
	a.packed = false
}

func (a *FlatArrayStore) ensureIndex() {
	if a.idxLen == len(a.entries) && (a.packed || a.keyIndex != nil) {
		return
	}
	a.rebuildIndex()
}

func (a *FlatArrayStore) rebuildIndex() {
	n := len(a.entries)
	a.idxLen = n
	packed := true
	// idx 延迟分配：纯 packed 数组（$a[] = / 列表字面量，热路径上最常见）没有任何
	// 命名键，原先无条件 make(map, n) 出来的表当场就被丢弃 —— 占全部分配的 6.89%。
	var idx map[string]int
	max := math.MinInt
	hasInt := false
	for i, z := range a.entries {
		if z == nil {
			continue
		}
		if z.EmptyStrKey {
			packed = false
			if idx == nil {
				idx = make(map[string]int, n)
			}
			idx[""] = i
			continue
		}
		if z.Name != "" {
			packed = false
			if idx == nil {
				idx = make(map[string]int, n)
			}
			idx[z.Name] = i
			if k, ok := ParseIntArrayKeyName(z.Name); ok {
				hasInt = true
				if k > max {
					max = k
				}
			}
			continue
		}
		hasInt = true
		if i > max {
			max = i
		}
	}
	a.packed = packed
	if !a.appendKeyKnown {
		a.nextIntKey = 0
		a.intKeySeen = false
		a.appendKeyKnown = true
	}
	if hasInt {
		a.advanceAppendKey(max)
	}
	// 没有命名键时 idx 保持 nil，等价于原来的 a.keyIndex = nil
	a.keyIndex = idx
}

// advanceAppendKey follows PHP 8.3+: the first integer key may be negative.
func (a *FlatArrayStore) advanceAppendKey(key int) {
	if !a.intKeySeen || key >= a.nextIntKey {
		a.nextIntKey = key
		if key < math.MaxInt {
			a.nextIntKey++
		}
	}
	a.intKeySeen = true
}

// NextAppendIntKey includes previously used integer keys, even after unset.
func (a *FlatArrayStore) NextAppendIntKey() int {
	a.ensureIndex()
	return a.nextIntKey
}

func (a *FlatArrayStore) AppendValue(value Value) bool {
	i := a.NextAppendIntKey()
	if i == math.MaxInt {
		if slot, _ := a.FindSlotByIntKey(i); slot != nil {
			return false
		}
	}
	a.SetIntKey(i, value)
	return true
}

func (a *FlatArrayStore) AppendSlot(value Value) *ZVal {
	i := a.NextAppendIntKey()
	if !a.AppendValue(value) {
		return nil
	}
	z, _ := a.FindSlotByIntKey(i)
	return z
}

// LookupZValByStringKey 按字符串键查找槽位；纯数字字符串键会回退整数键查找（如 "0" → 列表下标 0）
func (a *FlatArrayStore) LookupZValByStringKey(key string) (*ZVal, bool) {
	a.ensureIndex()
	if a.keyIndex != nil {
		if i, ok := a.keyIndex[key]; ok && i >= 0 && i < len(a.entries) {
			z := a.entries[i]
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
func (a *FlatArrayStore) SetStringKey(key string, value Value) {
	if n, ok := ParseIntArrayKeyName(key); ok {
		a.SetIntKey(n, value)
		return
	}
	if z, ok := a.LookupZValByStringKey(key); ok {
		z.Value = value
		return
	}
	// 新增键一律追加到末尾，索引可增量维护（旧实现整表失效，下次读要重建整张索引）
	if key == "" {
		a.appendSlotIncremental(NewEmptyStringKeyZVal(value))
		return
	}
	a.appendSlotIncremental(NewNamedZVal(key, value))
}

// FindSlotByIntKey 按 PHP 整数键查找槽位（含稀疏键 Name=="6" 等）
func (a *FlatArrayStore) FindSlotByIntKey(i int) (*ZVal, int) {
	a.ensureIndex()
	if a.packed {
		if i >= 0 && i < len(a.entries) {
			if z := a.entries[i]; z.IsPackedIntSlot() {
				return z, i
			}
		}
		return nil, -1
	}
	keyStr := IntArrayKeyName(i)
	if a.keyIndex != nil {
		if j, ok := a.keyIndex[keyStr]; ok && j >= 0 && j < len(a.entries) {
			z := a.entries[j]
			if z != nil && z.Name == keyStr {
				return z, j
			}
		}
	}
	if i >= 0 && i < len(a.entries) {
		if z := a.entries[i]; z.IsPackedIntSlot() {
			return z, i
		}
	}
	return nil, -1
}

// SetIntKey 设置整数键（不将稀疏数组转为 ObjectValue）
func (a *FlatArrayStore) SetIntKey(i int, value Value) {
	a.ensureIndex()
	if z, _ := a.FindSlotByIntKey(i); z != nil {
		z.Value = value
		return
	}
	packed, n := a.packed, len(a.entries)
	// 未找到整数键 i 时只能追加。禁止用 List[i] 覆盖：该槽可能是 PHP 空字符串键 ''。
	if packed && i == n {
		// 纯追加（$a[] = ... 循环里最常见）：索引仍然有效，推进计数即可。
		// 旧实现在这里整表失效，于是下一次读要重建整张索引 —— O(n) 的重复劳动。
		a.appendSlotIncremental(NewZVal(value))
		return
	}
	a.appendSlotIncremental(NewNamedZVal(IntArrayKeyName(i), value))
}

// appendSlotIncremental 追加一个「纯追加」槽位并增量维护索引。
// 前提：调用前索引是最新的（ensureIndex 已跑过），且本次只追加、不改动已有槽位。
// 不满足前提时（索引已脏）退回整表失效，交给下次重建 —— 宁可慢，也不能留半张索引。
func (a *FlatArrayStore) appendSlotIncremental(z *ZVal) {
	n := len(a.entries)
	if n == 0 {
		a.iterator = 0
	}
	a.entries = append(a.entries, z)
	if z.EmptyStrKey || z.Name != "" {
		wasPacked := a.packed
		a.packed = false
		if a.keyIndex == nil {
			if !wasPacked {
				a.invalidateIndex()
				return
			}
			// packed 数组的槽位按下标查、不进键表，跨到非 packed 时建新表即可
			a.keyIndex = make(map[string]int, n+1)
		}
		if z.EmptyStrKey {
			a.keyIndex[""] = n
		} else {
			a.keyIndex[z.Name] = n
			if k, ok := ParseIntArrayKeyName(z.Name); ok {
				a.advanceAppendKey(k)
			}
		}
	} else {
		a.advanceAppendKey(n)
	}
	a.idxLen = n + 1
}

// normalizeDenseIntKeys 将 Name=="" 的连续槽位转为显式整数字符串键，避免 unset 中间元素时误压缩后续键
func (a *FlatArrayStore) normalizeDenseIntKeys() {
	for j, z := range a.entries {
		if z != nil && z.IsPackedIntSlot() {
			z.Name = IntArrayKeyName(j)
		}
	}
	a.invalidateIndex()
}

// UnsetKey preserves PHP integer keys, insertion order, and automatic-key state.
func (a *FlatArrayStore) UnsetKey(index Value) {
	var slot *ZVal
	switch key := index.(type) {
	case *NullValue:
		slot, _ = a.LookupZValByStringKey("")
	case *StringValue:
		slot, _ = a.LookupZValByStringKey(key.AsString())
	case AsInt:
		if i, err := key.AsInt(); err == nil {
			slot, _ = a.FindSlotByIntKey(i)
		}
	}
	if slot == nil {
		return
	}
	// Compact only after a successful lookup. Explicit keys prevent other
	// integer entries from changing identity when a string entry is removed.
	for position, entry := range a.entries {
		if entry == slot {
			a.normalizeDenseIntKeys()
			copy(a.entries[position:], a.entries[position+1:])
			a.entries[len(a.entries)-1] = nil
			a.entries = a.entries[:len(a.entries)-1]
			if position < a.iterator {
				a.iterator--
			}
			if a.iterator >= len(a.entries) {
				a.iterator = -1
			}
			return
		}
	}
}

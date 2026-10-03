package data

// ZVal owns bucket metadata. Only explicit references allocate a shared cell;
// copying or renaming a bucket never changes another array's key. The ordinary
// slot remains 48 bytes on amd64 and stores its value inline.
type ZVal struct {
	Name         string
	InitialValue Value // Construction only; execution uses ReadValue/StoreRaw.
	reference    *ReferenceCell
	// Defined 表示该槽是否已赋值。CreateContext 预留的符号表槽为 false，
	// 与 PHP「仅引用未赋值不算已存在」对齐（EXTR_SKIP / isset 等）。
	Defined bool
	// EmptyStrKey 为 true 时 Name=="" 表示 PHP 数组键 ''，而不是 packed 整数键。
	// Filament NavigationManager::groupBy('') 依赖二者区分。
	EmptyStrKey bool
}

type ReferenceCell struct {
	value Value
	guard *ReferenceGuard
	count int32
}

func (z *ZVal) ReadValue() Value {
	if z.reference != nil {
		return z.reference.value
	}
	return z.InitialValue
}

// StoreRaw is for writes whose PHP checks are performed by the caller.
func (z *ZVal) StoreRaw(value Value) {
	if z.reference != nil {
		z.reference.value = value
	} else {
		z.InitialValue = value
	}
}

func (z *ZVal) ensureReference() *ReferenceCell {
	if z.reference == nil {
		z.reference = &ReferenceCell{value: z.InitialValue}
		z.InitialValue = nil
	}
	return z.reference
}

func (z *ZVal) Guard() *ReferenceGuard {
	if z.reference == nil {
		return nil
	}
	return z.reference.guard
}

func (z *ZVal) setGuard(guard *ReferenceGuard) { z.ensureReference().guard = guard }

func (z *ZVal) RefCount() int32 {
	if z.reference == nil {
		return 0
	}
	return z.reference.count
}

func (z *ZVal) ReleaseRefSlot() {
	if z != nil && z.reference != nil && z.reference.count > 0 {
		z.reference.count--
	}
}

func (z *ZVal) ReferenceIdentity() *ReferenceCell { return z.reference }

func CopyReferenceBucket(z *ZVal) *ZVal {
	copy := *z
	return &copy
}

// BindContextReference owns one reference binding independently of the source
// bucket. Rebinding releases the prior local owner without changing its value.
func BindContextReference(ctx Context, index int, source *ZVal) {
	if source == nil {
		return
	}
	previous := ctx.GetIndexZVal(index)
	if previous == source {
		return
	}
	source.AddRefSlot()
	bucket := CopyReferenceBucket(source)
	if previous != nil {
		bucket.Name = previous.Name
		previous.ReleaseRefSlot()
	}
	ctx.SetIndexZVal(index, bucket)
}

// MarkDefined 标记槽位已赋值（写入 Value 后调用）。
func (z *ZVal) MarkDefined() {
	if z != nil {
		z.Defined = true
	}
}

// AddRefSlot 标记该数组槽位被引用绑定（如 $x =& $arr[0]）
func (z *ZVal) AddRefSlot() {
	if z != nil {
		z.ensureReference().count++
	}
}

// NewZVal 创建一个新的 ZVal（视为已赋值）
func NewZVal(v Value) *ZVal {
	return &ZVal{
		InitialValue: v,
		Defined:      true,
	}
}

// NewNamedZVal 创建一个带名称的 ZVal（视为已赋值）
func NewNamedZVal(name string, v Value) *ZVal {
	return &ZVal{
		Name:         name,
		InitialValue: v,
		Defined:      true,
	}
}

// NewNamedZValSlot 创建仅占位的命名槽（未赋值，供函数/闭包符号表预分配）
func NewNamedZValSlot(name string) *ZVal {
	return &ZVal{
		Name:         name,
		InitialValue: NewNullValue(),
		Defined:      false,
	}
}

// NewEmptyStringKeyZVal 创建 PHP 空字符串键 ” 的数组槽。
func NewEmptyStringKeyZVal(v Value) *ZVal {
	return &ZVal{
		InitialValue: v,
		Defined:      true,
		EmptyStrKey:  true,
	}
}

// CopyZValKeepName 复制槽位的键身份（含空字符串键），替换 Value。
func CopyZValKeepName(z *ZVal, value Value) *ZVal {
	if z == nil {
		return NewZVal(value)
	}
	return &ZVal{
		Name:         z.Name,
		InitialValue: value,
		Defined:      true,
		EmptyStrKey:  z.EmptyStrKey,
	}
}

// IsPackedIntSlot 是否为 packed 整数键（Name 空且不是 PHP ”）。
func (z *ZVal) IsPackedIntSlot() bool {
	return z != nil && z.Name == "" && !z.EmptyStrKey
}

// PHPArrayKey 返回该槽对应的 PHP 数组键。
func (z *ZVal) PHPArrayKey(slot int) Value {
	if z == nil {
		return NewIntValue(slot)
	}
	if z.EmptyStrKey {
		return NewStringValue("")
	}
	if z.Name != "" {
		if n, ok := ParseIntArrayKeyName(z.Name); ok {
			return NewIntValue(n)
		}
		return NewStringValue(z.Name)
	}
	return NewIntValue(slot)
}

// SameArrayKey compares key identity without allocating temporary PHP values.
func (z *ZVal) SameArrayKey(position int, other *ZVal, otherPosition int) bool {
	if z == nil || other == nil {
		return z == other && position == otherPosition
	}
	if z.EmptyStrKey != other.EmptyStrKey {
		return false
	}
	if z.EmptyStrKey {
		return true
	}
	if z.Name == other.Name {
		return z.Name != "" || position == otherPosition
	}
	if z.Name == "" {
		key, ok := ParseIntArrayKeyName(other.Name)
		return ok && key == position
	}
	if other.Name == "" {
		key, ok := ParseIntArrayKeyName(z.Name)
		return ok && key == otherPosition
	}
	return false
}

type ZValGetter interface {
	GetZVal(v Variable) *ZVal
}

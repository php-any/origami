package data

// ZVal 模仿 PHP 的 zval 结构
type ZVal struct {
	Name  string // 变量名称，用于 extract 等按名称操作
	Value Value
	// Defined 表示该槽是否已赋值。CreateContext 预留的符号表槽为 false，
	// 与 PHP「仅引用未赋值不算已存在」对齐（EXTR_SKIP / isset 等）。
	Defined bool
	// RefSlotCount 表示有多少变量通过 &$arr[i] 等方式绑定到该槽位（用于 COW / 写穿）
	RefSlotCount int
	// EmptyStrKey 为 true 时 Name=="" 表示 PHP 数组键 ''，而不是 packed 整数键。
	// Filament NavigationManager::groupBy('') 依赖二者区分。
	EmptyStrKey bool
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
		z.RefSlotCount++
	}
}

// NewZVal 创建一个新的 ZVal（视为已赋值）
func NewZVal(v Value) *ZVal {
	return &ZVal{
		Value:   v,
		Defined: true,
	}
}

// NewNamedZVal 创建一个带名称的 ZVal（视为已赋值）
func NewNamedZVal(name string, v Value) *ZVal {
	return &ZVal{
		Name:    name,
		Value:   v,
		Defined: true,
	}
}

// NewNamedZValSlot 创建仅占位的命名槽（未赋值，供函数/闭包符号表预分配）
func NewNamedZValSlot(name string) *ZVal {
	return &ZVal{
		Name:    name,
		Value:   NewNullValue(),
		Defined: false,
	}
}

// NewEmptyStringKeyZVal 创建 PHP 空字符串键 '' 的数组槽。
func NewEmptyStringKeyZVal(v Value) *ZVal {
	return &ZVal{
		Value:       v,
		Defined:     true,
		EmptyStrKey: true,
	}
}

// CopyZValKeepName 复制槽位的键身份（含空字符串键），替换 Value。
func CopyZValKeepName(z *ZVal, value Value) *ZVal {
	if z == nil {
		return NewZVal(value)
	}
	return &ZVal{
		Name:        z.Name,
		Value:       value,
		Defined:     true,
		EmptyStrKey: z.EmptyStrKey,
	}
}

// IsPackedIntSlot 是否为 packed 整数键（Name 空且不是 PHP ''）。
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

type ZValGetter interface {
	GetZVal(v Variable) *ZVal
}

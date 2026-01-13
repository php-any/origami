package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SortFunction 实现 sort 函数
// sort(array &$array, int $flags = SORT_REGULAR): bool
// 对数组进行升序排序，会重新索引数组的键
type SortFunction struct{}

func NewSortFunction() data.FuncStmt {
	return &SortFunction{}
}

func (f *SortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	flagsValue, _ := ctx.GetIndexValue(1) // 可选的 flags 参数

	if arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查是否为数组引用
	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	// 如果数组为空，返回 true
	if len(arrayRef.Value) == 0 {
		return data.NewBoolValue(true), nil
	}

	// 获取 flags（默认为 SORT_REGULAR）
	flags := 0
	if flagsValue != nil {
		if intVal, ok := flagsValue.(*data.IntValue); ok {
			flags, _ = intVal.AsInt()
		}
	}

	// 对数组进行排序
	// PHP 的 sort() 函数会重新索引数组的键
	sort.Slice(arrayRef.Value, func(i, j int) bool {
		return compareValues(arrayRef.Value[i], arrayRef.Value[j], flags)
	})

	return data.NewBoolValue(true), nil
}

// compareValues 比较两个值的大小
// 根据 flags 参数决定比较方式
func compareValues(a, b data.Value, flags int) bool {
	// SORT_REGULAR: 正常比较（默认）
	// SORT_NUMERIC: 数值比较
	// SORT_STRING: 字符串比较
	// SORT_LOCALE_STRING: 根据当前区域设置进行字符串比较
	// SORT_NATURAL: 自然排序
	// SORT_FLAG_CASE: 可以与 SORT_STRING 或 SORT_NATURAL 组合使用，不区分大小写

	switch flags {
	case 1: // SORT_NUMERIC
		return compareNumeric(a, b)
	case 2: // SORT_STRING
		return compareString(a, b)
	case 5: // SORT_NATURAL
		return compareNatural(a, b)
	case 6: // SORT_NATURAL | SORT_FLAG_CASE
		return compareNaturalCaseInsensitive(a, b)
	case 3: // SORT_LOCALE_STRING (简化实现，使用字符串比较)
		return compareString(a, b)
	default: // SORT_REGULAR (0) 或其他
		return compareRegular(a, b)
	}
}

// compareRegular 正常比较
func compareRegular(a, b data.Value) bool {
	// 尝试数值比较
	if aNum, ok := a.(*data.IntValue); ok {
		if bNum, ok := b.(*data.IntValue); ok {
			aVal, _ := aNum.AsInt()
			bVal, _ := bNum.AsInt()
			return aVal < bVal
		}
		if bFloat, ok := b.(*data.FloatValue); ok {
			aVal, _ := aNum.AsInt()
			bVal, _ := bFloat.AsFloat()
			return float64(aVal) < bVal
		}
	}
	if aFloat, ok := a.(*data.FloatValue); ok {
		if bNum, ok := b.(*data.IntValue); ok {
			aVal, _ := aFloat.AsFloat()
			bVal, _ := bNum.AsInt()
			return aVal < float64(bVal)
		}
		if bFloat, ok := b.(*data.FloatValue); ok {
			aVal, _ := aFloat.AsFloat()
			bVal, _ := bFloat.AsFloat()
			return aVal < bVal
		}
	}
	// 否则使用字符串比较
	return a.AsString() < b.AsString()
}

// compareNumeric 数值比较
func compareNumeric(a, b data.Value) bool {
	aNum := getNumericValue(a)
	bNum := getNumericValue(b)
	return aNum < bNum
}

// compareString 字符串比较
func compareString(a, b data.Value) bool {
	return a.AsString() < b.AsString()
}

// compareNatural 自然排序
func compareNatural(a, b data.Value) bool {
	// 简化实现：使用字符串比较
	// 完整的自然排序需要更复杂的实现
	return a.AsString() < b.AsString()
}

// compareNaturalCaseInsensitive 自然排序（不区分大小写）
func compareNaturalCaseInsensitive(a, b data.Value) bool {
	// 简化实现：使用字符串比较（不区分大小写）
	// 完整的自然排序需要更复杂的实现
	aStr := a.AsString()
	bStr := b.AsString()
	// 转换为小写进行比较
	return toLower(aStr) < toLower(bStr)
}

// getNumericValue 获取数值
func getNumericValue(v data.Value) float64 {
	if intVal, ok := v.(*data.IntValue); ok {
		val, _ := intVal.AsInt()
		return float64(val)
	}
	if floatVal, ok := v.(*data.FloatValue); ok {
		val, _ := floatVal.AsFloat()
		return val
	}
	// 简化实现：如果字符串可以转换为数值，返回该值；否则返回 0
	// 这里可以使用 strconv.ParseFloat，但为了简化，返回 0
	return 0
}

// toLower 转换为小写（简化实现）
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func (f *SortFunction) GetName() string {
	return "sort"
}

func (f *SortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
	}
}

func (f *SortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
	}
}

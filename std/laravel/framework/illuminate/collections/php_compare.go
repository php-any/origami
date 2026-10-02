package collections

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

// PHP 值比较语义，供 Collection 的 where / whereIn / unique / sortBy 等使用。
//
// 关系比较规则与 node/binary_compare.go 的 phpCompareValues 一致：
//   - 数字与数字、数字与数字字符串、两个数字字符串：按数值比较
//   - 两个非数字字符串：按字典序比较
//
// 这里单独实现而不复用 node 包，是为了不给语言运算符的热路径加固定税。

// phpCompareValues 返回 -1/0/1；ok=false 表示无法比较（调用方应视为不匹配）。
func phpCompareValues(a, b data.Value) (int, bool) {
	a = unwrapValue(a)
	b = unwrapValue(b)

	aNum, aHasNum := phpToNumber(a)
	bNum, bHasNum := phpToNumber(b)
	aPure := phpIsPureNumber(a)
	bPure := phpIsPureNumber(b)
	as, aStr := phpToString(a)
	bs, bStr := phpToString(b)

	switch {
	case aPure && bPure:
		return cmpFloat(aNum, bNum), true
	case aPure && bHasNum, bPure && aHasNum:
		return cmpFloat(aNum, bNum), true
	case aPure && bStr:
		return cmpFloat(aNum, 0), true
	case bPure && aStr:
		return cmpFloat(0, bNum), true
	case aStr && bStr && aHasNum && bHasNum:
		return cmpFloat(aNum, bNum), true
	case aStr && bStr:
		return strings.Compare(as, bs), true
	default:
		return 0, false
	}
}

// phpLooseEquals 对齐 PHP 的 ==（PHP 8 规则）。
func phpLooseEquals(a, b data.Value) bool {
	a = unwrapValue(a)
	b = unwrapValue(b)

	if isNull(a) || isNull(b) {
		if isNull(a) && isNull(b) {
			return true
		}
		other := a
		if isNull(a) {
			other = b
		}
		// null == false|0|0.0|'0'? PHP: '' 与 '0' 不成立；null == '' 成立、null == '0' 不成立
		switch v := other.(type) {
		case *data.BoolValue:
			return !v.Value
		case *data.IntValue:
			return v.Value == 0
		case *data.FloatValue:
			return v.Value == 0
		case *data.StringValue:
			return v.Value == ""
		case *data.ArrayValue:
			return v.Len() == 0
		default:
			return false
		}
	}

	// 布尔参与比较时，两侧都转布尔
	if isPHPBool(a) || isPHPBool(b) {
		return phpToBool(a) == phpToBool(b)
	}

	// 数字与数字/数字字符串按数值比较；数字与非数字字符串按字符串比较（PHP 8）
	na, aNum := phpToNumber(a)
	nb, bNum := phpToNumber(b)
	if aNum && bNum && (phpIsPureNumber(a) || phpIsPureNumber(b) || phpIsNumericString(a) || phpIsNumericString(b)) {
		return na == nb
	}

	return phpStringify(a) == phpStringify(b)
}

// phpStrictEquals 对齐 PHP 的 ===。数组按结构（键序 + 严格值）比较，对象按引用比较。
func phpStrictEquals(a, b data.Value) bool {
	return phpStrictEqualsDepth(a, b, 0)
}

func phpStrictEqualsDepth(a, b data.Value, depth int) bool {
	if depth > 32 {
		return false
	}
	a = unwrapValue(a)
	b = unwrapValue(b)

	switch av := a.(type) {
	case *data.NullValue:
		return isNull(b)
	case *data.BoolValue:
		bv, ok := b.(*data.BoolValue)
		return ok && av.Value == bv.Value
	case *data.IntValue:
		bv, ok := b.(*data.IntValue)
		return ok && av.Value == bv.Value
	case *data.FloatValue:
		bv, ok := b.(*data.FloatValue)
		return ok && av.Value == bv.Value
	case *data.StringValue:
		bv, ok := b.(*data.StringValue)
		return ok && av.Value == bv.Value
	case *data.ArrayValue:
		bv, ok := b.(*data.ArrayValue)
		if !ok || av.Len() != bv.Len() {
			return false
		}
		for arraySlots46, i := av.View(), 0; i < arraySlots46.Len(); i++ {
			az := arraySlots46.At(i)
			bz := bv.At(i)
			if az == nil || bz == nil {
				if az != bz {
					return false
				}
				continue
			}
			if az.Name != bz.Name || az.EmptyStrKey != bz.EmptyStrKey {
				return false
			}
			if !phpStrictEqualsDepth(az.Value, bz.Value, depth+1) {
				return false
			}
		}
		return true
	case *data.ObjectValue:
		bv, ok := b.(*data.ObjectValue)
		return ok && av == bv
	default:
		return a == b
	}
}

func cmpFloat(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func phpIsPureNumber(v data.Value) bool {
	switch v.(type) {
	case *data.IntValue, *data.FloatValue, *data.BoolValue:
		return true
	default:
		return false
	}
}

func phpIsNumericString(v data.Value) bool {
	if s, ok := v.(*data.StringValue); ok {
		_, numeric := phpNumericString(s.Value)
		return numeric
	}
	return false
}

func isPHPBool(v data.Value) bool {
	_, ok := v.(*data.BoolValue)
	return ok
}

func phpToBool(v data.Value) bool {
	v = unwrapValue(v)
	switch t := v.(type) {
	case nil:
		return false
	case *data.NullValue:
		return false
	case *data.BoolValue:
		return t.Value
	case *data.IntValue:
		return t.Value != 0
	case *data.FloatValue:
		return t.Value != 0
	case *data.StringValue:
		return t.Value != "" && t.Value != "0"
	case *data.ArrayValue:
		return t.Len() > 0
	case *data.ObjectValue:
		return true
	default:
		return kit.Truthy(v)
	}
}

func phpToNumber(v data.Value) (float64, bool) {
	switch n := v.(type) {
	case *data.IntValue:
		return float64(n.Value), true
	case *data.FloatValue:
		return n.Value, true
	case *data.BoolValue:
		if n.Value {
			return 1, true
		}
		return 0, true
	case *data.StringValue:
		return phpNumericString(n.Value)
	default:
		return 0, false
	}
}

func phpToString(v data.Value) (string, bool) {
	switch n := v.(type) {
	case *data.StringValue:
		return n.Value, true
	case *data.NullValue:
		return "", true
	default:
		return "", false
	}
}

func phpStringify(v data.Value) string {
	if v == nil || isNull(v) {
		return ""
	}
	return v.AsString()
}

// phpNumericString 判断是否为 PHP 数字字符串（允许前导空白，不允许尾随空白/杂讯）。
func phpNumericString(s string) (float64, bool) {
	i := 0
	for i < len(s) && isPHPSpace(s[i]) {
		i++
	}
	if i >= len(s) {
		return 0, false
	}
	rest := s[i:]
	for _, r := range rest {
		if r <= unicode.MaxASCII && isPHPSpace(byte(r)) {
			return 0, false
		}
	}
	f, err := strconv.ParseFloat(rest, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func isPHPSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	default:
		return false
	}
}

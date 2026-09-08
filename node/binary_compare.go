package node

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
)

// phpCompareValues 按 PHP 关系比较语义比较两个已求值操作数。
// 返回 -1 / 0 / 1；ok=false 表示无法比较（调用方应得到 false）。
//
// 关键规则（与 PHP 8 对齐，供 Livewire Content-Length 等数字字符串使用）：
//   - 数字与数字、数字与数字字符串、两个数字字符串：按数值比较
//   - 两个非数字字符串：按字典序比较
func phpCompareValues(ctx data.Context, lv, rv data.GetValue) (int, bool) {
	lv = unwrapValue(lv)
	rv = unwrapValue(rv)

	if lts, ok := getDateTimeTimestamp(ctx, lv); ok {
		if rts, ok2 := getDateTimeTimestamp(ctx, rv); ok2 {
			return cmpFloat(float64(lts), float64(rts)), true
		}
	}

	lNum, lHasNum := phpToNumber(lv)
	rNum, rHasNum := phpToNumber(rv)
	lPure := phpIsPureNumber(lv)
	rPure := phpIsPureNumber(rv)
	ls, lStr := phpToString(lv)
	rs, rStr := phpToString(rv)

	switch {
	case lPure && rPure:
		return cmpFloat(lNum, rNum), true
	case lPure && rHasNum, rPure && lHasNum:
		return cmpFloat(lNum, rNum), true
	case lPure && rStr:
		return cmpFloat(lNum, 0), true
	case rPure && lStr:
		return cmpFloat(0, rNum), true
	case lStr && rStr && lHasNum && rHasNum:
		return cmpFloat(lNum, rNum), true
	case lStr && rStr:
		return strings.Compare(ls, rs), true
	default:
		return 0, false
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

func phpIsPureNumber(v data.GetValue) bool {
	switch v.(type) {
	case *data.IntValue, *data.FloatValue, *data.BoolValue:
		return true
	default:
		return false
	}
}

func phpToNumber(v data.GetValue) (float64, bool) {
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

func phpToString(v data.GetValue) (string, bool) {
	switch n := v.(type) {
	case *data.StringValue:
		return n.Value, true
	case *data.NullValue:
		return "", true
	default:
		return "", false
	}
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

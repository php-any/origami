package core

import (
	"sort"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrtrFunction 实现 strtr 函数
type StrtrFunction struct{}

// strtrPair 是数组形式 strtr 的一条替换规则。
type strtrPair struct{ key, value string }

func NewStrtrFunction() data.FuncStmt {
	return &StrtrFunction{}
}

// strtrToString 尝试将任意值转换为字符串：
// - 对普通标量/数组等，使用 AsString
// - 对对象，若存在 __toString 方法，则优先调用该方法获取字符串结果
func strtrToString(ctx data.Context, v data.Value) (string, data.Control) {
	if v == nil {
		return "", nil
	}

	switch sv := v.(type) {
	case *data.StringValue:
		return sv.Value, nil
	case *data.ClassValue:
		result, accepted, ctl := data.ObjectToStringValue(sv, ctx)
		if ctl != nil {
			return "", ctl
		}
		if accepted {
			return result.AsString(), nil
		}
		return sv.AsString(), nil
	case *data.ThisValue:
		// ThisValue 内部持有 ClassValue，共享同样的 __toString 逻辑
		return strtrToString(ctx, sv.ClassValue)
	default:
		return v.AsString(), nil
	}
}

func (f *StrtrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strValue, _ := ctx.GetIndexValue(0)
	str, ctl := strtrToString(ctx, strValue)
	if ctl != nil {
		return nil, ctl
	}

	fromValue, _ := ctx.GetIndexValue(1)

	var pairs []strtrPair

	if v, ok := fromValue.(*data.ArrayValue); ok {
		// PHP 数组形式按「每个位置最长匹配、替换结果不再扫描」执行。原实现走
		// strings.NewReplacer：它会为整组键建字典树，再叠加 map + keys 切片 +
		// args 切片，每次调用都是十余次分配，而这里只需要一张小表。
		pairs = make([]strtrPair, 0, v.View().Len())
		for slots, i := v.View(), 0; i < slots.Len(); i++ {
			z := slots.At(i)
			if z == nil || z.ReadValue() == nil {
				continue
			}
			key := z.PHPArrayKey(i).AsString()
			if key == "" {
				// PHP 忽略空键：空键会让每个位置都命中。
				continue
			}
			pairs = append(pairs, strtrPair{key: key, value: z.ReadValue().AsString()})
		}
		// 等长的键在同一位置至多命中一个，因此只需按长度降序保证最长优先。
		sort.Slice(pairs, func(i, j int) bool { return len(pairs[i].key) > len(pairs[j].key) })
	}

	if pairs != nil {
		if len(pairs) == 0 {
			return data.NewStringValue(str), nil
		}
		matched := false
		var builder strings.Builder
		builder.Grow(len(str))
		for i := 0; i < len(str); {
			replaced := false
			for _, pair := range pairs {
				if strings.HasPrefix(str[i:], pair.key) {
					builder.WriteString(pair.value)
					i += len(pair.key)
					replaced = true
					matched = true
					break
				}
			}
			if !replaced {
				builder.WriteByte(str[i])
				i++
			}
		}
		if !matched {
			// 没有任何替换命中时复用入参，避免为新字符串值单独分配。
			return data.NewStringValue(str), nil
		}
		return data.NewStringValue(builder.String()), nil
	}

	// strtr(string $string, string $from, string $to): string
	from := fromValue.AsString()
	toValue, _ := ctx.GetIndexValue(2)

	if toValue == nil {
		// If toValue is missing, and we are here, it means fromValue was NOT an array/object.
		// So this is invalid usage of strtr(string, string) -> missing 3rd arg.
		// We should throw error or return empty string?
		// PHP Warning: strtr() expects at least 3 parameters, 2 given
		// But here we can just return str or throw.
		// Let's return str (no replacement) or better, throw error.
		// But for now, let's just return str to avoid crash.
		return data.NewStringValue(str), nil
	}

	to := toValue.AsString()

	l := len(from)
	if len(to) < l {
		l = len(to)
	}

	mapping := make(map[byte]byte)
	for i := 0; i < l; i++ {
		mapping[from[i]] = to[i]
	}

	var builder strings.Builder
	builder.Grow(len(str))
	for i := 0; i < len(str); i++ {
		if val, ok := mapping[str[i]]; ok {
			builder.WriteByte(val)
		} else {
			builder.WriteByte(str[i])
		}
	}

	return data.NewStringValue(builder.String()), nil
}

func (f *StrtrFunction) GetName() string {
	return "strtr"
}

var strtrFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
	node.NewParameter(nil, "from", 1, nil, nil),
	node.NewParameter(nil, "to", 2, node.NewNullLiteral(nil), nil),
}

func (f *StrtrFunction) GetParams() []data.GetValue {
	return strtrFunctionGetParams
}

var strtrFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	node.NewVariable(nil, "from", 1, data.NewBaseType("string|array")),
	node.NewVariable(nil, "to", 2, data.NewBaseType("string")),
}

func (f *StrtrFunction) GetVariables() []data.Variable {
	return strtrFunctionGetVariables
}

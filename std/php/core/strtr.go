package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// StrtrFunction 实现 strtr 函数
type StrtrFunction struct{}

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
		if m, ok := sv.GetMethod("__toString"); ok && m != nil {
			fnCtx := sv.CreateContext(m.GetVariables())
			fnCtx.SetCallArgs([]data.GetValue{})
			ret, ctl := m.Call(fnCtx)
			if ctl != nil {
				return "", ctl
			}
			if ret == nil {
				return "", nil
			}
			if s, ok := ret.(*data.StringValue); ok {
				return s.Value, nil
			}
			if val, ok := ret.(data.Value); ok {
				return val.AsString(), nil
			}
			return "", nil
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

	var pairs map[string]string

	switch v := fromValue.(type) {
	case *data.ObjectValue:
		pairs = make(map[string]string)
		for k, val := range v.GetProperties() {
			pairs[k] = val.AsString()
		}
	case *data.ArrayValue:
		pairs = make(map[string]string)
		valueList := v.ToValueList()
		for i, val := range valueList {
			pairs[fmt.Sprintf("%d", i)] = val.AsString()
		}
	}

	if pairs != nil {
		// strtr(string $string, array $replace_pairs): string
		keys := make([]string, 0, len(pairs))
		for k := range pairs {
			keys = append(keys, k)
		}

		// Sort keys by length descending to ensure longest match first
		sort.Slice(keys, func(i, j int) bool {
			return len(keys[i]) > len(keys[j])
		})

		var args []string
		for _, k := range keys {
			args = append(args, k, pairs[k])
		}
		replacer := strings.NewReplacer(args...)
		return data.NewStringValue(replacer.Replace(str)), nil
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

func (f *StrtrFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "from", 1, nil, nil),
		node.NewParameter(nil, "to", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *StrtrFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "from", 1, data.NewBaseType("string|array")),
		node.NewVariable(nil, "to", 2, data.NewBaseType("string")),
	}
}

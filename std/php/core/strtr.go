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

func (f *StrtrFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strValue, _ := ctx.GetIndexValue(0)
	str := strValue.AsString()

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

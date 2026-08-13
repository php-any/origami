package php

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IntvalFunction 实现 PHP intval(mixed $value, int $base = 10): int
type IntvalFunction struct{}

func NewIntvalFunction() data.FuncStmt {
	return &IntvalFunction{}
}

func (f *IntvalFunction) GetName() string { return "intval" }

func (f *IntvalFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "base", 1, node.NewIntLiteral(nil, "10"), data.NewBaseType("int")),
	}
}

func (f *IntvalFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
		node.NewVariable(nil, "base", 1, data.NewBaseType("int")),
	}
}

func (f *IntvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	base := 10
	if b, ok := ctx.GetIndexValue(1); ok && b != nil {
		if ai, ok := b.(interface{ AsInt() (int, error) }); ok {
			if n, err := ai.AsInt(); err == nil {
				base = n
			}
		}
	}
	return data.NewIntValue(phpIntval(v, base)), nil
}

func phpIntval(v data.Value, base int) int {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case *data.IntValue:
		return t.Value
	case *data.FloatValue:
		f, _ := t.AsFloat()
		return int(f)
	case *data.BoolValue:
		if t.Value {
			return 1
		}
		return 0
	case *data.NullValue:
		return 0
	case *data.ArrayValue:
		if len(t.List) == 0 {
			return 0
		}
		return 1
	case *data.ObjectValue:
		n := 0
		t.RangeProperties(func(string, data.Value) bool {
			n++
			return false
		})
		if n == 0 {
			return 0
		}
		return 1
	case *data.ClassValue:
		return 1
	case *data.StringValue:
		return parseIntvalString(t.Value, base)
	default:
		return parseIntvalString(v.AsString(), base)
	}
}

// parseIntvalString 对齐 PHP intval 对字符串的 strtol 风格解析。
func parseIntvalString(s string, base int) int {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	if s == "" {
		return 0
	}
	if base == 0 {
		lower := strings.ToLower(s)
		sign := ""
		rest := lower
		if lower[0] == '+' || lower[0] == '-' {
			sign = string(lower[0])
			rest = lower[1:]
		}
		switch {
		case strings.HasPrefix(rest, "0x"):
			base = 16
			s = sign + rest[2:]
		case strings.HasPrefix(rest, "0b"):
			base = 2
			s = sign + rest[2:]
		case strings.HasPrefix(rest, "0") && len(rest) > 1:
			base = 8
		default:
			base = 10
		}
	}
	if base < 2 || base > 36 {
		return 0
	}
	sign := 1
	i := 0
	if s[0] == '+' {
		i = 1
	} else if s[0] == '-' {
		sign = -1
		i = 1
	}
	start := i
	for i < len(s) {
		c := s[i]
		var dig int
		switch {
		case c >= '0' && c <= '9':
			dig = int(c - '0')
		case c >= 'a' && c <= 'z':
			dig = int(c-'a') + 10
		case c >= 'A' && c <= 'Z':
			dig = int(c-'A') + 10
		default:
			dig = -1
		}
		if dig < 0 || dig >= base {
			break
		}
		i++
	}
	if i == start {
		return 0
	}
	n, err := strconv.ParseInt(s[start:i], base, 64)
	if err != nil {
		u, uerr := strconv.ParseUint(s[start:i], base, 64)
		if uerr != nil {
			return 0
		}
		return sign * int(u)
	}
	return sign * int(n)
}

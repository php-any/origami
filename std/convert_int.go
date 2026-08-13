package std

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewIntFunction() data.FuncStmt { return &IntFunction{} }

type IntFunction struct{}

func (f *IntFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewIntValue(0), nil
	}

	switch tv := v.(type) {
	case *data.IntValue:
		return tv, nil
	case *data.FloatValue:
		f64, _ := tv.AsFloat()
		return data.NewIntValue(int(f64)), nil
	case *data.BoolValue:
		if tv.Value {
			return data.NewIntValue(1), nil
		}
		return data.NewIntValue(0), nil
	case *data.NullValue:
		return data.NewIntValue(0), nil
	case data.AsString:
		return data.NewIntValue(intFromString(tv.AsString())), nil
	case data.AsInt:
		if i, err := tv.AsInt(); err == nil {
			return data.NewIntValue(i), nil
		}
	case data.AsFloat:
		if f64, err := tv.AsFloat(); err == nil {
			return data.NewIntValue(int(f64)), nil
		}
	}

	if s, ok := v.(data.AsString); ok {
		return data.NewIntValue(intFromString(s.AsString())), nil
	}
	return data.NewIntValue(0), nil
}

// intFromString 对齐 PHP (int)/intval 对十进制字符串的 strtol 风格：
// "3.14" → 3，"3abc" → 3，非法 → 0。
func intFromString(s string) int {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	if s == "" {
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
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == start {
		return 0
	}
	n, err := strconv.ParseInt(s[start:i], 10, 64)
	if err != nil {
		return 0
	}
	return sign * int(n)
}

func (f *IntFunction) GetName() string { return "int" }

func (f *IntFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *IntFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}

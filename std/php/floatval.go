package php

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FloatvalFunction 实现 PHP floatval(mixed $value): float
type FloatvalFunction struct{}

func NewFloatvalFunction() data.FuncStmt {
	return &FloatvalFunction{}
}

func (f *FloatvalFunction) GetName() string { return "floatval" }

func (f *FloatvalFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *FloatvalFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

func (f *FloatvalFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return data.NewFloatValue(phpFloatval(v)), nil
}

// DoublevalFunction 是 floatval 的别名。
type DoublevalFunction struct {
	FloatvalFunction
}

func NewDoublevalFunction() data.FuncStmt {
	return &DoublevalFunction{}
}

func (f *DoublevalFunction) GetName() string { return "doubleval" }

func phpFloatval(v data.Value) float64 {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case *data.FloatValue:
		f, _ := t.AsFloat()
		return f
	case *data.IntValue:
		return float64(t.Value)
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
	case *data.ClassValue, *data.ObjectValue:
		return 1
	case *data.StringValue:
		return parseFloatvalString(t.Value)
	default:
		if af, ok := v.(data.AsFloat); ok {
			if n, err := af.AsFloat(); err == nil {
				return n
			}
		}
		return parseFloatvalString(v.AsString())
	}
}

// parseFloatvalString 对齐 PHP floatval / (float) 对字符串的 zend_strtod 前缀解析。
func parseFloatvalString(s string) float64 {
	s = strings.TrimLeftFunc(s, unicode.IsSpace)
	if s == "" {
		return 0
	}
	i := 0
	if s[0] == '+' || s[0] == '-' {
		i++
	}
	if i >= len(s) {
		return 0
	}
	sawDigit := false
	sawDot := false
	sawExp := false
	for i < len(s) {
		c := s[i]
		if c >= '0' && c <= '9' {
			sawDigit = true
			i++
			continue
		}
		if c == '.' && !sawDot && !sawExp {
			sawDot = true
			i++
			continue
		}
		if (c == 'e' || c == 'E') && sawDigit && !sawExp {
			sawExp = true
			i++
			if i < len(s) && (s[i] == '+' || s[i] == '-') {
				i++
			}
			continue
		}
		break
	}
	if !sawDigit {
		return 0
	}
	f, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return 0
	}
	return f
}

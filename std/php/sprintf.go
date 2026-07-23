package php

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SprintfFunction 实现 sprintf 函数
type SprintfFunction struct{}

func NewSprintfFunction() data.FuncStmt {
	return &SprintfFunction{}
}

func (f *SprintfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	formatValue, _ := ctx.GetIndexValue(0)
	format := formatValue.AsString()

	// Collect args
	args := []interface{}{}

	valuesValue, _ := ctx.GetIndexValue(1)
	if valuesValue != nil {
		if paramsArray, ok := valuesValue.(*data.ArrayValue); ok {
			valueList := paramsArray.ToValueList()
			for _, val := range valueList {
				// Convert data.Value to Go value for fmt.Sprintf

				// Simple conversion for common types
				if v, ok := val.(*data.IntValue); ok {
					args = append(args, v.Value)
				} else if v, ok := val.(*data.FloatValue); ok {
					args = append(args, v.Value)
				} else if v, ok := val.(*data.StringValue); ok {
					args = append(args, v.Value)
				} else if v, ok := val.(*data.BoolValue); ok {
					if v.Value {
						args = append(args, 1)
					} else {
						args = append(args, 0)
					}
				} else {
					args = append(args, val.AsString())
				}
			}
		}
	}

	// PHP sprintf 格式串与 Go fmt 略有差异：
	// - PHP 支持位置参数：%1$s、%2$d 等；Go 使用 %[1]s 这样的语法。
	// 这里做一次简单转换：%1$-20s -> %[1]-20s 等。
	goFormat := phpToGoFormat(format)
	args = coerceArgsForFormat(goFormat, args)

	// 残缺格式串（如 ProgressBar 误走到 sprintf('%', $x)）勿把 Go fmt 的 %!(NOVERB) 泄露给用户
	if hasIncompletePrintfVerb(goFormat) {
		return data.NewBoolValue(false), nil
	}

	result := fmt.Sprintf(goFormat, args...)
	if strings.Contains(result, "%!(NOVERB)") || strings.Contains(result, "%!(EXTRA") {
		return data.NewBoolValue(false), nil
	}

	return data.NewStringValue(result), nil
}

func (f *SprintfFunction) GetName() string {
	return "sprintf"
}

func (f *SprintfFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, nil),
		node.NewParameters(nil, "values", 1, nil, nil),
	}
}

func (f *SprintfFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "values", 1, data.NewBaseType("mixed")),
	}
}

// phpToGoFormat 将 PHP 的 %1$s / %2$-10s 风格占位符转换为 Go 的 %[1]s / %[2]-10s。
var phpPositionalRe = regexp.MustCompile(`%([0-9]+)\$([+#0\- ]*\d*(?:\.\d+)?[bcdeEfFgGosxX])`)

func phpToGoFormat(format string) string {
	// 将 %1$s / %2$-10s 转为 %[1]s / %[2]-10s
	return phpPositionalRe.ReplaceAllString(format, "%[$1]$2")
}

// goFmtSpecRe 匹配 Go fmt 占位符（含 phpToGoFormat 转换后的 %[n] 位置参数）。
var goFmtSpecRe = regexp.MustCompile(`%(?:%|(?:\[(\d+)\])?([+#0\- ]*)(\d*)(?:\.(\d+))?([bcdeEfFgGosxX]))`)

// coerceArgsForFormat 按各占位符类型单独转换实参，兼容 PHP 的弱类型 sprintf 行为。
func coerceArgsForFormat(format string, args []interface{}) []interface{} {
	if len(args) == 0 {
		return args
	}

	result := append([]interface{}(nil), args...)
	nextIndex := 0

	for _, m := range goFmtSpecRe.FindAllStringSubmatch(format, -1) {
		if m[0] == "%%" {
			continue
		}

		verb := m[5]
		argIdx := nextIndex
		if m[1] != "" {
			if n, err := strconv.Atoi(m[1]); err == nil && n > 0 {
				argIdx = n - 1
			}
		} else {
			nextIndex++
		}

		if argIdx < 0 || argIdx >= len(result) {
			continue
		}
		result[argIdx] = coerceArgForVerb(verb, result[argIdx])
	}

	return result
}

func coerceArgForVerb(verb string, arg interface{}) interface{} {
	switch verb {
	case "s":
		return coerceToString(arg)
	case "d", "b", "o", "x", "X":
		return coerceToInt(arg)
	case "f", "F", "e", "E", "g", "G":
		return coerceToFloat(arg)
	default:
		return arg
	}
}

func coerceToString(arg interface{}) interface{} {
	switch v := arg.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case bool:
		if v {
			return "1"
		}
		return ""
	default:
		return fmt.Sprint(arg)
	}
}

func coerceToInt(arg interface{}) interface{} {
	switch v := arg.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return arg
}

func coerceToFloat(arg interface{}) interface{} {
	switch v := arg.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case bool:
		if v {
			return 1.0
		}
		return 0.0
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return arg
}

// hasIncompletePrintfVerb 检测以孤立 % 结尾或非法残缺占位符（PHP 会报错，勿走 Go fmt）。
func hasIncompletePrintfVerb(format string) bool {
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		if i+1 >= len(format) {
			return true
		}
		if format[i+1] == '%' {
			i++
			continue
		}
		// 跳过可选的 %[n] flags width.precision
		j := i + 1
		if j < len(format) && format[j] == '[' {
			for j < len(format) && format[j] != ']' {
				j++
			}
			if j >= len(format) {
				return true
			}
			j++
		}
		for j < len(format) && strings.ContainsRune("+#0- ", rune(format[j])) {
			j++
		}
		for j < len(format) && format[j] >= '0' && format[j] <= '9' {
			j++
		}
		if j < len(format) && format[j] == '.' {
			j++
			for j < len(format) && format[j] >= '0' && format[j] <= '9' {
				j++
			}
		}
		if j >= len(format) {
			return true
		}
		verb := format[j]
		if !strings.ContainsRune("bcdeEfFgGosxXvtT", rune(verb)) {
			// 允许 Go 额外动词；若仍不是字母则视为残缺
			if (verb < 'a' || verb > 'z') && (verb < 'A' || verb > 'Z') {
				return true
			}
		}
		i = j
	}
	return false
}

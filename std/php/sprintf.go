package php

import (
	"fmt"
	"regexp"

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

	// 对于 %s 格式符，Go 不会自动将 int/bool 转为字符串
	// 需要将参数转为字符串以兼容 PHP 行为
	if needsStringConversion(goFormat) {
		for i, arg := range args {
			switch v := arg.(type) {
			case int:
				args[i] = fmt.Sprintf("%d", v)
			case int64:
				args[i] = fmt.Sprintf("%d", v)
			case float64:
				args[i] = fmt.Sprintf("%g", v)
			case bool:
				if v {
					args[i] = "1"
				} else {
					args[i] = ""
				}
			}
		}
	}

	result := fmt.Sprintf(goFormat, args...)

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

// needsStringConversion 检查格式串中是否包含 %s
// 如果有，Go 不会自动将 int/bool 转为字符串，需要手动转换
var stringSpecRe = regexp.MustCompile(`%[+#0\- ]*\d*(?:\.\d+)?s`)

func needsStringConversion(format string) bool {
	return stringSpecRe.MatchString(format)
}

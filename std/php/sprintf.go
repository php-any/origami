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
						args = append(args, 1) // PHP bool to int/string in sprintf?
					} else {
						args = append(args, 0) // or empty string?
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

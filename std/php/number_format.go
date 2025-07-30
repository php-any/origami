package php

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewNumberFormatFunction() data.FuncStmt {
	return &NumberFormatFunction{}
}

type NumberFormatFunction struct {
	data.Function
}

func (f *NumberFormatFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	params := f.GetParams()
	if len(params) == 0 {
		return data.NewStringValue("0"), nil
	}

	// 获取第一个参数（数字）
	numberParam := params[0]
	numberValue, _ := numberParam.GetValue(ctx)
	if numberValue == nil {
		return data.NewStringValue("0"), nil
	}

	// 获取数字值
	var number float64
	switch v := numberValue.(type) {
	case *data.IntValue:
		number = float64(v.Value)
	case *data.FloatValue:
		number = v.Value
	case *data.StringValue:
		if parsed, err := strconv.ParseFloat(v.Value, 64); err == nil {
			number = parsed
		} else {
			return data.NewStringValue("0"), nil
		}
	default:
		return data.NewStringValue("0"), nil
	}

	// 默认参数
	decimals := 0
	decimalSeparator := "."
	thousandsSeparator := ","

	// 处理可选参数
	if len(params) > 1 {
		if decParam := params[1]; decParam != nil {
			if decValue, _ := decParam.GetValue(ctx); decValue != nil {
				if decInt, ok := decValue.(*data.IntValue); ok {
					decimals = decInt.Value
				}
			}
		}
	}

	if len(params) > 2 {
		if decSepParam := params[2]; decSepParam != nil {
			if decSepValue, _ := decSepParam.GetValue(ctx); decSepValue != nil {
				if decSepStr, ok := decSepValue.(*data.StringValue); ok {
					decimalSeparator = decSepStr.Value
				}
			}
		}
	}

	if len(params) > 3 {
		if thouSepParam := params[3]; thouSepParam != nil {
			if thouSepValue, _ := thouSepParam.GetValue(ctx); thouSepValue != nil {
				if thouSepStr, ok := thouSepValue.(*data.StringValue); ok {
					thousandsSeparator = thouSepStr.Value
				}
			}
		}
	}

	// 格式化数字
	formatted := formatNumber(number, decimals, decimalSeparator, thousandsSeparator)
	return data.NewStringValue(formatted), nil
}

func (f *NumberFormatFunction) GetName() string {
	return "number_format"
}

func (f *NumberFormatFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "number", 0, nil, nil),
		node.NewParameters(nil, "decimals", 1, nil, nil),
		node.NewParameters(nil, "decimal_separator", 2, nil, nil),
		node.NewParameters(nil, "thousands_separator", 3, nil, nil),
	}
}

func (f *NumberFormatFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "number", 0, nil),
		node.NewVariable(nil, "decimals", 1, nil),
		node.NewVariable(nil, "decimal_separator", 2, nil),
		node.NewVariable(nil, "thousands_separator", 3, nil),
	}
}

// formatNumber 格式化数字
func formatNumber(number float64, decimals int, decimalSeparator, thousandsSeparator string) string {
	// 使用 fmt.Sprintf 格式化数字
	formatStr := fmt.Sprintf("%%.%df", decimals)
	formatted := fmt.Sprintf(formatStr, number)

	// 分割整数和小数部分
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// 处理负数
	isNegative := false
	if strings.HasPrefix(integerPart, "-") {
		isNegative = true
		integerPart = integerPart[1:]
	}

	// 添加千位分隔符
	if thousandsSeparator != "" {
		integerPart = addThousandsSeparator(integerPart, thousandsSeparator)
	}

	// 重新组合
	result := integerPart
	if decimalPart != "" {
		result += decimalSeparator + decimalPart
	}

	if isNegative {
		result = "-" + result
	}

	return result
}

// addThousandsSeparator 添加千位分隔符
func addThousandsSeparator(s, separator string) string {
	if len(s) <= 3 {
		return s
	}

	var result strings.Builder
	start := len(s) % 3
	if start == 0 {
		start = 3
	}

	result.WriteString(s[:start])

	for i := start; i < len(s); i += 3 {
		result.WriteString(separator)
		result.WriteString(s[i : i+3])
	}

	return result.String()
}

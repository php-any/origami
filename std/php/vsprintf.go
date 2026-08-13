package php

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// VsprintfFunction 实现 vsprintf 函数
// vsprintf(string $format, array $values): string|false
// 返回根据格式化字符串生成的字符串，参数以数组形式提供
type VsprintfFunction struct{}

func NewVsprintfFunction() data.FuncStmt {
	return &VsprintfFunction{}
}

func (f *VsprintfFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：格式化字符串
	formatValue, _ := ctx.GetIndexValue(0)
	if formatValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取第二个参数：参数数组
	argsValue, _ := ctx.GetIndexValue(1)
	if argsValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 将格式化字符串转换为字符串
	var format string
	if strVal, ok := formatValue.(*data.StringValue); ok {
		format = strVal.AsString()
	} else {
		format = formatValue.AsString()
	}

	// 将参数数组转换为 Go 的 []interface{}
	var args []interface{}
	if arrayVal, ok := argsValue.(*data.ArrayValue); ok {
		valueList := arrayVal.ToValueList()
		args = make([]interface{}, len(valueList))
		for i, val := range valueList {
			// 将值转换为 Go 类型
			switch v := val.(type) {
			case *data.StringValue:
				args[i] = v.AsString()
			case *data.IntValue:
				args[i] = v.Value
			case *data.FloatValue:
				args[i] = v.Value
			case *data.BoolValue:
				// 布尔值转换为整数（PHP 行为：true -> 1, false -> 0）
				if v.Value {
					args[i] = 1
				} else {
					args[i] = 0
				}
			case *data.NullValue:
				args[i] = nil
			default:
				// 其他类型转换为字符串
				args[i] = v.AsString()
			}
		}
	} else {
		return data.NewBoolValue(false), nil
	}

	// 使用 fmt.Sprintf 格式化字符串
	result := fmt.Sprintf(format, args...)

	return data.NewStringValue(result), nil
}

func (f *VsprintfFunction) GetName() string {
	return "vsprintf"
}

func (f *VsprintfFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "format", 0, nil, data.String{}),
		node.NewParameter(nil, "values", 1, nil, data.Arrays{}),
	}
}

func (f *VsprintfFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "format", 0, data.String{}),
		node.NewVariable(nil, "values", 1, data.Arrays{}),
	}
}

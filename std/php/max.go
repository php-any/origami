package php

import (
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MaxFunction 实现 PHP 内置函数 max
//
// 支持的用法子集：
//   - max(value1, value2, ...): 返回参数中最大的一个（按 PHP 的“数字优先，字符串其次”的转换规则近似处理：这里直接用 AsFloat 比较）
//   - max(array $values): 只传一个数组参数时，在数组元素中取最大值
//
// 对于空数组或未提供有效参数的情况，返回 null。
type MaxFunction struct{}

func NewMaxFunction() data.FuncStmt {
	return &MaxFunction{}
}

func (f *MaxFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 收集所有实参
	var args []data.Value
	for i := 0; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok {
			break
		}
		args = append(args, v)
	}

	if len(args) == 0 {
		return data.NewNullValue(), nil
	}

	// 如果只有一个参数且是数组，按数组模式处理
	if len(args) == 1 {
		if arr, ok := args[0].(*data.ArrayValue); ok {
			list := arr.ToValueList()
			if len(list) == 0 {
				return data.NewNullValue(), nil
			}
			maxVal := list[0]
			maxNum := toFloat(maxVal)
			for _, v := range list[1:] {
				if n := toFloat(v); n > maxNum {
					maxNum = n
					maxVal = v
				}
			}
			return maxVal, nil
		}
	}

	// 多参数模式：在参数中取最大值
	maxVal := args[0]
	maxNum := toFloat(maxVal)
	for _, v := range args[1:] {
		if n := toFloat(v); n > maxNum {
			maxNum = n
			maxVal = v
		}
	}
	return maxVal, nil
}

// toFloat 尝试将任意 Value 按 PHP 语义近似转换为 float，用于大小比较。
func toFloat(v data.Value) float64 {
	if v == nil {
		return 0
	}
	if asFloat, ok := v.(data.AsFloat); ok {
		if f, err := asFloat.AsFloat(); err == nil {
			return f
		}
	}
	// 退化为字符串再尝试解析，简化：无法解析时视为 0
	str := v.AsString()
	if n, err := strconv.ParseFloat(str, 64); err == nil {
		return n
	}
	return 0
}

func (f *MaxFunction) GetName() string {
	return "max"
}

func (f *MaxFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "values", 0, nil, nil),
	}
}

func (f *MaxFunction) GetVariables() []data.Variable {
	return []data.Variable{
		// values 可以是任意类型或数组，这里用 mixed 约束
		node.NewVariable(nil, "values", 0, data.NewBaseType("mixed")),
	}
}

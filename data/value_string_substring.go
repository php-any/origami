package data

import (
	"strconv"
)

type StringValueSubstring struct {
	source string
}

func (s *StringValueSubstring) Call(ctx Context) (GetValue, Control) {
	// 获取参数
	startParam, startOk := ctx.GetIndexValue(0)
	if !startOk {
		return NewStringValue(""), nil
	}

	// 将 start 参数转换为整数
	var start int
	switch v := startParam.(type) {
	case *IntValue:
		start = v.Value
	case *StringValue:
		if val, err := strconv.Atoi(v.AsString()); err == nil {
			start = val
		} else {
			start = 0
		}
	case *FloatValue:
		start = int(v.Value)
	case *BoolValue:
		if v.Value {
			start = 1
		} else {
			start = 0
		}
	case *NullValue:
		start = 0
	default:
		start = 0
	}

	// 处理负数索引
	if start < 0 {
		start = 0
	}
	if start >= len(s.source) {
		return NewStringValue(""), nil
	}

	// 检查是否有第二个参数（end）
	endParam, endOk := ctx.GetIndexValue(1)
	var end int
	if endOk {
		// 将 end 参数转换为整数
		switch v := endParam.(type) {
		case *IntValue:
			end = v.Value
		case *StringValue:
			if val, err := strconv.Atoi(v.AsString()); err == nil {
				end = val
			} else {
				end = len(s.source)
			}
		case *FloatValue:
			end = int(v.Value)
		case *BoolValue:
			if v.Value {
				end = 1
			} else {
				end = 0
			}
		case *NullValue:
			end = len(s.source)
		default:
			end = len(s.source)
		}

		// 处理负数索引
		if end < 0 {
			end = 0
		}
		if end > len(s.source) {
			end = len(s.source)
		}
	} else {
		// 如果没有第二个参数，end 为字符串长度
		end = len(s.source)
	}

	// 确保 start <= end
	if start > end {
		start, end = end, start
	}

	// 提取子字符串
	result := s.source[start:end]
	return NewStringValue(result), nil
}

func (s *StringValueSubstring) GetName() string {
	return "substring"
}

func (s *StringValueSubstring) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueSubstring) GetIsStatic() bool {
	return false
}

func (s *StringValueSubstring) GetParams() []GetValue {
	return []GetValue{
		NewParameter("start", 0),
		NewParameterDefault("end", 1, NewNullValue(), nil),
	}
}

func (s *StringValueSubstring) GetVariables() []Variable {
	return []Variable{
		NewVariable("start", 0, nil),
		NewVariable("end", 1, nil),
	}
}

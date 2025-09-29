package data

import (
	"strings"
)

type StringValueSplit struct {
	source string
}

func (s *StringValueSplit) Call(ctx Context) (GetValue, Control) {
	// 获取分隔符参数
	separatorParam, separatorOk := ctx.GetIndexValue(0)
	if !separatorOk {
		// 如果没有参数，使用空格作为默认分隔符
		parts := strings.Fields(s.source)
		values := make([]Value, len(parts))
		for i, part := range parts {
			values[i] = NewStringValue(part)
		}
		return NewArrayValue(values), nil
	}

	// 将分隔符参数转换为字符串
	var separator string
	switch v := separatorParam.(type) {
	case *StringValue:
		separator = v.AsString()
	case *IntValue:
		separator = v.AsString()
	case *FloatValue:
		separator = v.AsString()
	case *BoolValue:
		separator = v.AsString()
	case *NullValue:
		separator = ""
	case *ArrayValue:
		separator = v.AsString()
	default:
		// 对于其他类型，尝试使用 AsString 方法
		if strValue, ok := v.(AsString); ok {
			separator = strValue.AsString()
		} else {
			separator = ""
		}
	}

	// 使用分隔符分割字符串
	parts := strings.Split(s.source, separator)
	values := make([]Value, len(parts))
	for i, part := range parts {
		values[i] = NewStringValue(part)
	}

	return NewArrayValue(values), nil
}

func (s *StringValueSplit) GetName() string {
	return "split"
}

func (s *StringValueSplit) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueSplit) GetIsStatic() bool {
	return false
}

func (s *StringValueSplit) GetParams() []GetValue {
	return []GetValue{
		NewParameterDefault("separator", 0, NewNullValue(), nil),
	}
}

func (s *StringValueSplit) GetVariables() []Variable {
	return []Variable{
		NewVariable("separator", 0, nil),
	}
}

func (s *StringValueSplit) GetReturnType() Types {
	return Arrays{}
}

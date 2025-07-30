package data

import (
	"strings"
)

type StringValueIndexOf struct {
	source string
}

func (s *StringValueIndexOf) Call(ctx Context) (GetValue, Control) {
	// 获取参数
	searchParam, ok := ctx.GetIndexValue(0)
	if !ok {
		return NewIntValue(-1), nil
	}

	// 将参数转换为字符串
	var searchStr string
	switch v := searchParam.(type) {
	case *StringValue:
		searchStr = v.AsString()
	case *IntValue:
		searchStr = v.AsString()
	case *FloatValue:
		searchStr = v.AsString()
	case *BoolValue:
		searchStr = v.AsString()
	case *NullValue:
		searchStr = "null"
	case *ArrayValue:
		searchStr = v.AsString()
	default:
		// 对于其他类型，尝试使用 AsString 方法
		if strValue, ok := v.(AsString); ok {
			searchStr = strValue.AsString()
		} else {
			searchStr = "undefined"
		}
	}

	// 使用 strings.Index 查找子字符串的位置
	// strings.Index 在找不到时返回 -1，找到时返回索引位置
	index := strings.Index(s.source, searchStr)

	// 返回找到的索引（如果没找到返回 -1）
	return NewIntValue(index), nil
}

func (s *StringValueIndexOf) GetName() string {
	return "indexOf"
}

func (s *StringValueIndexOf) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueIndexOf) GetIsStatic() bool {
	return false
}

func (s *StringValueIndexOf) GetParams() []GetValue {
	return []GetValue{
		NewVariable("search", 0, nil),
	}
}

func (s *StringValueIndexOf) GetVariables() []Variable {
	return []Variable{
		NewVariable("search", 0, nil),
	}
}

package data

import "strings"

type StringValueReplace struct {
	source string
}

func (s *StringValueReplace) Call(ctx Context) (GetValue, Control) {
	// 获取参数
	searchParam, searchOk := ctx.GetIndexValue(0)
	if !searchOk {
		return NewStringValue(s.source), nil
	}

	replaceParam, replaceOk := ctx.GetIndexValue(1)
	if !replaceOk {
		return NewStringValue(s.source), nil
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

	var replaceStr string
	switch v := replaceParam.(type) {
	case *StringValue:
		replaceStr = v.AsString()
	case *IntValue:
		replaceStr = v.AsString()
	case *FloatValue:
		replaceStr = v.AsString()
	case *BoolValue:
		replaceStr = v.AsString()
	case *NullValue:
		replaceStr = "null"
	case *ArrayValue:
		replaceStr = v.AsString()
	default:
		// 对于其他类型，尝试使用 AsString 方法
		if strValue, ok := v.(AsString); ok {
			replaceStr = strValue.AsString()
		} else {
			replaceStr = "undefined"
		}
	}

	// 使用 strings.ReplaceAll 替换所有匹配的字符串
	result := strings.ReplaceAll(s.source, searchStr, replaceStr)
	return NewStringValue(result), nil
}

func (s *StringValueReplace) GetName() string {
	return "replace"
}

func (s *StringValueReplace) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueReplace) GetIsStatic() bool {
	return false
}

func (s *StringValueReplace) GetParams() []GetValue {
	return []GetValue{
		NewVariable("search", 0, nil),
		NewVariable("replace", 1, nil),
	}
}

func (s *StringValueReplace) GetVariables() []Variable {
	return []Variable{
		NewVariable("search", 0, nil),
		NewVariable("replace", 1, nil),
	}
}

func (s *StringValueReplace) GetReturnType() Types {
	return String{}
}

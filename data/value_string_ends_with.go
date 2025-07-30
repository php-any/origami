package data

import "strings"

type StringValueEndsWith struct {
	source string
}

func (s *StringValueEndsWith) Call(ctx Context) (GetValue, Control) {
	// 获取参数
	searchParam, searchOk := ctx.GetIndexValue(0)
	if !searchOk {
		return NewBoolValue(false), nil
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

	// 检查字符串是否以指定后缀结束
	result := strings.HasSuffix(s.source, searchStr)
	return NewBoolValue(result), nil
}

func (s *StringValueEndsWith) GetName() string {
	return "endsWith"
}

func (s *StringValueEndsWith) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueEndsWith) GetIsStatic() bool {
	return false
}

func (s *StringValueEndsWith) GetParams() []GetValue {
	return []GetValue{
		NewVariable("search", 0, nil),
	}
}

func (s *StringValueEndsWith) GetVariables() []Variable {
	return []Variable{
		NewVariable("search", 0, nil),
	}
}

package data

import "strings"

type StringValueToLowerCase struct {
	source string
}

func (s *StringValueToLowerCase) Call(ctx Context) (GetValue, Control) {
	// 将字符串转换为小写
	result := strings.ToLower(s.source)
	return NewStringValue(result), nil
}

func (s *StringValueToLowerCase) GetName() string {
	return "toLowerCase"
}

func (s *StringValueToLowerCase) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueToLowerCase) GetIsStatic() bool {
	return false
}

func (s *StringValueToLowerCase) GetParams() []GetValue {
	return []GetValue{}
}

func (s *StringValueToLowerCase) GetVariables() []Variable {
	return []Variable{}
}

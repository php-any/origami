package data

import "strings"

type StringValueToUpperCase struct {
	source string
}

func (s *StringValueToUpperCase) Call(ctx Context) (GetValue, Control) {
	// 将字符串转换为大写
	result := strings.ToUpper(s.source)
	return NewStringValue(result), nil
}

func (s *StringValueToUpperCase) GetName() string {
	return "toUpperCase"
}

func (s *StringValueToUpperCase) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueToUpperCase) GetIsStatic() bool {
	return false
}

func (s *StringValueToUpperCase) GetParams() []GetValue {
	return []GetValue{}
}

func (s *StringValueToUpperCase) GetVariables() []Variable {
	return []Variable{}
}

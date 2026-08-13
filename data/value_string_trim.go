package data

import "strings"

type StringValueTrim struct {
	source string
}

func (s *StringValueTrim) Call(ctx Context) (GetValue, Control) {
	// 去除字符串首尾的空白字符
	result := strings.TrimSpace(s.source)
	return NewStringValue(result), nil
}

func (s *StringValueTrim) GetName() string {
	return "trim"
}

func (s *StringValueTrim) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueTrim) GetIsStatic() bool {
	return false
}

func (s *StringValueTrim) GetParams() []GetValue {
	return []GetValue{}
}

func (s *StringValueTrim) GetVariables() []Variable {
	return []Variable{}
}

func (s *StringValueTrim) GetReturnType() Types {
	return String{}
}

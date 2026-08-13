package data

type StringValueLength struct {
	source string
}

func (s *StringValueLength) Call(ctx Context) (GetValue, Control) {
	// 返回字符串的长度
	return NewIntValue(len(s.source)), nil
}

func (s *StringValueLength) GetName() string {
	return "length"
}

func (s *StringValueLength) GetModifier() Modifier {
	return ModifierPublic
}

func (s *StringValueLength) GetIsStatic() bool {
	return false
}

func (s *StringValueLength) GetParams() []GetValue {
	return []GetValue{}
}

func (s *StringValueLength) GetVariables() []Variable {
	return []Variable{}
}

func (s *StringValueLength) GetReturnType() Types {
	return Int{}
}

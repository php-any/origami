package data

import "strconv"

func NewStringValue(s string) Value {
	return &StringValue{Value: s}
}

type AsString interface {
	AsString() string
}

type StringValue struct {
	Value string
}

func (s *StringValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *StringValue) AsString() string {
	return s.Value
}

func (s *StringValue) AsInt() (int64, error) {
	return strconv.ParseInt(s.Value, 10, 64)
}

func (s *StringValue) AsFloat() (float64, error) {
	return strconv.ParseFloat(s.Value, 64)
}

func (s *StringValue) GetMethod(name string) (Method, bool) {
	switch name {
	case "indexOf":
		return &StringValueIndexOf{s.Value}, true
	case "substring":
		return &StringValueSubstring{s.Value}, true
	case "length":
		return &StringValueLength{s.Value}, true
	case "toLowerCase":
		return &StringValueToLowerCase{s.Value}, true
	case "toUpperCase":
		return &StringValueToUpperCase{s.Value}, true
	case "trim":
		return &StringValueTrim{s.Value}, true
	case "replace":
		return &StringValueReplace{s.Value}, true
	case "split":
		return &StringValueSplit{s.Value}, true
	case "startsWith":
		return &StringValueStartsWith{s.Value}, true
	case "endsWith":
		return &StringValueEndsWith{s.Value}, true
	}

	return nil, false
}

func (s *StringValue) GetProperty(name string) (Value, bool) {
	switch name {
	case "length":
		return NewIntValue(len(s.Value)), true
	}
	return nil, false
}

func (s *StringValue) AsBool() (bool, error) {
	return s.Value != "", nil
}

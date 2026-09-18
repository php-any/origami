package data

import (
	"fmt"
	"strconv"
)

func NewStringValue(s string) Value {
	return &StringValue{Value: s}
}

// NewByteStringValue 把单个字节做成 PHP 字节串（长度 1）。
// 禁止 string(byte)：Go 会把 >=0x80 的字节当成 rune 再编成 UTF-8（变成 2 字节）。
func NewByteStringValue(b byte) Value {
	return &StringValue{Value: string([]byte{b})}
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

func (s *StringValue) AsInt() (int, error) {
	n, err := strconv.ParseInt(s.Value, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
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

func (s *StringValue) GetProperty(name string) (Value, Control) {
	switch name {
	case "length":
		return NewIntValue(len(s.Value)), nil
	}
	return nil, NewErrorThrow(nil, fmt.Errorf("StringValue.GetProperty called with name %s", name))
}

func (s *StringValue) AsBool() (bool, error) {
	return s.Value != "" && s.Value != "0", nil
}

func (s *StringValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalString(s)
}

func (s *StringValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalString(data, s)
}

func (s *StringValue) ToGoValue(_ Serializer) (any, error) {
	return s.Value, nil
}

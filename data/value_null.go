package data

import (
	"fmt"
)

func NewNullValue() Value {
	return &NullValue{}
}

type AsNull interface {
	AsInt() (int, error)
}

type NullValue struct {
	Value
}

func (s *NullValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *NullValue) AsString() string {
	return fmt.Sprint("null")
}

func (s *NullValue) AsInt() (int, error) {
	return 0, nil
}

func (s *NullValue) AsFloat() (float64, error) {
	return 0, nil
}

func (s *NullValue) AsBool() (bool, error) {
	return false, nil
}

func (s *NullValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalNull(s)
}
func (s *NullValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalNull(data, s)
}

func (s *NullValue) ToGoValue(_ Serializer) (any, error) {
	return nil, nil
}

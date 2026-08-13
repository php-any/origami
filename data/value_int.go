package data

import (
	"fmt"
)

func NewIntValue(v int) Value {
	return &IntValue{
		Value: v,
	}
}

type AsInt interface {
	Value
	AsInt() (int, error)
}

type IntValue struct {
	Value int
}

func (s *IntValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *IntValue) AsString() string {
	return fmt.Sprintf("%d", s.Value)
}

func (s *IntValue) AsInt() (int, error) {
	return s.Value, nil
}

func (s *IntValue) AsFloat() (float64, error) {
	return float64(s.Value), nil
}

func (s *IntValue) AsBool() (bool, error) {
	return s.Value > 0, nil
}

func (s *IntValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalInt(s)
}
func (s *IntValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalInt(data, s)
}

func (s *IntValue) ToGoValue(_ Serializer) (any, error) {
	return s.Value, nil
}

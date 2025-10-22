package data

import (
	"fmt"
)

func NewBoolValue(v bool) Value {
	return &BoolValue{
		Value: v,
	}
}

type AsBool interface {
	AsBool() (bool, error)
}

type BoolValue struct {
	Value bool
}

func (s *BoolValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *BoolValue) AsString() string {
	if s.Value {
		return fmt.Sprint("true")
	}
	return fmt.Sprint("false")
}

func (s *BoolValue) AsBool() (bool, error) {
	return s.Value, nil
}

func (s *BoolValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalBool(s)
}

func (s *BoolValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalBool(data, s)
}

func (s *BoolValue) ToGoValue(_ Serializer) (any, error) {
	return s.Value, nil
}

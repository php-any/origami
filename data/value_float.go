package data

import (
	"fmt"
)

func NewFloatValue(v float64) Value {
	return &FloatValue{
		Value: v,
	}
}

type AsFloat interface {
	AsFloat() (float64, error)
}

type AsFloat32 interface {
	AsFloat32() (float32, error)
}

type AsFloat64 interface {
	AsFloat() (float64, error)
}

type FloatValue struct {
	Value float64
}

func (s *FloatValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *FloatValue) AsString() string {
	return fmt.Sprintf("%f", s.Value)
}

func (s *FloatValue) AsInt() (int, error) {
	return int(s.Value), nil
}

func (s *FloatValue) AsFloat() (float64, error) {
	return s.Value, nil
}
func (s *FloatValue) AsFloat32() (float32, error) {
	return float32(s.Value), nil
}

func (s *FloatValue) AsBool() (bool, error) {
	return s.Value > 0, nil
}

func (s *FloatValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalFloat(s)
}

func (s *FloatValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalFloat(data, s)
}

func (s *FloatValue) ToGoValue(_ Serializer) (any, error) {
	return s.Value, nil
}

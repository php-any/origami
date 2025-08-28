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

func (s *FloatValue) AsBool() (bool, error) {
	return s.Value > 0, nil
}

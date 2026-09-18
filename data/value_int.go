package data

import (
	"strconv"
)

var internInts [256]*IntValue

func init() {
	for i := 0; i < 256; i++ {
		internInts[i] = &IntValue{Value: i}
	}
}

func NewIntValue(v int) Value {
	if v >= 0 && v < len(internInts) {
		return internInts[v]
	}
	return &IntValue{Value: v}
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
	return strconv.Itoa(s.Value)
}

func (s *IntValue) AsInt() (int, error) {
	return s.Value, nil
}

func (s *IntValue) AsFloat() (float64, error) {
	return float64(s.Value), nil
}

func (s *IntValue) AsBool() (bool, error) {
	// PHP：任意非零整数（含负数）均为 true
	return s.Value != 0, nil
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

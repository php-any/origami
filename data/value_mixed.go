package data

import "fmt"

func NewMixedValue(v interface{}) Value {
	return &MixedValue{
		Value: v,
	}
}

type MixedValue struct {
	Value interface{}
}

func (m *MixedValue) GetValue(ctx Context) (GetValue, Control) {
	return m, nil
}

func (m *MixedValue) AsString() string {
	return fmt.Sprintf("%v", m.Value)
}

func (m *MixedValue) AsInt() (int, error) {
	return 0, nil
}

func (m *MixedValue) AsFloat() (float64, error) {
	return float64(0), nil
}

func (m *MixedValue) AsBool() (bool, error) {
	return false, nil
}

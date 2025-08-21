package data

import (
	"fmt"
)

func NewReferenceValue(v Value, ctx Context) Value {
	return &ReferenceValue{
		Value: v,
		Ctx:   ctx,
	}
}

type ReferenceValue struct {
	Value Value
	Ctx   Context
}

func (s *ReferenceValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *ReferenceValue) AsString() string {
	return fmt.Sprintf("%d", s.Value)
}

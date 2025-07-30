package data

import "fmt"

func NewAnyValue(v any) *AnyValue {
	return &AnyValue{
		Value: v,
	}
}

type AnyValue struct {
	Value any
}

func (c *AnyValue) GetValue(ctx Context) (GetValue, Control) {
	return nil, nil
}

func (c *AnyValue) AsString() string {
	return fmt.Sprintf("%v", c.Value)
}

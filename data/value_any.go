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

func (c *AnyValue) Marshal(serializer Serializer) ([]byte, error) {
	return serializer.MarshalAny(c)
}

func (c *AnyValue) Unmarshal(data []byte, serializer Serializer) error {
	return serializer.UnmarshalAny(data, c)
}

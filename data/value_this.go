package data

func NewThisValue(v *ClassValue) *ThisValue {
	return &ThisValue{
		ClassValue: v,
	}
}

type ThisValue struct {
	*ClassValue
}

func (c *ThisValue) GetValue(ctx Context) (GetValue, Control) {
	return c.ClassValue, nil
}

package data

func NewThisValue(v *ClassValue) *ThisValue {
	if v == nil {
		return &ThisValue{}
	}
	if cached := v.thisValue.Load(); cached != nil {
		return cached
	}
	value := &ThisValue{ClassValue: v}
	if v.thisValue.CompareAndSwap(nil, value) {
		return value
	}
	return v.thisValue.Load()
}

type ThisValue struct {
	*ClassValue
}

func (c *ThisValue) GetName() string {
	return c.ClassValue.GetName()
}

func (c *ThisValue) GetValue(ctx Context) (GetValue, Control) {
	return c.ClassValue, nil
}

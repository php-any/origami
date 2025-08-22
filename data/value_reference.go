package data

func NewReferenceValue(v Variable, ctx Context) Value {
	return &ReferenceValue{
		Val: v,
		Ctx: ctx,
	}
}

type ReferenceValue struct {
	Val Variable
	Ctx Context
}

func (s *ReferenceValue) GetValue(ctx Context) (GetValue, Control) {
	return s, nil
}

func (s *ReferenceValue) AsString() string {
	v, _ := s.Val.GetValue(s.Ctx)
	return v.(Value).AsString()
}

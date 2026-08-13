package data

func NewReturnControl(v Value) ReturnControl {
	return &ReturnValue{
		V: v,
	}
}

// ReturnControl 表示返回语句控制流
type ReturnControl interface {
	Control

	ReturnValue() Value
}

type ReturnValue struct {
	V Value
}

func (t *ReturnValue) GetValue(ctx Context) (GetValue, Control) {
	return t.V.GetValue(ctx)
}

func (t *ReturnValue) AsString() string {
	return "return " + t.V.AsString()
}

func (t *ReturnValue) ReturnValue() Value {
	return t.V
}

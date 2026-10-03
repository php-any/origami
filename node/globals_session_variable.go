package node

import "github.com/php-any/origami/data"

// $_SESSION

type SessionVariable struct {
	*Node `pp:"-"`
}

// sessionValue 仅作无 VM 作用域时的回退；HTTP 请求必须走 SuperglobalArrayProvider。
var sessionValue *data.ArrayValue

func NewSessionVariable(from data.From) data.Variable {
	return &SessionVariable{Node: NewNode(from)}
}

func (v *SessionVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_SESSION", func() *data.ArrayValue { return data.NewArrayValueFromSlots(nil) }), nil
}

func (v *SessionVariable) GetIndex() int       { return -1 }
func (v *SessionVariable) GetName() string     { return "$_SESSION" }
func (v *SessionVariable) GetType() data.Types { return nil }
func (v *SessionVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_SESSION", value)
}

func (v *SessionVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_SESSION"), nil
}
func (*SessionVariable) SuperglobalName() string { return "_SESSION" }

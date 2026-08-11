package node

import "github.com/php-any/origami/data"

// $_SESSION

type SessionVariable struct {
	*Node `pp:"-"`
}

// sessionValue 仅作无 VM 作用域时的回退；HTTP 请求必须走 SuperglobalArrayProvider。
var sessionValue *data.ObjectValue

func NewSessionVariable(from data.From) data.Variable {
	return &SessionVariable{Node: NewNode(from)}
}

func (v *SessionVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return sessionArrayFromContext(ctx), nil
}

func (v *SessionVariable) GetIndex() int       { return 0 }
func (v *SessionVariable) GetName() string     { return "$_SESSION" }
func (v *SessionVariable) GetType() data.Types { return nil }
func (v *SessionVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}

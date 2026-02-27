package node

import "github.com/php-any/origami/data"

// $_SESSION

type SessionVariable struct {
	*Node `pp:"-"`
}

var sessionValue *data.ObjectValue

func NewSessionVariable(from data.From) data.Variable {
	return &SessionVariable{Node: NewNode(from)}
}

func (v *SessionVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if sessionValue == nil {
		sessionValue = data.NewObjectValue()
	}
	return sessionValue, nil
}

func (v *SessionVariable) GetIndex() int       { return 0 }
func (v *SessionVariable) GetName() string     { return "$_SESSION" }
func (v *SessionVariable) GetType() data.Types { return nil }
func (v *SessionVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}

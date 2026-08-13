package node

import (
	"github.com/php-any/origami/data"
)

// $_ENV

type EnvVariable struct {
	*Node `pp:"-"`
}

var envValue *data.ObjectValue

func NewEnvVariable(from data.From) data.Variable {
	return &EnvVariable{Node: NewNode(from)}
}

func (v *EnvVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return ensureEnvValue(), nil
}

func (v *EnvVariable) GetIndex() int       { return 0 }
func (v *EnvVariable) GetName() string     { return "$_ENV" }
func (v *EnvVariable) GetType() data.Types { return nil }
func (v *EnvVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}

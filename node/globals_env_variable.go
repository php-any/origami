package node

import (
	"github.com/php-any/origami/data"
	"strings"
)

// $_ENV

type EnvVariable struct {
	*Node `pp:"-"`
}

func NewEnvVariable(from data.From) data.Variable {
	return &EnvVariable{Node: NewNode(from)}
}

func (v *EnvVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_ENV", func() *data.ArrayValue {
		array := data.NewArrayValueFromSlots(nil)
		for _, entry := range EnvironmentEntries(ctx) {
			if name, value, ok := strings.Cut(entry, "="); ok {
				array.SetStringKey(name, data.NewStringValue(value))
			}
		}
		return array
	}), nil
}

func (v *EnvVariable) GetIndex() int       { return -1 }
func (v *EnvVariable) GetName() string     { return "$_ENV" }
func (v *EnvVariable) GetType() data.Types { return nil }
func (v *EnvVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_ENV", value)
}

func (v *EnvVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_ENV"), nil
}
func (v *EnvVariable) SuperglobalName() string { return "_ENV" }

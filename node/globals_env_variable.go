package node

import (
	"os"
	"strings"

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
	if envValue == nil {
		envValue = data.NewObjectValue()
	}

	// 每次调用时合并一次环境变量（与之前逻辑一致）
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envValue.SetProperty(parts[0], data.NewStringValue(parts[1]))
		}
	}

	return envValue, nil
}

func (v *EnvVariable) GetIndex() int       { return 0 }
func (v *EnvVariable) GetName() string     { return "$_ENV" }
func (v *EnvVariable) GetType() data.Types { return nil }
func (v *EnvVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}

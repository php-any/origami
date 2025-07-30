package os

import (
	"github.com/php-any/origami/data"
)

type OSHostnameMethod struct {
	source *OS
}

func (h *OSHostnameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	name, err := h.source.Hostname()
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewStringValue(name), nil
}

func (h *OSHostnameMethod) GetName() string {
	return "hostname"
}

func (h *OSHostnameMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *OSHostnameMethod) GetIsStatic() bool {
	return false
}

func (h *OSHostnameMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *OSHostnameMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (h *OSHostnameMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

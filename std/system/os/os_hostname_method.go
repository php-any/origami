package os

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type OSHostnameMethod struct {
	source *OS
}

func (h *OSHostnameMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	name, err := h.source.Hostname()
	if err != nil {
		return nil, utils.NewThrow(err)
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

var oSHostnameMethodGetParams = []data.GetValue{}

func (h *OSHostnameMethod) GetParams() []data.GetValue {
	return oSHostnameMethodGetParams
}

var oSHostnameMethodGetVariables = []data.Variable{}

func (h *OSHostnameMethod) GetVariables() []data.Variable {
	return oSHostnameMethodGetVariables
}

// GetReturnType 返回方法返回类型
func (h *OSHostnameMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

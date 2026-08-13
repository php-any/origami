package exception

import (
	"github.com/php-any/origami/data"
)

type ExceptionGetTraceAsStringMethod struct {
	source *Exception
}

func (h *ExceptionGetTraceAsStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(h.source.GetTraceAsString()), nil
}

func (h *ExceptionGetTraceAsStringMethod) GetName() string {
	return "getTraceAsString"
}

func (h *ExceptionGetTraceAsStringMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ExceptionGetTraceAsStringMethod) GetIsStatic() bool {
	return false
}

func (h *ExceptionGetTraceAsStringMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ExceptionGetTraceAsStringMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (h *ExceptionGetTraceAsStringMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

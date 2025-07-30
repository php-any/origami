package exception

import (
	"github.com/php-any/origami/data"
)

type ExceptionErrorMethod struct {
	source *Exception
}

func (h *ExceptionErrorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(h.source.Error()), nil
}

func (h *ExceptionErrorMethod) GetName() string {
	return "error"
}

func (h *ExceptionErrorMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ExceptionErrorMethod) GetIsStatic() bool {
	return false
}

func (h *ExceptionErrorMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ExceptionErrorMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (h *ExceptionErrorMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

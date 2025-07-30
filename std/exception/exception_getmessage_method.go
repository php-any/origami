package exception

import (
	"github.com/php-any/origami/data"
)

type ExceptionGetMessageMethod struct {
	source *Exception
}

func (h *ExceptionGetMessageMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(h.source.GetMessage()), nil
}

func (h *ExceptionGetMessageMethod) GetName() string {
	return "getMessage"
}

func (h *ExceptionGetMessageMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *ExceptionGetMessageMethod) GetIsStatic() bool {
	return false
}

func (h *ExceptionGetMessageMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *ExceptionGetMessageMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (h *ExceptionGetMessageMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

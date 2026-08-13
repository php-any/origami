package exception

import (
	"github.com/php-any/origami/data"
)

// ExceptionGetCodeMethod 实现 Exception::getCode(): int
type ExceptionGetCodeMethod struct {
	source *Exception
}

func (h *ExceptionGetCodeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if code, ok := instancePropertyInt(ctx, "code"); ok {
		return data.NewIntValue(code), nil
	}
	return data.NewIntValue(0), nil
}

func (h *ExceptionGetCodeMethod) GetName() string { return "getCode" }

func (h *ExceptionGetCodeMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (h *ExceptionGetCodeMethod) GetIsStatic() bool { return false }

func (h *ExceptionGetCodeMethod) GetParams() []data.GetValue { return []data.GetValue{} }

func (h *ExceptionGetCodeMethod) GetVariables() []data.Variable { return []data.Variable{} }

func (h *ExceptionGetCodeMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

package exception

import (
	"github.com/php-any/origami/data"
)

// ExceptionGetPreviousMethod 实现 Exception::getPrevious(): ?Throwable
type ExceptionGetPreviousMethod struct {
	source *Exception
}

func (h *ExceptionGetPreviousMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if prev, ok := instancePropertyValue(ctx, "previous"); ok {
		return prev, nil
	}
	return data.NewNullValue(), nil
}

func (h *ExceptionGetPreviousMethod) GetName() string { return "getPrevious" }

func (h *ExceptionGetPreviousMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (h *ExceptionGetPreviousMethod) GetIsStatic() bool { return false }

func (h *ExceptionGetPreviousMethod) GetParams() []data.GetValue { return []data.GetValue{} }

func (h *ExceptionGetPreviousMethod) GetVariables() []data.Variable { return []data.Variable{} }

func (h *ExceptionGetPreviousMethod) GetReturnType() data.Types {
	return data.NewNullableType(data.NewBaseType("Throwable"))
}

func instancePropertyValue(ctx data.Context, name string) (data.Value, bool) {
	ov := instanceObjectValue(ctx)
	if ov == nil || !ov.HasProperty(name) {
		return nil, false
	}
	v, _ := ov.GetProperty(name)
	if v == nil {
		return nil, false
	}
	return v, true
}

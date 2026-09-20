package exception

import (
	"github.com/php-any/origami/data"
)

// ExceptionGetSeverityMethod 实现 ErrorException::getSeverity(): int
type ExceptionGetSeverityMethod struct {
	source *Exception
}

func (h *ExceptionGetSeverityMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if severity, ok := instancePropertyInt(ctx, "severity"); ok {
		return data.NewIntValue(severity), nil
	}
	return data.NewIntValue(1), nil
}

func (h *ExceptionGetSeverityMethod) GetName() string { return "getSeverity" }

func (h *ExceptionGetSeverityMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (h *ExceptionGetSeverityMethod) GetIsStatic() bool { return false }

func (h *ExceptionGetSeverityMethod) GetParams() []data.GetValue { return []data.GetValue{} }

func (h *ExceptionGetSeverityMethod) GetVariables() []data.Variable { return []data.Variable{} }

func (h *ExceptionGetSeverityMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

package exception

import (
	"github.com/php-any/origami/data"
)

type ExceptionGetFileMethod struct {
	source *Exception
}

func (h *ExceptionGetFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if file, ok := instancePropertyString(ctx, "file"); ok && file != "" {
		return data.NewStringValue(file), nil
	}
	return data.NewStringValue(h.source.GetFile()), nil
}

func (h *ExceptionGetFileMethod) GetName() string            { return "getFile" }
func (h *ExceptionGetFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ExceptionGetFileMethod) GetIsStatic() bool          { return false }
func (h *ExceptionGetFileMethod) GetParams() []data.GetValue { return []data.GetValue{} }
func (h *ExceptionGetFileMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *ExceptionGetFileMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

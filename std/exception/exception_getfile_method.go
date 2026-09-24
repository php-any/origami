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
	if h.source != nil {
		return data.NewStringValue(h.source.GetFile()), nil
	}
	return data.NewStringValue(""), nil
}

func (h *ExceptionGetFileMethod) GetName() string            { return "getFile" }
func (h *ExceptionGetFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ExceptionGetFileMethod) GetIsStatic() bool          { return false }
var exceptionGetFileMethodGetParams = []data.GetValue{}

func (h *ExceptionGetFileMethod) GetParams() []data.GetValue { return exceptionGetFileMethodGetParams }
var exceptionGetFileMethodGetVariables = []data.Variable{}

func (h *ExceptionGetFileMethod) GetVariables() []data.Variable {
	return exceptionGetFileMethodGetVariables
}
func (h *ExceptionGetFileMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}

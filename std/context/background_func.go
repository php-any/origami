package context

import (
	"context"
	"github.com/php-any/origami/data"
)

type BackgroundFunction struct{}

func NewBackgroundFunction() data.FuncStmt {
	return &BackgroundFunction{}
}

func (h *BackgroundFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := context.Background()
	return data.NewClassValue(NewContextClassFrom(ret0), ctx), nil
}

func (h *BackgroundFunction) GetName() string            { return "context\\background" }
func (h *BackgroundFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *BackgroundFunction) GetIsStatic() bool          { return true }
func (h *BackgroundFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (h *BackgroundFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (h *BackgroundFunction) GetReturnType() data.Types { return data.NewBaseType("void") }

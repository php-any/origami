package signal

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	signalsrc "os/signal"
)

type ResetFunction struct{}

func NewResetFunction() data.FuncStmt { return &ResetFunction{} }

func (f *ResetFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	sigs, ctl := parseSignalsFromContext(ctx, 0)
	if ctl != nil {
		return nil, ctl
	}
	signalsrc.Reset(sigs...)
	return nil, nil
}

func (f *ResetFunction) GetName() string            { return "Signal\\reset" }
func (f *ResetFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *ResetFunction) GetIsStatic() bool          { return true }
func (f *ResetFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "signals", 0, nil, data.NewBaseType("int")),
	}
}
func (f *ResetFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "signals", 0, data.NewBaseType("int")),
	}
}
func (f *ResetFunction) GetReturnType() data.Types { return data.NewBaseType("void") }

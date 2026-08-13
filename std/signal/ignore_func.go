package signal

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	signalsrc "os/signal"
)

type IgnoreFunction struct{}

func NewIgnoreFunction() data.FuncStmt { return &IgnoreFunction{} }

func (f *IgnoreFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	sigs, ctl := parseSignalsFromContext(ctx, 0)
	if ctl != nil {
		return nil, ctl
	}
	signalsrc.Ignore(sigs...)
	return nil, nil
}

func (f *IgnoreFunction) GetName() string            { return "Signal\\ignore" }
func (f *IgnoreFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *IgnoreFunction) GetIsStatic() bool          { return true }
func (f *IgnoreFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameters(nil, "signals", 0, nil, data.NewBaseType("int")),
	}
}
func (f *IgnoreFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "signals", 0, data.NewBaseType("int")),
	}
}
func (f *IgnoreFunction) GetReturnType() data.Types { return data.NewBaseType("void") }

package signal

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	signalsrc "os/signal"
)

type StopFunction struct{}

func NewStopFunction() data.FuncStmt { return &StopFunction{} }

func (f *StopFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	chVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: channel"))
	}
	ch, ctl := extractSignalChannel(chVal)
	if ctl != nil {
		return nil, ctl
	}
	signalsrc.Stop(ch.channel)
	return nil, nil
}

func (f *StopFunction) GetName() string            { return "Signal\\stop" }
func (f *StopFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *StopFunction) GetIsStatic() bool          { return true }
func (f *StopFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "channel", 0, nil, nil),
	}
}
func (f *StopFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "channel", 0, nil),
	}
}
func (f *StopFunction) GetReturnType() data.Types { return data.NewBaseType("void") }

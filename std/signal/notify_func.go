package signal

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	signalsrc "os/signal"
)

type NotifyFunction struct{}

func NewNotifyFunction() data.FuncStmt { return &NotifyFunction{} }

func (f *NotifyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	chVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: channel"))
	}
	ch, ctl := extractSignalChannel(chVal)
	if ctl != nil {
		return nil, ctl
	}
	sigs, ctl := parseSignalsFromContext(ctx, 1)
	if ctl != nil {
		return nil, ctl
	}
	if len(sigs) == 0 {
		return nil, utils.NewThrow(errors.New("至少需要一个信号"))
	}
	signalsrc.Notify(ch.channel, sigs...)
	return nil, nil
}

func (f *NotifyFunction) GetName() string            { return "Signal\\notify" }
func (f *NotifyFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *NotifyFunction) GetIsStatic() bool          { return true }
func (f *NotifyFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "channel", 0, nil, nil),
		node.NewParameters(nil, "signals", 1, nil, data.NewBaseType("int")),
	}
}
func (f *NotifyFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "channel", 0, nil),
		node.NewVariable(nil, "signals", 1, data.NewBaseType("int")),
	}
}
func (f *NotifyFunction) GetReturnType() data.Types { return data.NewBaseType("void") }

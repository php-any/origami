package signal

import (
	"errors"
	"os"
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	signalsrc "os/signal"
)

type WaitFunction struct{}

func NewWaitFunction() data.FuncStmt { return &WaitFunction{} }

func (f *WaitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	sigVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数: signals"))
	}
	sigs, ctl := parseSignalsFromArray(sigVal)
	if ctl != nil {
		return nil, ctl
	}
	if len(sigs) == 0 {
		return nil, utils.NewThrow(errors.New("信号列表不能为空"))
	}

	ch := make(chan os.Signal, 1)
	signalsrc.Notify(ch, sigs...)
	defer signalsrc.Stop(ch)

	sig := <-ch
	if s, ok := sig.(syscall.Signal); ok {
		return data.NewIntValue(int(s)), nil
	}
	return data.NewIntValue(0), nil
}

func (f *WaitFunction) GetName() string            { return "Signal\\wait" }
func (f *WaitFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (f *WaitFunction) GetIsStatic() bool          { return true }
func (f *WaitFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "signals", 0, nil, data.NewBaseType("array")),
	}
}
func (f *WaitFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "signals", 0, data.NewBaseType("array")),
	}
}
func (f *WaitFunction) GetReturnType() data.Types { return data.NewBaseType("int") }

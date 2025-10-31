package driver

import (
	driversrc "database/sql/driver"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ConnPrepareMethod struct {
	source driversrc.Conn
}

func (h *ConnPrepareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0 := a0.(*data.StringValue).AsString()

	ret0, err := h.source.Prepare(arg0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewStmtClassFrom(ret0), ctx), nil
}

func (h *ConnPrepareMethod) GetName() string            { return "prepare" }
func (h *ConnPrepareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnPrepareMethod) GetIsStatic() bool          { return true }
func (h *ConnPrepareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}

func (h *ConnPrepareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}

func (h *ConnPrepareMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

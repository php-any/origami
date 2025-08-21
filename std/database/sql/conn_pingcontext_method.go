package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ConnPingContextMethod struct {
	source *sqlsrc.Conn
}

func (h *ConnPingContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0 := a0.(*data.AnyValue).Value.(context.Context)

	if err := h.source.PingContext(arg0); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *ConnPingContextMethod) GetName() string            { return "pingContext" }
func (h *ConnPingContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnPingContextMethod) GetIsStatic() bool          { return true }
func (h *ConnPingContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
	}
}

func (h *ConnPingContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
	}
}

func (h *ConnPingContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

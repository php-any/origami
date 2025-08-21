package sql

import (
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ConnRawMethod struct {
	source *sqlsrc.Conn
}

func (h *ConnRawMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0 := a0.(*data.AnyValue).Value.(func(interface{}) error)

	if err := h.source.Raw(arg0); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *ConnRawMethod) GetName() string            { return "raw" }
func (h *ConnRawMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ConnRawMethod) GetIsStatic() bool          { return true }
func (h *ConnRawMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "f", 0, nil, nil),
	}
}

func (h *ConnRawMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "f", 0, nil),
	}
}

func (h *ConnRawMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

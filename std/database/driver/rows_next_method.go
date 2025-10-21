package driver

import (
	driversrc "database/sql/driver"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type RowsNextMethod struct {
	source driversrc.Rows
}

func (h *RowsNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	arg0 := make([]driversrc.Value, 0)
	for _, v := range a0.(*data.ArrayValue).Value {
		arg0 = append(arg0, v.(driversrc.Value))
	}

	if err := h.source.Next(arg0); err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return nil, nil
}

func (h *RowsNextMethod) GetName() string            { return "next" }
func (h *RowsNextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowsNextMethod) GetIsStatic() bool          { return true }
func (h *RowsNextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}

func (h *RowsNextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}

func (h *RowsNextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

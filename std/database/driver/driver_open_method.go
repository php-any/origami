package driver

import (
	driversrc "database/sql/driver"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DriverOpenMethod struct {
	source driversrc.Driver
}

func (h *DriverOpenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0 := a0.(*data.StringValue).AsString()

	ret0, err := h.source.Open(arg0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewConnClassFrom(ret0), ctx), nil
}

func (h *DriverOpenMethod) GetName() string            { return "open" }
func (h *DriverOpenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DriverOpenMethod) GetIsStatic() bool          { return true }
func (h *DriverOpenMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}

func (h *DriverOpenMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}

func (h *DriverOpenMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type TxQueryRowMethod struct {
	source *sqlsrc.Tx
}

func (h *TxQueryRowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 1"))
	}

	arg0 := a0.(*data.StringValue).AsString()
	arg1 := make([]any, 0)
	for _, v := range a1.(*data.ArrayValue).Value {
		arg1 = append(arg1, ConvertValueToGoType(v))
	}

	ret0 := h.source.QueryRow(arg0, arg1...)
	return data.NewClassValue(NewRowClassFrom(ret0), ctx), nil
}

func (h *TxQueryRowMethod) GetName() string            { return "queryRow" }
func (h *TxQueryRowMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxQueryRowMethod) GetIsStatic() bool          { return true }
func (h *TxQueryRowMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "query", 0, nil, nil),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (h *TxQueryRowMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "query", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *TxQueryRowMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

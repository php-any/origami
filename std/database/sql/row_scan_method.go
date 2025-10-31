package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RowScanMethod struct {
	source *sqlsrc.Row
}

func (h *RowScanMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	arg0 := make([]any, 0)
	for _, v := range a0.(*data.ArrayValue).Value {
		arg0 = append(arg0, ConvertValueToGoType(v))
	}

	if err := h.source.Scan(arg0...); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *RowScanMethod) GetName() string            { return "scan" }
func (h *RowScanMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RowScanMethod) GetIsStatic() bool          { return true }
func (h *RowScanMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParametersReference(nil, "dest", 0, nil, nil),
	}
}

func (h *RowScanMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "dest", 0, nil),
	}
}

func (h *RowScanMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

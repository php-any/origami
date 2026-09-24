package driver

import (
	driversrc "database/sql/driver"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type TxRollbackMethod struct {
	source driversrc.Tx
}

func (h *TxRollbackMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Rollback(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *TxRollbackMethod) GetName() string            { return "rollback" }
func (h *TxRollbackMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxRollbackMethod) GetIsStatic() bool          { return true }
var txRollbackMethodGetParams = []data.GetValue{}

func (h *TxRollbackMethod) GetParams() []data.GetValue {
	return txRollbackMethodGetParams
}

var txRollbackMethodGetVariables = []data.Variable{}

func (h *TxRollbackMethod) GetVariables() []data.Variable {
	return txRollbackMethodGetVariables
}

func (h *TxRollbackMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

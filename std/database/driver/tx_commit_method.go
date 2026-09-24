package driver

import (
	driversrc "database/sql/driver"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type TxCommitMethod struct {
	source driversrc.Tx
}

func (h *TxCommitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Commit(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *TxCommitMethod) GetName() string            { return "commit" }
func (h *TxCommitMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxCommitMethod) GetIsStatic() bool          { return true }
var txCommitMethodGetParams = []data.GetValue{}

func (h *TxCommitMethod) GetParams() []data.GetValue {
	return txCommitMethodGetParams
}

var txCommitMethodGetVariables = []data.Variable{}

func (h *TxCommitMethod) GetVariables() []data.Variable {
	return txCommitMethodGetVariables
}

func (h *TxCommitMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

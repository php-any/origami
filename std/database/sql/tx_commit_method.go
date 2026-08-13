package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type TxCommitMethod struct {
	source *sqlsrc.Tx
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
func (h *TxCommitMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *TxCommitMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *TxCommitMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type StmtCloseMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	if err := h.source.Close(); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *StmtCloseMethod) GetName() string            { return "close" }
func (h *StmtCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtCloseMethod) GetIsStatic() bool          { return true }
func (h *StmtCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *StmtCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *StmtCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

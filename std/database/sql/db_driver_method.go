package sql

import (
	sqlsrc "database/sql"
	"github.com/php-any/origami/data"
)

type DBDriverMethod struct {
	source *sqlsrc.DB
}

func (h *DBDriverMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.Driver()
	return data.NewAnyValue(ret0), nil
}

func (h *DBDriverMethod) GetName() string            { return "driver" }
func (h *DBDriverMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBDriverMethod) GetIsStatic() bool          { return true }
func (h *DBDriverMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (h *DBDriverMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (h *DBDriverMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

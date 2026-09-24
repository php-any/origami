package sql

import (
	sqlsrc "database/sql"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/database/driver"
)

type DBDriverMethod struct {
	source *sqlsrc.DB
}

func (h *DBDriverMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	ret0 := h.source.Driver()
	return data.NewClassValue(driver.NewDriverClassFrom(ret0), ctx), nil
}

func (h *DBDriverMethod) GetName() string            { return "driver" }
func (h *DBDriverMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBDriverMethod) GetIsStatic() bool          { return true }
var dBDriverMethodGetParams = []data.GetValue{}

func (h *DBDriverMethod) GetParams() []data.GetValue {
	return dBDriverMethodGetParams
}

var dBDriverMethodGetVariables = []data.Variable{}

func (h *DBDriverMethod) GetVariables() []data.Variable {
	return dBDriverMethodGetVariables
}

func (h *DBDriverMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

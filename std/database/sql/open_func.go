package sql

import (
	"database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type OpenFunction struct{}

func NewOpenFunction() data.FuncStmt {
	return &OpenFunction{}
}

func (h *OpenFunction) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	arg0 := a0.(*data.StringValue).AsString()
	arg1 := a1.(*data.StringValue).AsString()
	ret0, err := sql.Open(arg0, arg1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewDBClassFrom(ret0), ctx), nil
}

func (h *OpenFunction) GetName() string            { return "database\\sql\\open" }
func (h *OpenFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *OpenFunction) GetIsStatic() bool          { return true }
func (h *OpenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "driverName", 0, nil, nil),
		node.NewParameter(nil, "dataSourceName", 1, nil, nil),
	}
}
func (h *OpenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "driverName", 0, nil),
		node.NewVariable(nil, "dataSourceName", 1, nil),
	}
}
func (h *OpenFunction) GetReturnType() data.Types { return data.NewBaseType("void") }

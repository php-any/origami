package sql

import (
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DBExecMethod struct {
	source *sqlsrc.DB
}

func (h *DBExecMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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

	ret0, err := h.source.Exec(arg0, arg1...)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewResultClassFrom(ret0), ctx), nil
}

func (h *DBExecMethod) GetName() string            { return "exec" }
func (h *DBExecMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBExecMethod) GetIsStatic() bool          { return true }
func (h *DBExecMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "query", 0, nil, nil),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (h *DBExecMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "query", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *DBExecMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

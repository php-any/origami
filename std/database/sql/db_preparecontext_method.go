package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DBPrepareContextMethod struct {
	source *sqlsrc.DB
}

func (h *DBPrepareContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 1"))
	}

	var arg0 context.Context
	switch v := a0.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			arg0 = p.GetSource().(context.Context)
		} else {
			return nil, utils.NewThrow(errors.New("参数类型不支持, index: 0"))
		}
	case *data.AnyValue:
		arg0 = v.Value.(context.Context)
	default:
		return nil, utils.NewThrow(errors.New("参数类型不支持, index: 0"))
	}
	arg1 := a1.(*data.StringValue).AsString()

	ret0, err := h.source.PrepareContext(arg0, arg1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewStmtClassFrom(ret0), ctx), nil
}

func (h *DBPrepareContextMethod) GetName() string            { return "prepareContext" }
func (h *DBPrepareContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBPrepareContextMethod) GetIsStatic() bool          { return true }
func (h *DBPrepareContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "query", 1, nil, nil),
	}
}

func (h *DBPrepareContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "query", 1, nil),
	}
}

func (h *DBPrepareContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

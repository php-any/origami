package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type TxPrepareContextMethod struct {
	source *sqlsrc.Tx
}

func (h *TxPrepareContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	var arg0 context.Context
	switch v := a0.(type) {
	case *data.ClassValue:
		if p, ok := v.Class.(interface{ GetSource() any }); ok {
			arg0 = p.GetSource().(context.Context)
		} else {
			return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 0"))
		}
	case *data.AnyValue:
		arg0 = v.Value.(context.Context)
	default:
		return nil, data.NewErrorThrow(nil, errors.New("参数类型不支持, index: 0"))
	}
	arg1 := a1.(*data.StringValue).AsString()

	ret0, err := h.source.PrepareContext(arg0, arg1)
	if err != nil {
		return nil, data.NewErrorThrow(nil, err)
	}
	return data.NewClassValue(NewStmtClassFrom(ret0), ctx), nil
}

func (h *TxPrepareContextMethod) GetName() string            { return "prepareContext" }
func (h *TxPrepareContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxPrepareContextMethod) GetIsStatic() bool          { return true }
func (h *TxPrepareContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "query", 1, nil, nil),
	}
}

func (h *TxPrepareContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "query", 1, nil),
	}
}

func (h *TxPrepareContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

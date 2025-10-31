package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type TxStmtContextMethod struct {
	source *sqlsrc.Tx
}

func (h *TxStmtContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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
	arg1Class := a1.(*data.ClassValue).Class.(*StmtClass)
	arg1 := arg1Class.source

	ret0 := h.source.StmtContext(arg0, arg1)
	return data.NewClassValue(NewStmtClassFrom(ret0), ctx), nil
}

func (h *TxStmtContextMethod) GetName() string            { return "stmtContext" }
func (h *TxStmtContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *TxStmtContextMethod) GetIsStatic() bool          { return true }
func (h *TxStmtContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameter(nil, "stmt", 1, nil, nil),
	}
}

func (h *TxStmtContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "stmt", 1, nil),
	}
}

func (h *TxStmtContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type StmtQueryRowContextMethod struct {
	source *sqlsrc.Stmt
}

func (h *StmtQueryRowContextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

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
	arg1 := make([]any, 0)
	for _, v := range a1.(*data.ArrayValue).Value {
		arg1 = append(arg1, ConvertValueToGoType(v))
	}

	ret0 := h.source.QueryRowContext(arg0, arg1...)
	return data.NewClassValue(NewRowClassFrom(ret0), ctx), nil
}

func (h *StmtQueryRowContextMethod) GetName() string            { return "queryRowContext" }
func (h *StmtQueryRowContextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *StmtQueryRowContextMethod) GetIsStatic() bool          { return true }
func (h *StmtQueryRowContextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (h *StmtQueryRowContextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

func (h *StmtQueryRowContextMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

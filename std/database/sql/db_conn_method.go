package sql

import (
	"context"
	sqlsrc "database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type DBConnMethod struct {
	source *sqlsrc.DB
}

func (h *DBConnMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
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

	ret0, err := h.source.Conn(arg0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	return data.NewClassValue(NewConnClassFrom(ret0), ctx), nil
}

func (h *DBConnMethod) GetName() string            { return "conn" }
func (h *DBConnMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *DBConnMethod) GetIsStatic() bool          { return true }
func (h *DBConnMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "ctx", 0, nil, nil),
	}
}

func (h *DBConnMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "ctx", 0, nil),
	}
}

func (h *DBConnMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

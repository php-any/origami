package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ArrayFunction struct{}

func NewArrayFunction() data.FuncStmt { return &ArrayFunction{} }

func (f *ArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	a1, has := ctx.GetIndexValue(0)
	if !has {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	switch v := a1.(type) {
	case *data.ArrayValue:
		return v, nil
	case *data.StringValue, *data.IntValue, *data.FloatValue, *data.BoolValue, *data.NullValue:
		// PHP: (array) 标量 => array(0 => 标量)
		return data.NewArrayValue([]data.Value{a1}), nil
	default:
		return data.NewArrayValue([]data.Value{a1}), nil
	}
}

func (f *ArrayFunction) GetName() string { return "array" }

func (f *ArrayFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
	}
}

func (f *ArrayFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.Mixed{}),
	}
}

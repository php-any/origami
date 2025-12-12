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

	switch f := a1.(type) {
	case *data.ArrayValue:
		return f, nil
	}

	return nil, utils.NewThrow(errors.New("无法转化数组"))
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

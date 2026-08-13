package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DieFunction is an alias for exit/die
type DieFunction struct{}

func NewDieFunction() data.FuncStmt {
	return &DieFunction{}
}

func (f *DieFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return (&ExitFunction{}).Call(ctx)
}

func (f *DieFunction) GetName() string {
	return "die"
}

func (f *DieFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "status", 0, node.NewNullLiteral(nil), data.NewBaseType("mixed")),
	}
}

func (f *DieFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "status", 0, data.NewBaseType("mixed")),
	}
}

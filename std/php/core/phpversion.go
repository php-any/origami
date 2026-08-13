package core

import (
	"github.com/php-any/origami/data"
)

type PhpVersionFunction struct{}

func NewPhpVersionFunction() data.FuncStmt {
	return &PhpVersionFunction{}
}

func (f *PhpVersionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if v, ok := vm.GetConstant("PHP_VERSION"); ok {
		return v, nil
	}
	return data.NewStringValue("8.2.0"), nil
}

func (f *PhpVersionFunction) GetName() string { return "phpversion" }

func (f *PhpVersionFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *PhpVersionFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
)

// SplAutoloadFunctionsFunction 实现 spl_autoload_functions
type SplAutoloadFunctionsFunction struct{}

func NewSplAutoloadFunctionsFunction() data.FuncStmt {
	return &SplAutoloadFunctionsFunction{}
}

func (f *SplAutoloadFunctionsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	fns := runtime.GetAutoLoad()
	list := make([]*data.ZVal, len(fns))
	for i, fn := range fns {
		list[i] = data.NewZVal(fn)
	}
	return &data.ArrayValue{List: list}, nil
}

func (f *SplAutoloadFunctionsFunction) GetName() string { return "spl_autoload_functions" }

func (f *SplAutoloadFunctionsFunction) GetParams() []data.GetValue { return nil }

func (f *SplAutoloadFunctionsFunction) GetVariables() []data.Variable { return nil }

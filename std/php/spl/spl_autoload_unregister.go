package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

// SplAutoloadUnregisterFunction 实现 spl_autoload_unregister
type SplAutoloadUnregisterFunction struct{}

func NewSplAutoloadUnregisterFunction() data.FuncStmt { return &SplAutoloadUnregisterFunction{} }

func (f *SplAutoloadUnregisterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	callback, _ := ctx.GetIndexValue(0)
	if _, ctl := node.ResolveCallback(ctx, callback); ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(runtime.UnregisterAutoloadInContext(ctx, callback)), nil
}
func (f *SplAutoloadUnregisterFunction) GetName() string { return "spl_autoload_unregister" }

var splAutoloadUnregisterFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "callback", 0, nil, nil),
}

func (f *SplAutoloadUnregisterFunction) GetParams() []data.GetValue {
	return splAutoloadUnregisterFunctionGetParams
}

var splAutoloadUnregisterFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "callback", 0, data.Mixed{}),
}

func (f *SplAutoloadUnregisterFunction) GetVariables() []data.Variable {
	return splAutoloadUnregisterFunctionGetVariables
}

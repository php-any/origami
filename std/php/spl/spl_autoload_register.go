package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

// SplAutoloadRegisterFunction 实现 spl_autoload_register
type SplAutoloadRegisterFunction struct{}

func NewSplAutoloadRegisterFunction() data.FuncStmt { return &SplAutoloadRegisterFunction{} }

func (f *SplAutoloadRegisterFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	callback, _ := ctx.GetIndexValue(0)
	if callback == nil {
		callback = data.NewStringValue("spl_autoload")
	}
	if _, isNull := callback.(*data.NullValue); isNull {
		callback = data.NewStringValue("spl_autoload")
	}
	resolved, ctl := node.ResolveCallback(ctx, callback)
	if ctl != nil {
		return nil, ctl
	}
	prepend := false
	if value, ok := ctx.GetIndexValue(2); ok {
		if flag, ok := value.(*data.BoolValue); ok {
			prepend = flag.Value
		}
	}
	runtime.RegisterAutoloadInContext(ctx, callback, resolved, prepend)
	return data.NewBoolValue(true), nil
}
func (f *SplAutoloadRegisterFunction) GetName() string { return "spl_autoload_register" }

var splAutoloadRegisterFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "callback", 0, data.NewNullValue(), nil),
	node.NewParameter(nil, "throw", 1, data.NewBoolValue(true), data.TypeBool),
	node.NewParameter(nil, "prepend", 2, data.NewBoolValue(false), data.TypeBool),
}

func (f *SplAutoloadRegisterFunction) GetParams() []data.GetValue {
	return splAutoloadRegisterFunctionGetParams
}

var splAutoloadRegisterFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "callback", 0, data.Mixed{}),
	node.NewVariable(nil, "throw", 1, data.TypeBool),
	node.NewVariable(nil, "prepend", 2, data.TypeBool),
}

func (f *SplAutoloadRegisterFunction) GetVariables() []data.Variable {
	return splAutoloadRegisterFunctionGetVariables
}

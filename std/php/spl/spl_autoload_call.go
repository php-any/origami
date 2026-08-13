package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

// SplAutoloadCallFunction 实现 spl_autoload_call
type SplAutoloadCallFunction struct{}

func NewSplAutoloadCallFunction() data.FuncStmt {
	return &SplAutoloadCallFunction{}
}

func (f *SplAutoloadCallFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	classVal, ok := ctx.GetIndexValue(0)
	if !ok || classVal == nil {
		return data.NewBoolValue(false), nil
	}
	name := classVal.AsString()
	okLoaded, acl := runtime.CallAutoLoad(name, ctx)
	if acl != nil {
		return nil, acl
	}
	return data.NewBoolValue(okLoaded), nil
}

func (f *SplAutoloadCallFunction) GetName() string { return "spl_autoload_call" }

func (f *SplAutoloadCallFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "class_name", 0, nil, data.NewBaseType("string")),
	}
}

func (f *SplAutoloadCallFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "class_name", 0, data.NewBaseType("string")),
	}
}

package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsCallableFunction 实现 is_callable 函数
type IsCallableFunction struct{}

func NewIsCallableFunction() data.FuncStmt {
	return &IsCallableFunction{}
}

func (f *IsCallableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	syntaxValue, _ := ctx.GetIndexValue(1)
	syntaxOnly := false
	if flag, ok := syntaxValue.(*data.BoolValue); ok {
		syntaxOnly = flag.Value
	}
	name, syntaxValid := node.CallableName(value)
	if ctl := ctx.SetVariableValue(isCallableFunctionGetVariables[2], data.NewStringValue(name)); ctl != nil {
		return nil, ctl
	}
	if syntaxOnly {
		return data.NewBoolValue(syntaxValid), nil
	}
	if !syntaxValid {
		return data.NewBoolValue(false), nil
	}
	resolved, ctl := node.ResolveCallback(ctx, value)
	if ctl != nil {
		if thrown, ok := ctl.(*data.ThrowValue); !ok || thrown.Name != "TypeError" {
			return nil, ctl
		}
	}
	return data.NewBoolValue(resolved != nil && ctl == nil), nil
}
func (f *IsCallableFunction) GetName() string {
	return "is_callable"
}

var isCallableFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "value", 0, nil, nil),
	node.NewParameter(nil, "syntax_only", 1, node.NewBooleanLiteral(nil, false), nil),
	node.NewParameterReference(nil, "callable_name", 2, node.NewNullLiteral(nil), nil),
}

func (f *IsCallableFunction) GetParams() []data.GetValue {
	return isCallableFunctionGetParams
}

var isCallableFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	node.NewVariable(nil, "syntax_only", 1, data.NewBaseType("bool")),
	node.NewVariable(nil, "callable_name", 2, data.NewBaseType("string")),
}

func (f *IsCallableFunction) GetVariables() []data.Variable {
	return isCallableFunctionGetVariables
}

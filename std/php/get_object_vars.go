package php

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type GetObjectVarsFunction struct{}

func (*GetObjectVarsFunction) GetName() string { return "get_object_vars" }
func (*GetObjectVarsFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewCallerContextParameter(nil)}
}
func (*GetObjectVarsFunction) GetVariables() []data.Variable { return nil }
func (*GetObjectVarsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	args := ctx.GetCallArgs()
	if len(args) != 1 {
		return nil, data.NewTypeError(nil, fmt.Errorf("get_object_vars() expects exactly one object"))
	}
	value, ctl := args[0].GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	var object *data.ClassValue
	switch v := value.(type) {
	case *data.ClassValue:
		object = v
	case *data.ThisValue:
		object = v.ClassValue
	}
	if object == nil {
		return nil, data.NewTypeError(nil, fmt.Errorf("get_object_vars(): argument must be an object"))
	}
	return node.VisibleObjectVars(ctx, object), nil
}

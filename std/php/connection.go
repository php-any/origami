package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
	"net/http"
)

// Connection state belongs to the PHP request, including CLI's default normal
// state. A server execution deadline is independent of the ignore-abort flag.
type ConnectionFunction struct{ name string }

func (f *ConnectionFunction) GetName() string { return f.name }
func (f *ConnectionFunction) GetParams() []data.GetValue {
	if f.name == "ignore_user_abort" {
		return connectionIgnoreParams
	}
	return nil
}
func (f *ConnectionFunction) GetVariables() []data.Variable {
	if f.name == "ignore_user_abort" {
		return connectionIgnoreVariables
	}
	return nil
}

var connectionIgnoreParams = []data.GetValue{node.NewParameter(nil, "enable", 0, data.NewNullValue(), data.NewNullableType(data.Bool{}))}
var connectionIgnoreVariables = []data.Variable{node.NewVariable(nil, "enable", 0, data.NewNullableType(data.Bool{}))}

func (f *ConnectionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if f.name == "ignore_user_abort" {
		previous, _ := core.IniGetInContext(ctx, "ignore_user_abort")
		value := 0
		if previous == "1" || previous == "on" || previous == "true" {
			value = 1
		}
		if enable, ok := ctx.GetIndexValue(0); ok {
			if flag, ok := enable.(*data.BoolValue); ok {
				next := "0"
				if flag.Value {
					next = "1"
				}
				core.IniSetInContext(ctx, "ignore_user_abort", next)
			}
		}
		return data.NewIntValue(value), nil
	}
	status := 0
	if host, ok := ctx.GetVM().(interface{ HTTPRequest() *http.Request }); ok {
		if request := host.HTTPRequest(); request != nil && request.Context().Err() != nil {
			status = 1
		}
	}
	return data.NewIntValue(status), nil
}

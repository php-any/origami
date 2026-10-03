package php

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DateTimeCreateFromInterfaceMethod struct{}

func (*DateTimeCreateFromInterfaceMethod) GetName() string            { return "createFromInterface" }
func (*DateTimeCreateFromInterfaceMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (*DateTimeCreateFromInterfaceMethod) GetIsStatic() bool          { return true }
func (*DateTimeCreateFromInterfaceMethod) GetReturnType() data.Types  { return data.TypeStatic }

var dateTimeCreateFromInterfaceParams = []data.GetValue{node.NewParameter(nil, "object", 0, nil, data.NewBaseType("DateTimeInterface"))}
var dateTimeCreateFromInterfaceVariables = []data.Variable{node.NewVariable(nil, "object", 0, data.NewBaseType("DateTimeInterface"))}

func (*DateTimeCreateFromInterfaceMethod) GetParams() []data.GetValue {
	return dateTimeCreateFromInterfaceParams
}
func (*DateTimeCreateFromInterfaceMethod) GetVariables() []data.Variable {
	return dateTimeCreateFromInterfaceVariables
}
func (*DateTimeCreateFromInterfaceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	source, ok := toClassValue(value)
	if !ok {
		return nil, data.NewTypeError(nil, fmt.Errorf("DateTime::createFromInterface() expects DateTimeInterface"))
	}
	frame := ctx.(*data.ClassMethodContext)
	class := frame.StaticClass
	if class == nil {
		class = frame.Class
	}
	copy := data.NewClassValue(class, ctx.CreateBaseContext())
	for _, name := range []string{"timestamp", "timezone", "microsecond"} {
		if value, ctl := source.GetProperty(name); ctl == nil && value != nil {
			copy.SetProperty(name, value)
		}
	}
	return copy, nil
}

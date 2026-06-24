package compile

import (
	"reflect"

	"github.com/php-any/origami/data"
)

type dataValueEmitter func(g *Generator, v data.GetValue) error

var dataValueEmitters map[reflect.Type]dataValueEmitter

func init() {
	dataValueEmitters = map[reflect.Type]dataValueEmitter{
		reflect.TypeOf((*data.IntValue)(nil)):    emitIntValue,
		reflect.TypeOf((*data.StringValue)(nil)): emitStringValue,
		reflect.TypeOf((*data.BoolValue)(nil)):   emitBoolValue,
		reflect.TypeOf((*data.NullValue)(nil)):   emitNullValue,
		reflect.TypeOf((*data.FloatValue)(nil)):  emitFloatValue,
	}
}

func emitIntValue(g *Generator, v data.GetValue) error {
	iv := v.(*data.IntValue)
	g.printf("data.NewIntValue(%d)", iv.Value)
	return nil
}

func emitStringValue(g *Generator, v data.GetValue) error {
	sv := v.(*data.StringValue)
	g.printf("data.NewStringValue(%q)", sv.Value)
	return nil
}

func emitBoolValue(g *Generator, v data.GetValue) error {
	bv := v.(*data.BoolValue)
	g.printf("data.NewBoolValue(%v)", bv.Value)
	return nil
}

func emitNullValue(g *Generator, v data.GetValue) error {
	g.printf("data.NewNullValue()")
	return nil
}

func emitFloatValue(g *Generator, v data.GetValue) error {
	fv := v.(*data.FloatValue)
	g.printf("data.NewFloatValue(%g)", fv.Value)
	return nil
}

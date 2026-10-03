package php

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"strconv"
	"time"
)

func dateTimeWireLocation(zone string) (*time.Location, error) {
	if len(zone) == 6 && (zone[0] == '+' || zone[0] == '-') && zone[3] == ':' {
		h, e1 := strconv.Atoi(zone[1:3])
		m, e2 := strconv.Atoi(zone[4:6])
		if e1 != nil || e2 != nil || h > 23 || m > 59 {
			return nil, errors.New("Invalid timezone offset")
		}
		offset := h*3600 + m*60
		if zone[0] == '-' {
			offset = -offset
		}
		return time.FixedZone(zone, offset), nil
	}
	return time.LoadLocation(zone)
}

type DateTimeSerializeMethod struct{ DateTimeConstructMethod }

func (*DateTimeSerializeMethod) GetName() string               { return "__serialize" }
func (*DateTimeSerializeMethod) GetParams() []data.GetValue    { return nil }
func (*DateTimeSerializeMethod) GetVariables() []data.Variable { return nil }
func (*DateTimeSerializeMethod) GetReturnType() data.Types     { return data.Arrays{} }
func (*DateTimeSerializeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	t, ctl := getDateTime(ctx)
	if ctl != nil {
		return nil, ctl
	}
	object := ctx.(*data.ClassMethodContext).ClassValue
	zone := "UTC"
	if value, ctl := object.GetProperty("timezone"); ctl == nil && value != nil {
		if name, ok := value.(*data.StringValue); ok && name.Value != "" {
			zone = name.Value
		}
	}
	zoneType := 3
	if len(zone) > 0 && (zone[0] == '+' || zone[0] == '-') {
		zoneType = 1
	}
	result := data.NewArrayValueFromSlots(nil)
	result.SetStringKey("date", data.NewStringValue(t.Format("2006-01-02 15:04:05.000000")))
	result.SetStringKey("timezone_type", data.NewIntValue(zoneType))
	result.SetStringKey("timezone", data.NewStringValue(zone))
	object.RangeProperties(func(name string, value data.Value) bool {
		if name != "timestamp" && name != "microsecond" && name != "timezone" {
			result.SetStringKey(name, value)
		}
		return true
	})
	return result, nil
}

type DateTimeUnserializeMethod struct{ DateTimeConstructMethod }

func (*DateTimeUnserializeMethod) GetName() string { return "__unserialize" }
func (*DateTimeUnserializeMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "data", 0, nil, data.Arrays{})}
}
func (*DateTimeUnserializeMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "data", 0, data.Arrays{})}
}
func (*DateTimeUnserializeMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
func (*DateTimeUnserializeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	raw, _ := ctx.GetIndexValue(0)
	values, ok := raw.(*data.ArrayValue)
	if !ok {
		return nil, data.NewTypeError(nil, errors.New("DateTime::__unserialize(): data must be an array"))
	}
	date, hasDate := values.LookupZValByStringKey("date")
	zone, hasZone := values.LookupZValByStringKey("timezone")
	kind, hasKind := values.LookupZValByStringKey("timezone_type")
	invalid := func() (data.GetValue, data.Control) {
		return nil, data.NewErrorThrowByName(nil, errors.New("Invalid serialization data for DateTime object"), "Error")
	}
	if !hasDate || !hasZone || !hasKind {
		return invalid()
	}
	d, okDate := date.ReadValue().(*data.StringValue)
	z, okZone := zone.ReadValue().(*data.StringValue)
	k, okKind := kind.ReadValue().(*data.IntValue)
	if !okDate || !okZone || !okKind || k.Value < 1 || k.Value > 3 {
		return invalid()
	}
	loc, err := dateTimeWireLocation(z.Value)
	if err != nil {
		return invalid()
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05.000000", d.Value, loc)
	if err != nil {
		return invalid()
	}
	if ctl := setDateTime(ctx, t); ctl != nil {
		return nil, ctl
	}
	object := ctx.(*data.ClassMethodContext).ClassValue
	object.SetProperty("timezone", data.NewStringValue(z.Value))
	view := values.View()
	for i := 0; i < view.Len(); i++ {
		slot := view.At(i)
		name := slot.PHPArrayKey(i).AsString()
		if name != "date" && name != "timezone" && name != "timezone_type" {
			object.SetProperty(name, slot.ReadValue())
		}
	}
	return data.NewNullValue(), nil
}

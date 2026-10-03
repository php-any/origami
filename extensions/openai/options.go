package openai

import (
	"errors"
	"github.com/php-any/origami/data"
)

func clientOptions(ctx data.Context, index int) (map[string]any, data.Control) {
	value, found := ctx.GetIndexValue(index)
	if !found || value == nil {
		return nil, nil
	}
	if _, null := value.(*data.NullValue); null {
		return nil, nil
	}
	array, ok := value.(*data.ArrayValue)
	if !ok {
		return nil, data.NewTypeError(nil, errors.New("options must be an array"))
	}
	result := make(map[string]any, array.Len())
	view := array.View()
	for i := 0; i < view.Len(); i++ {
		slot := view.At(i)
		result[slot.PHPArrayKey(i).AsString()] = optionValue(slot.ReadValue())
	}
	return result, nil
}

func optionValue(value data.Value) any {
	if array, ok := value.(*data.ArrayValue); ok {
		view := array.View()
		list := make([]any, view.Len())
		object := make(map[string]any, view.Len())
		isList := true
		for i := 0; i < view.Len(); i++ {
			slot := view.At(i)
			key := slot.PHPArrayKey(i)
			if number, ok := key.(*data.IntValue); !ok || number.Value != i {
				isList = false
			}
			list[i] = optionValue(slot.ReadValue())
			object[key.AsString()] = list[i]
		}
		if isList {
			return list
		}
		return object
	}
	switch value := value.(type) {
	case *data.NullValue:
		return nil
	case *data.StringValue:
		return value.Value
	case *data.IntValue:
		return value.Value
	case *data.FloatValue:
		return value.Value
	case *data.BoolValue:
		return value.Value
	default:
		return value.AsString()
	}
}

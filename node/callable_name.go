package node

import "github.com/php-any/origami/data"

// CallableName performs syntax validation without loading a class or function.
func CallableName(value data.Value) (string, bool) {
	switch callback := value.(type) {
	case *data.StringValue:
		return callback.Value, true
	case *data.FuncValue, *data.BoundFuncValue:
		return "Closure::__invoke", true
	case *data.ThisValue:
		return callback.Class.GetName() + "::__invoke", true
	case *data.ClassValue:
		return callback.Class.GetName() + "::__invoke", true
	case *data.ArrayValue:
		if callback.Len() != 2 {
			return "Array", false
		}
		first, _ := callback.FindSlotByIntKey(0)
		second, _ := callback.FindSlotByIntKey(1)
		if first == nil || second == nil {
			return "Array", false
		}
		method, ok := second.ReadValue().(*data.StringValue)
		if !ok {
			return "Array", false
		}
		var class string
		switch receiver := first.ReadValue().(type) {
		case *data.StringValue:
			class = receiver.Value
		case *data.ThisValue:
			class = receiver.Class.GetName()
		case *data.ClassValue:
			class = receiver.Class.GetName()
		default:
			return "Array", false
		}
		return class + "::" + method.Value, true
	case *data.NullValue:
		return "", false
	default:
		if value == nil {
			return "", false
		}
		return value.AsString(), false
	}
}

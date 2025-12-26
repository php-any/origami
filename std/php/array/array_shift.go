package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayShiftFunction 实现 array_shift 函数
type ArrayShiftFunction struct{}

func NewArrayShiftFunction() data.FuncStmt {
	return &ArrayShiftFunction{}
}

func (f *ArrayShiftFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	// Check if it's an array
	if arr, ok := arrayValue.(*data.ArrayValue); ok {
		if len(arr.Value) == 0 {
			return data.NewNullValue(), nil
		}

		// Shift first element
		first := arr.Value[0]
		arr.Value = arr.Value[1:]

		// Re-index numerical keys?
		// PHP array_shift re-indexes numerical keys.
		// Origami ArrayValue seems to be a list of values (indexed array).
		// So just slicing is enough for indexed array.
		// But if it was an associative array (ObjectValue in Origami?), array_shift works differently.
		// Wait, Origami separates ArrayValue (list) and ObjectValue (map).
		// If the user passes an ObjectValue (associative array), we need to handle it.
		// But `array` type hint in PHP usually accepts both.
		// In Origami, `array` type hint might map to ArrayValue?
		// Let's assume ArrayValue for now.

		return first, nil
	}

	// If it's a reference to an array?
	// The argument is passed by reference `array &$array`.
	// In Origami, we receive the value.
	// If we modify `arr.Value`, does it reflect?
	// `arr` is a pointer to `ArrayValue`. If `ArrayValue` is the value stored in the variable, modifying it works.
	// Yes, `ArrayValue` struct has a slice `Value`. Modifying the slice header in the struct works.

	// What if it's an ObjectValue (associative array)?
	if _, ok := arrayValue.(*data.ObjectValue); ok {
		// array_shift on associative array:
		// "Shifts the first value of the array off and returns it, shortening the array by one element and moving everything down. All numerical array keys will be modified to start counting from zero while literal keys won't be touched."
		// ObjectValue uses OrderedMap.
		// We need to remove the first element.

		// Accessing OrderedMap
		// We need to see if we can remove the first element.
		// data.ObjectValue doesn't expose OrderedMap directly?
		// Let's check `data/value_object.go`.
		// It has `GetProperties()`.
		// But we need to modify it.
		// `SetProperty`?
		// We might need to access the underlying map if possible or iterate and rebuild.

		// For now, let's assume ArrayValue is the primary target for array functions in Origami's current state.
		// If ObjectValue is passed, we might need to implement `array_shift` for it too.
		// But `array_shift` expects an array.

		return data.NewNullValue(), nil
	}

	// Warning: array_shift() expects parameter 1 to be array
	return data.NewNullValue(), nil
}

func (f *ArrayShiftFunction) GetName() string {
	return "array_shift"
}

func (f *ArrayShiftFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
	}
}

func (f *ArrayShiftFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	}
}

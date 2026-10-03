package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayReplaceRecursiveFunction() data.FuncStmt {
	return &ArrayReplaceRecursiveFunction{}
}

type ArrayReplaceRecursiveFunction struct{}

func (f *ArrayReplaceRecursiveFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramsValue, _ := ctx.GetIndexValue(0)
	if paramsValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsArray, ok := paramsValue.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}

	paramsList := paramsArray.ToValueList()
	if len(paramsList) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	result := deepCopyPreserveKeys(paramsList[0])
	for _, replacement := range paramsList[1:] {
		result = recursiveReplaceByKeys(result, replacement)
	}

	return result, nil
}

func recursiveReplaceByKeys(base, replacement data.Value) data.Value {
	a, ok := base.(*data.ArrayValue)
	if !ok {
		return replacement
	}
	b, ok := replacement.(*data.ArrayValue)
	if !ok {
		return replacement
	}
	out := data.CloneArrayValue(a)
	for slots, i := b.View(), 0; i < slots.Len(); i++ {
		slot := slots.At(i)
		key := slot.PHPArrayKey(i)
		if existing, ok := out.LookupZValByStringKey(key.AsString()); ok && bothArraysOrObjects(existing.ReadValue(), slot.ReadValue()) {
			out.SetKey(key, recursiveReplaceByKeys(existing.ReadValue(), slot.ReadValue()))
		} else if slot.RefCount() > 0 {
			out.BindReference(key, slot)
		} else {
			out.SetKey(key, slot.ReadValue())
		}
	}
	return out
}

func bothArraysOrObjects(a, b data.Value) bool {
	_, left := a.(*data.ArrayValue)
	_, right := b.(*data.ArrayValue)
	return left && right
}

func deepCopyPreserveKeys(v data.Value) data.Value {
	switch val := v.(type) {

	case *data.ArrayValue:
		cloned := data.CloneArrayValue(val)
		for arraySlots103, arrayPosition103 := cloned.View(), 0; arrayPosition103 < arraySlots103.Len(); arrayPosition103++ {
			z := arraySlots103.At(arrayPosition103)
			if z != nil {
				z.StoreRaw(deepCopyPreserveKeys(z.ReadValue()))
			}
		}
		return cloned
	default:
		return v
	}
}

func (f *ArrayReplaceRecursiveFunction) GetName() string {
	return "array_replace_recursive"
}

var arrayReplaceRecursiveFunctionGetParams = []data.GetValue{
	node.NewParameters(nil, "arrays", 0, nil, nil),
}

func (f *ArrayReplaceRecursiveFunction) GetParams() []data.GetValue {
	return arrayReplaceRecursiveFunctionGetParams
}

var arrayReplaceRecursiveFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
}

func (f *ArrayReplaceRecursiveFunction) GetVariables() []data.Variable {
	return arrayReplaceRecursiveFunctionGetVariables
}

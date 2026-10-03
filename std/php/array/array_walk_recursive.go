package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayWalkRecursiveFunction 实现 array_walk_recursive
// array_walk_recursive(array|object &$array, callable $callback, mixed $arg = null): bool
// PHP：回调返回值忽略；只有 callback(&$value, ...) 才会改写数组元素。
type ArrayWalkRecursiveFunction struct{}

func NewArrayWalkRecursiveFunction() data.FuncStmt {
	return &ArrayWalkRecursiveFunction{}
}

func (fn *ArrayWalkRecursiveFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrZVal := ctx.GetIndexZVal(0)
	cbVal, _ := ctx.GetIndexValue(1)
	userdata, _ := ctx.GetIndexValue(2)

	if arrZVal == nil || cbVal == nil {
		return data.NewBoolValue(false), nil
	}
	data.CowSeparateZVal(arrZVal)

	if ctl := walkRecursive(ctx, cbVal, userdata, arrZVal.ReadValue()); ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(true), nil
}

func walkRecursive(ctx data.Context, cbVal, userdata, val data.Value) data.Control {
	switch v := val.(type) {
	case *data.ArrayValue:
		for arraySlots107, arrayPosition107 := v.View(), 0; arrayPosition107 < arraySlots107.Len(); arrayPosition107++ {
			z := arraySlots107.At(arrayPosition107)
			if z == nil {
				continue
			}
			if isNestedArrayValue(z.ReadValue()) {
				z.StoreRaw(separateNestedWalkArray(z.ReadValue()))
				if ctl := walkRecursive(ctx, cbVal, userdata, z.ReadValue()); ctl != nil {
					return ctl
				}
				continue
			}
			key := z.PHPArrayKey(arrayPosition107)
			if ctl := invokeWalkCallback(ctx, cbVal, z, key, userdata); ctl != nil {
				return ctl
			}
		}
		return nil

	default:
		return nil
	}
}

func separateNestedWalkArray(value data.Value) data.Value {
	switch array := value.(type) {
	case *data.ArrayValue:
		return data.CloneArrayValue(array)

	default:
		return value
	}
}

func (fn *ArrayWalkRecursiveFunction) GetName() string { return "array_walk_recursive" }

var arrayWalkRecursiveFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.NewBaseType("array")),
	node.NewParameter(nil, "callback", 1, nil, nil),
	node.NewParameter(nil, "arg", 2, nil, nil),
}

func (fn *ArrayWalkRecursiveFunction) GetParams() []data.GetValue {
	return arrayWalkRecursiveFunctionGetParams
}

var arrayWalkRecursiveFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	node.NewVariable(nil, "callback", 1, data.Mixed{}),
	node.NewVariable(nil, "arg", 2, data.Mixed{}),
}

func (fn *ArrayWalkRecursiveFunction) GetVariables() []data.Variable {
	return arrayWalkRecursiveFunctionGetVariables
}

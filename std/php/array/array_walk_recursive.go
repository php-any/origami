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

	if ctl := walkRecursive(ctx, cbVal, userdata, arrZVal.Value); ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(true), nil
}

func walkRecursive(ctx data.Context, cbVal, userdata, val data.Value) data.Control {
	switch v := val.(type) {
	case *data.ArrayValue:
		for _, z := range v.List {
			if z == nil {
				continue
			}
			if isNestedArrayValue(z.Value) {
				if ctl := walkRecursive(ctx, cbVal, userdata, z.Value); ctl != nil {
					return ctl
				}
				continue
			}
			key := data.NewIntValue(0)
			if z.Name != "" {
				if n, ok := data.ParseIntArrayKeyName(z.Name); ok {
					key = data.NewIntValue(n)
				} else {
					key = data.NewStringValue(z.Name)
				}
			}
			if ctl := invokeWalkCallback(ctx, cbVal, z, key, userdata); ctl != nil {
				return ctl
			}
		}
		return nil
	case *data.ObjectValue:
		keys := make([]string, 0)
		v.RangeProperties(func(key string, _ data.Value) bool {
			keys = append(keys, key)
			return true
		})
		for _, key := range keys {
			zv, ctl := v.GetZVal(key)
			if ctl != nil {
				return ctl
			}
			if zv == nil {
				continue
			}
			if isNestedArrayValue(zv.Value) {
				if ctl := walkRecursive(ctx, cbVal, userdata, zv.Value); ctl != nil {
					return ctl
				}
				continue
			}
			if ctl := invokeWalkCallback(ctx, cbVal, zv, data.NewStringValue(key), userdata); ctl != nil {
				return ctl
			}
		}
		return nil
	default:
		return nil
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

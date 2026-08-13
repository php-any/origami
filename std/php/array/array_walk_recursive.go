package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayWalkRecursiveFunction 实现 array_walk_recursive
// array_walk_recursive(array|object &$array, callable $callback, mixed $arg = null): bool
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

	_, ctl := walkRecursive(ctx, cbVal, userdata, arrZVal.Value)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(true), nil
}

func walkRecursive(ctx data.Context, cbVal, userdata, val data.Value) (data.Value, data.Control) {
	switch v := val.(type) {
	case *data.ArrayValue:
		for _, z := range v.List {
			if z == nil {
				continue
			}
			if isNestedArrayValue(z.Value) {
				_, ctl := walkRecursive(ctx, cbVal, userdata, z.Value)
				if ctl != nil {
					return nil, ctl
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
			args := []data.Value{z.Value, key}
			if userdata != nil {
				args = append(args, userdata)
			}
			ret, ctl := invokeCallback(ctx, cbVal, args)
			if ctl != nil {
				return nil, ctl
			}
			z.Value = ret
		}
		return v, nil
	case *data.ObjectValue:
		keys := make([]string, 0)
		v.RangeProperties(func(key string, _ data.Value) bool {
			keys = append(keys, key)
			return true
		})
		for _, key := range keys {
			propVal, ctl := v.GetProperty(key)
			if ctl != nil {
				return nil, ctl
			}
			if isNestedArrayValue(propVal) {
				_, ctl := walkRecursive(ctx, cbVal, userdata, propVal)
				if ctl != nil {
					return nil, ctl
				}
				continue
			}
			args := []data.Value{propVal, data.NewStringValue(key)}
			if userdata != nil {
				args = append(args, userdata)
			}
			ret, ctl := invokeCallback(ctx, cbVal, args)
			if ctl != nil {
				return nil, ctl
			}
			v.SetProperty(key, ret)
		}
		return v, nil
	default:
		return val, nil
	}
}

func (fn *ArrayWalkRecursiveFunction) GetName() string { return "array_walk_recursive" }

func (fn *ArrayWalkRecursiveFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "callback", 1, nil, nil),
		node.NewParameter(nil, "arg", 2, nil, nil),
	}
}

func (fn *ArrayWalkRecursiveFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "callback", 1, data.Mixed{}),
		node.NewVariable(nil, "arg", 2, data.Mixed{}),
	}
}

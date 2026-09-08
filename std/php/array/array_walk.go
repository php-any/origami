package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayWalkFunction() data.FuncStmt {
	return &ArrayWalkFunction{}
}

type ArrayWalkFunction struct{}

func (fn *ArrayWalkFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	cbVal, _ := ctx.GetIndexValue(1)
	if cbVal == nil {
		return data.NewBoolValue(false), nil
	}

	arrZVal := ctx.GetIndexZVal(0)
	if arrZVal == nil {
		return data.NewBoolValue(false), nil
	}

	userdata, _ := ctx.GetIndexValue(2)

	if ctl := walkFlat(ctx, cbVal, userdata, arrZVal.Value); ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(true), nil
}

func walkFlat(ctx data.Context, cbVal, userdata, val data.Value) data.Control {
	switch arr := val.(type) {
	case *data.ArrayValue:
		for i, z := range arr.List {
			if z == nil {
				continue
			}
			key := data.NewIntValue(i)
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
		arr.RangeProperties(func(key string, _ data.Value) bool {
			keys = append(keys, key)
			return true
		})
		for _, key := range keys {
			zv, ctl := arr.GetZVal(key)
			if ctl != nil {
				return ctl
			}
			if zv == nil {
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

func (fn *ArrayWalkFunction) GetName() string {
	return "array_walk"
}

func (fn *ArrayWalkFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "callback", 1, nil, nil),
		node.NewParameter(nil, "arg", 2, nil, nil),
	}
}

func (fn *ArrayWalkFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "callback", 1, data.Mixed{}),
		node.NewVariable(nil, "arg", 2, data.Mixed{}),
	}
}

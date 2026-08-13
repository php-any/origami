package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ArrayFlipFunction struct{}

func NewArrayFlipFunction() data.FuncStmt { return &ArrayFlipFunction{} }

func (f *ArrayFlipFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}
	arr, ok := v.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}
	flipped := make(map[string]int)
	for _, zv := range arr.List {
		key := ""
		if sv, ok := zv.Value.(data.AsString); ok {
			key = sv.AsString()
		} else if iv, ok := zv.Value.(data.AsInt); ok {
			if val, err := iv.AsInt(); err == nil {
				key = data.NewIntValue(val).AsString()
			}
		}
		if key != "" {
			flipped[key] = len(flipped)
		}
	}
	result := data.NewArrayValue([]data.Value{})
	for k := range flipped {
		zv := data.NewZVal(data.NewIntValue(flipped[k]))
		zv.Name = k
		result.(*data.ArrayValue).List = append(result.(*data.ArrayValue).List, zv)
	}
	return result, nil
}

func (f *ArrayFlipFunction) GetName() string { return "array_flip" }

func (f *ArrayFlipFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array"))}
}

func (f *ArrayFlipFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}
}

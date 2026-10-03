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

	pairs := make([][2]data.Value, 0)
	switch arr := v.(type) {
	case *data.ArrayValue:
		for arraySlots89, idx := arr.View(), 0; idx < arraySlots89.Len(); idx++ {
			zv := arraySlots89.At(idx)
			if zv == nil || zv.ReadValue() == nil {
				continue
			}
			oldKey := zv.PHPArrayKey(idx)
			pairs = append(pairs, [2]data.Value{zv.ReadValue(), oldKey})
		}

	default:
		return data.NewArrayValue([]data.Value{}), nil
	}

	result := data.NewArrayValue([]data.Value{}).(*data.ArrayValue)
	for _, p := range pairs {
		switch k := p[0].(type) {
		case *data.IntValue:
			result.SetIntKey(k.Value, p[1])
		case *data.StringValue:
			result.SetStringKey(k.Value, p[1])
		default:
			continue
		}
	}
	return result, nil
}

func (f *ArrayFlipFunction) GetName() string { return "array_flip" }

var arrayFlipFunctionGetParams = []data.GetValue{node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array"))}

func (f *ArrayFlipFunction) GetParams() []data.GetValue {
	return arrayFlipFunctionGetParams
}

var arrayFlipFunctionGetVariables = []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}

func (f *ArrayFlipFunction) GetVariables() []data.Variable {
	return arrayFlipFunctionGetVariables
}

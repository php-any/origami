package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayKeyLastFunction 实现 array_key_last
// array_key_last(array $array): int|string|null
type ArrayKeyLastFunction struct{}

func NewArrayKeyLastFunction() data.FuncStmt {
	return &ArrayKeyLastFunction{}
}

func (f *ArrayKeyLastFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewNullValue(), nil
	}

	switch v := val.(type) {
	case *data.ArrayValue:
		if len(v.List) == 0 {
			return data.NewNullValue(), nil
		}
		last := v.List[len(v.List)-1]
		if last != nil && last.Name != "" {
			if n, ok := data.ParseIntArrayKeyName(last.Name); ok {
				return data.NewIntValue(n), nil
			}
			return data.NewStringValue(last.Name), nil
		}
		return data.NewIntValue(len(v.List) - 1), nil
	case *data.ObjectValue:
		var lastKey string
		found := false
		v.RangeProperties(func(key string, _ data.Value) bool {
			lastKey = key
			found = true
			return true
		})
		if !found {
			return data.NewNullValue(), nil
		}
		return data.NewStringValue(lastKey), nil
	default:
		return data.NewNullValue(), nil
	}
}

func (f *ArrayKeyLastFunction) GetName() string { return "array_key_last" }

func (f *ArrayKeyLastFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "array", 0, nil, nil)}
}

func (f *ArrayKeyLastFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}
}

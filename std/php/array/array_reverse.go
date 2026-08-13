package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayReverseFunction 实现 array_reverse
func NewArrayReverseFunction() data.FuncStmt {
	return &ArrayReverseFunction{}
}

type ArrayReverseFunction struct{}

func (f *ArrayReverseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	preserveKeysVal, _ := ctx.GetIndexValue(1)
	if arrayValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	_ = preserveKeysVal // preserve_keys 暂未用于 ObjectValue 顺序

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		n := len(v.List)
		out := make([]data.Value, n)
		for i := 0; i < n; i++ {
			out[n-1-i] = v.List[i].Value
		}
		return data.NewArrayValue(out), nil
	case *data.ObjectValue:
		keys := make([]string, 0)
		v.RangeProperties(func(key string, _ data.Value) bool {
			keys = append(keys, key)
			return true
		})
		n := len(keys)
		out := make([]data.Value, n)
		for i := 0; i < n; i++ {
			val, _ := v.GetProperty(keys[n-1-i])
			out[i] = val
		}
		return data.NewArrayValue(out), nil
	default:
		return data.NewArrayValue([]data.Value{}), nil
	}
}

func (f *ArrayReverseFunction) GetName() string {
	return "array_reverse"
}

func (f *ArrayReverseFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "preserve_keys", 1, nil, nil),
	}
}

func (f *ArrayReverseFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "preserve_keys", 1, data.NewBaseType("bool")),
	}
}

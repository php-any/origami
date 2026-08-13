package array

import (
	"math/rand"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayRandFunction 实现 array_rand 函数
// 从数组中随机取出一个或多个键
type ArrayRandFunction struct{}

func NewArrayRandFunction() data.FuncStmt {
	return &ArrayRandFunction{}
}

func (f *ArrayRandFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	numReq := 1
	if nv, _ := ctx.GetIndexValue(1); nv != nil {
		if iv, ok := nv.(*data.IntValue); ok {
			n, _ := iv.AsInt()
			if n > 0 {
				numReq = n
			}
		}
	}

	switch v := arrayValue.(type) {
	case *data.ArrayValue:
		length := len(v.List)
		if length == 0 {
			return data.NewNullValue(), nil
		}

		// 收集所有键
		keys := make([]data.Value, 0, length)
		for i, zval := range v.List {
			if zval != nil && zval.Name != "" {
				keys = append(keys, data.NewStringValue(zval.Name))
			} else {
				keys = append(keys, data.NewIntValue(i))
			}
		}

		if numReq >= length {
			// 返回所有键（打乱顺序）
			for i := len(keys) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				keys[i], keys[j] = keys[j], keys[i]
			}
			return data.NewArrayValue(keys), nil
		}

		if numReq == 1 {
			return keys[rand.Intn(len(keys))], nil
		}

		// 随机选取 numReq 个
		selected := make([]data.Value, 0, numReq)
		used := make(map[int]bool)
		for len(selected) < numReq {
			idx := rand.Intn(length)
			if !used[idx] {
				used[idx] = true
				selected = append(selected, keys[idx])
			}
		}
		return data.NewArrayValue(selected), nil

	case *data.ObjectValue:
		props := v.GetProperties()
		length := len(props)
		if length == 0 {
			return data.NewNullValue(), nil
		}

		keys := make([]data.Value, 0, length)
		for k := range props {
			keys = append(keys, data.NewStringValue(k))
		}

		if numReq >= length {
			for i := len(keys) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				keys[i], keys[j] = keys[j], keys[i]
			}
			return data.NewArrayValue(keys), nil
		}

		if numReq == 1 {
			return keys[rand.Intn(len(keys))], nil
		}

		selected := make([]data.Value, 0, numReq)
		used := make(map[int]bool)
		for len(selected) < numReq {
			idx := rand.Intn(length)
			if !used[idx] {
				used[idx] = true
				selected = append(selected, keys[idx])
			}
		}
		return data.NewArrayValue(selected), nil
	}

	return data.NewNullValue(), nil
}

func (f *ArrayRandFunction) GetName() string {
	return "array_rand"
}

func (f *ArrayRandFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "num", 1, data.NewIntValue(1), data.Int{}),
	}
}

func (f *ArrayRandFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "num", 1, data.Int{}),
	}
}

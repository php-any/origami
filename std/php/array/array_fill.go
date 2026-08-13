package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayFillFunction 实现 array_fill 函数
// array_fill(int $start_index, int $count, mixed $value): array
type ArrayFillFunction struct{}

func NewArrayFillFunction() data.FuncStmt {
	return &ArrayFillFunction{}
}

func (f *ArrayFillFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	startVal, _ := ctx.GetIndexValue(0)
	countVal, _ := ctx.GetIndexValue(1)
	fillVal, _ := ctx.GetIndexValue(2)

	start := 0
	if iv, ok := startVal.(data.AsInt); ok {
		if v, err := iv.AsInt(); err == nil {
			start = v
		}
	}

	count := 0
	if iv, ok := countVal.(data.AsInt); ok {
		if v, err := iv.AsInt(); err == nil {
			count = v
		}
	}

	if count <= 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	var val data.Value
	if fillVal != nil {
		val = fillVal
	} else {
		val = data.NewNullValue()
	}

	list := make([]*data.ZVal, 0, count)
	if start == 0 {
		for i := 0; i < count; i++ {
			list = append(list, data.NewZVal(val))
		}
	} else {
		for i := 0; i < count; i++ {
			key := data.IntArrayKeyName(start + i)
			list = append(list, data.NewNamedZVal(key, val))
		}
	}

	return &data.ArrayValue{List: list}, nil
}

func (f *ArrayFillFunction) GetName() string {
	return "array_fill"
}

func (f *ArrayFillFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "start_index", 0, nil, nil),
		node.NewParameter(nil, "count", 1, nil, nil),
		node.NewParameter(nil, "value", 2, nil, nil),
	}
}

func (f *ArrayFillFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "start_index", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "count", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "value", 2, data.Mixed{}),
	}
}

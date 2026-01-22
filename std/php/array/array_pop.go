package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayPopFunction() data.FuncStmt {
	return &ArrayPopFunction{}
}

type ArrayPopFunction struct{}

func (f *ArrayPopFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 检查是否为数组引用
	arrayRef, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return data.NewNullValue(), nil
	}
	if len(arrayRef.List) == 0 {
		return data.NewNullValue(), nil
	}

	// 弹出最后一个元素
	lastElement := arrayRef.List[len(arrayRef.List)-1].Value
	arrayRef.List = arrayRef.List[:len(arrayRef.List)-1]

	return lastElement, nil
}

func (f *ArrayPopFunction) GetName() string {
	return "array_pop"
}

func (f *ArrayPopFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *ArrayPopFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}

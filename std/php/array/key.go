package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// KeyFunction 实现 key 函数
// 返回数组中当前元素的键名
type KeyFunction struct{}

func NewKeyFunction() data.FuncStmt {
	return &KeyFunction{}
}

func (f *KeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		return val.PointerKey(), nil

	case *data.ClassValue:
		if ctl := data.EmitPHPError(ctx, 8192, "key(): Calling key() on an object is deprecated", nil); ctl != nil {
			return nil, ctl
		}
		return val.ObjectValue.PointerKey(), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *KeyFunction) GetName() string {
	return "key"
}

var keyFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, data.Mixed{}),
}

func (f *KeyFunction) GetParams() []data.GetValue {
	return keyFunctionGetParams
}

var keyFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
}

func (f *KeyFunction) GetVariables() []data.Variable {
	return keyFunctionGetVariables
}

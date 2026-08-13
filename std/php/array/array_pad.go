package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayPadFunction() data.FuncStmt {
	return &ArrayPadFunction{}
}

type ArrayPadFunction struct{}

func (fn *ArrayPadFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrVal, _ := ctx.GetIndexValue(0)
	lengthVal, _ := ctx.GetIndexValue(1)
	padVal, _ := ctx.GetIndexValue(2)

	targetLen := 0
	if l, ok := lengthVal.(data.AsInt); ok {
		if v, err := l.AsInt(); err == nil {
			targetLen = v
		}
	}

	var source []*data.ZVal
	switch arr := arrVal.(type) {
	case *data.ArrayValue:
		source = arr.List
	case *data.ObjectValue:
		// 转换 ObjectValue 为简单列表
		result := make([]*data.ZVal, 0)
		arr.RangeProperties(func(key string, value data.Value) bool {
			result = append(result, data.NewZVal(value))
			return true
		})
		source = result
	default:
		source = []*data.ZVal{}
	}

	absLen := targetLen
	if absLen < 0 {
		absLen = -absLen
	}

	if absLen <= len(source) {
		// 如果目标长度小于等于当前长度，返回原数组
		if targetLen > 0 || targetLen < 0 && len(source) > absLen {
			return data.NewArrayValue(listToValues(source)), nil
		}
	}

	result := make([]*data.ZVal, 0, absLen)

	if targetLen > 0 {
		// 正数：padding 在右边
		for _, z := range source {
			result = append(result, z)
		}
		for i := len(source); i < targetLen; i++ {
			result = append(result, data.NewZVal(padVal))
		}
	} else {
		// 负数：padding 在左边
		padCount := absLen - len(source)
		for i := 0; i < padCount; i++ {
			result = append(result, data.NewZVal(padVal))
		}
		for _, z := range source {
			result = append(result, z)
		}
	}

	return data.NewArrayValue(listToValues(result)), nil
}

func listToValues(list []*data.ZVal) []data.Value {
	result := make([]data.Value, len(list))
	for i, z := range list {
		result[i] = z.Value
	}
	return result
}

func (fn *ArrayPadFunction) GetName() string {
	return "array_pad"
}

func (fn *ArrayPadFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "length", 1, nil, nil),
		node.NewParameter(nil, "value", 2, nil, nil),
	}
}

func (fn *ArrayPadFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "length", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "value", 2, data.Mixed{}),
	}
}
